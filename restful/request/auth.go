package request

// swagger:model RefreshTokenRequest
type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type EmailLoginRequest struct {

	// required: true
	Email string `json:"email"`
	// required: true
	Password string `json:"password"`
}

// swagger:parameters  loginUser
type Req struct {
	// desc
	// in: body
	// required: true
	Body EmailLoginRequest `json:",inline"`
}
