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

package supertokens

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

// Middleware for verifying and refreshing session. ExtraParams are: bool, func(error, http.ResponseWriter)
func Middleware(theirHandler http.HandlerFunc, extraParams ...interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" || r.Method == "TRACE" {
			theirHandler.ServeHTTP(w, r)
			return
		}
		var path = r.URL.Path
		handshakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
		if handshakeInfoError != nil {
			if len(extraParams) != 2 {
				HandleErrorAndRespond(handshakeInfoError, w)
			} else {
				extraParams[1].(func(err error, w http.ResponseWriter))(handshakeInfoError, w)
			}
			return
		}
		refreshTokenPath := handshakeInfo.RefreshTokenPath
		if configMap != nil && configMap.RefreshAPIPath != "" {
			refreshTokenPath = configMap.RefreshAPIPath
		}
		if (refreshTokenPath == path ||
			(refreshTokenPath+"/") == path ||
			refreshTokenPath == (path+"/")) &&
			r.Method == "POST" {
			session, sessionError := RefreshSession(w, r)
			if sessionError != nil {
				if len(extraParams) != 2 {
					HandleErrorAndRespond(sessionError, w)
				} else {
					extraParams[1].(func(err error, w http.ResponseWriter))(sessionError, w)
				}
				return
			}
			ctx := context.WithValue(r.Context(), sessionContext, session)
			theirHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			var actualDoAntiCsrfCheck = r.Method != "GET"
			if len(extraParams) != 0 && extraParams[0] != nil {
				actualDoAntiCsrfCheck = extraParams[0].(bool)
			}
			session, sessionError := GetSession(w, r, actualDoAntiCsrfCheck)
			if sessionError != nil {
				if len(extraParams) != 2 {
					HandleErrorAndRespond(sessionError, w)
				} else {
					extraParams[1].(func(err error, w http.ResponseWriter))(sessionError, w)
				}
				return
			}
			ctx := context.WithValue(r.Context(), sessionContext, session)
			theirHandler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// HandleErrorAndRespond if error handlers are provided, then uses those, else does default error handling depending on the type of error
func HandleErrorAndRespond(err error, w http.ResponseWriter) {
	errorHandlers := core.GetErrorHandlersInstance()
	if errors.IsUnauthorizedError(err) {
		errorHandlers.OnUnauthorizedErrorHandler(err, w)
	} else if errors.IsTryRefreshTokenError(err) {
		errorHandlers.OnTryRefreshTokenErrorHandler(err, w)
	} else if errors.IsTokenTheftDetectedError(err) {
		actualError := err.(errors.TokenTheftDetectedError)
		errorHandlers.OnTokenTheftDetectedErrorHandler(actualError.SessionHandle, actualError.UserID, w)
	} else {
		errorHandlers.OnGeneralErrorHandler(err, w)
	}
}
