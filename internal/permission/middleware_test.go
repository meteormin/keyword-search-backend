package permission_test

import (
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/utils"
	"strings"
	"testing"
)

var hasPerm = []permission.Permission{
	{
		Actions: []permission.Action{
			{
				Methods: []permission.Method{
					"GET",
				},
				Resource: "/test",
			},
		},
		GroupId: 1,
		Name:    "TEST",
	},
}

func TestCheckPermission(t *testing.T) {
	pass := false
	method := "GET"
	utils.NewCollection(hasPerm).For(func(perm permission.Permission, i int) {
		utils.NewCollection(perm.Actions).For(func(action permission.Action, j int) {
			routePath := "/api/test/1"
			if strings.Contains(routePath, action.Resource) {

				if method == "OPTION" {
					method = "GET"
				}

				filtered := utils.NewCollection(action.Methods).Filter(func(v permission.Method, i int) bool {
					return string(v) == method
				})

				if filtered.Count() != 0 {
					pass = true
				}
			}
		})
	})

	if !pass {
		t.Error(pass)
	}
}
