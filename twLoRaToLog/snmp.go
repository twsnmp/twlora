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

	snmp := &gosnmp.GoSNMP{
		Target:    snmpTrapDst,
		Port:      162,
		Community: snmpCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   gosnmp.Default.Timeout,
		Retries:   gosnmp.Default.Retries,
	}

	if err := snmp.Connect(); err != nil {
		log.Printf("SNMP connect error: %v", err)
		return
	}
	defer snmp.Conn.Close()

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

	if _, err := snmp.SendTrap(trap); err != nil {
		log.Printf("SNMP send trap error: %v", err)
		return
	}
	lastTrapSent[id] = time.Now()
}
