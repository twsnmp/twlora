#include <HardwareSerial.h>

#define LORA_RX D5
#define LORA_TX D4

HardwareSerial LoRaSerial(1);

void setup() {
  Serial.begin(115200);
  // 運用モードも9600bps（設定に合わせてください）
  LoRaSerial.begin(9600, SERIAL_8N1, LORA_RX, LORA_TX);
  Serial.println("LoRa Normal Mode Ready...");
}

void loop() {
  // PCからLoRaへ送信
  if (Serial.available()) {
    LoRaSerial.write(Serial.read());
  }
  // LoRaから受信してPCへ表示
  if (LoRaSerial.available()) {
    Serial.write(LoRaSerial.read());
  }
}