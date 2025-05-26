package usecase

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
	"user-service/internal/entity"
	"user-service/pkg/jwt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, user entity.User, passwordHash []byte) (int64, error)
}

type UserProvider interface {
	GetUserByName(ctx context.Context, email string) (entity.User, error)
}

type UserService struct {
	userSaver    UserSaver
	userProvider UserProvider
	log          entity.Logger
	tokenTTL     time.Duration
}

func NewUserService(us UserSaver, up UserProvider, log entity.Logger, tokenTTL time.Duration) *UserService {
	return &UserService{
		userSaver:    us,
		userProvider: up,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (us *UserService) Registre(ctx context.Context, user entity.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(user.PassHash, bcrypt.DefaultCost)
	if err != nil {
		us.log.Error(ctx, "failed to hash password", "error", err)
		return 0, err
	}
	
	//TODO: remove log
	us.log.Info(ctx, "hashed password", "hashedPassword", string(hashedPassword))
	//TODO: remove log

	id, err := us.userSaver.SaveUser(ctx, user, hashedPassword)
	if err != nil {
		us.log.Error(ctx, "failed to create user", "userID", id, "error", err)
		return 0, err
	}

	us.log.Info(ctx, "created user", "userID", id)
	return id, nil
}

func (us *UserService) Login(ctx context.Context, userName string, password []byte) (string, error) {
	user, err := us.userProvider.GetUserByName(ctx, userName)
	if err != nil {
		us.log.Error(ctx, "failed to get user", "error", err)
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, password); err != nil {
		us.log.Error(ctx, "failed to compare password", "error", err)
		return "", err
	}
	us.log.Info(ctx, "user authenticated", "userID", user.ID)

	token, err := jwt.NewToken(user, us.tokenTTL)
	if err != nil {
		us.log.Error(ctx, "failed to create token", "error", err)
		return "", err
	}
	return token, nil
}
