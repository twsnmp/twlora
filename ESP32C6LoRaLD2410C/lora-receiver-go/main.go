package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go.bug.st/serial"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go list          - List available serial ports")
		fmt.Println("  go run main.go <device_path> - Start receiving from the specified port")
		return
	}

	arg := os.Args[1]

	if arg == "list" {
		listPorts()
		return
	}

	startReceiver(arg)
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

func startReceiver(portName string) {
	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatalf("Error opening serial port %s: %v", portName, err)
	}
	defer port.Close()

	log.Printf("Started listening on %s (9600 bps)\n", portName)

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log.Println(line)

		// Parse format: RM,ID,status,dist1,dist2 (e.g., RM,1,T,45,20)
		parts := strings.Split(line, ",")
		if len(parts) >= 3 && parts[2] == "T" {
			// Trigger BEEP sound
			fmt.Print("\a") // ASCII Bell
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from serial port: %v\n", err)
	}
}
