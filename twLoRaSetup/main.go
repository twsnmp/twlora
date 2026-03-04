package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

// main はプログラムのエントリポイントです。
// コマンドライン引数からシリアルポート、チャンネル、アドレスを取得し、
// LoRaモジュールの設定を表示または書き込みを行います。
func main() {
	var err error
	ch := 920    // デフォルトの周波数 (MHz)
	addr := 0    // デフォルトのアドレス (0-65535)
	set := false // 設定書き込みを行うかどうかのフラグ

	// 引数チェック: 少なくともポート名の指定が必要
	if len(os.Args) < 2 {
		log.Fatalln("usage: loraconf <Port> [Channel] [Address]")
	}

	// 引数のパース: チャンネルとアドレスが指定されているか確認
	if len(os.Args) >= 3 {
		ch, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid channel: %v\n", err)
		}
		set = true
	}
	if len(os.Args) >= 4 {
		addr, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("invalid address: %v\n", err)
		}
	}

	// シリアルポートの設定
	// Name: ポートパス, Baud: 9600, ReadTimeout: 応答待ちのタイムアウト
	config := &serial.Config{Name: os.Args[1], Baud: 9600, ReadTimeout: time.Second * 2}
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// バッファのクリア
	// 接続直後にシリアルバッファに残っている可能性のある不要なデータを捨てる
	dummy := make([]byte, 128)
	_, _ = s.Read(dummy)
	s.Flush()

	if set {
		// 指定された値で設定を書き込む
		if err = setConfig(s, addr, ch); err != nil {
			log.Fatalln(err)
		}
	} else {
		// 現在の設定を読み出して表示する
		showCurrentConfig(s)
	}
}

// showCurrentConfig はモジュールから現在の設定値を読み出し、解析結果を表示します。
func showCurrentConfig(s *serial.Port) {
	// 設定読み出しコマンド (C1 C1 C1) を送信
	fmt.Println("--- 設定読み出し開始 ---")
	readCmd := []byte{0xC1, 0xC1, 0xC1}
	if _, err := s.Write(readCmd); err != nil {
		fmt.Printf("コマンド送信エラー: %v\n", err)
		return
	}

	// モジュールが内部処理を終えてレスポンスを返すまで待機
	time.Sleep(100 * time.Millisecond)

	// レスポンスの読み取り (設定データは6バイト固定)
	buf := make([]byte, 6)
	n, err := s.Read(buf)
	if err != nil || n < 6 {
		fmt.Printf("読み取り不完全です (受信バイト数: %d). M0/M1の状態を確認してください。\n", n)
	} else {
		fmt.Printf("現在の設定 (HEX): %X\n", buf)
		printConfig(buf)
	}
}

// setConfig は指定されたアドレスと周波数をモジュールに書き込みます。
// 書き込み後、正しく反映されたか確認するために再度読み込みを行います。
func setConfig(s *serial.Port, addr, ch int) error {
	// 設定コマンドの構築
	// 構造: [HEAD, ADDH, ADDL, SPED, CHAN, OPTION]
	// CHAN: 862MHzからのオフセット (例: 920MHzの場合は 920-862=58 (0x3A))
	newConfig := []byte{
		0xC0,                // HEAD: 0xC0は電源切断後も保存、0xC1は一時保存
		byte(addr >> 8),     // ADDH: アドレス上位バイト
		byte(addr & 0x00ff), // ADDL: アドレス下位バイト
		0x1A,                // SPED: 通信速度 (UART: 9600bps, Air: 2.4kbps)
		byte(ch - 862),      // CHAN: 周波数設定 (862 + CHAN MHz)
		0x44,                // OPTION: 動作モード (固定送信モード、100mW)
	}

	fmt.Printf("新しい設定を書き込み中: %X\n", newConfig)
	if _, err := s.Write(newConfig); err != nil {
		return err
	}

	// モジュールのフラッシュメモリへの書き込みと反映を待機
	time.Sleep(100 * time.Millisecond)

	// 書き込み後に返ってくる現在の設定値を確認
	buf := make([]byte, 6)
	n, err := s.Read(buf)
	if err != nil || n != 6 {
		return fmt.Errorf("書き込み後の確認に失敗しました (受信バイト数: %d)", n)
	}

	fmt.Printf("書き込み完了。現在の設定: %X\n", buf)
	printConfig(buf)
	return nil
}

// printConfig は設定バイナリの内容を解析し、人間が読みやすい形式で表示します。
func printConfig(buf []byte) {
	fmt.Printf("アドレス: %d (0x%X)\n", (int(buf[1])<<8)|int(buf[2]), buf[1:3])
	fmt.Printf("SPEED: 0x%X\n", buf[3])
	fmt.Printf("周波数: %dMHz (CHAN: 0x%X)\n", uint(buf[4])+862, buf[4])
	fmt.Printf("OPTION: 0x%X\n", buf[5])
}
