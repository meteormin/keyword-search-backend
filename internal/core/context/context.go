package context

type Key string

// context constants
// ctx.Locals() 메서드에서 주로 사용됨
const (
	Container    Key = "container"
	App          Key = "app"
	DB           Key = "db"
	Config       Key = "config"
	Logger       Key = "logger"
	AuthUser     Key = "authUser"
	JwtGenerator Key = "jwtGenerator"
	Permissions  Key = "permissions"
)
