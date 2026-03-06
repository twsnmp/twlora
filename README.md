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

## Building and Flashing Firmware

After setting up the environment, you can build and flash the firmware using the following methods.

### 1. Arduino IDE (Easiest for Beginners)
1. Install [Arduino IDE](https://www.arduino.cc/en/software).
2. Follow the [XIAO ESP32C6 Setup Guide](https://wiki.seeedstudio.com/xiao_esp32c6_getting_started/) to add the ESP32 board support.
3. Install the required libraries via **Library Manager**:
   - `ld2410` by Trevor Shannon
   - `EspSoftwareSerial`
4. Open the `.ino` file from the project directory (e.g., `ESP32C6LoRaLD2410C/ESP32C6LoRaLD2410C.ino`).
5. Select **Tools > Board > esp32 > Seeed Studio XIAO ESP32C6**.
6. Connect your device, select the correct port in **Tools > Port**.
7. Click the **Upload** button.

### 2. Arduino CLI (Recommended for Advanced Users)
If you have `arduino-cli` installed (via `mise install`):

```bash
# For Radar Sensor
arduino-cli upload -p <PORT> --fqbn esp32:esp32:XIAO_ESP32C6 ESP32C6LoRaLD2410C

# For Setup Tool
arduino-cli upload -p <PORT> --fqbn esp32:esp32:XIAO_ESP32C6 ESP32LoRaSetup
```
*Note: Replace `<PORT>` with your actual serial port (e.g., `/dev/tty.usbmodem...` on macOS or `COMx` on Windows).*

### 3. ESP32 Flash Download Tool (Windows GUI)
If you are using the pre-compiled binaries in the `dist/` directory:
1. Download the [ESP32 Flash Download Tool](https://www.espressif.com/en/support/download/other-tools).
2. Select **ChipType: ESP32-C6**.
3. Load the binary files from `dist/<component>/`:
   - `...bootloader.bin` @ `0x0`
   - `...partitions.bin` @ `0x8000`
   - `...ino.bin` @ `0x10000`
4. Set **SPI SPEED** to **80MHz** and **SPI MODE** to **DIO**.
5. Select the COM port and click **START**.

### 4. esptool (Command Line)
You can also use [esptool.py](https://github.com/espressif/esptool) (installable via `pip install esptool`):

```bash
esptool.py --chip esp32c6 --port <PORT> --baud 921600 write_flash 0x0 dist/radar/ESP32C6LoRaLD2410C.ino.bootloader.bin 0x8000 dist/radar/ESP32C6LoRaLD2410C.ino.partitions.bin 0x10000 dist/radar/ESP32C6LoRaLD2410C.ino.bin
```

## License
[Apache2](./LICENSE)
