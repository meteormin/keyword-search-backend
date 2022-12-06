package hosts

type CreateHost struct {
	UserId      uint   `json:"user_id"`
	Host        string `json:"host" validate:"required"`
	Subject     string `json:"subject" validate:"required,uri"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required,dir"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type UpdateHost struct {
	UserId      uint   `json:"user_id"`
	Host        string `json:"host" validate:"required"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required,dir"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type HostResponse struct {
	Id          uint   `json:"id"`
	UserId      uint   `json:"user_id"`
	Host        string `json:"host"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Publish     bool   `json:"publish"`
}
