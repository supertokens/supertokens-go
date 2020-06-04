package core

import (
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
	lastTriedIndex int16
	apiVersion     *string
}

var querierInstantiated *querier
var querierLock sync.Mutex

// GetQuerierInstance function used to get querier struct
func GetQuerierInstance() *querier {
	if querierInstantiated == nil {
		querierLock.Lock()
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
		querierLock.Unlock()
	}
	return querierInstantiated
}

// InitQuerier set hosts
func InitQuerier(hostsStr string) error {
	if querierInstantiated == nil {
		querierLock.Lock()
		if querierInstantiated == nil {

			// convert "hostname1:port1;hostname2:port2" to proper data type
			var hostsArr = make([]hosts, 0)
			var splitted = strings.Split(hostsStr, ";")
			for i := 0; i < len(splitted); i++ {
				var curr = splitted[i]
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
		querierLock.Unlock()
	}
	return nil
}
