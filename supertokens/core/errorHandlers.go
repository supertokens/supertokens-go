package core

import (
	"net/http"
	"sync"
)

type errorHandlers struct {
	OnTokenTheftDetectedErrorHandler func(sessionHandle string, userID string, response http.ResponseWriter)
	OnUnauthorisedErrorHandler       func(error, http.ResponseWriter)
	OnTryRefreshTokenErrorHandler    func(error, http.ResponseWriter)
	OnGeneralErrorHandler            func(error, http.ResponseWriter)
}

func defaultTokenTheftDetectedErrorHandler(sessionHandle string, userID string, w http.ResponseWriter) {
	handshakeInfo, handshakeInfoError := GetHandshakeInfoInstance()
	if handshakeInfoError != nil {
		GetErrorHandlersInstance().OnGeneralErrorHandler(handshakeInfoError, w)
		return
	}
	w.WriteHeader(handshakeInfo.SessionExpiredStatusCode)
	w.Write([]byte("token theft detected"))
	_, _ = RevokeSession(sessionHandle)
}

func defaultUnauthorisedErrorHandler(err error, w http.ResponseWriter) {
	handshakeInfo, handshakeInfoError := GetHandshakeInfoInstance()
	if handshakeInfoError != nil {
		GetErrorHandlersInstance().OnGeneralErrorHandler(handshakeInfoError, w)
		return
	}
	w.WriteHeader(handshakeInfo.SessionExpiredStatusCode)
	w.Write([]byte("unauthorised: " + err.Error()))
}

func defaultTryRefreshTokenErrorHandler(err error, w http.ResponseWriter) {
	handshakeInfo, handshakeInfoError := GetHandshakeInfoInstance()
	if handshakeInfoError != nil {
		GetErrorHandlersInstance().OnGeneralErrorHandler(handshakeInfoError, w)
		return
	}
	w.WriteHeader(handshakeInfo.SessionExpiredStatusCode)
	w.Write([]byte("try refresh token: " + err.Error()))
}

func defaultGeneralErrorHandler(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Internal error: " + err.Error()))
}

var errorHandlerInstantiated *errorHandlers
var errorHandlersOnce sync.Once

// GetErrorHandlersInstance returns all the error handlers.
func GetErrorHandlersInstance() *errorHandlers {
	errorHandlersOnce.Do(func() {
		errorHandlerInstantiated = &errorHandlers{
			OnTokenTheftDetectedErrorHandler: defaultTokenTheftDetectedErrorHandler,
			OnUnauthorisedErrorHandler:       defaultUnauthorisedErrorHandler,
			OnTryRefreshTokenErrorHandler:    defaultTryRefreshTokenErrorHandler,
			OnGeneralErrorHandler:            defaultGeneralErrorHandler,
		}
	})
	return errorHandlerInstantiated
}
