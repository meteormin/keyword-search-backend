package users

type Service interface {
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) *ServiceStruct {
	return &ServiceStruct{repo: repo}
}
