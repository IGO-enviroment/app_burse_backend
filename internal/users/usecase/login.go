package users_usecase

import (
	"app_burse_backend/internal/app"
	"app_burse_backend/internal/service"
	users_entity "app_burse_backend/internal/users/entity"
	users_repository "app_burse_backend/internal/users/repo"
	tokenservice "app_burse_backend/pkg/token_service"
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type LoginUsecase struct {
	entity users_entity.LoginEntity
	repo   users_repository.Repository

	app app.AppContext
}

func NewLoginUsecase(repo users_repository.Repository, entity users_entity.LoginEntity, app app.AppContext) *LoginUsecase {
	return &LoginUsecase{repo: repo, entity: entity, app: app}
}

func (u *LoginUsecase) Call() service.Result {
	user, err := u.repo.GetUserByEmail(u.entity.Email)
	if err != nil {
		return service.NewErrorResult(
			errors.New(u.app.Locales().MustLocalize(&i18n.LocalizeConfig{MessageID: "errors.unknown_error"})),
		)
	}

	// Check if password matches
	if !user.CheckPassword(u.entity.Password, user.DigestPassword) {
		return service.NewErrorResult(
			errors.New(u.app.Locales().MustLocalize(&i18n.LocalizeConfig{
				MessageID: "users.errors.login.invalid_credentials",
			})),
		)
	}

	token, err := u.generateHeaderToken(user.ID)
	if err != nil {
		return service.NewErrorResult(
			errors.New(u.app.Locales().MustLocalize(&i18n.LocalizeConfig{MessageID: "errors.unknown_error"})),
		)
	}

	return service.NewSuccessResult(token)
}

func (u *LoginUsecase) generateHeaderToken(userID int) (string, error) {
	сfg := u.app.Configs().Web
	return tokenservice.NewTokenService(сfg.TokenSecret, сfg.TokenExpiration).GenerateToken(userID)
}
