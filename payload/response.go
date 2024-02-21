package payload

type Response200 struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
