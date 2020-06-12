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

import (
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

func TestDriverInfoCheckWithoutFrontendSDK(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")
	info, err := core.GetHandshakeInfoInstance()
	if err != nil {
		t.Error(err)
	}

	if info.AccessTokenPath != "/" {
		t.Error("AccessToken path is not /")
	}
	if info.CookieDomain != "supertokens.io" {
		if info.CookieDomain != "localhost" {
			t.Error("incorrect cookie domain")
		}
	}
	if info.CookieSecure {
		t.Error("cookie secure set as true")
	}
	if info.RefreshTokenPath != "/refresh" {
		t.Error("incorrect refresh token path")
	}
	if !info.EnableAntiCsrf {
		t.Error("enable anticsrf set to false")
	}
	if info.AccessTokenBlacklistingEnabled {
		t.Error("accessTokenBlacklisting enabled")
	}
	info.UpdateJwtSigningPublicKeyInfo("hello", 100)

	info2, err2 := core.GetHandshakeInfoInstance()
	if err2 != nil {
		t.Error(err2)
	}

	if info2.JwtSigningPublicKey != "hello" {
		t.Error("JwtSigningPublicKey value does not match")
	}
	if info2.JwtSigningPublicKeyExpiryTime != 100 {
		t.Error("JwtSigningPublicKeyExpiryTime does not match")
	}
}
