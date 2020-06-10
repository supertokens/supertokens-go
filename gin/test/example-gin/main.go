package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/gin/supertokens"
)

var noOfTimesGetSessionCalledDuringTest int = 0
var noOfTimesRefreshCalledDuringTest int = 0

func main() {
	supertokens.Config("localhost:9000;")
	r := gin.Default()
	r.Any("/login", login)
	r.Any("/testUserConfig", testUserConfig)
	r.Any("/multipleInterceptors", multipleInterceptors)
	r.Any("/", supertokensMiddleware(), defaultHandler)
	r.Any("/beforeeach", beforeeach)
	r.Any("/testing", testing)
	r.Any("/logout", supertokensMiddleware(), logout)
	r.Any("/revokeAll", supertokensMiddleware(), revokeAll)
	r.Any("/refresh", supertokensMiddleware(), refresh)
	r.Any("/refreshCalledTime", refreshCalledTime)
	r.Any("/getSessionCalledTime", getSessionCalledTime)
	r.Any("/ping", ping)
	r.Any("/testHeader", testHeader)
	r.Any("/checkDeviceInfo", checkDeviceInfo)
	r.Any("/checkAllowCredentials", checkAllowCredentials)
	r.Any("/testError", testError)
	r.Any("/index.html", index)
	r.Any("/fail", fail)
	r.Any("/update-jwt", supertokensMiddleware(), updateJwt)
	supertokens.OnTryRefreshToken(customOnTryRefreshTokenError)
	supertokens.OnUnauthorized(customOnUnauthorizedError)
	supertokens.OnGeneralError(customOnGeneralError)
	r.Run("0.0.0.0:8080")
}

func fail(c *gin.Context) {
	w := c.Writer
	w.WriteHeader(404)
	w.Write([]byte(""))
}

func index(c *gin.Context) {
	w := c.Writer
	dat, _ := ioutil.ReadFile("./static/index.html")
	w.Header().Set("Content-Type", "text/html")
	w.Write(dat)
}

func options(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Headers", "content-type")
	response.Header().Set("Access-Control-Allow-Methods", "*")
	supertokens.SetRelevantHeadersForOptionsAPI(response)
	response.Write([]byte(""))
}

func login(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}

	var body map[string]interface{}
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		response.Write([]byte("error when parsing body"))
		return
	}
	userID := body["userId"].(string)
	_, err = supertokens.CreateNewSession(response, userID)

	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
		return
	}
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte(userID))

}

func testUserConfig(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	response.Write([]byte(""))

}
func multipleInterceptors(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	interceptorheader2 := request.Header.Get("interceptorheader2")
	interceptorheader1 := request.Header.Get("interceptorheader1")

	var resp string
	if interceptorheader2 != "" && interceptorheader1 != "" {
		resp = "success"
	} else {
		resp = "failure"
	}
	response.Write([]byte(resp))
}

func defaultHandler(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	noOfTimesGetSessionCalledDuringTest++
	value := c.MustGet(supertokens.GinContext)
	session := value.(supertokens.Session)
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte(session.GetUserID()))
}

func updateJwt(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
	} else if request.Method == "GET" {
		response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		session := c.MustGet(supertokens.GinContext).(supertokens.Session)
		json.NewEncoder(response).Encode(session.GetJWTPayload())
	} else if request.Method == "POST" {
		var body map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			response.Write([]byte("error when parsing the body"))
			return
		}
		session := c.MustGet(supertokens.GinContext).(supertokens.Session)
		session.UpdateJWTPayload(body)
		response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		json.NewEncoder(response).Encode(session.GetJWTPayload())
	} else {
		response.Write([]byte("incorrect Method, requires POST or GET"))
	}
}

func beforeeach(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	noOfTimesRefreshCalledDuringTest = 0
	noOfTimesGetSessionCalledDuringTest = 0
	core.ResetHandshakeInfo()
	response.Write([]byte(""))
}

func testing(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	}
	value := request.Header.Get("testing")
	if value != "" {
		response.Header().Set("testing", value)
	}
	response.Write([]byte("success"))
}

func logout(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}

	value := c.MustGet(supertokens.GinContext)
	session := value.(supertokens.Session)
	err := session.RevokeSession()
	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
		return
	}
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte("success"))

}

func revokeAll(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	value := c.MustGet(supertokens.GinContext)
	session := value.(supertokens.Session)
	userID := session.GetUserID()
	supertokens.RevokeAllSessionsForUser(userID)
	response.Write([]byte("success"))
}

func refresh(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	noOfTimesRefreshCalledDuringTest++
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte("refresh success"))
}

func refreshCalledTime(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Write([]byte(strconv.Itoa(noOfTimesRefreshCalledDuringTest)))
}

func getSessionCalledTime(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Write([]byte(strconv.Itoa(noOfTimesGetSessionCalledDuringTest)))
}

func ping(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Write([]byte(""))
}

func testHeader(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	testheader := request.Header.Get("st-custom-header")
	success := testheader != ""
	json.NewEncoder(response).Encode(map[string]interface{}{
		"success": success,
	})
}

func checkDeviceInfo(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	sdkName := request.Header.Get("supertokens-sdk-name")
	sdkVersion := request.Header.Get("supertokens-sdk-version")
	response.Write([]byte(strconv.FormatBool(sdkName == "website" && sdkVersion != "")))
}

func checkAllowCredentials(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	response.Write([]byte(strconv.FormatBool(request.Header.Get("allow-credentials") != "")))
}

func testError(c *gin.Context) {
	response := c.Writer
	request := c.Request
	if request.Method == "OPTIONS" {
		options(response, request)
		return
	} else if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("test error message"))
}

func customOnTryRefreshTokenError(err error, response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.WriteHeader(440)
	response.Write([]byte(""))

}

func customOnUnauthorizedError(err error, response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.WriteHeader(440)
	response.Write([]byte(""))
}

func customOnGeneralError(err error, response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("Something went wrong"))
}
