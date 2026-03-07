package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

var snmpCh = make(chan string, 2000)
var lastTrapSent = make(map[string]time.Time)

func startSNMP(ctx context.Context) {
	log.Println("start snmp trap sender")
	for {
		select {
		case <-ctx.Done():
			log.Println("stop snmp trap sender")
			return
		case msg := <-snmpCh:
			sendTrapMessage(msg)
		}
	}
}

func sendSNMPTrap(msg string) {
	select {
	case snmpCh <- msg:
	default:
		if debug {
			log.Println("snmp channel full, skipping message")
		}
	}
}

func sendTrapMessage(msg string) {
	// Parse format: RM,ID,status,dist1,dist2 (e.g., RM,1,T,45,20)
	a := strings.Split(msg, ",")
	if len(a) != 5 || a[0] != "RM" {
		return
	}
	if a[2] != "T" {
		return // Only send trap when detected
	}
	id := a[1]
	if snmpInterval > 0 {
		if last, ok := lastTrapSent[id]; ok {
			if time.Since(last) < time.Duration(snmpInterval)*time.Minute {
				return
			}
		}
	}

	movingDist, err := strconv.Atoi(a[3])
	if err != nil {
		return
	}
	stationaryDist, err := strconv.Atoi(a[4])
	if err != nil {
		return
	}
	port := 162
	ta := strings.SplitN(snmpTrapDst, ":", 2)
	if len(ta) > 1 {
		snmpTrapDst = ta[0]
		if v, err := strconv.ParseInt(ta[1], 10, 64); err == nil && v > 0 && v < 0xfffe {
			port = int(v)
		}
	}
	gosnmp.Default.Target = snmpTrapDst
	gosnmp.Default.Port = uint16(port)
	gosnmp.Default.Timeout = time.Duration(3) * time.Second
	switch snmpMode {
	case "v3auth":
		gosnmp.Default.Version = gosnmp.Version3
		gosnmp.Default.SecurityModel = gosnmp.UserSecurityModel
		gosnmp.Default.MsgFlags = gosnmp.AuthNoPriv
		gosnmp.Default.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 snmpUser,
			AuthoritativeEngineID:    snmpEngineID,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: snmpPassword,
		}
	case "v3authpriv":
		gosnmp.Default.Version = gosnmp.Version3
		gosnmp.Default.SecurityModel = gosnmp.UserSecurityModel
		gosnmp.Default.MsgFlags = gosnmp.AuthPriv
		gosnmp.Default.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 snmpUser,
			AuthoritativeEngineID:    snmpEngineID,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: snmpPassword,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        snmpPassword,
		}
	default:
		gosnmp.Default.Version = gosnmp.Version2c
		gosnmp.Default.Community = snmpCommunity
	}
	err = gosnmp.Default.Connect()
	if err != nil {
		log.Printf("SNMP connect error: %v", err)
		return
	}
	defer gosnmp.Default.Conn.Close()

	// OIDs from twLoRaToLogTrap.txt
	// .1.3.6.1.4.1.17861.1.11
	baseOID := ".1.3.6.1.4.1.17861.1.11"
	trapOID := baseOID + ".0.1" // twLoRaToLogDetectedTrap
	objOID := baseOID + ".1"

	trap := gosnmp.SnmpTrap{
		Variables: []gosnmp.SnmpPDU{
			{
				Name:  ".1.3.6.1.6.3.1.1.4.1.0",
				Type:  gosnmp.ObjectIdentifier,
				Value: trapOID,
			},
			{
				Name:  objOID + ".0", // twLoRaToLogSensorID
				Type:  gosnmp.OctetString,
				Value: id,
			},
			{
				Name:  objOID + ".1", // twLoRaToLogMovingDistance
				Type:  gosnmp.Integer,
				Value: movingDist,
			},
			{
				Name:  objOID + ".2", // twLoRaToLogStationaryDistance
				Type:  gosnmp.Integer,
				Value: stationaryDist,
			},
		},
	}

	if _, err := gosnmp.Default.SendTrap(trap); err != nil {
		log.Printf("SNMP send trap error: %v", err)
		return
	}
	lastTrapSent[id] = time.Now()
}
