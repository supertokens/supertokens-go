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
	JwtSigningPublicKeyExpiryTime  uint64
	CookieSameSite                 string
	IDRefreshTokenPath             string
	SessionExpiredStatusCode       int8
}

var handshakeInfoInstantiated *handshakeInfo
var handshakeInfoOnce sync.Once

// GetHandshakeInfoInstance returns handshake info.
func GetHandshakeInfoInstance() *handshakeInfo {
	if handshakeInfoInstantiated == nil {
		handshakeInfoLock.Lock()
		if handshakeInfoInstantiated == nil {
			// TODO: fetch from querier
			handshakeInfoInstantiated = &handshakeInfo{}
		}
		handshakeInfoLock.Unlock()
	}
	return handshakeInfoInstantiated
}

var handshakeInfoLock sync.Mutex

func (info *handshakeInfo) UpdateJwtSigningPublicKeyInfo(newKey string, newExpiry uint64) {
	handshakeInfoLock.Lock()
	info.JwtSigningPublicKey = newKey
	info.JwtSigningPublicKeyExpiryTime = newExpiry
	handshakeInfoLock.Unlock()
}
