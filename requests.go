package main

import (
	"sync"
	"net"
)

type requestSessionData struct {
	sessionID uint32
	fromAddr net.UDPAddr
}
var requestSessionDatas []requestSessionData
var requestSessionDatasMutex = &sync.Mutex{}

func RequestAdd(sessionID uint32, fromAddr *net.UDPAddr) {
	requestSessionDatasMutex.Lock()
	rsd := requestSessionData{sessionID, *fromAddr}
	requestSessionDatas = append(requestSessionDatas, rsd)
	requestSessionDatasMutex.Unlock()
}

func requestGetIndex(sessionID uint32, fromAddr *net.UDPAddr) int {
    for i, v := range requestSessionDatas {
        if v.sessionID == sessionID && v.fromAddr.String() == fromAddr.String() {
            return i
        }
    }
    return -1
}

func RequestIsAdded(sessionID uint32, fromAddr *net.UDPAddr) bool {
	requestSessionDatasMutex.Lock()
	defer requestSessionDatasMutex.Unlock()
	return requestGetIndex(sessionID, fromAddr) >= 0
}

// http://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang
func removeFromSlice(s []requestSessionData, i int) []requestSessionData {
	if len(s) == 0 {
		return s
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func RequestRemove(sessionID uint32, fromAddr *net.UDPAddr) {
	requestSessionDatasMutex.Lock()
	requestSessionDatas = removeFromSlice(requestSessionDatas, requestGetIndex(sessionID, fromAddr))
	requestSessionDatasMutex.Unlock()
}
