package bookmarks

type Service interface {
	All()
	Find(pk uint)
	Create()
	Update()
	Delete()
}

type ServiceStruct struct {
	repo Repository
}

func (s *ServiceStruct) All() {

}

func (s *ServiceStruct) Find() {

}

func (s *ServiceStruct) Create() {

}

func (s *ServiceStruct) Update() {

}

func (s *ServiceStruct) Delete() {

}
