# twlora - TWSNMP LoRa Toolset

[English README](./README.md)

TWSNMPシリーズにおけるLoRa通信実験および、LoRaを用いたセンサーネットワーク構築のためのプログラム群です。

## 含まれるツール・プログラム

### 🛠 設定ツール
- **LoRa Config Tool (Go)**: PCからUSBシリアル経由でLoRaモジュールの周波数や出力を設定するツール。
- **LoRa Config Tool (Arduino)**: ESP32-C6を使用してLoRaモジュールの設定を行うツール。

### 📡 ファームウェア (Arduino)
- **Radar Sensor Transmitter (ESP32-C6)**: LD2410C 24GHz mmWave レーダーセンサーのデータを取得し、LoRa経由で送信するプログラム。
- **LoRa TxRx Test**: LoRa通信の疎通確認用ブリッジプログラム。

### 🖥 受信・統合 (Go)
- **twLoRaToLog**: LoRa経由で受信したデータを、syslogやMQTTなどのプロトコルへ中継・転送するプログラム。TWSNMPでの分析を想定しています。

## ビルドとセットアップ

このプロジェクトではツールの管理とビルドタスクの実行に [mise](https://mise.jdx.dev/) を使用しています。

### 準備
- `mise` をシステムにインストールしてください。
- このリポジトリの設定を信頼するように設定します：
  ```bash
  mise trust
  ```

### ビルド手順

1. **ツールのインストールとArduino環境のセットアップ**:
   ```bash
   mise install
   mise run arduino:setup
   ```

2. **全コンポーネントのビルド**:
   ```bash
   mise run build
   ```
   ビルドされたバイナリとファームウェアは `dist/` ディレクトリに保存されます。

3. **個別のタスク**:
   - Goプロジェクトのビルド: `mise run go:build` （`twLoRaToLog`, `twLoRaSetup` を生成）
   - レーダー用ファームウェアのコンパイル: `mise run arduino:compile:radar`
   - 設定用ファームウェアのコンパイル: `mise run arduino:compile:setup`

## ライセンス
[Apache2](./LICENSE)
