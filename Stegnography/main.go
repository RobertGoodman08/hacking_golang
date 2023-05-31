package main

import (
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
	"os"
)

func main() {

	carrieFile := "go.jpg"
	MaliciousFile := "mal.exe"
	encodedFile := "encoded_test.jpg"

	err := steg.EncodeByFileNames(carrieFile, MaliciousFile, encodedFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("file has been successfully encoded")
	}

}
