package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	primary   []string
	fallbacks [][]string
}

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
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			continue
		}
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			runCommandWithFallbacks("SMBIOS (UUID)", Command{
				primary: []string{"wmic", "csproduct", "get", "uuid"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_ComputerSystemProduct).UUID"},
				},
			})
		case "2":
			runCommandWithFallbacks("BIOS (Serial Number)", Command{
				primary: []string{"wmic", "bios", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_BIOS | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_BIOS).SerialNumber"},
				},
			})
		case "3":
			runCommandWithFallbacks("Motherboard (Serial Number)", Command{
				primary: []string{"wmic", "baseboard", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_BaseBoard | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_BaseBoard).SerialNumber"},
				},
			})
		case "4":
			runCommandWithFallbacks("Chassis (Serial Number)", Command{
				primary: []string{"wmic", "systemenclosure", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_SystemEnclosure | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_SystemEnclosure).SerialNumber"},
				},
			})
		case "5":
			runCommandWithFallbacks("CPU (Serial Number)", Command{
				primary: []string{"wmic", "cpu", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_Processor).ProcessorId"},
				},
			})
		case "6":
			runCommandWithFallbacks("HDD/SSD (Serial Number)", Command{
				primary: []string{"wmic", "diskdrive", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_DiskDrive | Select-Object Model, SerialNumber"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_DiskDrive | Select-Object Model, SerialNumber"},
					{"powershell", "-Command", "Get-PhysicalDisk | Select-Object FriendlyName, SerialNumber"},
				},
			})
		case "7":
			runCommandWithFallbacks("Volume Information", Command{
				primary: []string{"vol"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-Volume"},
					{"powershell", "-Command", "Get-WmiObject Win32_LogicalDisk | Select-Object DeviceID, VolumeName, VolumeSerialNumber"},
				},
			})
		case "8":
			runCommandWithFallbacks("RAM (Serial Number)", Command{
				primary: []string{"wmic", "memorychip", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
				},
			})
		case "9":
			runCommandWithFallbacks("Windows Product ID", Command{
				primary: []string{"wmic", "os", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "(Get-WmiObject -Class Win32_OperatingSystem).SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_OperatingSystem).SerialNumber"},
					{"powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' -Name ProductId | Select-Object -ExpandProperty ProductId"},
				},
			})
			runCommandWithFallbacks("Windows Product ID (Alternative)", Command{
				primary: []string{"systeminfo", "|", "findstr", "/B", "/C:\"OS Serial Number\""},
				fallbacks: [][]string{
					{"powershell", "-Command", "systeminfo | Select-String 'OS Serial Number'"},
				},
			})
		case "10":
			runCommandWithFallbacks("MAC Addresses", Command{
				primary: []string{"getmac", "/v"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-NetAdapter | Select-Object Name, Status, MacAddress"},
				},
			})
			runCommandWithFallbacks("MAC Addresses", Command{
				primary: []string{"wmic", "path", "Win32_NetworkAdapter", "where", `"MacAddress like '%%:%%:%%:%%:%%:%%'"`, "get", "Name, MacAddress"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
				},
			})
			runCommandWithFallbacks("MAC Addresses", Command{
				primary: []string{"wmic", "nic", "get", "Name, MACAddress"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object { $_.IPEnabled -eq $true } | Select-Object Description, MACAddress"},
				},
			})
			runCommandWithFallbacks("MAC Addresses", Command{
				primary: []string{"ipconfig", "/all", "|", "findstr", `"Physical Address"`},
				fallbacks: [][]string{
					{"powershell", "-Command", "ipconfig /all | Select-String 'Physical Address'"},
				},
			})
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

func runCommandWithFallbacks(description string, command Command) bool {
	success := executeCommand(description, command.primary)
	if success {
		return true
	}

	for i, fallback := range command.fallbacks {
		fmt.Printf("Primary command failed. Trying fallback %d...\n", i+1)
		success := executeCommand(description+" (Fallback)", fallback)
		if success {
			return true
		}
	}

	fmt.Fprintf(os.Stderr, "All commands for '%s' failed\n", description)
	return false
}
func executeCommand(description string, args []string) bool {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command provided for '%s'\n", description)
		return false
	}

	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))
	
	if containsPipe(args) {
		return executePipedCommand(description, args)
	}

	_, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command '%s' not found: %s\n", args[0], err)
		return false
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing '%s': %s\n", description, err)
		return false
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		fmt.Fprintf(os.Stderr, "Command '%s' returned empty output\n", description)
		return false
	}

	fmt.Println(string(output))
	return true
}

func containsPipe(args []string) bool {
	for _, arg := range args {
		if arg == "|" {
			return true
		}
	}
	return false
}

func executePipedCommand(description string, args []string) bool {
	fullCommand := strings.Join(args, " ")
	cmd := exec.Command("cmd.exe", "/C", fullCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing piped command '%s': %s\n", description, err)
		return false
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		fmt.Fprintf(os.Stderr, "Piped command '%s' returned empty output\n", description)
		return false
	}

	fmt.Println(string(output))
	return true
}

func saveAllToFile() {
	file, err := os.Create("hwid_info.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %s\n", err)
		return
	}
	defer file.Close()

	type FileCommand struct {
		description string
		command     Command
	}

	commands := []FileCommand{
		{"SMBIOS (UUID)", Command{
			primary: []string{"wmic", "csproduct", "get", "uuid"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_ComputerSystemProduct).UUID"},
			},
		}},
		{"BIOS (Serial Number)", Command{
			primary: []string{"wmic", "bios", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_BIOS | Select-Object -ExpandProperty SerialNumber"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_BIOS).SerialNumber"},
			},
		}},
		{"Motherboard (Serial Number)", Command{
			primary: []string{"wmic", "baseboard", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_BaseBoard | Select-Object -ExpandProperty SerialNumber"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_BaseBoard).SerialNumber"},
			},
		}},
		{"Chassis (Serial Number)", Command{
			primary: []string{"wmic", "systemenclosure", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_SystemEnclosure | Select-Object -ExpandProperty SerialNumber"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_SystemEnclosure).SerialNumber"},
			},
		}},
		{"CPU (Serial Number)", Command{
			primary: []string{"wmic", "cpu", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_Processor).ProcessorId"},
			},
		}},
		{"HDD/SSD (Serial Number)", Command{
			primary: []string{"wmic", "diskdrive", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_DiskDrive | Select-Object Model, SerialNumber"},
				{"powershell", "-Command", "Get-CimInstance -Class Win32_DiskDrive | Select-Object Model, SerialNumber"},
				{"powershell", "-Command", "Get-PhysicalDisk | Select-Object FriendlyName, SerialNumber"},
			},
		}},
		{"Volume Information", Command{
			primary: []string{"vol"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-Volume"},
				{"powershell", "-Command", "Get-WmiObject Win32_LogicalDisk | Select-Object DeviceID, VolumeName, VolumeSerialNumber"},
			},
		}},
		{"RAM (Serial Number)", Command{
			primary: []string{"wmic", "memorychip", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
				{"powershell", "-Command", "Get-CimInstance -Class Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
			},
		}},
		{"Windows Product ID", Command{
			primary: []string{"wmic", "os", "get", "serialnumber"},
			fallbacks: [][]string{
				{"powershell", "-Command", "(Get-WmiObject -Class Win32_OperatingSystem).SerialNumber"},
				{"powershell", "-Command", "(Get-CimInstance -Class Win32_OperatingSystem).SerialNumber"},
				{"powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' -Name ProductId | Select-Object -ExpandProperty ProductId"},
			},
		}},
		{"Windows Product ID (Alternative)", Command{
			primary: []string{"systeminfo", "|", "findstr", "/B", "/C:\"OS Serial Number\""},
			fallbacks: [][]string{
				{"powershell", "-Command", "systeminfo | Select-String 'OS Serial Number'"},
			},
		}},
		{"MAC Addresses (GetMac)", Command{
			primary: []string{"getmac", "/v"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-NetAdapter | Select-Object Name, Status, MacAddress"},
			},
		}},
		{"MAC Addresses (WMIC Path)", Command{
			primary: []string{"wmic", "path", "Win32_NetworkAdapter", "where", `"MacAddress like '%%:%%:%%:%%:%%:%%'"`, "get", "Name, MacAddress"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
				{"powershell", "-Command", "Get-CimInstance -Class Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
			},
		}},
		{"MAC Addresses (WMIC NIC)", Command{
			primary: []string{"wmic", "nic", "get", "Name, MACAddress"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object { $_.IPEnabled -eq $true } | Select-Object Description, MACAddress"},
			},
		}},
		{"MAC Addresses (IPConfig)", Command{
			primary: []string{"ipconfig", "/all", "|", "findstr", `"Physical Address"`},
			fallbacks: [][]string{
				{"powershell", "-Command", "ipconfig /all | Select-String 'Physical Address'"},
			},
		}},
	}

	for _, cmd := range commands {
		fmt.Fprintf(file, "\n%s:\n", cmd.description)
		fmt.Fprintf(file, "Executing command: %s\n", strings.Join(cmd.command.primary, " "))
		
		success := false
		
		if containsPipe(cmd.command.primary) {
			fullCommand := strings.Join(cmd.command.primary, " ")
			execCmd := exec.Command("cmd.exe", "/C", fullCommand)
			output, err := execCmd.CombinedOutput()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				fmt.Fprintf(file, "%s\n", output)
				success = true
			}
		} else {
			_, err := exec.LookPath(cmd.command.primary[0])
			if err == nil {
				execCmd := exec.Command(cmd.command.primary[0], cmd.command.primary[1:]...)
				output, err := execCmd.CombinedOutput()
				if err == nil && len(strings.TrimSpace(string(output))) > 0 {
					fmt.Fprintf(file, "%s\n", output)
					success = true
				}
			}
		}
		
		if !success {
			for i, fallback := range cmd.command.fallbacks {
				fmt.Fprintf(file, "Primary command failed. Trying fallback %d:\n", i+1)
				fmt.Fprintf(file, "Executing command: %s\n", strings.Join(fallback, " "))
				
				if containsPipe(fallback) {
					fullCommand := strings.Join(fallback, " ")
					execCmd := exec.Command("cmd.exe", "/C", fullCommand)
					output, err := execCmd.CombinedOutput()
					if err == nil && len(strings.TrimSpace(string(output))) > 0 {
						fmt.Fprintf(file, "%s\n", output)
						success = true
						break
					}
				} else {
					_, err := exec.LookPath(fallback[0])
					if err == nil {
						execCmd := exec.Command(fallback[0], fallback[1:]...)
						output, err := execCmd.CombinedOutput()
						if err == nil && len(strings.TrimSpace(string(output))) > 0 {
							fmt.Fprintf(file, "%s\n", output)
							success = true
							break
						}
					}
				}
			}
		}
		
		if !success {
			fmt.Fprintf(file, "All commands for '%s' failed\n", cmd.description)
		}
	}
	
	fmt.Println("All information saved to hwid_info.txt")
}
