package request

type Login struct {
	Username string `json:"username" valid:"username,required"`
	Password string `json:"password" valid:"password,required"`
}
