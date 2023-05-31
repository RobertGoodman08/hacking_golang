package main

import (
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
)

func main() {

	encodedFile := "encoded_test.jpg"
	ResultFile := "malicious.exe"

	err := steg.DecodeByFileNames(encodedFile, ResultFile)

	if err != nil {
		fmt.Errorf("Cant decode File")
	} else {
		fmt.Println("Successfully decoded file")
	}
}
