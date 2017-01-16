package main

import (
	"sync"
)

var requestSessionIDs []uint32
var requestSessionIDsMutex = &sync.Mutex{}

func RequestAdd(sessionID uint32) {
	requestSessionIDsMutex.Lock()
	requestSessionIDs = append(requestSessionIDs, sessionID)
	requestSessionIDsMutex.Unlock()
}

func requestGetIndex(sessionID uint32) int {
    for i, v := range requestSessionIDs {
        if v == sessionID {
            return i
        }
    }
    return -1
}

func RequestIsAdded(sessionID uint32) bool {
	requestSessionIDsMutex.Lock()
	defer requestSessionIDsMutex.Unlock()
	return requestGetIndex(sessionID) >= 0
}

// http://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang
func removeFromSlice(s []uint32, i int) []uint32 {
	if len(s) == 0 {
		return s
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func RequestRemove(sessionID uint32) {
	requestSessionIDsMutex.Lock()
	requestSessionIDs = removeFromSlice(requestSessionIDs, requestGetIndex(sessionID))
	requestSessionIDsMutex.Unlock()
}
