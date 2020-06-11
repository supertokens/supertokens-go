package core

import "sync"

type handshakeInfo struct {
	JwtSigningPublicKey            string
	CookieDomain                   string
	CookieSecure                   bool
	AccessTokenPath                string
	RefreshTokenPath               string
	EnableAntiCsrf                 bool
	AccessTokenBlacklistingEnabled bool
	JwtSigningPublicKeyExpiryTime  int64
	CookieSameSite                 string
	IDRefreshTokenPath             string
	SessionExpiredStatusCode       int
}

var handshakeInfoInstantiated *handshakeInfo

// GetHandshakeInfoInstance returns handshake info.
func GetHandshakeInfoInstance() (*handshakeInfo, error) {
	if handshakeInfoInstantiated == nil {
		handshakeInfoLock.Lock()
		defer handshakeInfoLock.Unlock()
		if handshakeInfoInstantiated == nil {
			response, err := GetQuerierInstance().SendPostRequest("handshake", "/handshake", map[string]interface{}{})
			if err != nil {
				return nil, err
			}
			handshakeInfoInstantiated = &handshakeInfo{
				JwtSigningPublicKey:            response["jwtSigningPublicKey"].(string),
				CookieDomain:                   response["cookieDomain"].(string),
				CookieSecure:                   response["cookieSecure"].(bool),
				AccessTokenPath:                response["accessTokenPath"].(string),
				RefreshTokenPath:               response["refreshTokenPath"].(string),
				EnableAntiCsrf:                 response["enableAntiCsrf"].(bool),
				AccessTokenBlacklistingEnabled: response["accessTokenBlacklistingEnabled"].(bool),
				JwtSigningPublicKeyExpiryTime:  int64(response["jwtSigningPublicKeyExpiryTime"].(float64)),
				CookieSameSite:                 response["cookieSameSite"].(string),
				IDRefreshTokenPath:             response["idRefreshTokenPath"].(string),
				SessionExpiredStatusCode:       int(response["sessionExpiredStatusCode"].(float64)),
			}
		}
	}
	return handshakeInfoInstantiated, nil
}

var handshakeInfoLock sync.Mutex

func (info *handshakeInfo) UpdateJwtSigningPublicKeyInfo(newKey string, newExpiry int64) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	info.JwtSigningPublicKey = newKey
	info.JwtSigningPublicKeyExpiryTime = newExpiry
}

// ResetHandshakeInfo to be used for testing only
func ResetHandshakeInfo() {
	handshakeInfoInstantiated = nil
}
