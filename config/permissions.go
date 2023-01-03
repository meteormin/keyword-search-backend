package config

type PermissionConfig struct {
	Name      string
	GroupId   uint
	Methods   []PermissionMethod
	Resources []string
}

type PermissionMethod string

const (
	GET    PermissionMethod = "GET"
	POST   PermissionMethod = "POST"
	PUT    PermissionMethod = "PUT"
	PATCH  PermissionMethod = "PATCH"
	DELETE PermissionMethod = "DELETE"
)

func getPermissions() []PermissionConfig {
	return []PermissionConfig{
		{
			Name:    "ADMIN",
			GroupId: 1,
			Methods: []PermissionMethod{GET, POST, PUT, PATCH, DELETE},
			Resources: []string{
				"/users",
				"/hosts",
				"/search",
				"/redirect",
				"/groups",
			},
		},
	}
}
