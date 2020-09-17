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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/supertokens/supertokens-go/supertokens/core"
)

// GetInstallationDir install dir for supertokens
func GetInstallationDir() string {
	// INSTALL_DIR=../supertokens-root go test ./... -count=1
	var installDir = os.Getenv("INSTALL_DIR")
	if installDir == "" {
		installDir = "../../com-root"
	}
	return installDir
}

func setKeyValueInConfig(key string, value string) {
	f, err := os.OpenFile(GetInstallationDir()+"/config.yaml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(key + ": " + value + "\n"); err != nil {
		panic(err)
	}
}

func executeCommand(waitFor bool, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = GetInstallationDir()
	cmd.Start()
	if waitFor {
		cmd.Wait()
	}
}

func setupST() {
	executeCommand(true, "cp", "temp/licenseKey", "./licenseKey")
	executeCommand(true, "cp", "temp/config.yaml", "./config.yaml")
	setKeyValueInConfig("refresh_api_path", "/refresh")
	setKeyValueInConfig("enable_anti_csrf", "true")
}

func cleanST() {
	executeCommand(true, "rm", "licenseKey")
	executeCommand(true, "rm", "config.yaml")
	executeCommand(true, "rm", "-rf", ".webserver-temp-*")
	executeCommand(true, "rm", "-rf", ".started")
}

func getListOfPids() []string {
	result := []string{}
	installationDir := GetInstallationDir()
	if _, err := os.Stat(installationDir + "/.started"); os.IsNotExist(err) {
		return result
	}
	items, _ := ioutil.ReadDir(installationDir + "/.started")
	for _, item := range items {
		dat, _ := ioutil.ReadFile(installationDir + "/.started/" + item.Name())
		strData := string(dat)
		result = append(result, strData)
	}
	return result
}

func stopST(pid string) {
	pidsBefore := getListOfPids()
	if len(pidsBefore) == 0 {
		return
	}
	executeCommand(true, "kill", pid)
	startTime := getCurrTimeInMS()
	for getCurrTimeInMS()-startTime < 10000 {
		pidsAfter := getListOfPids()
		if itemExists(pidsAfter, pid) {
			time.Sleep(100 * time.Millisecond)
		} else {
			return
		}
	}
	panic("Could not stop ST")
}

func killAllST() {
	pids := getListOfPids()
	for i := 0; i < len(pids); i++ {
		stopST(pids[i])
	}
	core.ResetDeviceDriverInfo()
	core.ResetError()
	core.ResetHandshakeInfo()
	core.ResetQuerier()
	core.ResetProcessState()
	core.ResetHTTPMocking()
}

func startST(host string, port string) string {
	pidsBefore := getListOfPids()
	executeCommand(false, "bash", "-c", "java -Djava.security.egd=file:/dev/urandom -classpath \"./core/*:./plugin-interface/*\" io.supertokens.Main ./ DEV host="+host+" port="+port)
	startTime := getCurrTimeInMS()
	for getCurrTimeInMS()-startTime < 20000 {
		pidsAfter := getListOfPids()
		if len(pidsAfter) <= len(pidsBefore) {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		nonIntersection := getNonIntersection(pidsAfter, pidsBefore)
		if len(nonIntersection) < 1 {
			panic("something went wrong while starting ST")
		} else {
			return nonIntersection[0]
		}
	}
	panic("could not start ST process")
}

func getCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func itemExists(arr []string, item string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return true
		}
	}
	return false
}

func getNonIntersection(a1 []string, a2 []string) []string {
	var result = []string{}
	for i := 0; i < len(a1); i++ {
		there := false
		for y := 0; y < len(a2); y++ {
			if a1[i] == a2[y] {
				there = true
			}
		}
		if !there {
			result = append(result, a1[i])
		}
	}
	return result
}
func containsHost(hostsAlive []string, host string) bool {
	if len(hostsAlive) == 0 {
		return false
	}
	for _, value := range hostsAlive {
		if value == host {
			return true
		}
	}
	return false
}

func extractInfoFromResponseHeader(res *http.Response) map[string]string {

	headerInfo := extractInfoFromCookies(res.Cookies())
	if res.Header.Get("Anti-Csrf") != "" {
		headerInfo["antiCsrf"] = res.Header.Get("Anti-Csrf")
	}
	if res.Header.Get("Id-Refresh-Token") != "" {
		headerInfo["idRefreshTokenFromHeader"] = res.Header.Get("Id-Refresh-Token")
	}

	return headerInfo
}

func extractInfoFromCookies(cookies []*http.Cookie) map[string]string {
	var response = map[string]string{}
	if len(cookies) == 0 {
		return response
	}

	for _, cookie := range cookies {
		if cookie.Name == "sAccessToken" {
			response["accessToken"] = cookie.Value
			response["accessTokenPath"] = cookie.Path
			response["accessTokenDomain"] = cookie.Domain
			response["accessTokenExpiry"] = cookie.RawExpires
			response["accessTokenSecure"] = strconv.FormatBool(cookie.Secure)
		} else if cookie.Name == "sRefreshToken" {
			response["refreshToken"] = cookie.Value
			response["refreshTokenPath"] = cookie.Path
			response["refreshTokenDomain"] = cookie.Domain
			response["refreshTokenExpiry"] = cookie.RawExpires
		} else {
			response["idRefreshTokenFromCookie"] = cookie.Value
			response["idRefreshTokenExpiry"] = cookie.RawExpires
			response["idRefreshTokenDomain"] = cookie.Domain
		}
	}
	return response
}

func beforeEach() {
	killAllST()
	setupST()
}
