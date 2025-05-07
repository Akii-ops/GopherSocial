package main

import (
	"backend/internal/auth"
	"backend/internal/db"
	"backend/internal/env"
	"backend/internal/mailer"
	"backend/internal/ratelimiter"
	"backend/internal/store"
	"backend/internal/store/cache"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title						GGM API
//	@host						localhost:9999
//	@version					0.0.2
//	@description				API for  API
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:         env.GetString("ENV", "development"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			ethereal: EtherealConfig{
				host:     env.GetString("ETHEREAL_EMAIL_HOST", ""),
				port:     env.GetInt("ETHEREAL_EMAIL_PORT", 25),
				username: env.GetString("ETHEREAL_EMAIL_USERNAME", ""),
				password: env.GetString("ETHEREAL_EMAIL_PASSWORD", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: "aki",
				pwd:  "aki",
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKENN_SECRET", "example"),
				exp:    time.Hour * 24 * 3,
				iss:    "GopherSocial",
			},
		},
		redis: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "loaclhost:6379"),
			password: env.GetString("REDIS_PWD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATE_LIMITER_REQUEST_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// logger

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	// Redis cache
	var redis_client *redis.Client

	if cfg.redis.enabled {
		redis_client = cache.NewRedisClient(cfg.redis.addr, cfg.redis.password, cfg.redis.db)
		logger.Info("redis cache connection established")
	}

	redis_storage := cache.NewRedisStorage(redis_client)

	store := store.NewStorage(db)

	mailer := mailer.NewEtherealMailer(cfg.mail.fromEmail,
		cfg.mail.ethereal.host,
		cfg.mail.ethereal.username,
		cfg.mail.ethereal.password,
		cfg.mail.ethereal.port,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(cfg.rateLimiter.RequestPerTimeFrame, cfg.rateLimiter.TimeFrame)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cachestore:    redis_storage,
		rateLimiter:   rateLimiter,
	}

	//

	mux := app.mount()

	logger.Fatal(app.run(mux))

}
