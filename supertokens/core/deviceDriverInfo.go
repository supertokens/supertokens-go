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

import "sync"

type device struct {
	name    string
	version string
}

type deviceInfo struct {
	frontendSDK []device
}

var deviceInfoInstantiated *deviceInfo
var deviceInfoOnce *sync.Once = new(sync.Once)
var deviceInfoLock sync.Mutex

// GetDeviceInfoInstance get device info struct - singleton
func GetDeviceInfoInstance() *deviceInfo {
	deviceInfoOnce.Do(func() {
		deviceInfoInstantiated = &deviceInfo{
			frontendSDK: []device{},
		}
	})
	return deviceInfoInstantiated
}

// AddToFrontendSDKs add a device's info to array
func (info *deviceInfo) AddToFrontendSDKs(name string, version string) {
	deviceInfoLock.Lock()
	defer deviceInfoLock.Unlock()
	for i := 0; i < len(info.frontendSDK); i++ {
		curr := info.frontendSDK[i]
		if curr.name == name && curr.version == version {
			return
		}
	}
	info.frontendSDK = append(info.frontendSDK, device{
		name:    name,
		version: version,
	})
}

// GetFrontendSDKs get info about devices that have queried
func (info *deviceInfo) GetFrontendSDKs() []map[string]string {
	result := []map[string]string{}
	for i := 0; i < len(info.frontendSDK); i++ {
		result = append(result, map[string]string{
			"name":    info.frontendSDK[i].name,
			"version": info.frontendSDK[i].version,
		})
	}
	return result
}

// ResetDeviceDriverInfo to be used for testing only
func ResetDeviceDriverInfo() {
	deviceInfoInstantiated = nil
	deviceInfoOnce = new(sync.Once)
}
