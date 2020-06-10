package testing

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/supertokens/supertokens-go/supertokens/core"
)

func getInstallationDir() string {
	// INSTALL_DIR=../supertokens-root go test ./... -count=1
	var installDir = os.Getenv("INSTALL_DIR")
	if installDir == "" {
		installDir = "../../com-root"
	}
	return installDir
}

func setKeyValueInConfig(key string, value string) {
	f, err := os.OpenFile(getInstallationDir()+"/config.yaml", os.O_APPEND|os.O_WRONLY, 0600)
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
	cmd.Dir = getInstallationDir()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	if waitFor {
		err := cmd.Wait()
		if err != nil {
			panic(err)
		}
	}
}

func setupST() {
	executeCommand(true, "cp", "temp/licenseKey", "./licenseKey")
	executeCommand(true, "cp", "temp/config.yaml", "./config.yaml")
}

func cleanST() {
	executeCommand(true, "rm", "licenseKey")
	executeCommand(true, "rm", "config.yaml")
	executeCommand(true, "rm", "-rf", ".webserver-temp-*")
	executeCommand(true, "rm", "-rf", ".started")
}

func getListOfPids() []string {
	result := []string{}
	installationDir := getInstallationDir()
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
	// TODO: process state reset
}

func startST(host string, port string) string {
	pidsBefore := getListOfPids()
	executeCommand(false, "bash", "-c", "java -Djava.security.egd=file:/dev/urandom -classpath \"./core/*:./plugin-interface/*\" io.supertokens.Main ./ DEV host="+host+" port="+port)
	startTime := getCurrTimeInMS()
	for getCurrTimeInMS()-startTime < 10000 {
		pidsAfter := getListOfPids()
		if len(pidsAfter) <= len(pidsBefore) {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		nonIntersection := getNonIntersection(pidsAfter, pidsBefore)
		if len(nonIntersection) != 1 {
			panic("something went wrong while starting ST")
		} else {
			return nonIntersection[0]
		}
	}
	panic("could not start ST process")
}

func getCurrTimeInMS() int64 {
	return time.Now().UnixNano() / 1000000
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

func beforeEach() {
	killAllST()
	setupST()
}
