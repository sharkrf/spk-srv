package main

import (
	"fmt"
)

const SPK_PACKET_MAGIC						= "SRFSPK"

const SPK_PACKET_TYPE_RESPONSE_TERMINATOR	= 0
const SPK_PACKET_TYPE_RESPONSE				= 1
const SPK_PACKET_TYPE_REQUEST				= 2
type spkPacketType uint8

const SPK_ANNOUNCE_TYPE_DEFAULT				= 0
const SPK_ANNOUNCE_TYPE_CONNECTING			= 1
const SPK_ANNOUNCE_TYPE_CONNECTED			= 2
const SPK_ANNOUNCE_TYPE_STATUS				= 3
const SPK_ANNOUNCE_TYPE_STARTUP				= 4
type spkAnnounceType uint8

const SPK_MODEM_MODE_DMR					= 2
const SPK_MODEM_MODE_DSTAR					= 3
const SPK_MODEM_MODE_C4FM					= 4
type spkModemMode uint8

const SPK_CONNECTOR_ID_UNKNOWN				= 0
const SPK_CONNECTOR_ID_DMRPLUS				= 1
const SPK_CONNECTOR_ID_HOMEBREW				= 2
const SPK_CONNECTOR_ID_HOMEBREW_MMDVM		= 3
const SPK_CONNECTOR_ID_DCS					= 4
const SPK_CONNECTOR_ID_FCS					= 5
const SPK_CONNECTOR_ID_SRFIPCONN_CLIENT		= 6
const SPK_CONNECTOR_ID_SRFIPCONN_SERVER		= 7
const SPK_CONNECTOR_ID_REF					= 8
const SPK_CONNECTOR_ID_YSFREF				= 9
type spkConnectorId uint8

const SPK_ANNOUNCE_DATA_MAX_LENGTH			= 33
const SPK_REQUEST_PACKET_SIZE				= 23 + SPK_ANNOUNCE_DATA_MAX_LENGTH

type spkRequestPacket struct {
	Magic [6]byte
	Version uint8
	PacketType spkPacketType
	SessionID uint32
	ConnectorID spkConnectorId
	AnnounceType spkAnnounceType
	AnnounceTypeData [2]uint32
	ModemMode spkModemMode
	CodeStr [SPK_ANNOUNCE_DATA_MAX_LENGTH]byte
}

const SPK_RESPONSE_PACKET_SIZE				= 41

type spkResponsePacket struct {
	Magic [6]byte
	Version uint8
	PacketType spkPacketType
	SessionID uint32
	SeqNum uint8
	AMBEFrameCount uint8
	AMBEFrames [3][9]byte
}

func getModemModeNameStr(modemMode spkModemMode) string {
	switch (modemMode) {
		case SPK_MODEM_MODE_DMR: return "dmr"
		case SPK_MODEM_MODE_DSTAR: return "dstar"
		case SPK_MODEM_MODE_C4FM: return "c4fm"
		default: return "unknown"
	}
}

func decodeAnnounceTypeAndDataToStr(at spkAnnounceType, atd [2]uint32) (string, string) {
	var res string
	var resData string

	switch (at) {
		default:
			res = "default"
			resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
		case SPK_ANNOUNCE_TYPE_CONNECTING:
			res = "connecting"
			resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0] >> 24, (atd[0] >> 16) & 0xff, (atd[0] >> 8) & 0xff, atd[0] & 0xff, atd[1])
		case SPK_ANNOUNCE_TYPE_CONNECTED:
			res = "connected"
			resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0] >> 24, (atd[0] >> 16) & 0xff, (atd[0] >> 8) & 0xff, atd[0] & 0xff, atd[1])
		case SPK_ANNOUNCE_TYPE_STATUS:
			res = "status"
			resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0] >> 24, (atd[0] >> 16) & 0xff, (atd[0] >> 8) & 0xff, atd[0] & 0xff, atd[1])
		case SPK_ANNOUNCE_TYPE_STARTUP:
			res = "startup"
			resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	}
	return res, resData
}

func getConnectorIdNameStr(id spkConnectorId) string {
	switch (id) {
		case SPK_CONNECTOR_ID_DMRPLUS: return "dmp"
		case SPK_CONNECTOR_ID_HOMEBREW: return "hbr"
		case SPK_CONNECTOR_ID_HOMEBREW_MMDVM: return "mmd"
		case SPK_CONNECTOR_ID_DCS: return "dcs"
		case SPK_CONNECTOR_ID_FCS: return "fcs"
		case SPK_CONNECTOR_ID_SRFIPCONN_CLIENT: return "sfc"
		case SPK_CONNECTOR_ID_SRFIPCONN_SERVER: return "sfs"
		case SPK_CONNECTOR_ID_REF: return "ref"
		case SPK_CONNECTOR_ID_YSFREF: return "ysf"
		default: return "unk"
	}
}
