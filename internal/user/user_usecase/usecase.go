package user_usecase

import (
	"context"
	"errors"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"gofermart/internal/balance"
	"gofermart/internal/token"
	"gofermart/internal/token/token_models"
	"gofermart/internal/user"
	"gofermart/internal/user/user_models"
	"gofermart/pkg/errs"
	"gofermart/pkg/jwt"
	"gofermart/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo    user.Repository
	tokenRepo   token.Repository
	balanceRepo balance.Repository
	trManager   *manager.Manager
}

func NewUserUC(
	userRepo user.Repository,
	tokenRepo token.Repository,
	balanceRepo balance.Repository,
	trManager *manager.Manager,
) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		balanceRepo: balanceRepo,
		trManager:   trManager,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, req user_models.UserRegisterRequest) error {

	if err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		data, err := uc.userRepo.GetUserByLogin(ctx, req.Login)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if data != nil {
			return errs.HttpErrUserAlreadyExists
		}

		data = &user_models.UserStorageData{Login: req.Login}
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return err
		}
		data.Password = string(hashedPassword)

		userId, err := uc.userRepo.CreateUser(ctx, *data)
		if err != nil {
			return err
		}

		return uc.balanceRepo.CreateBalance(ctx, userId)
	}); err != nil {
		logger.Log.Error("rollback transaction: trManager.Do() failed", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (uc *UserUseCase) Login(ctx context.Context, req user_models.UserLoginRequest) (string, error) {
	var tokenValue string
	if err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		userData, err := uc.userRepo.GetUserByLogin(ctx, req.Login)
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Error("There is no user with such login")
			return errs.HttpErrUserIncorrectLogin
		}
		if err != nil {
			return err
		}

		if !validatePassword(req.Password, userData.Password) {
			logger.Log.Error("Password for this user is incorrect")
			return errs.HttpErrUserIncorrectLogin
		}

		newToken, err := uc.generateNewToken(ctx, userData.UserId)
		if err != nil {
			return err
		}

		err = uc.tokenRepo.CreateToken(ctx, *newToken)
		if err != nil {
			return err
		}

		tokenValue = newToken.Value
		return nil
	}); err != nil {
		return "", err
	}

	return tokenValue, nil
}

func (uc *UserUseCase) generateNewToken(ctx context.Context, userId int64) (*token_models.TokenStorageData, error) {
	key, err := uc.tokenRepo.GetSecretKey(ctx)
	if err != nil {
		return nil, err
	}

	value, err := jwt.GenerateTokenValue(userId, key)
	if err != nil {
		return nil, err
	}
	newToken := token_models.TokenStorageData{
		UserId: userId,
		Value:  value,
	}

	return &newToken, nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Error to hash password", zap.Error(err))
		return nil, err
	}

	return hashedPassword, nil
}

func validatePassword(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
