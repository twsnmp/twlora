#include <HardwareSerial.h>

// ピン定義
#define LORA_RX D5  // D5
#define LORA_TX D4  // D4

// UART1を使用
HardwareSerial LoRaSerial(1);

void setup() {
  Serial.begin(115200);
  while (!Serial); 

  // XIAO ESP32-C6のピンを明示的に指定して初期化
  // 9600bps, 8N1, RX=7, TX=6
  LoRaSerial.begin(9600, SERIAL_8N1, LORA_RX, LORA_TX);
  
  Serial.println("\n--- Debug: LoRa Config Mode ---");
  delay(1000); // モジュールが安定するまで長めに待機

  // 送信前にバッファを空にする
  while(LoRaSerial.available()) LoRaSerial.read();

  // 1. 読み出しコマンド
  Serial.println("Reading current config...");
  uint8_t readCmd[] = {0xC1, 0xC1, 0xC1};
  LoRaSerial.write(readCmd, 3);
  LoRaSerial.flush(); // 送信完了まで待機

  delay(1000); // 応答を待つ時間を延長

  if (LoRaSerial.available()) {
    Serial.print("Current Data: ");
    while (LoRaSerial.available()) {
      Serial.printf("%02X ", LoRaSerial.read());
    }
    Serial.println();
  } else {
    Serial.println("No response for Read Command.");
  }

  // 2. 書き込みコマンド
  uint8_t set920MHz[] = {0xC0, 0x00, 0x00, 0x1A, 0x3A, 0x44};
  Serial.println("Writing 920MHz config...");
  LoRaSerial.write(set920MHz, 6);
  LoRaSerial.flush();

  delay(1000);

  if (LoRaSerial.available()) {
    Serial.print("Success Response: ");
    while (LoRaSerial.available()) {
      Serial.printf("%02X ", LoRaSerial.read());
    }
    Serial.println();
  } else {
    Serial.println("No response for Write Command.");
  }
}

void loop() {}