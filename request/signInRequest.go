package request

type SignIn struct {
	UserNameOrEmail string `json:"userNameOrEmail"`
	Password        string `json:"password"`
}
