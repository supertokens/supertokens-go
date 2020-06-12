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
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/supertokens/supertokens-go/supertokens/errors"
)

// Hosts to host location of SuperTokens instances.
type hosts struct {
	hostname string
	port     int
}

type querier struct {
	hosts          []hosts
	lastTriedIndex int
	apiVersion     *string
}

var querierInstantiated *querier
var querierLock sync.Mutex
var hostsAliveForTesting = []string{}

// ResetQuerier to be used for testing only
func ResetQuerier() {
	querierInstantiated = nil
	hostsAliveForTesting = []string{}
}

// GetQuerierInstance function used to get querier struct
func GetQuerierInstance() *querier {
	if querierInstantiated == nil {
		querierLock.Lock()
		defer querierLock.Unlock()
		if querierInstantiated == nil {
			querierInstantiated = &querier{
				hosts: []hosts{
					hosts{
						hostname: "localhost",
						port:     3567,
					},
				},
				lastTriedIndex: 0,
				apiVersion:     nil,
			}
		}
	}
	return querierInstantiated
}

// InitQuerier set hosts
func InitQuerier(hostsStr string) error {
	if querierInstantiated == nil {
		querierLock.Lock()
		defer querierLock.Unlock()
		if querierInstantiated == nil {

			// convert "hostname1:port1;hostname2:port2" to proper data type
			var hostsArr = make([]hosts, 0)
			var splitted = strings.Split(hostsStr, ";")
			for i := 0; i < len(splitted); i++ {
				var curr = splitted[i]
				if curr == "" {
					continue
				}
				var hostname = strings.Split(curr, ":")[0]
				var port, err = strconv.Atoi(strings.Split(curr, ":")[1])
				if err != nil {
					return errors.GeneralError{
						Msg:         "Invalid syntax for connection string",
						ActualError: nil,
					}
				}
				hostsArr = append(hostsArr, hosts{
					hostname: hostname,
					port:     port,
				})
			}

			querierInstantiated = &querier{
				hosts:          hostsArr,
				lastTriedIndex: 0,
				apiVersion:     nil,
			}
		}
	}
	return nil
}

func (querierInstance *querier) getAPIVersion() (string, error) {
	if querierInstance.apiVersion != nil {
		return *(querierInstance.apiVersion), nil
	}
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierInstance.apiVersion != nil {
		return *(querierInstance.apiVersion), nil
	}
	response, err := querierInstance.sendRequestHelper("/apiversion", func(url string) (*http.Response, error) {
		return http.Get(url)
	}, len(querierInstance.hosts))

	if err != nil {
		return "", err
	}

	cdiSupportedByServerInterface := response["versions"].([]interface{})
	cdiSupportedByServer := []string{}
	for i := 0; i < len(cdiSupportedByServerInterface); i++ {
		cdiSupportedByServer = append(cdiSupportedByServer, cdiSupportedByServerInterface[i].(string))
	}

	supportedVersion := getLargestVersionFromIntersection(cdiSupportedByServer, CdiVersion)

	if supportedVersion == nil {
		return "", errors.GeneralError{
			Msg: "The running SuperTokens core version is not compatible with this Golang SDK. Please visit https://supertokens.io/docs/community/compatibility to find the right version",
		}
	}

	querierInstance.apiVersion = supportedVersion

	return *(querierInstance.apiVersion), nil
}
func (querierInstance *querier) GetHostsAliveForTesting() []string {
	return hostsAliveForTesting
}

func (querierInstance *querier) SendPostRequest(requestID string, path string, data map[string]interface{}) (map[string]interface{}, error) {
	if path == "/session" || path == "/session/verify" || path == "/session/refresh" || path == "/handshake" {
		data["frontendSDK"] = GetDeviceInfoInstance().GetFrontendSDKs()
		data["drive"] = map[string]interface{}{
			"name":    "go",
			"version": VERSION,
		}
	}
	return querierInstance.sendRequestHelper(path, func(url string) (*http.Response, error) {
		jsonData, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, apiVersionError := querierInstance.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)

		client := getHTTPClient(requestID)
		return client.Do(req)
	}, len(querierInstance.hosts))
}

func getHTTPClient(requestID string) MockedHTTPClient {
	mock := GetMockedHTTPClient(requestID)
	if mock == nil {
		return &http.Client{}
	}
	return mock
}

func (querierInstance *querier) SendDeleteRequest(requestID string, path string, data map[string]interface{}) (map[string]interface{}, error) {
	return querierInstance.sendRequestHelper(path, func(url string) (*http.Response, error) {
		jsonData, _ := json.Marshal(data)
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, apiVersionError := querierInstance.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)

		client := getHTTPClient(requestID)
		return client.Do(req)
	}, len(querierInstance.hosts))
}

func (querierInstance *querier) SendGetRequest(requestID string, path string, params map[string]string) (map[string]interface{}, error) {
	return querierInstance.sendRequestHelper(path, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()

		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()

		apiVerion, apiVersionError := querierInstance.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}
		req.Header.Set("cdi-version", apiVerion)

		client := getHTTPClient(requestID)
		return client.Do(req)
	}, len(querierInstance.hosts))
}

func (querierInstance *querier) SendPutRequest(requestID string, path string, data map[string]interface{}) (map[string]interface{}, error) {
	return querierInstance.sendRequestHelper(path, func(url string) (*http.Response, error) {
		jsonData, _ := json.Marshal(data)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, apiVersionError := querierInstance.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)

		client := getHTTPClient(requestID)
		return client.Do(req)
	}, len(querierInstance.hosts))
}

type httpRequestFunction func(url string) (*http.Response, error)

func (querierInstance *querier) sendRequestHelper(path string, httpRequest httpRequestFunction,
	numberOfTries int) (map[string]interface{}, error) {
	if numberOfTries == 0 {
		return nil, errors.GeneralError{
			Msg:         "No SuperTokens core available to query",
			ActualError: nil,
		}
	}
	var currentHost = querierInstance.hosts[querierInstance.lastTriedIndex]
	hostPortString := currentHost.hostname + ":" + strconv.Itoa(currentHost.port)
	querierInstance.lastTriedIndex = (querierInstance.lastTriedIndex + 1) % len(querierInstance.hosts)
	var resp, err = httpRequest("http://" + hostPortString + path)

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return querierInstance.sendRequestHelper(path, httpRequest, numberOfTries-1)
		}
		if resp != nil {
			resp.Body.Close()
		}
		return nil, errors.GeneralError{
			Msg:         "Error while querying SuperTokens core",
			ActualError: err,
		}
	}

	if flag.Lookup("test.v") != nil && !containsHost(hostsAliveForTesting, hostPortString) {
		hostsAliveForTesting = append(hostsAliveForTesting, hostPortString)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.GeneralError{
			Msg:         resp.Status,
			ActualError: nil,
		}
	}

	var body, readErr = ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, errors.GeneralError{
			Msg:         "Error while querying SuperTokens core",
			ActualError: readErr,
		}
	}

	finalResult := make(map[string]interface{})
	jsonError := json.Unmarshal(body, &finalResult)
	if jsonError != nil {
		return map[string]interface{}{
			"result": string(body),
		}, nil
	}
	return finalResult, nil
}
