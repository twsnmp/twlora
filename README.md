# twlora - TWSNMP LoRa Toolset

[日本語のREADME](./README_ja.md)

`twlora` is a comprehensive suite of tools and programs for LoRa communication experiments and sensor network integration within the TWSNMP ecosystem.

## Overview
This repository provides a collection of software to bridge the gap between physical LoRa modules and network management systems. It covers everything from hardware configuration to data forwarding via modern protocols.

## Included Components

### 🛠 Configuration Tools
* **Serial/USB Configurator (Go)**: A utility to configure LoRa module parameters via USB-Serial interface.
* **ESP32 Configurator (Arduino)**: A firmware tool to manage LoRa module settings directly from an ESP32-C6.

### 📡 Firmware (Arduino)
* **Radar Sensor Transmitter (ESP32-C6)**: A program that captures LD2410C 24GHz mmWave Radar data and broadcasts it over LoRa.
* **LoRa TxRx Test**: A simple bridge program to test LoRa communication.

### 🖥 Receiver & Integration (Go)
* **twLoRaToLog**: A high-performance Go-based receiver that processes incoming LoRa packets and forwards them to Syslog or MQTT for TWSNMP analysis.

## Build and Setup

This project uses [mise](https://mise.jdx.dev/) for tool and task management.

### Prerequisites
- Install `mise` on your system.
- Trust the configuration in this repository:
  ```bash
  mise trust
  ```

### Build Steps

1. **Install Tools and Setup Arduino Environment**:
   ```bash
   mise install
   mise run arduino:setup
   ```

2. **Build All Components**:
   ```bash
   mise run build
   ```
   Built binaries and firmware will be saved in the `dist/` directory.

3. **Individual Tasks**:
   - Build Go projects: `mise run go:build` (Outputs: `twLoRaToLog`, `twLoRaSetup`)
   - Compile Radar firmware: `mise run arduino:compile:radar`
   - Compile Setup firmware: `mise run arduino:compile:setup`

## License
[Apache2](./LICENSE)
