package session

type Session struct {
	Sub   string   `json:"sub"`
	Roles []string `json:"roles"`
}
