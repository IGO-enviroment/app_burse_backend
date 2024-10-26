package usecase

type Usecase struct {
}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Create(email string) error {
	// Создание пользователя
	// from domain.User{Email: email}
	// to raitingRecord
	// validate feedback
	// save to db

	return nil
}
