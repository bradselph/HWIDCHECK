package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nHWID Checker")
		fmt.Println("Select an option:")
		fmt.Println("1. SMBIOS (UUID)")
		fmt.Println("2. BIOS (Serial Number)")
		fmt.Println("3. Motherboard (Serial Number)")
		fmt.Println("4. Chassis (Serial Number)")
		fmt.Println("5. CPU (Serial Number)")
		fmt.Println("6. HDD/SSD (Serial Number)")
		fmt.Println("7. Volume Information")
		fmt.Println("8. RAM (Serial Number)")
		fmt.Println("9. Windows Product ID")
		fmt.Println("10. MAC Addresses")
		fmt.Println("11. Print All to File and Save")
		fmt.Println("12. Exit")

		fmt.Print("Enter your choice: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			runCommand("SMBIOS (UUID)", "wmic", "csproduct", "get", "uuid")
		case "2":
			runCommand("BIOS (Serial Number)", "wmic", "bios", "get", "serialnumber")
		case "3":
			runCommand("Motherboard (Serial Number)", "wmic", "baseboard", "get", "serialnumber")
		case "4":
			runCommand("Chassis (Serial Number)", "wmic", "systemenclosure", "get", "serialnumber")
		case "5":
			runCommand("CPU (Serial Number)", "wmic", "cpu", "get", "serialnumber")
		case "6":
			runCommand("HDD/SSD (Serial Number)", "wmic", "diskdrive", "get", "serialnumber")
		case "7":
			runCommand("Volume Information", "vol")
		case "8":
			runCommand("RAM (Serial Number)", "wmic", "memorychip", "get", "serialnumber")
		case "9":
			runCommand("Windows Product ID", "wmic", "os", "get", "serialnumber")
			runCommand("Windows Product ID (Alternative)", "systeminfo", "|", "findstr", "/B", "/C:\"OS Serial Number\"")
		case "10":
			runCommand("MAC Addresses", "getmac", "/v")
			runCommand("MAC Addresses", "powershell", "-Command", "Get-NetAdapter")
			runCommand("MAC Addresses", "wmic", "path", "Win32_NetworkAdapter", "where", `"MacAddress like '%%:%%:%%:%%:%%:%%'"`, "get", "Name, MacAddress")
			runCommand("MAC Addresses", "wmic", "nic", "get", "Name, MACAddress")
			runCommand("MAC Addresses", "ipconfig", "/all", "|", "findstr", `"Physical Address"`)
		case "11":
			saveAllToFile()
		case "12":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
func runCommand(description, command string, args ...string) {
	fmt.Printf("Executing command: %s %s\n", command, strings.Join(args, " "))
	_, err := exec.LookPath(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command '%s' not found: %s\n", command, err)
		return
	}
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing '%s': %s\n", description, err)
		return
	}
	fmt.Println(string(output))
}
func saveAllToFile() {
	file, err := os.Create("hwid_info.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %s\n", err)
		return
	}
	defer file.Close()
	commands := map[string][]string{
		"SMBIOS (UUID)":                    {"wmic", "csproduct", "get", "uuid"},
		"BIOS (Serial Number)":             {"wmic", "bios", "get", "serialnumber"},
		"Motherboard (Serial Number)":      {"wmic", "baseboard", "get", "serialnumber"},
		"Chassis (Serial Number)":          {"wmic", "systemenclosure", "get", "serialnumber"},
		"CPU (Serial Number)":              {"wmic", "cpu", "get", "serialnumber"},
		"HDD/SSD (Serial Number)":          {"wmic", "diskdrive", "get", "serialnumber"},
		"Volume Information":               {"vol"},
		"RAM (Serial Number)":              {"wmic", "memorychip", "get", "serialnumber"},
		"Windows Product ID":               {"wmic", "os", "get", "serialnumber"},
		"Windows Product ID (Alternative)": {"systeminfo", "|", "findstr", "/B", "/C:\"OS Serial Number\""},
		"MAC Addresses":                    {"getmac", "/v"},
		"MAC Addresses (Powershell)":       {"powershell", "-Command", "Get-NetAdapter"},
		"MAC Addresses (WMIC Path)":        {"wmic", "path", "Win32_NetworkAdapter", "where", `"MacAddress like '%%:%%:%%:%%:%%:%%'"`, "get", "Name, MacAddress"},
		"MAC Addresses (WMIC NIC)":         {"wmic", "nic", "get", "Name, MACAddress"},
		"MAC Addresses (IPConfig)":         {"ipconfig", "/all", "|", "findstr", `"Physical Address"`},
	}
	for description, commandArgs := range commands {
		fmt.Fprintf(file, "Executing command: %s\n", strings.Join(commandArgs, " "))
		_, err := exec.LookPath(commandArgs[0])
		if err != nil {
			fmt.Fprintf(file, "Command '%s' not found: %s\n", commandArgs[0], err)
			continue
		}
		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(file, "Error executing '%s': %s\n", description, err)
			continue
		}
		fmt.Fprintf(file, "%s:\n%s\n", description, output)
	}
	fmt.Println("All information saved to hwid_info.txt")
}
