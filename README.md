# twlora - TWSNMP LoRa Toolset

[日本語のREADME](./README_ja.md)

`twlora` is a comprehensive suite of tools and programs for LoRa communication experiments and sensor network integration within the TWSNMP ecosystem.

## Overview
This repository provides a collection of software to bridge the gap between physical LoRa modules and network management systems. It covers everything from hardware configuration to data forwarding via modern protocols.

## Included Components

### 🛠 Configuration Tools
* **Serial/USB Configurator**: A utility to configure LoRa module parameters (frequency, spreading factor, etc.) via USB-Serial interface.
* **ESP32 Configurator**: A firmware tool to manage LoRa module settings directly from an ESP32.

### 📡 Firmware
* **Motion Sensor Transmitter (ESP32)**: A program that captures PIR sensor data and broadcasts it over LoRa.

### 🖥 Receiver & Integration (Go)
* **LoRa Data Gateway**: A high-performance Go-based receiver that processes incoming LoRa packets and forwards them to:
    * **Syslog**: For centralized logging and TWSNMP analysis.
    * **MQTT**: For IoT dashboard integration and real-time monitoring.

## Getting Started
(You can add specific build or installation instructions here later)

## License
[Apache2](./LICENSE)
