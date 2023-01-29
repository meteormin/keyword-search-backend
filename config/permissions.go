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

// getPermissions
// 최초 기본 값 저장을 위한 설정입니다.
func getPermissions() []PermissionConfig {
	return []PermissionConfig{
		{
			Name:    "admin",
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
		{
			Name:    "owner",
			GroupId: 1,
			Methods: []PermissionMethod{GET, POST, PUT, PATCH, DELETE},
			Resources: []string{
				"/users",
				"/hosts",
				"/search",
				"/redirect",
			},
		},
		{
			Name:    "manager",
			GroupId: 1,
			Methods: []PermissionMethod{GET, POST, PUT, PATCH},
			Resources: []string{
				"/users",
				"/hosts",
				"/search",
				"/redirect",
				"/groups",
			},
		},
		{
			Name:    "member",
			GroupId: 1,
			Methods: []PermissionMethod{GET, POST},
			Resources: []string{
				"/hosts",
				"/search",
				"/redirect",
			},
		},
	}
}
