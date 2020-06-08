package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

var noOfTimesGetSessionCalledDuringTest int = 0
var noOfTimesRefreshCalledDuringTest int = 0

func main() {
	// TODO: create API according to https://github.com/supertokens/supertokens-javalin/blob/master/Example/src/main/java/example/Main.java
	supertokens.Config("localhost:9000;")
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/testUserConfig", testUserConfig)
	http.HandleFunc("/multipleInterceptors", multipleInterceptors)

}

func options(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Headers", "content-type")
	response.Header().Set("Access-Control-Allow-Methods", "*")
	supertokens.SetRelevantHeadersForOptionsAPI(response)
	response.Write([]byte(""))
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	var body map[string]interface{}
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		//TODO: error when parsing body
	}
	userID := body["userId"].(string)
	_, err = supertokens.CreateNewSession(response, userID)

	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
	}
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte(userID))

}

func testUserConfig(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(""))

}
func multipleInterceptors(response http.ResponseWriter, request *http.Request) {
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

func defaultHandler(response http.ResponseWriter, request *http.Request) {
	noOfTimesGetSessionCalledDuringTest++
	value := request.Context().Value(supertokens.SessionContext)
	session := value.(supertokens.Session)
	session.GetUserID()
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte("success"))
}

func beforeeach(response http.ResponseWriter, request *http.Request) {
	noOfTimesRefreshCalledDuringTest = 0
	noOfTimesGetSessionCalledDuringTest = 0
	core.ResetHandshakeInfo()
	response.Write([]byte(""))
}

func testing(response http.ResponseWriter, request *http.Request) {
	value := request.Header.Get("testing")
	if value != "" {
		response.Header().Set("testing", value)
	}
	response.Write([]byte("success"))
}

func logout(response http.ResponseWriter, request *http.Request) {
	value := request.Context().Value(supertokens.SessionContext)
	session := value.(supertokens.Session)
	err := session.RevokeSession()
	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
	}
	response.Write([]byte("success"))

}

func revokeAll(response http.ResponseWriter, request *http.Request) {
	value := request.Context().Value(supertokens.SessionContext)
	session := value.(supertokens.Session)
	userID := session.GetUserID()
	supertokens.RevokeAllSessionsForUser(userID)
	response.Write([]byte("success"))
}

func refresh(response http.ResponseWriter, request *http.Request) {
	noOfTimesRefreshCalledDuringTest++
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Write([]byte("refresh success"))
}

func refreshCalledTime(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Write([]byte(strconv.Itoa(noOfTimesRefreshCalledDuringTest)))
}

func getSessionCalledTime(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Write([]byte(strconv.Itoa(noOfTimesGetSessionCalledDuringTest)))
}

func getPackageVersion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	response.Write([]byte("4.1.3"))
}

func ping(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(""))
}

func testHeader(response http.ResponseWriter, request *http.Request) {
	testheader := request.Header.Get("st-custom-header")
	success := testheader != ""
	json.NewEncoder(response).Encode(map[string]interface{}{
		"success": success,
	})
}

func checkDeviceInfo(response http.ResponseWriter, request *http.Request) {
	sdkName := request.Header.Get("supertokens-sdk-name")
	sdkVersion := request.Header.Get("supertokens-sdk-version")
	response.Write([]byte(strconv.FormatBool(sdkName == "website" && sdkVersion == "4.1.3")))
}

func checkAllowCredentials(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(strconv.FormatBool(request.Header.Get("allow-credentials") != "")))
}

func testError(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("test error message"))
}
