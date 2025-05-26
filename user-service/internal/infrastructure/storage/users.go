package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"user-service/internal/entity"
	"user-service/pkg/errWrap"
)

func (s *Storage) SaveUser(ctx context.Context, user entity.User, passwordHash []byte) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error(ctx, "failed to begin transaction", "error", err)
		return 0, errWrap.WrapError(err)
	}
	defer tx.Rollback(ctx)

	var userId int64
	if err = s.pool.QueryRow(ctx,
		`INSERT INTO users (username, pass_hash) VALUES ($1, $2) RETURNING id`,
		user.Username, passwordHash).Scan(&userId); err != nil {
		s.log.Error(ctx, "failed to save user", "error", err)
		return 0, errWrap.WrapError(err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO balances (user_id, current, withdrawn) VALUES ($1, 0, 0)`,
		userId,
	)
	if err != nil {
		s.log.Error(ctx, "failed to create balance", "error", err)
		return 0, errWrap.WrapError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		s.log.Error(ctx, "failed to commit transaction", "error", err)
		return 0, errWrap.WrapError(err)
	}
	return userId, nil
}

func (s *Storage) GetUserByName(ctx context.Context, userName string) (entity.User, error) {
	query := `SELECT id, pass_hash
			  FROM users
			  WHERE username = $1`
	var user entity.User
	if err := s.pool.QueryRow(ctx, query, userName).Scan(
		&user.ID,
		&user.PassHash,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Error(ctx, "user not found", "userName", userName)
			return entity.User{}, errWrap.NewAppError(errWrap.ErrUnauthorized, "invalid credentials", err)
		}
		return entity.User{}, errWrap.WrapError(err)
	}
	user.Username = userName
	return user, nil
}
