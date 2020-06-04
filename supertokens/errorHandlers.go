package supertokens

import (
	"net/http"
	"sync"
)

type errorHandlers struct {
	onTokenTheftDetectedErrorHandler onTokenTheftDetectedErrorHandlerType
	onUnauthorisedErrorHandler       onUnauthorisedErrorHandlerType
	onTryRefreshTokenErrorHandler    onTryRefreshTokenErrorHandlerType
	onGeneralErrorHandler            onGeneralErrorHandlerType
}

type onTokenTheftDetectedErrorHandlerType func(string, string, http.ResponseWriter)
type onUnauthorisedErrorHandlerType func(error, http.ResponseWriter)
type onTryRefreshTokenErrorHandlerType func(error, http.ResponseWriter)
type onGeneralErrorHandlerType func(error, http.ResponseWriter)

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
			onTokenTheftDetectedErrorHandler: defaultTokenTheftDetectedErrorHandler,
			onUnauthorisedErrorHandler:       defaultUnauthorisedErrorHandler,
			onTryRefreshTokenErrorHandler:    defaultTryRefreshTokenErrorHandler,
			onGeneralErrorHandler:            defaultGeneralErrorHandler,
		}
	})
	return instantiated
}

// OnTokenTheftDetected function to override default behaviour of handling token thefts
func OnTokenTheftDetected(handler onTokenTheftDetectedErrorHandlerType) {
	(*GetInstance()).onTokenTheftDetectedErrorHandler = handler
}

// OnUnauthorised function to override default behaviour of handling unauthorised error
func OnUnauthorised(handler onUnauthorisedErrorHandlerType) {
	(*GetInstance()).onUnauthorisedErrorHandler = handler
}

// OnTryRefreshToken function to override default behaviour of handling try refresh token errors
func OnTryRefreshToken(handler onTryRefreshTokenErrorHandlerType) {
	(*GetInstance()).onTryRefreshTokenErrorHandler = handler
}

// OnGeneralError function to override default behaviour of handling general errors
func OnGeneralError(handler onGeneralErrorHandlerType) {
	(*GetInstance()).onGeneralErrorHandler = handler
}
