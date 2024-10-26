package feedback_usecase

import (
	"app_burse_backend/internal/errors"
	feedback_repository "app_burse_backend/internal/feedback/repo"
	feedback_validates "app_burse_backend/internal/feedback/validates"
	practice_repo "app_burse_backend/internal/practice/repo"
	profile_repo "app_burse_backend/internal/profile/repo"
	"app_burse_backend/internal/service"
)

type CreateEntity struct {
	EnterpriseId int `json:"enterpriseId"`

	ProfileId int `json:"profileId"`

	PracticeId int `json:"practiceId"`

	Rating  int    `json:"rating"`
	Comment string `json:"comment"`

	Recommendation string `json:"recommendation"`
}

type CreateService struct {
	feedbackRepo *feedback_repository.Repository
	practiceRepo *practice_repo.Repository
	profileRepo  *profile_repo.Repository
	entity       *CreateEntity
}

func NewUsecase(
	feedbackRepo *feedback_repository.Repository,
	practiceRepo *practice_repo.Repository,
	profileRepo *profile_repo.Repository,
	entity CreateEntity,
) *CreateService {
	return &CreateService{
		feedbackRepo: feedbackRepo,
		practiceRepo: practiceRepo,
		profileRepo:  profileRepo,
		entity:       &entity,
	}
}

// Предприятие создаёт отзыв по указанному студенту.
func (u *CreateService) CreateForStudentByEnterprise() error {
	err := u.validate()
	if err != nil {
		return service.NewErrorResult(err)
	}

	err = u.feedbackRepo.Create()
	if err != nil {
		return service.NewErrorResult(errors.NewErrorItem(errors.WithMessage("Ошибка при создании отзыва")))
	}

	return service.NewSuccessResult("Отзыв успешно создан")
}

func (u *CreateService) validate() error {
	// Проверка валидности входных данных
	if !feedback_validates.RaitingValidate(u.entity.Rating) {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "rating", Message: "Недопустимый рейтинг"},
				},
			),
		)
	}
	if !feedback_validates.CommentValidate(u.entity.Comment) {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "comment", Message: "Комментарий должен быть от 10 до 3000 символов"},
				},
			),
		)
	}
	if !feedback_validates.RecommendationsValidate(&u.entity.Recommendation) {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "recommendation", Message: "Максимальная длина рекомендаций - 2000 символов"},
				},
			),
		)
	}

	practice, err := u.practiceRepo.GetById(u.entity.PracticeId)
	if err != nil || practice == nil {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "practiceId", Message: "Нет такой практики"},
				},
			),
		)
	}

	if practice.EnterpriseId != u.entity.EnterpriseId {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "practiceId", Message: "Отзыв может быть создан только для практики вашего предприятия"},
				},
			),
		)
	}

	profile, err := u.profileRepo.GetById(u.entity.ProfileId)
	if err != nil || profile == nil {
		return errors.NewErrorItem(
			errors.WithErrors(
				[]errors.ErrorField{
					errors.ErrorField{Field: "profileId", Message: "Нет такого профиля"},
				},
			),
		)
	}

	return nil
}
