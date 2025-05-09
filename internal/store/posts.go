package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts 
		WHERE ID = $1;
	
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	ret := &Post{}
	err := s.db.QueryRowContext(ctx, query,
		id).Scan(
		&ret.ID,
		&ret.UserID,
		&ret.Title,
		&ret.Content,
		&ret.CreatedAt,
		&ret.UpdatedAt,
		pq.Array(&ret.Tags),
		&ret.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return ret, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
		delete from posts where posts.id = $1;
		
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil

}

func (s *PostStore) Update(ctx context.Context, post *Post) error {

	query := `
		UPDATE posts

		SET title = $1, content = $2, version = version + 1 

		WHERE id = $3 AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := s.db.QueryRowContext(ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err

		}
	}

	return nil

}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {

	query := `
		select p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
		u.username,
		count(c.id) as comments_count

		from posts p
		left join users u on  p.user_id = u.id
		left join comments c on c.post_id=p.id
		join followers f on f.follower_id = p.user_id or p.user_id = $1

		where  f.user_id = $1 or p.user_id = $1 and
			(p.title ilike '%'  || $4 || '%' or p.content ilike '%'  || $4 || '%')
			and
			(p.tags @> $5 or $5 =  '{}')

		group by p.id, u.username

		order by p.created_at ` + fq.Sort + `
		limit $2
		offset $3
	`

	var postMetadata = []PostWithMetadata{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	defer rows.Close()

	for rows.Next() {

		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		postMetadata = append(postMetadata, p)

	}

	return postMetadata, nil
}
