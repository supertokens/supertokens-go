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

package testing

// func TestMiddleware(t *testing.T) {
// 	beforeEach()
// 	startST("localhost", "8080")
// 	supertokens.Config("localhost:8080")
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/create", func(response http.ResponseWriter, requeset *http.Request) {
// 		supertokens.CreateNewSession(response, "testing-userID")
// 	})

// 	mux.HandleFunc("/user/id", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
// 		session := supertokens.GetSessionFromRequest(request)

// 		if session != nil {
// 			// json.NewEncoder(response).Encode(map[string]interface{}{
// 			// 	"message": session.GetUserID,
// 			// })
// 			response.Write([]byte(""))
// 			return
// 		}
// 		response.Write([]byte(""))
// 	}))

// 	mux.HandleFunc("/user/handle", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
// 		session := supertokens.GetSessionFromRequest(request)
// 		if session != nil {
// 			json.NewEncoder(response).Encode(map[string]interface{}{
// 				"message": session.GetHandle,
// 			})
// 			return
// 		}
// 		response.Write([]byte(""))
// 	}))

// 	mux.HandleFunc("/refresh", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {

// 		json.NewEncoder(response).Encode(map[string]interface{}{
// 			"message": true,
// 		})

// 	}))

// 	mux.HandleFunc("/logout", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
// 		session := supertokens.GetSessionFromRequest(request)
// 		if session != nil {
// 			err := session.RevokeSession()
// 			if err != nil {
// 				response.Write([]byte(""))
// 				return
// 			}
// 			json.NewEncoder(response).Encode(map[string]interface{}{
// 				"message": true,
// 			})
// 			return
// 		}
// 		response.Write([]byte(""))
// 	}))
// 	supertokens.OnTryRefreshToken(func(err error, response http.ResponseWriter) {
// 		response.WriteHeader(401)
// 		json.NewEncoder(response).Encode(map[string]interface{}{
// 			"message": " try refresh token",
// 		})
// 	})
// 	supertokens.OnTokenTheftDetected(func(val1 string, val2 string, response http.ResponseWriter) {
// 		response.WriteHeader(401)
// 		json.NewEncoder(response).Encode(map[string]interface{}{
// 			"message": " token theft detected",
// 		})
// 	})
// 	supertokens.OnGeneralError(func(err error, response http.ResponseWriter) {
// 		response.WriteHeader(401)
// 		json.NewEncoder(response).Encode(map[string]interface{}{
// 			"message": " general error",
// 		})
// 	})

// 	ts := httptest.NewServer(mux)
// 	defer ts.Close()

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
// 	res, _ := client.Do(req)

// 	response := extractInfoFromResponseHeader(res)

// 	req, _ = http.NewRequest("POST", ts.URL+"/user/id", nil)
// 	req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
// 	req.Header.Add("anti-csrf", response["antiCsrf"])
// 	res, _ = client.Do(req)

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	res.Body.Close()
// 	// var responseBody map[string]string
// 	// json.Unmarshal(body, &responseBody)
// 	fmt.Println(string(body))
// 	// if responseBody["message"] != "testing-userID" {
// 	// 	t.Error("incorrect response body")
// 	// }

// }

/* {
    Map<String, String> headers = new HashMap<>();
    headers.put("Cookie", "sAccessToken=" + response.get("accessToken") + ";sIdRefreshToken=" +
            response.get("idRefreshTokenFromCookie"));
    headers.put("anti-csrf", response.get("antiCsrf"));
    HttpURLConnection con = HttpRequest.sendGETRequest("http://localhost:8081/user/handle", new HashMap<>(), headers);

    assert (con.getResponseCode() == 200);
}*/
