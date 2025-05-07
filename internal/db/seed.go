package db

import (
	"backend/internal/store"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

var blogUsernames = []string{
	"Emily", "Liam", "Sophia", "Noah", "Olivia",
	"Ava", "Elijah", "Mia", "James", "Charlotte",
	"Benjamin", "Amelia", "Lucas", "Harper", "Mason",
	"Evelyn", "Ethan", "Abigail", "Logan", "Ella",
	"Alexander", "Grace", "Jacob", "Chloe", "Michael",
	"Penelope", "Daniel", "Lily", "Henry", "Aria",
	"Jackson", "Zoey", "Sebastian", "Riley", "Aiden",
	"Nora", "Matthew", "Hazel", "Samuel", "Ellie",
	"David", "Violet", "Joseph", "Luna", "Carter",
	"Stella", "Owen", "Lucy", "Wyatt", "Claire",
	"Gabriel", "Scarlett", "Christopher", "Hannah",
	"Anthony", "Zoe", "Dylan", "Elizabeth", "Isaac",
	"Addison", "Ryan", "Grace", "Nathan", "Layla",
	"Isaiah", "Natalie", "Adam", "Savannah", "Caleb",
	"Victoria", "Thomas", "Aurora", "Eli", "Bella",
	"Levi", "Nevaeh", "Aaron", "Stella", "Hunter",
}

var blogTitles = []string{
	// Original 20
	"AI Art Revolution: MidJourney Masterclass",
	"Building React 18 Performance Optimization Systems",
	"Metaverse Social: First Impressions of Horizon Worlds",
	"Decoding the Dark Forest Theory in 'The Three-Body Problem'",
	"Go Concurrency: Advanced Channel Patterns",
	"Vintage Loft Makeover: Industrial Chic Transformation",
	"Web3 Fundamentals: ENS Domain Investment Guide",
	"Coffee Science: Mastering Pour-Over Variables",
	"Flutter Cross-Platform: iOS/Android Production Patterns",
	"Carbon Neutral Era: ESG Investment Strategies",
	"Low-Key Photography: Dramatic Lighting Techniques",
	"Smart Home Automation: Advanced HomeKit Scenes",
	"Murder Mystery Design: Plot Twist Architecture",
	"Succulent Care Guide: Summer Survival Tactics",
	"Classical Deep Dive: Beethoven's Moonlight Sonata",
	"Capsule Wardrobe: Minimalist Fashion System",
	"Deep Learning Optimization: PyTorch Pro Tips",
	"Urban Exploration: Abandoned Factory Photography",
	"Biohacking 101: Quantified Self Experiments",
	"Indie Dev Journey: From MVP to Product Hunt",

	// New additions
	"Rust Memory Safety: Ownership Patterns Explained",
	"Kubernetes Cost Optimization: Cluster Autoscaling Deep Dive",
	"AR Development: RealityKit Spatial Anchors Implementation",
	"Quantum Computing Basics: Qubit Entanglement Demystified",
	"Serverless Security: AWS Lambda IAM Best Practices",
	"Digital Twin Architecture: IoT Sensor Fusion Techniques",
	"Ethical Hacking: SSRF Vulnerability Exploitation & Prevention",
	"FinTech Architecture: Payment Gateway Resilience Design",
	"3D Printing: Voronoi Structure Optimization for Lightweighting",
	"Climate Tech: Carbon Credit Verification Blockchain Solutions",
	"GameFi Economics: Play-to-Earn Tokenomics Analysis",
	"Cybersecurity: Zero Trust Network Access Implementation",
	"Biotech: CRISPR Gene Editing Workflow Automation",
	"Robotics: ROS2 Navigation Stack Customization",
	"DevOps: GitOps Pipeline Security Hardening",
	"Data Engineering: Delta Lake Schema Evolution Strategies",
	"UI/UX: Dark Mode Accessibility Considerations",
	"Space Tech: CubeSat Communication Protocols",
	"Edge Computing: WASM MicroRuntime Deployment",
	"Voice Tech: Alexa Skill Multi-Language Localization",
}

var blogContents = []string{
	// Original 20
	"Master MidJourney commands and parameters for precise artistic style generation through prompt engineering",
	"Optimize initial load performance with code splitting and lazy loading in React's concurrent rendering architecture",
	"Comparing social dynamics between Horizon Worlds and VRChat, with avatar creation walkthrough",
	"Analyzing game theory foundations of Dark Forest theory and potential Fermi paradox solutions",
	"Practical patterns for goroutine lifecycle management using select statements and context cancellation",
	"Documenting industrial element preservation in loft conversions with modern minimalist design principles",
	"Evaluating ENS valuation models and step-by-step domain trading on OpenSea",
	"Quantifying grind size/water temperature effects and building reproducible brew parameter matrices",
	"Implementing cross-platform state management solutions for platform-specific UI challenges",
	"Breaking down ESG metrics and analyzing renewable energy's evolving role in corporate ratings",
	"Single-light setup demonstrations using Rembrandt lighting for portrait drama",
	"Geofence-based automation configurations for energy-efficient smart home ecosystems",
	"Multi-threaded narrative structures using red herrings for mystery storytelling",
	"Species-specific soil mixtures and standardized summer watering/ventilation protocols",
	"Comparative analysis of tempo interpretations by renowned pianists in 3rd movement",
	"33-item core wardrobe system with seasonal combination algorithms",
	"30%+ inference speed gains through mixed precision training and knowledge distillation",
	"Industrial composition techniques and safety protocols for hazardous location shoots",
	"Cold exposure/IF protocols with HRV tracking for biohacking beginners",
	"Launch checklist covering ASO optimization and press kit preparation workflows",

	// New additions
	"Implementing Rust's borrow checker patterns for safe concurrent memory access in high-performance systems",
	"Cluster capacity planning using vertical pod autoscaling and spot instance integration strategies",
	"Implementing persistent AR experiences with plane detection and scene understanding in RealityKit",
	"Quantum circuit simulations demonstrating superdense coding and quantum teleportation protocols",
	"Least privilege principle implementation using service-linked roles and permission boundaries",
	"Time-series sensor data synchronization using Kalman filters for digital twin accuracy",
	"Burp Suite exploitation lab demonstrating SSRF chaining to cloud metadata endpoints",
	"Circuit breaker implementation for payment processing with automated failover testing",
	"Topology optimization algorithms for strength-to-weight ratio maximization in Fusion 360",
	"ERC-1155 smart contract development for transparent carbon offset tracking on Ethereum",
	"Analyzing Axie Infinity's dual-token model and sustainability challenges in P2E ecosystems",
	"ZTNA implementation guide using Cloudflare Access and mutual TLS authentication",
	"Automated guide RNA design pipelines using Benchling API and CRISPR library management",
	"Customizing ROS2 Nav2 stack for omnidirectional mobile robots in warehouse environments",
	"Supply chain security implementation with signed commits and binary authorization in Argo CD",
	"Schema migration patterns using Delta Lake time travel and merge operations",
	"Contrast ratio optimization and color perception testing for dark theme implementations",
	"LoRaWAN protocol implementation for ground station communication in CubeSat deployments",
	"WebAssembly System Interface (WASI) integration for edge computing function portability",
	"Language model switching techniques using Alexa Conversations multilingual skill toolkit",
}

var blogTags = [][]string{
	// Original 20
	{"AI Art", "Digital Art", "Tutorial", "AIGC"},
	{"Frontend", "React", "Optimization", "Web Dev"},
	{"Metaverse", "VR", "Social Media", "Tech Trends"},
	{"Sci-Fi", "Game Theory", "Physics", "Book Analysis"},
	{"Go", "Concurrency", "Systems Design", "Backend"},
	{"Interior Design", "Renovation", "Industrial", "Lifestyle"},
	{"Blockchain", "Web3", "NFTs", "Investing"},
	{"Coffee", "Brewing", "Specialty", "Lifestyle"},
	{"Mobile", "Cross-Platform", "Flutter", "UI/UX"},
	{"Sustainability", "ESG", "Carbon Neutral", "Finance"},
	{"Photography", "Portraits", "Lighting", "Low-Key"},
	{"Smart Home", "HomeKit", "IoT", "Automation"},
	{"Game Design", "Storytelling", "Social Games", "Entertainment"},
	{"Gardening", "Succulents", "Plant Care", "Botany"},
	{"Classical", "Piano", "Music Theory", "Beethoven"},
	{"Minimalism", "Fashion", "Lifestyle", "Organization"},
	{"Deep Learning", "PyTorch", "MLOps", "AI"},
	{"Urban", "Photography", "Industrial", "Documentary"},
	{"Health", "Biohacking", "Quantified Self", "Wellness"},
	{"Startups", "Product Design", "Entrepreneurship", "MVP"},

	// New additions
	{"Rust", "Memory Safety", "Systems Programming", "Concurrency"},
	{"Kubernetes", "Cloud", "Cost Optimization", "DevOps"},
	{"AR", "RealityKit", "Spatial Computing", "iOS Dev"},
	{"Quantum", "Physics", "Qubits", "Algorithms"},
	{"AWS", "Serverless", "Security", "IAM"},
	{"IoT", "Digital Twin", "Sensor Fusion", "Industry 4.0"},
	{"Cybersecurity", "Ethical Hacking", "Web Security", "SSRF"},
	{"FinTech", "Microservices", "Resilience", "Payments"},
	{"3D Printing", "Topology Optimization", "CAD", "Engineering"},
	{"Blockchain", "Climate Tech", "Carbon Credits", "Sustainability"},
	{"GameFi", "NFTs", "Tokenomics", "Web3"},
	{"Zero Trust", "Network Security", "ZTNA", "Cloudflare"},
	{"Biotech", "CRISPR", "Automation", "Bioinformatics"},
	{"Robotics", "ROS2", "Navigation", "Warehouse Tech"},
	{"GitOps", "CI/CD", "Security", "Supply Chain"},
	{"Data Engineering", "Delta Lake", "Spark", "Big Data"},
	{"Accessibility", "UI Design", "Dark Mode", "Color Theory"},
	{"Space Tech", "CubeSat", "LoRaWAN", "Aerospace"},
	{"WASM", "Edge Computing", "Containers", "Microservices"},
	{"Voice Tech", "Alexa", "NLP", "Localization"},
}

var blogComments = []string{
	// Original 20
	"The prompt engineering section is spot-on! Finally understand how to use the seed parameter properly.",
	"Are there real-world benchmarks for React 18's concurrent features?",
	"Cross-platform avatar compatibility would be a game-changer!",
	"Brilliant connection between Dark Forest and prisoner's dilemma!",
	"Context timeout handling is crucial for microservice resilience.",
	"Can we get more details on exposed ductwork finishing techniques?",
	"Would love to see more case studies on ENS investment risks.",
	"Brew variable matrix now permanently taped to my coffee station!",
	"Could this state management approach work for e-commerce apps?",
	"Any recommended ESG implementation guides for manufacturing?",
	"Need camera angle diagrams for the Rembrandt lighting demo!",
	"How precise is geofencing in high-rise apartment buildings?",
	"What's the ideal percentage of red herrings in mystery plots?",
	"Specific brand recommendations for inorganic soil components?",
	"Spectrogram comparisons add scientific rigor to music analysis!",
	"Is there open-source code for the wardrobe algorithm?",
	"How to choose teacher models for distillation effectively?",
	"Essential safety gear list saved me multiple times!",
	"What's the most accurate HRV monitoring device?",
	"Would you share your Product Hunt email outreach template?",

	// New additions
	"Lifetime annotations in Rust could use more practical examples",
	"Cluster autoscaler configuration thresholds need deeper explanation",
	"ARKit vs RealityKit performance comparison would be valuable",
	"Quantum circuit visualization tools recommendation needed",
	"Lambda cold start mitigation techniques missing in IAM context",
	"Digital twin accuracy metrics framework request",
	"SSRF prevention in serverless architectures needs expansion",
	"Payment gateway circuit breaker pattern implementation details?",
	"Voroni structure FEA simulation parameters unclear",
	"Carbon credit verification oracle design challenges?",
	"GameFi token inflation control mechanisms analysis?",
	"ZTNA vs VPN performance overhead comparison data?",
	"CRISPR library scale management best practices?",
	"ROS2 navigation stack customization tutorials?",
	"GitOps pipeline CVE scanning implementation guide?",
	"Delta Lake schema evolution conflict resolution strategies?",
	"Dark mode contrast verification tools recommendation?",
	"CubeSat deployment frequency coordination challenges?",
	"WASI filesystem access security implications?",
	"Alexa Conversations dialog management complexity analysis?",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(1000)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {

		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}

	}

	tx.Commit()

	posts := generatePosts(2000, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(10000, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("seeding complete")

}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)
	pwd := &store.Password{}
	pwd.Set("123123")
	for i := 0; i < num; i++ {

		users[i] = &store.User{
			Username: blogUsernames[i%len(blogUsernames)] + fmt.Sprintf("%d", i),
			Email:    blogUsernames[i%len(blogUsernames)] + fmt.Sprintf("%d", i) + "@example2.com",
			Role: store.Role{
				Name: "user",
			},

			Password: *pwd,
		}

	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Version: 0,
			Title:   blogTitles[rand.Intn(len(blogTitles))],
			Content: blogContents[rand.Intn(len(blogContents))],
			Tags:    blogTags[rand.Intn(len(blogTags))],
		}
	}

	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: blogComments[rand.Intn(len(blogComments))],
		}
	}

	return comments
}
