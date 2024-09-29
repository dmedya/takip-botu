package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const telegramToken = "TOKEN-GİR"
const chatID = "ID-GİR"

const targetURL = "https://dmedya.github.io/"

const responseTimeLimit = 100 * time.Millisecond

func main() {
	for {

		start := time.Now()
		resp, err := http.Get(targetURL)
		elapsed := time.Since(start)

		if err != nil {
			log.Println("Site erişilemedi:", err)
			sendTelegramMessage("❌ Web sitesi down: " + targetURL)
		} else {
			defer resp.Body.Close()

			if elapsed > responseTimeLimit {
				message := fmt.Sprintf("⚠️  Shrek site yavaş yanıt veriyor: %s (Süre: %d ms)", targetURL, elapsed.Milliseconds())
				sendTelegramMessage(message)
			} else {
				message := fmt.Sprintf("✅ Web sitesi erişilebilir: %s (Süre: %d ms)", targetURL, elapsed.Milliseconds())
				sendTelegramMessage(message)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

func sendTelegramMessage(message string) {
	url := "https://api.telegram.org/bot" + telegramToken + "/sendMessage"
	text := fmt.Sprintf(`{"chat_id":"%s","text":"%s"}`, chatID, message)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(text)))
	if err != nil {
		log.Println("Telegram isteği oluşturulamadı:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Telegram mesajı gönderilemedi:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Telegram yanıtı:", string(body))
}
