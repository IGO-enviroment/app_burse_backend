package feedback_usecase

import (
	feedback_repository "app_burse_backend/internal/feedback/repo"
	feedback_validates "app_burse_backend/internal/feedback/validates"
	practice_repo "app_burse_backend/internal/practice/repo"
	profile_repo "app_burse_backend/internal/profile/repo"
	types_item "app_burse_backend/internal/types/error_item"
	types_result "app_burse_backend/internal/types/result"
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
	errItem := u.validate()
	if errItem != nil {
		return types_result.NewErrorResult(types_result.ErrorWithErrorItem(errItem))
	}

	err := u.feedbackRepo.Create()
	if err != nil {
		return types_result.NewErrorResult(types_result.ErrorWithMsg("Ошибка при создании отзыва"))
	}

	return types_result.NewSuccessResult("Отзыв успешно создан")
}

func (u *CreateService) validate() *types_item.ErrorItem {
	// Проверка валидности входных данных
	if !feedback_validates.RaitingValidate(u.entity.Rating) {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"rating", "Недопустимый рейтинг"}},
			),
		)
	}
	if !feedback_validates.CommentValidate(u.entity.Comment) {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"comment", "Комментарий должен быть от 10 до 3000 символов"}},
			),
		)
	}
	if !feedback_validates.RecommendationsValidate(&u.entity.Recommendation) {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"recommendation", "Максимальная длина рекомендаций - 2000 символов"}},
			),
		)
	}

	practice, err := u.practiceRepo.GetById(u.entity.PracticeId)
	if err != nil || practice == nil {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"practiceId", "Нет такой практики"}},
			),
		)
	}

	if practice.EnterpriseId != u.entity.EnterpriseId {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"practiceId", "Отзыв может быть создан только для практики вашего предприятия"}},
			),
		)
	}

	profile, err := u.profileRepo.GetById(u.entity.ProfileId)
	if err != nil || profile == nil {
		return types_item.NewErrorItem(
			types_item.WithErrors(
				[][]string{{"profileId", "Нет такого профиля"}},
			),
		)
	}

	return nil
}
