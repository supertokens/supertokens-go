package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens"
)

var noOfTimesGetSessionCalledDuringTest int = 0
var noOfTimesRefreshCalledDuringTest int = 0

func main() {
	// TODO: create API according to https://github.com/supertokens/supertokens-javalin/blob/master/Example/src/main/java/example/Main.java
	fmt.Println("Hello World!")
	supertokens.Config("localhost:8000;")

}
func loginHandler(response http.ResponseWriter, request *http.Request) {
	var body map[string]interface{}
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		// error when parsing body
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

func beforeeach(response http.ResponseWriter, request *http.Request) {
	noOfTimesRefreshCalledDuringTest = 0
	noOfTimesGetSessionCalledDuringTest = 0
	supertokens.ResetHandshakeInfo()
	response.Write([]byte(""))
}
