package auth

// RequestRegister representa la solicitud de registro
type RequestRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RequestLogin representa la solicitud de login
type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"` // true -> "Recordarme" (expiración larga)
}

// ResponseUser representa la respuesta de usuario
type ResponseUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// ResponseMessage representa una respuesta simple con mensaje
type ResponseMessage struct {
	Message string `json:"message"`
}
