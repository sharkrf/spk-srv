//go:generate sh generate.sh

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func sendAMBEAnswer(udpConn *net.UDPConn, toAddr *net.UDPAddr, res *spkAMBEResponsePacket) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, res); err != nil {
		log.Printf("send answer error: %v", err)
		return
	}
	writtenBytes, err := udpConn.WriteToUDP(buf.Bytes(), toAddr)
	if writtenBytes != SPK_AMBE_RESPONSE_PACKET_SIZE || err != nil {
		log.Printf("warning: can't send udp packet to %s\n", toAddr.String())
	}

	if res.PacketType != SPK_PACKET_TYPE_RESPONSE_TERMINATOR {
		res.SeqNum++
		time.Sleep(time.Duration(res.FrameCount) * 20 * time.Millisecond)
	}
}

func sendIMBEAnswer(udpConn *net.UDPConn, toAddr *net.UDPAddr, res *spkIMBEResponsePacket) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, res); err != nil {
		log.Printf("send answer error: %v", err)
		return
	}
	writtenBytes, err := udpConn.WriteToUDP(buf.Bytes(), toAddr)
	if writtenBytes != SPK_IMBE_RESPONSE_PACKET_SIZE || err != nil {
		log.Printf("warning: can't send udp packet to %s\n", toAddr.String())
	}

	if res.PacketType != SPK_PACKET_TYPE_RESPONSE_TERMINATOR {
		res.SeqNum++
		time.Sleep(time.Duration(res.FrameCount) * 20 * time.Millisecond)
	}
}

func main() {
	var bindIp = ""
	var bindPort = 65200
	var silent bool
	var logToFile bool

	flag.IntVar(&bindPort, "p", bindPort, "bind to port")
	flag.StringVar(&bindIp, "i", bindIp, "bind to ip address")
	flag.BoolVar(&silent, "s", false, "disable logging")
	flag.BoolVar(&logToFile, "f", false, "log to file spk-srv.log")
	flag.Parse()

	if logToFile && !silent {
		logFile, err := os.OpenFile("spk-srv.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Println("warning: can't open spk-srv.log for writing: ", err)
		} else {
			defer logFile.Close()
			log.SetOutput(io.MultiWriter(os.Stdout, logFile))
		}
	}

	log.Printf("spk-srv start, binding to %s:%d\n", bindIp, bindPort)

	if silent {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(bindIp),
		Port: bindPort,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	go BMProcess()

	log.Println("starting listening loop")
	buffer := make([]byte, 128)
	for {
		readBytes, fromAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Did we read at least magic + version number of bytes? Does packet magic match?
		if readBytes >= 7 && strings.Compare(SPK_PACKET_MAGIC, string(buffer[:6])) == 0 {
			log.Printf("got %d byte packet from %s\n", readBytes, fromAddr.String())

			switch buffer[6] {
			default:
				log.Printf("ignoring packet with version 0x%.2x\n", buffer[6])
			case 0:
				v0processPacket(udpConn, fromAddr, buffer, readBytes)
			case 1:
				v1processPacket(udpConn, fromAddr, buffer, readBytes)
			}
		}
	}
}
