package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var syslogCh = make(chan string, 2000)

func startSyslog(ctx context.Context) {
	dstList := strings.Split(syslogDst, ",")
	dst := []net.Conn{}
	for _, d := range dstList {
		if !strings.Contains(d, ":") {
			d += ":514"
		}
		s, err := net.Dial("udp", d)
		if err != nil {
			log.Printf("failed to dial syslog %s: %v", d, err)
			continue
		}
		syslogCh <- fmt.Sprintf("start send syslog to %s", d)
		dst = append(dst, s)
	}
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	defer func() {
		for _, d := range dst {
			d.Close()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("stop syslog")
			return
		case msg := <-syslogCh:
			if s := makeSyslogMsg(host, msg); s != "" {
				for _, d := range dst {
					d.Write([]byte(s))
				}
			}
		}
	}
}

func makeSyslogMsg(host string, msg string) string {
	// Parse format: RM,ID,detected,movingDistance,stationaryDistance (e.g., RM,1,T,45,20)
	a := strings.Split(msg, ",")
	if len(a) != 5 || a[0] != "RM" {
		return fmt.Sprintf("<%d>%s %s twLoRa: %s", 21*8+6, time.Now().Format("2006-01-02T15:04:05-07:00"), host, msg)
	}
	return fmt.Sprintf("<%d>%s %s twLoRa: type=%s,id=%s,detected=%s,movingDistance=%s,stationaryDistance=%s", 21*8+6, time.Now().Format("2006-01-02T15:04:05-07:00"), host, a[0], a[1], a[2], a[3], a[4])
}

func sendSyslog(msg string) {
	select {
	case syslogCh <- msg:
	default:
		if debug {
			log.Println("syslog channel full, skipping message")
		}
	}
}
