package auth

type authError struct {
	ErrorMsg string `json:"errorMessage"`
}

type verifyEmailResponse struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Username string `json:"username"`
}

type refreshTokenResponse struct {
	Username string `json:"username"`
}
