package config

import "github.com/miniyus/keyword-search-backend/internal/permission"

// getPermissions
// 최초 기본 값 저장을 위한 설정입니다.
func permissionConfig() []permission.Config {
	return []permission.Config{
		{
			Name:    "admin",
			GroupId: 1,
			Methods: []permission.Method{permission.GET, permission.PATCH, permission.DELETE, permission.POST, permission.PUT},
			Resources: []string{
				"/users",
				"/hosts",
				"/search",
				"/redirect",
				"/groups",
				"/worker",
			},
		},
		{
			Name:    "owner",
			GroupId: 1,
			Methods: []permission.Method{permission.GET, permission.POST, permission.PUT, permission.PATCH, permission.DELETE},
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
			Methods: []permission.Method{permission.GET, permission.POST, permission.PUT, permission.PATCH},
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
			Methods: []permission.Method{permission.GET, permission.POST},
			Resources: []string{
				"/hosts",
				"/search",
				"/redirect",
			},
		},
	}
}
