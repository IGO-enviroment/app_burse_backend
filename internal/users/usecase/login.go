package users_usecase

import (
	"app_burse_backend/internal/app"
	types_result "app_burse_backend/internal/types/result"
	users_entity "app_burse_backend/internal/users/entity"
	users_repository "app_burse_backend/internal/users/repo"

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

func (u *LoginUsecase) Call() types_result.Result {
	user, err := u.repo.GetUserByEmail(u.entity.Email)
	if err != nil {
		return types_result.NewErrorResult(
			types_result.ErrorWithMsg(u.app.Locales().MustLocalize(&i18n.LocalizeConfig{MessageID: "errors.unknown_error"})),
		)
	}

	// Check if password matches
	if !user.CheckPassword(u.entity.Password, user.DigestPassword) {
		return types_result.NewErrorResult(
			types_result.ErrorWithMsg(u.app.Locales().MustLocalize(&i18n.LocalizeConfig{
				MessageID: "users.errors.login.invalid_credentials",
			})),
		)
	}

	return types_result.NewSuccessResult(user.ID)
}
