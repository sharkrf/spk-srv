package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"log"
	"net"
	"sync"
	"strings"
)

type bmServerData struct {
	Network string
	Name string
	Host string
}

type bmReflectorData struct {
	Active int
}

type bmSubscription struct {
	Talkgroup int
}

type bmClientData struct {
	Reflector bmReflectorData
	StaticSubscriptions []bmSubscription
	DynamicSubscriptions []bmSubscription
}

type bmServerIP string
var bmServerIPHosts = make(map[bmServerIP]bmServerData)
var bmServerIPHostsMutex = &sync.Mutex{}

func getJson(url string, target interface{}) error {
	var httpClient = &http.Client{ Timeout: 2000 * time.Millisecond }
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func BMGetClientData(clientId uint32, result *bmClientData, finished chan bool) {
	url := fmt.Sprintf("https://api.brandmeister.network/v1.0/repeater/?action=PROFILE&q=%d", clientId)
	err := getJson(url, result)
	if err != nil {
		log.Println("getjson error: ", err)
	} else {
		finished <- true
	}
}

func BMGenerateCodeStrFromClientData(cd *bmClientData, sd *bmServerData) string {
	var networkIDStr string
	var refStr string
	var stgStr string
	var dtgStr string

	if lastIndex := strings.LastIndex(sd.Name, "/"); lastIndex >= 0 {
		for i := lastIndex+1; i < len(sd.Name); i++ {
			networkIDStr += "0" + string(sd.Name[i])
		}
	}

	if len(cd.StaticSubscriptions) > 0 {
		if len(cd.StaticSubscriptions) == 1 {
			stgStr = "LKSTTG"
		} else {
			stgStr = "LKSTGS"
		}
		for i := 0; i < len(cd.StaticSubscriptions); i++ {
			if i > 0 {
				stgStr += "AN"
			}
			tg := fmt.Sprintf("%d", cd.StaticSubscriptions[i].Talkgroup)
			for j := 0; j < len(tg); j++ {
				stgStr += "0" + string(tg[j])
			}
		}
	}

	if len(cd.DynamicSubscriptions) > 0 {
		if len(cd.DynamicSubscriptions) == 1 {
			dtgStr = "LKDNTG"
		} else {
			dtgStr = "LKDNGS"
		}
		for i := 0; i < len(cd.DynamicSubscriptions); i++ {
			if i > 0 {
				dtgStr += "AN"
			}
			tg := fmt.Sprintf("%d", cd.DynamicSubscriptions[i].Talkgroup)
			for j := 0; j < len(tg); j++ {
				dtgStr += "0" + string(tg[j])
			}
		}
	}

	if cd.Reflector.Active != 4000 && cd.Reflector.Active != 0 {
		refStr = "LKRF"
		ref := fmt.Sprintf("%d", cd.Reflector.Active)
		for i := 0; i < len(ref); i++ {
			refStr += "0" + string(ref[i])
		}
	}

	return "BM" + networkIDStr + stgStr + dtgStr + refStr
}

func BMGetServerDataForServerIP(addr string) (bmServerData, bool) {
	bmServerIPHostsMutex.Lock()
	defer bmServerIPHostsMutex.Unlock()
	val, ok := bmServerIPHosts[bmServerIP(addr)]
	return val, ok
}

func BMUpdateServerList() {
	log.Println("updating bm server list")

	var bmServers []bmServerData
	err := getJson("http://x.sharkrf.com/db/homebrew/servers.json", &bmServers)
	if err != nil {
		log.Println("update bm server list getjson error: ", err)
		return
	}

	newList := make(map[bmServerIP]bmServerData)
	for _, bmServer := range bmServers {
		if bmServer.Network != "BrandMeister" {
			continue
		}

		addrs, err := net.LookupHost(bmServer.Host)
		if err == nil {
			// Storing all addresses so we can get server data for an IP later.
			for _, addr := range addrs {
				newList[bmServerIP(addr)] = bmServer
			}
		}
	}

	if len(newList) > 3 {
		bmServerIPHostsMutex.Lock()
		bmServerIPHosts = newList
		bmServerIPHostsMutex.Unlock()
		log.Println("updating bm server list finished")
	} else {
		log.Println("updating bm server list failed")
	}
}

func BMProcess() {
	for {
		BMUpdateServerList()
		time.Sleep(time.Hour)
	}
}
