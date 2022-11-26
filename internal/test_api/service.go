package test_api

type Service interface {
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) ServiceImpl {
	return ServiceImpl{repo: repo}
}
