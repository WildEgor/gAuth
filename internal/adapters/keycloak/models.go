package keycloak_adapter

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type KeycloakUserInfo struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type KeycloakErr struct {
	Code int    `json:"error"`
	Msg  string `json:"message"`
	Type string `json:"type"`
}
