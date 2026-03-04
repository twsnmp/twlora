# ESP32C6 LoRa LD2410C Project

## Project Overview
This is an Arduino-based project for an ESP32-C6 microcontroller that integrates an **LD2410C 24GHz mmWave Radar sensor** and a **LoRa communication module**. The system detects human presence (moving and stationary) and bridges serial data between the PC and the LoRa module.

### Core Technologies
- **Microcontroller:** ESP32-C6
- **Sensors:** LD2410C (Human Presence Radar)
- **Communication:** LoRa (Serial via SoftwareSerial), Hardware Serial
- **Framework:** Arduino (ESP32 Core)

## Hardware Configuration
| Peripheral | Pin (ESP32-C6) | Connection |
| :--- | :--- | :--- |
| **LoRa RX** | D5 | LoRa Module TX |
| **LoRa TX** | D4 | LoRa Module RX |
| **Radar RX** | D0 | LD2410C TX |
| **Radar TX** | D1 | LD2410C RX |
| **Radar OUT** | D2 | LD2410C OUT (Digital Presence) |

## Dependencies
- `ld2410` by Trevor Shannon (for the radar sensor)
- `SoftwareSerial` (for ESP32-C6 LoRa communication)
- `HardwareSerial` (standard ESP32 library)

## Building and Running
1. **Board Selection:** Select "ESP32C6 Dev Module" (or specific board like Seeed Studio XIAO ESP32C6) in the Arduino IDE or PlatformIO.
2. **Library Installation:** Ensure the `ld2410` library is installed via the Library Manager.
3. **Upload:** Use the standard upload process at 115200 baud.
4. **Monitoring:** Open the Serial Monitor at **115200 baud**.

## Development Conventions
- **Radar Data:** The radar is sampled every loop using `radar.read()`.
- **Reporting:** Presence status and distance are printed to the console every 5 seconds.
- **LoRa Bridge:** Any data sent via the Serial Monitor is forwarded to the LoRa module, and vice versa.
- **Pin Mapping:** Uses the `D0-D5` naming convention common in compact ESP32-C6 development boards.
