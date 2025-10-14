package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	primary   []string
	fallbacks [][]string
}

type CommandResult struct {
	success bool
	output  string
	error   string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========================================")
		fmt.Println("           HWID Checker")
		fmt.Println("========================================")
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
		fmt.Println("========================================")

		fmt.Print("Enter your choice: ")
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			logError(fmt.Sprintf("Error reading input: %s", err))
			fmt.Println("Press Enter to continue...")
			if _, readErr := reader.ReadString('\n'); readErr != nil {
				logError(fmt.Sprintf("Failed to read continuation: %s", readErr))
			}
			continue
		}
		choiceStr = strings.TrimSpace(choiceStr)

		if choiceStr == "" {
			fmt.Println("No input provided. Please enter a valid option.")
			continue
		}

		switch choiceStr {
		case "1":
			fmt.Println("\n[Starting] SMBIOS UUID Check...")
			runCommandWithFallbacks("SMBIOS (UUID)", Command{
				primary: []string{"wmic", "csproduct", "get", "uuid"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_ComputerSystemProduct).UUID"},
				},
			})
			fmt.Println("[Complete] SMBIOS UUID Check finished")
		case "2":
			fmt.Println("\n[Starting] BIOS Serial Number Check...")
			runCommandWithFallbacks("BIOS (Serial Number)", Command{
				primary: []string{"wmic", "bios", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_BIOS | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_BIOS).SerialNumber"},
				},
			})
			fmt.Println("[Complete] BIOS Serial Number Check finished")
		case "3":
			fmt.Println("\n[Starting] Motherboard Serial Number Check...")
			runCommandWithFallbacks("Motherboard (Serial Number)", Command{
				primary: []string{"wmic", "baseboard", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_BaseBoard | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_BaseBoard).SerialNumber"},
				},
			})
			fmt.Println("[Complete] Motherboard Serial Number Check finished")
		case "4":
			fmt.Println("\n[Starting] Chassis Serial Number Check...")
			runCommandWithFallbacks("Chassis (Serial Number)", Command{
				primary: []string{"wmic", "systemenclosure", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_SystemEnclosure | Select-Object -ExpandProperty SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_SystemEnclosure).SerialNumber"},
				},
			})
			fmt.Println("[Complete] Chassis Serial Number Check finished")
		case "5":
			fmt.Println("\n[Starting] CPU Serial Number Check...")
			runCommandWithFallbacks("CPU (Serial Number)", Command{
				primary: []string{"wmic", "cpu", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_Processor).ProcessorId"},
				},
			})
			fmt.Println("[Complete] CPU Serial Number Check finished")
		case "6":
			fmt.Println("\n[Starting] HDD/SSD Serial Number Check...")
			runCommandWithFallbacks("HDD/SSD (Serial Number)", Command{
				primary: []string{"wmic", "diskdrive", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_DiskDrive | Select-Object Model, SerialNumber"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_DiskDrive | Select-Object Model, SerialNumber"},
					{"powershell", "-Command", "Get-PhysicalDisk | Select-Object FriendlyName, SerialNumber"},
				},
			})
			fmt.Println("[Complete] HDD/SSD Serial Number Check finished")
		case "7":
			fmt.Println("\n[Starting] Volume Information Check...")
			runCommandWithFallbacks("Volume Information", Command{
				primary: []string{"vol"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-Volume"},
					{"powershell", "-Command", "Get-WmiObject Win32_LogicalDisk | Select-Object DeviceID, VolumeName, VolumeSerialNumber"},
				},
			})
			fmt.Println("[Complete] Volume Information Check finished")
		case "8":
			fmt.Println("\n[Starting] RAM Serial Number Check...")
			runCommandWithFallbacks("RAM (Serial Number)", Command{
				primary: []string{"wmic", "memorychip", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_PhysicalMemory | Select-Object DeviceLocator, SerialNumber"},
				},
			})
			fmt.Println("[Complete] RAM Serial Number Check finished")
		case "9":
			fmt.Println("\n[Starting] Windows Product ID Check...")
			runCommandWithFallbacks("Windows Product ID", Command{
				primary: []string{"wmic", "os", "get", "serialnumber"},
				fallbacks: [][]string{
					{"powershell", "-Command", "(Get-WmiObject -Class Win32_OperatingSystem).SerialNumber"},
					{"powershell", "-Command", "(Get-CimInstance -Class Win32_OperatingSystem).SerialNumber"},
					{"powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' -Name ProductId | Select-Object -ExpandProperty ProductId"},
				},
			})
			fmt.Println("\n[Starting] Windows Product ID Check (Alternative)...")
			runCommandWithFallbacks("Windows Product ID (Alternative)", Command{
				primary: []string{"systeminfo", "|", "findstr", "/B", "/C:\"OS Serial Number\""},
				fallbacks: [][]string{
					{"powershell", "-Command", "systeminfo | Select-String 'OS Serial Number'"},
				},
			})
			fmt.Println("[Complete] Windows Product ID Check finished")
		case "10":
			fmt.Println("\n[Starting] MAC Addresses Check (1/4)...")
			runCommandWithFallbacks("MAC Addresses (GetMac)", Command{
				primary: []string{"getmac", "/v"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-NetAdapter | Select-Object Name, Status, MacAddress"},
				},
			})
			fmt.Println("\n[Starting] MAC Addresses Check (2/4)...")
			runCommandWithFallbacks("MAC Addresses (WMIC Path)", Command{
				primary: []string{"wmic", "path", "Win32_NetworkAdapter", "where", `"MacAddress like '%%:%%:%%:%%:%%:%%'"`, "get", "Name, MacAddress"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
					{"powershell", "-Command", "Get-CimInstance -Class Win32_NetworkAdapter | Where-Object { $_.MacAddress -ne $null } | Select-Object Name, MacAddress"},
				},
			})
			fmt.Println("\n[Starting] MAC Addresses Check (3/4)...")
			runCommandWithFallbacks("MAC Addresses (WMIC NIC)", Command{
				primary: []string{"wmic", "nic", "get", "Name, MACAddress"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object { $_.IPEnabled -eq $true } | Select-Object Description, MACAddress"},
				},
			})
			fmt.Println("\n[Starting] MAC Addresses Check (4/4)...")
			runCommandWithFallbacks("MAC Addresses (IPConfig)", Command{
				primary: []string{"ipconfig", "/all", "|", "findstr", `"Physical Address"`},
				fallbacks: [][]string{
					{"powershell", "-Command", "ipconfig /all | Select-String 'Physical Address'"},
				},
			})
			fmt.Println("[Complete] MAC Addresses Check finished")
		case "11":
			saveAllToFile()
		case "12":
			fmt.Println("\nExiting HWID Checker...")
			return
		default:
			fmt.Printf("Invalid choice '%s'. Please enter a number between 1-12.\n", choiceStr)
		}

		fmt.Println("\nPress Enter to continue...")
		if _, err := reader.ReadString('\n'); err != nil {
			logError(fmt.Sprintf("Error waiting for input: %s", err))
		}
	}
}

func logError(msg string) {
	if _, err := fmt.Fprintln(os.Stderr, msg); err != nil {
		fmt.Println("Error logging failed:", err)
	}
}

func runCommandWithFallbacks(description string, command Command) bool {
	if len(command.primary) == 0 {
		logError(fmt.Sprintf("[Error] No primary command provided for '%s'", description))
		return false
	}

	fmt.Printf("[Attempt 1] Executing primary command...\n")
	result := executeCommandWithResult(command.primary)

	if result.success {
		fmt.Printf("[Success] Command completed successfully\n")
		displayCommandOutput(result.output)
		return true
	}

	logError(fmt.Sprintf("[Failed] Primary command error: %s", result.error))

	if len(command.fallbacks) == 0 {
		logError(fmt.Sprintf("[Failed] No fallback commands available for '%s'", description))
		return false
	}

	for i, fallback := range command.fallbacks {
		if len(fallback) == 0 {
			logError(fmt.Sprintf("[Warning] Skipping empty fallback %d", i+1))
			continue
		}

		fmt.Printf("[Attempt %d] Trying fallback command...\n", i+2)
		result = executeCommandWithResult(fallback)

		if result.success {
			fmt.Printf("[Success] Fallback %d completed successfully\n", i+1)
			displayCommandOutput(result.output)
			return true
		}

		logError(fmt.Sprintf("[Failed] Fallback %d error: %s", i+1, result.error))
	}

	logError(fmt.Sprintf("[Failed] All commands for '%s' failed", description))
	return false
}

func executeCommandWithResult(args []string) CommandResult {
	if len(args) == 0 {
		return CommandResult{
			success: false,
			error:   "No command provided",
		}
	}

	fmt.Printf("Command: %s\n", strings.Join(args, " "))

	if containsPipe(args) {
		return executePipedCommandWithResult(args)
	}

	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		return CommandResult{
			success: false,
			error:   fmt.Sprintf("Command '%s' not found in system PATH: %v", args[0], err),
		}
	}

	cmd := exec.Command(cmdPath, args[1:]...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		errMsg := fmt.Sprintf("Execution failed: %v", err)
		if len(outputStr) > 0 {
			errMsg += fmt.Sprintf(" | Output: %s", outputStr)
		}
		return CommandResult{
			success: false,
			error:   errMsg,
			output:  outputStr,
		}
	}

	if len(outputStr) == 0 {
		return CommandResult{
			success: false,
			error:   "Command returned empty output",
		}
	}

	return CommandResult{
		success: true,
		output:  outputStr,
	}
}

func executePipedCommandWithResult(args []string) CommandResult {
	if len(args) == 0 {
		return CommandResult{
			success: false,
			error:   "Empty piped command",
		}
	}

	fullCommand := strings.Join(args, " ")
	cmd := exec.Command("cmd.exe", "/C", fullCommand)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		errMsg := fmt.Sprintf("Piped command execution failed: %v", err)
		if len(outputStr) > 0 {
			errMsg += fmt.Sprintf(" | Output: %s", outputStr)
		}
		return CommandResult{
			success: false,
			error:   errMsg,
			output:  outputStr,
		}
	}

	if len(outputStr) == 0 {
		return CommandResult{
			success: false,
			error:   "Piped command returned empty output",
		}
	}

	return CommandResult{
		success: true,
		output:  outputStr,
	}
}

func displayCommandOutput(output string) {
	if len(output) == 0 {
		return
	}
	fmt.Println("----------------------------------------")
	fmt.Println(output)
	fmt.Println("----------------------------------------")
}

func containsPipe(args []string) bool {
	for _, arg := range args {
		if arg == "|" {
			return true
		}
	}
	return false
}

func generateTimestampedFilename() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("hwid_info_%s.txt", timestamp)
}

func saveAllToFile() {
	filename := generateTimestampedFilename()

	fmt.Printf("\n========================================\n")
	fmt.Printf("Starting full system scan...\n")
	fmt.Printf("Output file: %s\n", filename)
	fmt.Printf("========================================\n\n")

	file, err := os.Create(filename)
	if err != nil {
		logError(fmt.Sprintf("[Error] Failed to create file '%s': %s", filename, err))
		fmt.Println("\nPress Enter to continue...")
		reader := bufio.NewReader(os.Stdin)
		if _, readErr := reader.ReadString('\n'); readErr != nil {
			logError(fmt.Sprintf("Error reading input: %s", readErr))
		}
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logError(fmt.Sprintf("[Warning] Error closing file: %s", closeErr))
		}
	}()

	if err := writeFileHeader(file); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write file header: %s", err))
		return
	}

	commands := buildCommandList()
	totalCommands := len(commands)
	successCount := 0
	failureCount := 0
	startTime := time.Now()

	for idx, cmdEntry := range commands {
		progress := fmt.Sprintf("[%d/%d]", idx+1, totalCommands)
		percentComplete := float64(idx) / float64(totalCommands) * 100

		fmt.Printf("%s (%.1f%%) Processing: %s\n", progress, percentComplete, cmdEntry.description)

		success := processCommandForFile(file, cmdEntry, progress)

		if success {
			successCount++
			fmt.Printf("%s SUCCESS - %s\n", progress, cmdEntry.description)
		} else {
			failureCount++
			fmt.Printf("%s FAILED - %s\n", progress, cmdEntry.description)
		}
	}

	elapsed := time.Since(startTime)

	if err := writeFileSummary(file, totalCommands, successCount, failureCount, elapsed); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write summary: %s", err))
	}

	displayScanSummary(totalCommands, successCount, failureCount, elapsed, filename)

	fmt.Println("\nPress Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		logError(fmt.Sprintf("Error reading input: %s", err))
	}
}

func writeFileHeader(file *os.File) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	header := fmt.Sprintf(
		"========================================\n"+
			"  Hardware ID Information Report\n"+
			"========================================\n"+
			"Generated: %s\n"+
			"System: Windows\n"+
			"========================================\n\n",
		time.Now().Format("2006-01-02 15:04:05"),
	)

	_, err := fmt.Fprint(file, header)
	return err
}

func writeFileSummary(file *os.File, total, success, failure int, elapsed time.Duration) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	summary := fmt.Sprintf(
		"\n========================================\n"+
			"           Report Summary\n"+
			"========================================\n"+
			"Total Commands Executed: %d\n"+
			"Successful: %d\n"+
			"Failed: %d\n"+
			"Success Rate: %.1f%%\n"+
			"Execution Time: %s\n"+
			"Completion Time: %s\n"+
			"========================================\n",
		total,
		success,
		failure,
		float64(success)/float64(total)*100,
		elapsed.Round(time.Millisecond),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	_, err := fmt.Fprint(file, summary)
	return err
}

func displayScanSummary(total, success, failure int, elapsed time.Duration, filename string) {
	fmt.Printf("\n========================================\n")
	fmt.Printf("           Scan Complete\n")
	fmt.Printf("========================================\n")
	fmt.Printf("Total Commands: %d\n", total)
	fmt.Printf("Successful: %d\n", success)
	fmt.Printf("Failed: %d\n", failure)
	fmt.Printf("Success Rate: %.1f%%\n", float64(success)/float64(total)*100)
	fmt.Printf("Execution Time: %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Output saved to: %s\n", filename)
	fmt.Printf("========================================\n")
}

type FileCommandEntry struct {
	description string
	command     Command
}

func buildCommandList() []FileCommandEntry {
	return []FileCommandEntry{
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
}

func processCommandForFile(file *os.File, cmdEntry FileCommandEntry, progress string) bool {
	if file == nil {
		return false
	}

	header := fmt.Sprintf("\n%s %s\n"+
		"========================================\n"+
		"Primary Command: %s\n",
		progress,
		cmdEntry.description,
		strings.Join(cmdEntry.command.primary, " "))

	if _, err := fmt.Fprint(file, header); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write command header: %s", err))
		return false
	}

	result := executeCommandWithResult(cmdEntry.command.primary)

	if result.success {
		if err := writeCommandResult(file, result, true); err != nil {
			logError(fmt.Sprintf("[Error] Failed to write success result: %s", err))
		}
		return true
	}

	if err := writeCommandResult(file, result, false); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write failure result: %s", err))
	}

	for i, fallback := range cmdEntry.command.fallbacks {
		if len(fallback) == 0 {
			if err := writeSkippedFallback(file, i+1); err != nil {
				logError(fmt.Sprintf("[Error] Failed to write skipped fallback: %s", err))
			}
			continue
		}

		if err := writeFallbackAttempt(file, i+1, fallback); err != nil {
			logError(fmt.Sprintf("[Error] Failed to write fallback attempt: %s", err))
		}

		result = executeCommandWithResult(fallback)

		if result.success {
			if err := writeCommandResult(file, result, true); err != nil {
				logError(fmt.Sprintf("[Error] Failed to write fallback success: %s", err))
			}
			return true
		}

		if err := writeCommandResult(file, result, false); err != nil {
			logError(fmt.Sprintf("[Error] Failed to write fallback failure: %s", err))
		}
	}

	if err := writeFinalFailure(file); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write final failure: %s", err))
	}
	return false
}

func writeCommandResult(file *os.File, result CommandResult, success bool) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	var output string
	if success {
		output = fmt.Sprintf("Status: SUCCESS\nOutput:\n%s\n", result.output)
	} else {
		output = fmt.Sprintf("Status: FAILED\nError: %s\n", result.error)
		if len(result.output) > 0 {
			output += fmt.Sprintf("Partial Output: %s\n", result.output)
		}
	}

	_, err := fmt.Fprint(file, output)
	return err
}

func writeSkippedFallback(file *os.File, index int) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	_, err := fmt.Fprintf(file, "\nFallback %d: SKIPPED (empty command)\n", index)
	return err
}

func writeFallbackAttempt(file *os.File, index int, fallback []string) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	_, err := fmt.Fprintf(file, "\nAttempting Fallback %d:\nCommand: %s\n", index, strings.Join(fallback, " "))
	return err
}

func writeFinalFailure(file *os.File) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	_, err := fmt.Fprint(file, "\nFinal Status: ALL ATTEMPTS FAILED\n")
	return err
}
