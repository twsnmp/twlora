#include <HardwareSerial.h>
#include <Preferences.h>
#include <SoftwareSerial.h>
#include <WiFi.h>
#include <ld2410.h>

// --- ピン定義 (Seeed Studio XIAO ESP32-C6 等の D0-D5 命名規則に準拠) ---
#define LORA_RX D5  // LoRaモジュールの TX ピンと接続
#define LORA_TX D4  // LoRaモジュールの RX ピンと接続
#define RADAR_RX D0 // LD2410C センサーの TX ピンと接続 (ESP側のRX)
#define RADAR_TX D1 // LD2410C センサーの RX ピンと接続 (ESP側のTX)
#define PIN_OUT D2  // LD2410C のデジタル出力 (OUT: 感知時にHIGH) 監視用

// --- インスタンス・グローバル変数 ---
SoftwareSerial LoRaSerial;     // LoRa通信用ソフトウェアシリアル
HardwareSerial SerialRadar(1); // LD2410用ハードウェアシリアル (UART1を使用)
ld2410 radar;                  // レーダーライブラリのインスタンス
Preferences preferences;       // 不揮発性メモリ (NVM) への設定保存用
String sensorID = "1";         // このデバイスの識別ID。起動時に読み込まれる

void setup() {
  // --- シリアルモニタ初期化 (デバッグ・コマンド入力用) ---
  Serial.begin(115200);
  unsigned long start = millis();
  while (!Serial && (millis() - start < 1000)) {
    delay(10); // シリアルモニタが準備できるまで最大1秒間待機
  }

  // --- LoRaモジュールの初期化 ---
  // ボーレートはモジュールの仕様 (通常9600bps) に合わせて設定
  LoRaSerial.begin(9600, SWSERIAL_8N1, LORA_RX, LORA_TX, false);

  // --- 設定 (センサーID) の読み込み ---
  // "sensor-config" 名前空間を開き、保存されているIDを取得
  preferences.begin("sensor-config", false);
  sensorID = preferences.getString("id", "1"); // 保存がなければデフォルトの "1"

  // --- 乱数シードの初期化 ---
  // 送信間隔のゆらぎを作るために、未接続のアナログピンからノイズを読み取る
  randomSeed(analogRead(0));

  Serial.println("システム準備完了 (ESP32-C6 LoRa/Radar Bridge)");
  Serial.print("現在の Sensor ID: ");
  Serial.println(sensorID);

  // --- LD2410Cレーダーの初期化 ---
  // デフォルトの通信速度は 256000bps
  SerialRadar.begin(256000, SERIAL_8N1, RADAR_RX, RADAR_TX);

  if (radar.begin(SerialRadar, false)) {
    Serial.println("LD2410C レーダー: 接続成功");
  } else {
    Serial.println(
        "LD2410C レーダー: 接続失敗 (配線またはボーレートを確認してください)");
  }
}

void loop() {
  // --- PC (USB Serial) からの入力を処理 ---
  if (Serial.available()) {
    String input = Serial.readStringUntil('\n');
    input.trim();

    if (input.length() > 0) {
      // 受信した文字列を新しいセンサーIDとして保存し、以降このIDで報告する
      sensorID = input;
      preferences.putString("id", sensorID);
      Serial.print("Sensor ID を更新しました: ");
      Serial.println(sensorID);
    }

    // 入力されたデータはそのまま LoRa モジュールへも転送する (ブリッジ機能)
    LoRaSerial.println(input);
  }

  // --- LoRa モジュールからの受信データを PC へ転送 ---
  if (LoRaSerial.available()) {
    Serial.write(LoRaSerial.read());
  }

  // --- レーダーの状態更新 ---
  // ループごとに呼び出してセンサーからのフレームを解析する
  radar.read();

  static unsigned long lastTime = 0;
  static unsigned long currentInterval = 15000; // 送信間隔（初期値 15秒）

  // --- タイマー処理: 一定間隔でセンサーの状態を報告 ---
  if (millis() - lastTime > currentInterval) {
    lastTime = millis();

    // 送信衝突を避けるため、次回の間隔にランダムなゆらぎ (±1秒) を持たせる
    currentInterval = 14000 + random(2001); // 14.0秒 〜 16.0秒

    char statusBuf[128];
    // 報告フォーマット: RM,ID,状態,移動体距離,静止体距離
    // 状態: T (検知あり), F (検知なし), E (エラー/未接続)

    if (!radar.isConnected()) {
      // センサー通信エラー
      snprintf(statusBuf, sizeof(statusBuf), "RM,%s,E,-1,-1", sensorID.c_str());
    } else if (radar.presenceDetected()) {
      // 人体検知あり
      snprintf(statusBuf, sizeof(statusBuf), "RM,%s,T,%d,%d", sensorID.c_str(),
               radar.movingTargetDistance(), radar.stationaryTargetDistance());
    } else {
      // 検知なし
      snprintf(statusBuf, sizeof(statusBuf), "RM,%s,F,-1,-1", sensorID.c_str());
    }

    // 各出力先に送信
    Serial.println(statusBuf);     // PCデバッグ用
    LoRaSerial.println(statusBuf); // 遠隔地への通信用
  }
}