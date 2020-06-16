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
	"fmt"
	"os"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

func TestMain(m *testing.M) {
	code := m.Run()
	killAllST()
	cleanST()
	os.Exit(code)
}

func TestQuerierCalledWithoutInit(t *testing.T) {
	beforeEach()
	core.GetQuerierInstance()
}

func TestCoreNotAvailable(t *testing.T) {
	beforeEach()
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080;http://localhost:8081",
	})
	q := core.GetQuerierInstance()
	_, err := q.SendGetRequest("", "/", map[string]string{})
	if err == nil && err.Error() != "Error while querying SuperTokens core" {
		t.Error("failed")
	}
}

func TestThreeCoresAndRoundRobin(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	startST("localhost", "8081")
	startST("localhost", "8082")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080;http://localhost:8081;http://localhost:8082",
	})
	q := core.GetQuerierInstance()
	response, _ := q.SendGetRequest("", "/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	response, _ = q.SendDeleteRequest("", "/hello", map[string]interface{}{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	hostAlive := q.GetHostsAliveForTesting()

	if len(hostAlive) != 3 {
		t.Error("failed")
	}

	if !(containsHost(hostAlive, "http://localhost:8080") &&
		containsHost(hostAlive, "http://localhost:8081") && containsHost(hostAlive, "http://localhost:8082")) {
		fmt.Println(hostAlive)
		t.Error("failed")
	}
}

func TestThreeCoresOneDeadRoundRobin(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	startST("localhost", "8082")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080;http://localhost:8081;http://localhost:8082",
	})
	q := core.GetQuerierInstance()
	response, _ := q.SendGetRequest("", "/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	response, _ = q.SendDeleteRequest("", "/hello", map[string]interface{}{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	hostAlive := q.GetHostsAliveForTesting()
	if len(hostAlive) != 2 {
		t.Error("failed")
	}

	response, _ = q.SendGetRequest("", "/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}

	hostAlive = q.GetHostsAliveForTesting()
	if len(hostAlive) != 2 {
		t.Error("failed")
		return
	}
	if !(containsHost(hostAlive, "http://localhost:8080") &&
		containsHost(hostAlive, "http://localhost:8082")) {
		t.Error("failed")
	}
	if containsHost(hostAlive, "http://localhost:8081") {
		t.Error("failed")
	}
}
