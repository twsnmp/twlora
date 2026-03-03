# twlora - TWSNMP LoRa Toolset

TWSNMPシリーズにおけるLoRa通信実験および、LoRaを用いたセンサーネットワーク構築のためのプログラム群です。

## 含まれるツール・プログラム
- **LoRa Config Tool (Serial/USB)**: PCからUSBシリアル経由でLoRaモジュールの周波数や出力を設定するツール。
- **LoRa Config Tool (ESP32)**: ESP32を使用してLoRaモジュールの設定を行うツール。
- **Sensor Transmitter (ESP32)**: 人感センサー等のデータを取得し、LoRa経由で送信するESP32用プログラム。
- **LoRa Receiver (Go)**: LoRa経由で受信したデータを、syslogやMQTTなどの上位プロトコルへ中継・転送するGo言語製プログラム。

## 特徴
- TWSNMPシリーズとの連携を想定した設計
- ESP32とLoRaモジュールを組み合わせたプロトタイピングの支援
- 多様な受信通知機能（syslog, MQTT等）

