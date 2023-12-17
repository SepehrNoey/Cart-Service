package request

type AuthUser struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
