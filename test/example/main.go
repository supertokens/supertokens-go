package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens"
)

// noOfTimesGetSessionCalledDuringTest := 0
// noOfTimesRefreshCalledDuringTest := 0

func main() {
	// TODO: create API according to https://github.com/supertokens/supertokens-javalin/blob/master/Example/src/main/java/example/Main.java
	fmt.Println("Hello World!")
	supertokens.Config("localhost:8000;")

}
func loginHandler(response http.ResponseWriter, request *http.Request) {
	var userID map[string]interface{}
	json.NewDecoder(request.Body).Decode(&userID)

	_, err := supertokens.CreateNewSession(&response, userID["userId"].(string))

	if err != nil {
		supertokens.HandleErrorAndRespond(err, response)
	}

	json.NewEncoder(response).Encode([]byte("{\"message\": \"success\"}"))

}

/* JsonObject body = new JsonParser().parse(ctx.body()).getAsJsonObject();
   String userId = body.get("userId").getAsString();
   Session session = SuperTokens.newSession(ctx, userId).create();
   ctx.header("Access-Control-Allow-Origin", "http://127.0.0.1:8080");
   ctx.header("Access-Control-Allow-Credentials", "true");
   ctx.result(userId);*/
