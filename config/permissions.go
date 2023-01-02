package config

type Permission struct {
	Name     string
	GroupId  uint
	Method   Methods
	Resource string
}

type Methods string

const (
	GET    Methods = "GET"
	POST   Methods = "POST"
	PUT    Methods = "PUT"
	PATCH  Methods = "PATCH"
	DELETE Methods = "DELETE"
)

func (m Methods) ToString() string {
	return string(m)
}

func getPermissions() []Permission {
	return []Permission{
		{
			Name:     "name",
			GroupId:  1,
			Method:   GET,
			Resource: "/users",
		},
	}
}
