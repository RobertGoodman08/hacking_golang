package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {

	devices, err := pcap.FindAllDevs()

	if err != nil {
		log.Panicln(err)
	}

	for index, dev := range devices {
		fmt.Println(index, " ", dev.Name)

		for _, add := range dev.Addresses {
			fmt.Println("\tIP: ", add.IP)
			fmt.Println("\tNM: ", add.Netmask)
		}
	}
}
