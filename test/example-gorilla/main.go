/*
 * Copyright (c) 2020, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

var noOfTimesGetSessionCalledDuringTest int = 0
var noOfTimesRefreshCalledDuringTest int = 0

func main() {
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:9000",
	})
	r := mux.NewRouter()
	r.HandleFunc("/login", login)
	r.HandleFunc("/testUserConfig", testUserConfig)
	r.HandleFunc("/multipleInterceptors", multipleInterceptors)
	r.HandleFunc("/", supertokens.Middleware(defaultHandler))
	r.HandleFunc("/beforeeach", beforeeach)
	r.HandleFunc("/testing", testing)
	r.HandleFunc("/logout", supertokens.Middleware(logout))
	r.HandleFunc("/revokeAll", supertokens.Middleware(revokeAll))
	r.HandleFunc("/refresh", supertokens.Middleware(refresh))
	r.HandleFunc("/refreshCalledTime", refreshCalledTime)
	r.HandleFunc("/getSessionCalledTime", getSessionCalledTime)
	r.HandleFunc("/ping", ping)
	r.HandleFunc("/testHeader", testHeader)
	r.HandleFunc("/checkDeviceInfo", checkDeviceInfo)
	r.HandleFunc("/checkAllowCredentials", checkAllowCredentials)
	r.HandleFunc("/testError", testError)
	r.HandleFunc("/index.html", index)
	r.HandleFunc("/fail", fail)
	r.HandleFunc("/update-jwt", supertokens.Middleware(updateJwt))
	supertokens.OnTryRefreshToken(customOnTryRefreshTokenError)
	supertokens.OnUnauthorized(customOnUnauthorizedError)
	supertokens.OnGeneralError(customOnGeneralError)
	http.ListenAndServe("0.0.0.0:8080", handlers.CORS(
		handlers.AllowedHeaders(append([]string{"Content-Type"}, supertokens.GetCORSAllowedHeaders()...)),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8080"}),
		handlers.AllowCredentials(),
	)(r))
}

func fail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(""))
}

func index(w http.ResponseWriter, r *http.Request) {
	dat, _ := ioutil.ReadFile("./static/index.html")
	w.Header().Set("Content-Type", "text/html")
	w.Write(dat)
}

func login(response http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
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
	response.Write([]byte(userID))

}

func testUserConfig(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	response.Write([]byte(""))

}
func multipleInterceptors(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
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

func defaultHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	noOfTimesGetSessionCalledDuringTest++
	session := supertokens.GetSessionFromRequest(request)
	response.Write([]byte(session.GetUserID()))
}

func updateJwt(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		session := supertokens.GetSessionFromRequest(request)
		json.NewEncoder(response).Encode(session.GetJWTPayload())
	} else if request.Method == "POST" {
		var body map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			response.Write([]byte("error when parsing the body"))
			return
		}
		session := supertokens.GetSessionFromRequest(request)
		session.UpdateJWTPayload(body)
		json.NewEncoder(response).Encode(session.GetJWTPayload())
	} else {
		response.Write([]byte("incorrect Method, requires POST or GET"))
	}
}

func beforeeach(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
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
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}

	session := supertokens.GetSessionFromRequest(request)
	err := session.RevokeSession()
	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
		return
	}
	response.Write([]byte("success"))

}

func revokeAll(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	session := supertokens.GetSessionFromRequest(request)
	userID := session.GetUserID()
	supertokens.RevokeAllSessionsForUser(userID)
	response.Write([]byte("success"))
}

func refresh(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	noOfTimesRefreshCalledDuringTest++
	response.Write([]byte("refresh success"))
}

func refreshCalledTime(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Write([]byte(strconv.Itoa(noOfTimesRefreshCalledDuringTest)))
}

func getSessionCalledTime(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Write([]byte(strconv.Itoa(noOfTimesGetSessionCalledDuringTest)))
}

func ping(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.Write([]byte(""))
}

func testHeader(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	testheader := request.Header.Get("st-custom-header")
	success := testheader != ""
	json.NewEncoder(response).Encode(map[string]interface{}{
		"success": success,
	})
}

func checkDeviceInfo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	sdkName := request.Header.Get("supertokens-sdk-name")
	sdkVersion := request.Header.Get("supertokens-sdk-version")
	response.Write([]byte(strconv.FormatBool(sdkName == "website" && sdkVersion != "")))
}

func checkAllowCredentials(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Write([]byte("incorrect Method, requires POST"))
		return
	}
	response.Write([]byte(strconv.FormatBool(request.Header.Get("allow-credentials") != "")))
}

func testError(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.Write([]byte("incorrect Method, requires GET"))
		return
	}
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("test error message"))
}

func customOnTryRefreshTokenError(err error, response http.ResponseWriter) {
	response.WriteHeader(440)
	response.Write([]byte(""))

}

func customOnUnauthorizedError(err error, response http.ResponseWriter) {
	response.WriteHeader(440)
	response.Write([]byte(""))
}

func customOnGeneralError(err error, response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("Something went wrong"))
}
