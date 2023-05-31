package main

import (
	"log"
	"os"
)

func main() {
	dirPath := "."

	dir, err := os.Open(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			log.Println("Directory:", fileInfo.Name())
		} else {
			log.Println("File:", fileInfo.Name())
		}
	}
}
