package supertokens

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-go/supertokens"
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
			actualSession := supertokens.GetSessionFromRequest(r)
			if actualSession != nil {
				session := Session{
					actualSession: actualSession,
				}
				c.Set(sessionContext, &session)
			}
			c.Next()
		}, params...)
		handler(c.Writer, c.Request)
	}
}

// HandleErrorAndRespond if error handlers are provided, then uses those, else does default error handling depending on the type of error
func HandleErrorAndRespond(err error, c *gin.Context) {
	supertokens.HandleErrorAndRespond(err, c.Writer)
}
