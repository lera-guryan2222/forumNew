package usecase

type PostUsecase interface {
	// TODO: описать методы usecase
}

type postUsecase struct{}

func NewPostUsecase() PostUsecase {
	return &postUsecase{}
}
