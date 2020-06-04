package core

import (
	"net/http"
	"sync"
)

type errorHandlers struct {
	OnTokenTheftDetectedErrorHandler func(string, string, http.ResponseWriter)
	OnUnauthorisedErrorHandler       func(error, http.ResponseWriter)
	OnTryRefreshTokenErrorHandler    func(error, http.ResponseWriter)
	OnGeneralErrorHandler            func(error, http.ResponseWriter)
}

func defaultTokenTheftDetectedErrorHandler(sessionHandle string, userID string, response http.ResponseWriter) {
	// TODO:
}

func defaultUnauthorisedErrorHandler(err error, response http.ResponseWriter) {
	// TODO:
}

func defaultTryRefreshTokenErrorHandler(err error, response http.ResponseWriter) {
	// TODO:
}

func defaultGeneralErrorHandler(err error, response http.ResponseWriter) {
	// TODO:
}

var instantiated *errorHandlers
var once sync.Once

// GetInstance returns all the error handlers.
func GetInstance() *errorHandlers {
	once.Do(func() {
		instantiated = &errorHandlers{
			OnTokenTheftDetectedErrorHandler: defaultTokenTheftDetectedErrorHandler,
			OnUnauthorisedErrorHandler:       defaultUnauthorisedErrorHandler,
			OnTryRefreshTokenErrorHandler:    defaultTryRefreshTokenErrorHandler,
			OnGeneralErrorHandler:            defaultGeneralErrorHandler,
		}
	})
	return instantiated
}
