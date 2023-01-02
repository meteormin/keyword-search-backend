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
		if len(p.Methods) != len(p.Resources) {
			panic("resource length must be equals methods length")
		}

		var actions []Action
		for _, resource := range p.Resources {
			actions = append(actions, NewAction(p.Methods, resource))
		}

		permissions = append(permissions, NewPermission(p.GroupId, p.Name, actions))
	}

	return permissions
}
