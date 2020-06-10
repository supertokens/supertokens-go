package supertokens

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

// Middleware for verifying and refreshing session.
func Middleware(condition ...bool) func(*gin.Context) {
	return func(c *gin.Context) {
		var params = []interface{}{}
		if len(condition) == 1 {
			params = append(params, condition[0])
		} else {
			params = append(params, nil)
		}
		params = append(params, func(err error, w http.ResponseWriter) {
			c.Abort()
			supertokens.HandleErrorAndRespond(err, w)
		})
		handler := supertokens.Middleware(func(w http.ResponseWriter, r *http.Request) {
			session := supertokens.GetSessionFromRequest(r)
			if session != nil {
				c.Set(sessionContext, session)
			}
			c.Next()
		}, params...)
		handler(c.Writer, c.Request)
	}
}

// HandleErrorAndRespond if error handlers are provided, then uses those, else does default error handling depending on the type of error
func HandleErrorAndRespond(err error, c *gin.Context) {
	errorHandlers := core.GetErrorHandlersInstance()
	if errors.IsUnauthorizedError(err) {
		errorHandlers.OnUnauthorizedErrorHandler(err, c.Writer)
	} else if errors.IsTryRefreshTokenError(err) {
		errorHandlers.OnTryRefreshTokenErrorHandler(err, c.Writer)
	} else if errors.IsTokenTheftDetectedError(err) {
		actualError := err.(errors.TokenTheftDetectedError)
		errorHandlers.OnTokenTheftDetectedErrorHandler(actualError.SessionHandle, actualError.UserID, c.Writer)
	} else {
		errorHandlers.OnGeneralErrorHandler(err, c.Writer)
	}
}
