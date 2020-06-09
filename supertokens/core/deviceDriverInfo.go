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
func (info *deviceInfo) GetFrontendSDKs() []device {
	return info.frontendSDK
}

// ResetDeviceDriverInfo to be used for testing only
func ResetDeviceDriverInfo() {
	deviceInfoInstantiated = nil
	deviceInfoOnce = new(sync.Once)
}
