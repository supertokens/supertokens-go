package core

type accessTokenInfoStruct struct {
	sessionHandle           string
	userID                  string
	refreshTokenHash1       string
	parentRefreshTokenHash1 *string
	userData                map[string]interface{}
	antiCsrfToken           *string
	expiryTime              int64
	timeCreated             int64
}

func getInfoFromAccessToken(token string, jwtSigningPublicKey string, doAntiCsrfCheck bool) (accessTokenInfoStruct, error) {
	// TODO:
	return accessTokenInfoStruct{}, nil
}
