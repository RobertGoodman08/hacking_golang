package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func executeCommand(command string, args_arr []string) (err error) {
	args := args_arr
	cmd_obj := exec.Command(command, args...)
	cmd_obj.Stdout = os.Stdout
	cmd_obj.Stderr = os.Stderr
	cmd_obj.Stdin = os.Stdin

	err = cmd_obj.Run()

	if err != nil {
		log.Fatal(err)
		return
	}

	return nil
}

func main() {
	iface := flag.String("iface", "eth0", "Interface for which you want to change the MAC")
	newMac := flag.String("newMac", "", "provide the new MAC address")

	flag.Parse()

	executeCommand("powershell", []string{"-Command", "Get-NetAdapter | Where-Object { $_.Name -eq '" + *iface + "' }"})
	executeCommand("powershell", []string{"-Command", "Set-NetAdapter -Name '" + *iface + "' -MacAddress '" + *newMac + "'"})
	executeCommand("powershell", []string{"-Command", "Enable-NetAdapter -Name '" + *iface + "'"})
}
