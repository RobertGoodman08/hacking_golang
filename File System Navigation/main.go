package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Указываем начальную директорию для навигации
	rootDir := "C:/"

	// Рекурсивно обходит директории, начиная с rootDir
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Ошибка при обходе пути %s: %v\n", path, err)
			return err
		}

		// Печатаем полный путь к файлу или директории
		fmt.Println(path)

		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при навигации по файловой системе: %v\n", err)
	}
}
