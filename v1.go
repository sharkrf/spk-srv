package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"
)

func v1getFilePathForCodePair(modemMode spkModemMode, voiceID spkVoiceID, codePair string) string {
	dir := "voices/v1/"

	switch voiceID {
	default:
		dir += "srf-male-en"
	case SPK_VOICE_ID_FEMALE_EN:
		dir += "srf-female-en"
	}

	switch modemMode {
	case SPK_MODEM_MODE_DMR, SPK_MODEM_MODE_C4FM, SPK_MODEM_MODE_C4FM_HALF_DEVIATION, SPK_MODEM_MODE_NXDN:
		dir += "/dmr/"
	case SPK_MODEM_MODE_DSTAR:
		dir += "/dstar/"
	case SPK_MODEM_MODE_P25:
		dir += "/p25/"
	default:
		return ""
	}

	for filePath := range _bindata {
		if filepath.Dir(filePath) != filepath.Clean(dir) || filepath.Ext(filePath) != ".ambe" {
			continue
		}

		fileName := filepath.Base(filePath)
		// The requested code pair is stored in the first two characters of the filename.
		var fileCodePair = string(fileName[0]) + string(fileName[1])
		if fileCodePair == codePair {
			return filePath
		}
	}
	return ""
}

func v1StartSendAnswer(udpConn *net.UDPConn, toAddr net.UDPAddr, rp *spkRequestPacketv1) {
	defer RequestRemove(rp.SessionID, &toAddr)

	codeStr := strings.TrimRight(string(rp.CodeStr[:]), "\x00")

	// If the client is requesting a connect announce to a Homebrew server, we try to query a BM status from
	// the server's BM HTTP API to get linked talkgroups and reflector.
	bmGetClientDataFinished := make(chan bool)
	defer close(bmGetClientDataFinished)
	bmGetClientDataRunning := false
	var bmGetClientDataResult bmClientData
	var serverData bmServerData
	if rp.ConnectorID == SPK_CONNECTOR_ID_HOMEBREW &&
		(rp.AnnounceType == SPK_ANNOUNCE_TYPE_CONNECTED || rp.AnnounceType == SPK_ANNOUNCE_TYPE_CONNECTED_BRANDMEISTER_SHORTENED ||
			rp.AnnounceType == SPK_ANNOUNCE_TYPE_CONNECTOR_STATUS) {

		serverIP := fmt.Sprintf("%d.%d.%d.%d", rp.AnnounceTypeData[0]>>24, (rp.AnnounceTypeData[0]>>16)&0xff,
			(rp.AnnounceTypeData[0]>>8)&0xff, rp.AnnounceTypeData[0]&0xff)

		var ok bool
		if serverData, ok = BMGetServerDataForServerIP(serverIP); ok {
			clientId := rp.AnnounceTypeData[1]

			log.Printf("getting bm client data srv:%s cid:%d", serverIP, clientId)
			bmGetClientDataRunning = true
			go BMGetClientData(clientId, &bmGetClientDataResult, bmGetClientDataFinished)
		}
	}

	var res spkResponsePacket
	switch rp.ModemMode {
	default:
		copy(res.AMBE.Magic[:], SPK_PACKET_MAGIC)
		res.AMBE.Version = 1
		res.AMBE.PacketType = SPK_PACKET_TYPE_AMBE_RESPONSE
		res.AMBE.SessionID = rp.SessionID
	case SPK_MODEM_MODE_P25:
		copy(res.IMBE.Magic[:], SPK_PACKET_MAGIC)
		res.IMBE.Version = 1
		res.IMBE.PacketType = SPK_PACKET_TYPE_IMBE_RESPONSE
		res.IMBE.SessionID = rp.SessionID
	}

	// Stepping through each code char pair.
	for codeStrPos := 0; codeStrPos < len(codeStr); codeStrPos += 2 {
		if bmGetClientDataRunning {
			select {
			case finished := <-bmGetClientDataFinished:
				if finished && codeStrPos < 4 {
					toReplace := "HBSV"
					if strings.Contains(codeStr, "BMSV") {
						toReplace = "BMSV"
					}
					codeStr = strings.Replace(codeStr, toReplace, BMGenerateCodeStrFromClientData(&bmGetClientDataResult, &serverData,
						rp.AnnounceType == SPK_ANNOUNCE_TYPE_CONNECTED_BRANDMEISTER_SHORTENED), 1)
					log.Printf("code str modified for %s to %s", toAddr.String(), codeStr)
					bmGetClientDataRunning = false
				}
			default:
				break
			}
		}

		if codeStrPos+2 > len(codeStr) {
			log.Println("warning: last code pair is broken")
			break
		}

		var codePair = codeStr[codeStrPos : codeStrPos+2]

		filePath := v1getFilePathForCodePair(rp.ModemMode, rp.VoiceID, codePair)
		if filePath == "" {
			log.Printf("warning: file not found for modem mode %d code pair \"%s\", skipping\n", rp.ModemMode, codePair)
			continue
		}

		data, err := Asset(filePath)
		if err != nil {
			log.Printf("warning: can't find \"%s\", skipping\n", filePath)
			continue
		}

		log.Printf("playing %s to %s\n", filePath, toAddr.String())

		reader := bytes.NewReader(data)
		var fileFinished = false

		// Filling up frames from the file.
		for !fileFinished {
			switch rp.ModemMode {
			default:
				for ; res.AMBE.FrameCount < 3; res.AMBE.FrameCount++ {
					readBytes, err := reader.Read(res.AMBE.Frames[res.AMBE.FrameCount][:])
					if err != nil || readBytes != 9 {
						fileFinished = true
						break
					}
				}

				// Flushing if needed.
				if res.AMBE.FrameCount == 3 {
					sendAMBEAnswer(udpConn, &toAddr, &res.AMBE)
					res.AMBE.FrameCount = 0
				}
			case SPK_MODEM_MODE_P25:
				for ; res.IMBE.FrameCount < 3; res.IMBE.FrameCount++ {
					readBytes, err := reader.Read(res.IMBE.Frames[res.IMBE.FrameCount][:])
					if err != nil || readBytes != 18 {
						fileFinished = true
						break
					}
				}

				// Flushing if needed.
				if res.IMBE.FrameCount == 3 {
					sendIMBEAnswer(udpConn, &toAddr, &res.IMBE)
					res.IMBE.FrameCount = 0
				}
			}
		}
	}

	switch rp.ModemMode {
	default:
		res.AMBE.PacketType = SPK_PACKET_TYPE_RESPONSE_TERMINATOR
		sendAMBEAnswer(udpConn, &toAddr, &res.AMBE)
	case SPK_MODEM_MODE_P25:
		res.IMBE.PacketType = SPK_PACKET_TYPE_RESPONSE_TERMINATOR
		sendIMBEAnswer(udpConn, &toAddr, &res.IMBE)
	}

	if bmGetClientDataRunning {
		<-bmGetClientDataFinished
	}

	log.Printf("playing to %s finished\n", toAddr.String())
}

func v1processPacket(udpConn *net.UDPConn, fromAddr *net.UDPAddr, buffer []byte, readBytes int) {
	var packetType = buffer[7]

	switch packetType {
	default:
		log.Printf("ignoring packet with type 0x%.2x\n", packetType)
	case SPK_PACKET_TYPE_REQUEST:
		if readBytes != SPK_REQUEST_PACKET_V1_SIZE {
			log.Printf("ignoring packet with size %d\n", readBytes)
			return
		}

		// Reading the packet payload to our request struct.
		readBuf := bytes.NewReader(buffer)
		var rp spkRequestPacketv1
		err := binary.Read(readBuf, binary.BigEndian, &rp)
		if err != nil {
			log.Println("ignoring packet, binary parse error: ", err)
			return
		}

		switch rp.ModemMode {
		case SPK_MODEM_MODE_DMR:
			break
		case SPK_MODEM_MODE_DSTAR:
			break
		case SPK_MODEM_MODE_C4FM:
			break
		case SPK_MODEM_MODE_C4FM_HALF_DEVIATION:
			break
		case SPK_MODEM_MODE_NXDN:
			break
		case SPK_MODEM_MODE_P25:
			break
		default:
			log.Printf("ignoring packet, invalid modem mode %.2x\n", rp.ModemMode)
			return
		}

		rp.CodeStr[SPK_ANNOUNCE_DATA_MAX_LENGTH-1] = 0

		if RequestIsAdded(rp.SessionID, fromAddr) {
			//log.Printf("ignoring packet, request already under processing with sid:0x%.8x\n", rp.SessionID)
			return
		}
		RequestAdd(rp.SessionID, fromAddr)

		atStr, atdStr := decodeAnnounceTypeAndDataToStr(rp.AnnounceType, rp.AnnounceTypeData)
		log.Printf("sending \"%s\" to %s (sid:0x%.8x t:%s con:%s at:%s %s)\n",
			strings.TrimRight(string(rp.CodeStr[:]), "\x00"), fromAddr.String(), rp.SessionID, getModemModeNameStr(rp.ModemMode),
			getConnectorIdNameStr(rp.ConnectorID), atStr, atdStr)
		go v1StartSendAnswer(udpConn, *fromAddr, &rp)
	}
}
