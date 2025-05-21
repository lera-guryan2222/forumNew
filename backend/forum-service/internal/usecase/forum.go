package usecase

type ForumUsecase interface {
	// TODO: описать методы usecase
}

type forumUsecase struct{}

func NewForumUsecase() ForumUsecase {
	return &forumUsecase{}
}
