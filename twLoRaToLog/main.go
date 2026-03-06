package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.bug.st/serial"
)

var version = "v1.0.0"
var commit = ""
var syslogDst = ""
var mqttDst = ""
var mqttUser = ""
var mqttPassword = ""
var mqttClientID = "twlora-log"
var mqttTopic = "twlora/sensor"
var snmpTrapDst = ""
var snmpCommunity = "public"
var snmpInterval = 0
var portName = ""
var list = false
var debug = false

func init() {
	flag.BoolVar(&list, "list", false, "list available serial ports")
	flag.StringVar(&portName, "port", "", "serial port name")
	flag.StringVar(&syslogDst, "syslog", "", "syslog destnation list")
	flag.StringVar(&mqttDst, "mqtt", "", "mqtt broker destination")
	flag.StringVar(&mqttUser, "mqttuser", "", "mqtt username")
	flag.StringVar(&mqttPassword, "mqttpassword", "", "mqtt password")
	flag.StringVar(&mqttClientID, "mqttclientid", "twlora-log", "mqtt client id")
	flag.StringVar(&mqttTopic, "mqtttopic", "twlora/sensor", "mqtt topic")
	flag.StringVar(&snmpTrapDst, "snmp", "", "snmp trap destination")
	flag.StringVar(&snmpCommunity, "snmpcommunity", "public", "snmp community")
	flag.IntVar(&snmpInterval, "snmpinterval", 0, "snmp trap interval (minutes)")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv("TWLORATOLOG_" + strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()
}

func main() {
	if list {
		listPorts()
		return
	}
	if portName == "" {
		log.Fatalf("Error: No serial port specified")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if syslogDst != "" {
		go startSyslog(ctx)
	}
	if mqttDst != "" {
		go startMQTT(ctx)
	}
	if snmpTrapDst != "" {
		go startSNMP(ctx)
	}
	go startReceiver(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	time.Sleep(time.Second * 1)
}

func listPorts() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatalf("Error listing serial ports: %v", err)
	}

	if len(ports) == 0 {
		fmt.Println("No serial ports found.")
		return
	}

	fmt.Println("Available Serial Ports:")
	for _, port := range ports {
		fmt.Printf("  - %s\n", port)
	}
}

func startReceiver(ctx context.Context) {
	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatalf("Error opening serial port %s: %v", portName, err)
	}

	log.Printf("Started listening on %s (9600 bps)\n", portName)

	go func() {
		<-ctx.Done()
		log.Println("Closing serial port...")
		port.Close()
	}()

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if syslogDst != "" {
			sendSyslog(line)
		}
		if mqttDst != "" {
			publishMQTT(line)
		}
		if snmpTrapDst != "" {
			sendSNMPTrap(line)
		}
		if debug {
			log.Println(line)
			// Parse format: RM,ID,status,dist1,dist2 (e.g., RM,1,T,45,20)
			parts := strings.Split(line, ",")
			if len(parts) >= 3 && parts[2] == "T" {
				// Trigger BEEP sound
				fmt.Print("\a") // ASCII Bell
			}
		}
	}

	if err := scanner.Err(); err != nil {
		// Check if error is due to port being closed
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("Error reading from serial port: %v\n", err)
		}
	}
}
