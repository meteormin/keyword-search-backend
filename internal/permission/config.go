package permission

type Config struct {
	Name      string
	GroupId   uint
	Methods   []Method
	Resources []string
}

func NewPermissionsFromConfig(cfg []Config) []Permission {
	var permissions []Permission
	for _, p := range cfg {
		var actions []Action
		for _, resource := range p.Resources {
			actions = append(actions, NewAction(resource, p.Methods))
		}

		permissions = append(permissions, NewPermission(p.GroupId, p.Name, actions))
	}

	return permissions
}
