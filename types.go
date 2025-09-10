package main

import (
	"fmt"
)

const SPK_PACKET_MAGIC = "SRFSPK"

const SPK_PACKET_TYPE_RESPONSE_TERMINATOR = 0
const SPK_PACKET_TYPE_AMBE_RESPONSE = 1
const SPK_PACKET_TYPE_REQUEST = 2
const SPK_PACKET_TYPE_IMBE_RESPONSE = 3

type spkPacketType uint8

const SPK_ANNOUNCE_TYPE_DEFAULT = 0
const SPK_ANNOUNCE_TYPE_CONNECTING = 1
const SPK_ANNOUNCE_TYPE_CONNECTED = 2
const SPK_ANNOUNCE_TYPE_CONNECTOR_STATUS = 3
const SPK_ANNOUNCE_TYPE_STARTUP = 4
const SPK_ANNOUNCE_TYPE_CONNECTED_BRANDMEISTER_SHORTENED = 5
const SPK_ANNOUNCE_TYPE_DISCONNECTED = 6
const SPK_ANNOUNCE_TYPE_WIFI_DISCONNECTED = 7
const SPK_ANNOUNCE_TYPE_WIFI_CONNECTING = 8

type spkAnnounceType uint8

const SPK_MODEM_MODE_DMR = 2
const SPK_MODEM_MODE_DSTAR = 3
const SPK_MODEM_MODE_C4FM = 4
const SPK_MODEM_MODE_C4FM_HALF_DEVIATION = 5
const SPK_MODEM_MODE_NXDN = 6
const SPK_MODEM_MODE_P25 = 7

type spkModemMode uint8

const SPK_VOICE_ID_MALE_EN = 0
const SPK_VOICE_ID_FEMALE_EN = 1

type spkVoiceID uint8

const SPK_CONNECTOR_ID_UNKNOWN = 0
const SPK_CONNECTOR_ID_DMRPLUS = 1
const SPK_CONNECTOR_ID_HOMEBREW = 2
const SPK_CONNECTOR_ID_HOMEBREW_MMDVM = 3
const SPK_CONNECTOR_ID_DCS = 4
const SPK_CONNECTOR_ID_FCS = 5
const SPK_CONNECTOR_ID_SRFIPCONN_CLIENT = 6
const SPK_CONNECTOR_ID_SRFIPCONN_SERVER = 7
const SPK_CONNECTOR_ID_REF = 8
const SPK_CONNECTOR_ID_YSFREF = 9
const SPK_CONNECTOR_ID_NULL = 10
const SPK_CONNECTOR_ID_NXDNREF = 11
const SPK_CONNECTOR_ID_BCAST = 12
const SPK_CONNECTOR_ID_IAX2 = 13
const SPK_CONNECTOR_ID_P25REF = 14
const SPK_CONNECTOR_ID_ECHOLINK = 15

type spkConnectorId uint8

const SPK_ANNOUNCE_DATA_MAX_LENGTH = 33
const SPK_REQUEST_PACKET_V0_SIZE = 23 + SPK_ANNOUNCE_DATA_MAX_LENGTH

type spkRequestPacketv0 struct {
	Magic            [6]byte
	Version          uint8
	PacketType       spkPacketType
	SessionID        uint32
	ConnectorID      spkConnectorId
	AnnounceType     spkAnnounceType
	AnnounceTypeData [2]uint32
	ModemMode        spkModemMode
	CodeStr          [SPK_ANNOUNCE_DATA_MAX_LENGTH]byte
}

const SPK_REQUEST_PACKET_V1_SIZE = 24 + SPK_ANNOUNCE_DATA_MAX_LENGTH

type spkRequestPacketv1 struct {
	Magic            [6]byte
	Version          uint8
	PacketType       spkPacketType
	SessionID        uint32
	ConnectorID      spkConnectorId
	AnnounceType     spkAnnounceType
	AnnounceTypeData [2]uint32
	ModemMode        spkModemMode
	VoiceID          spkVoiceID
	CodeStr          [SPK_ANNOUNCE_DATA_MAX_LENGTH]byte
}

const SPK_AMBE_RESPONSE_PACKET_SIZE = 41

type spkAMBEResponsePacket struct {
	Magic      [6]byte
	Version    uint8
	PacketType spkPacketType
	SessionID  uint32
	SeqNum     uint8
	FrameCount uint8
	Frames     [3][9]byte
}

const SPK_IMBE_RESPONSE_PACKET_SIZE = 68

type spkIMBEResponsePacket struct {
	Magic      [6]byte
	Version    uint8
	PacketType spkPacketType
	SessionID  uint32
	SeqNum     uint8
	FrameCount uint8
	Frames     [3][18]byte
}

type spkResponsePacket struct {
	AMBE spkAMBEResponsePacket
	IMBE spkIMBEResponsePacket
}

func getModemModeNameStr(modemMode spkModemMode) string {
	switch modemMode {
	case SPK_MODEM_MODE_DMR:
		return "dmr"
	case SPK_MODEM_MODE_DSTAR:
		return "dstar"
	case SPK_MODEM_MODE_C4FM:
		return "c4fm"
	case SPK_MODEM_MODE_C4FM_HALF_DEVIATION:
		return "c4fm-half"
	case SPK_MODEM_MODE_NXDN:
		return "nxdn"
	case SPK_MODEM_MODE_P25:
		return "p25"
	default:
		return "unknown"
	}
}

func decodeAnnounceTypeAndDataToStr(at spkAnnounceType, atd [2]uint32) (string, string) {
	var res string
	var resData string

	switch at {
	default:
		res = "unknown"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	case SPK_ANNOUNCE_TYPE_DEFAULT:
		res = "default"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	case SPK_ANNOUNCE_TYPE_CONNECTING:
		res = "connecting"
		resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0]>>24, (atd[0]>>16)&0xff, (atd[0]>>8)&0xff, atd[0]&0xff, atd[1])
	case SPK_ANNOUNCE_TYPE_CONNECTED:
		res = "connected"
		resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0]>>24, (atd[0]>>16)&0xff, (atd[0]>>8)&0xff, atd[0]&0xff, atd[1])
	case SPK_ANNOUNCE_TYPE_CONNECTED_BRANDMEISTER_SHORTENED:
		res = "connected (bm shortened)"
		resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0]>>24, (atd[0]>>16)&0xff, (atd[0]>>8)&0xff, atd[0]&0xff, atd[1])
	case SPK_ANNOUNCE_TYPE_CONNECTOR_STATUS:
		res = "status"
		resData = fmt.Sprintf("srv:%d.%d.%d.%d cid:%d", atd[0]>>24, (atd[0]>>16)&0xff, (atd[0]>>8)&0xff, atd[0]&0xff, atd[1])
	case SPK_ANNOUNCE_TYPE_STARTUP:
		res = "startup"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	case SPK_ANNOUNCE_TYPE_DISCONNECTED:
		res = "disconnected"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	case SPK_ANNOUNCE_TYPE_WIFI_DISCONNECTED:
		res = "wi-fi disconnected"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	case SPK_ANNOUNCE_TYPE_WIFI_CONNECTING:
		res = "wi-fi connecting"
		resData = fmt.Sprintf("%.8x%.8x", atd[0], atd[1])
	}
	return res, resData
}

func getConnectorIdNameStr(id spkConnectorId) string {
	switch id {
	case SPK_CONNECTOR_ID_DMRPLUS:
		return "dmp"
	case SPK_CONNECTOR_ID_HOMEBREW:
		return "hbr"
	case SPK_CONNECTOR_ID_HOMEBREW_MMDVM:
		return "mmd"
	case SPK_CONNECTOR_ID_DCS:
		return "dcs"
	case SPK_CONNECTOR_ID_FCS:
		return "fcs"
	case SPK_CONNECTOR_ID_SRFIPCONN_CLIENT:
		return "sfc"
	case SPK_CONNECTOR_ID_SRFIPCONN_SERVER:
		return "sfs"
	case SPK_CONNECTOR_ID_REF:
		return "ref"
	case SPK_CONNECTOR_ID_YSFREF:
		return "ysf"
	case SPK_CONNECTOR_ID_NULL:
		return "nul"
	case SPK_CONNECTOR_ID_NXDNREF:
		return "nxd";
	case SPK_CONNECTOR_ID_BCAST:
		return "bca";
	case SPK_CONNECTOR_ID_IAX2:
		return "iax";
	case SPK_CONNECTOR_ID_P25REF:
		return "p25";
	case SPK_CONNECTOR_ID_ECHOLINK:
		return "echolink";
	default:
		return "unk"
	}
}
