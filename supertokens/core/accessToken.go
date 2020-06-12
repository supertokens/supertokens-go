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

package core

import (
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

type accessTokenInfoStruct struct {
	sessionHandle           string
	userID                  string
	refreshTokenHash1       string
	parentRefreshTokenHash1 *string
	userData                map[string]interface{}
	antiCsrfToken           *string
	expiryTime              uint64
	timeCreated             uint64
}

func getInfoFromAccessToken(token string, jwtSigningPublicKey string, doAntiCsrfCheck bool) (accessTokenInfoStruct, error) {
	payload, verifyError := verifyJWTAndGetPayload(token, jwtSigningPublicKey)
	if verifyError != nil {
		return accessTokenInfoStruct{}, errors.TryRefreshTokenError{
			Msg: verifyError.Error(),
		}
	}

	var sessionHandle *string = nil
	if payload["sessionHandle"] != nil {
		temp := payload["sessionHandle"].(string)
		sessionHandle = &temp
	}

	var userID *string = nil
	if payload["userId"] != nil {
		temp := payload["userId"].(string)
		userID = &temp
	}

	var refreshTokenHash1 *string = nil
	if payload["refreshTokenHash1"] != nil {
		temp := payload["refreshTokenHash1"].(string)
		refreshTokenHash1 = &temp
	}

	var parentRefreshTokenHash1 *string = nil
	if payload["parentRefreshTokenHash1"] != nil {
		temp := payload["parentRefreshTokenHash1"].(string)
		parentRefreshTokenHash1 = &temp
	}

	var userData *map[string]interface{} = nil
	if payload["userData"] != nil {
		temp := payload["userData"].(map[string]interface{})
		userData = &temp
	}

	var antiCsrfToken *string = nil
	if payload["antiCsrfToken"] != nil {
		temp := payload["antiCsrfToken"].(string)
		antiCsrfToken = &temp
	}

	var expiryTime *uint64 = nil
	if payload["expiryTime"] != nil {
		temp := uint64(payload["expiryTime"].(float64))
		expiryTime = &temp
	}

	var timeCreated *uint64 = nil
	if payload["timeCreated"] != nil {
		temp := uint64(payload["timeCreated"].(float64))
		timeCreated = &temp
	}

	if sessionHandle == nil ||
		userID == nil ||
		refreshTokenHash1 == nil ||
		userData == nil ||
		(antiCsrfToken == nil && doAntiCsrfCheck) ||
		expiryTime == nil ||
		timeCreated == nil {
		return accessTokenInfoStruct{}, errors.TryRefreshTokenError{
			Msg: "Access token does not contain all the information. Maybe the structure has changed?",
		}
	}

	if *expiryTime < getCurrTimeInMS() {
		return accessTokenInfoStruct{}, errors.TryRefreshTokenError{
			Msg: "Access token expired",
		}
	}

	return accessTokenInfoStruct{
		sessionHandle:           *sessionHandle,
		userID:                  *userID,
		refreshTokenHash1:       *refreshTokenHash1,
		parentRefreshTokenHash1: parentRefreshTokenHash1,
		userData:                *userData,
		antiCsrfToken:           antiCsrfToken,
		expiryTime:              *expiryTime,
		timeCreated:             *timeCreated,
	}, nil
}
