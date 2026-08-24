package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

type HWIDData struct {
	name  string
	value string
}

// isRunningAsAdmin reports whether the process has administrator privileges.
// Opening a physical drive handle requires elevation on Windows, so a
// successful open is a reliable signal without extra dependencies.
func isRunningAsAdmin() bool {
	f, err := os.Open(`\\.\PHYSICALDRIVE0`)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

// isPrivilegeError detects command output that reports missing admin rights
// even though the process itself exited successfully (e.g. Get-Tpm prints a
// localized "requires administrator privileges" message to stdout and still
// returns exit code 0). Without this check that output gets recorded as a
// SUCCESS with garbage content instead of a clear FAILED result.
func isPrivilegeError(output string) bool {
	lower := strings.ToLower(output)
	phrases := []string{
		"se requiere privilegios de administrador",
		"acceso denegado",
		"privilegios adecuados",
		"access is denied",
		"access denied",
		"administrator privileges are required",
		"run as administrator",
		"requires elevation",
		"you must run this cmdlet from an elevated",
	}
	for _, phrase := range phrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	if !isRunningAsAdmin() {
		fmt.Println("\n[Warning] Not running as Administrator — TPM, Secure Boot, and some")
		fmt.Println("          other checks will fail or return incomplete data.")
		fmt.Println("          Re-launch this program as Administrator for full results.")
	}

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
		fmt.Println("11. TPM and Secure Boot Status")
		fmt.Println("12. Print All to File (Detailed)")
		fmt.Println("13. Print Clean HWID List")
		fmt.Println("14. Compare with Previous Scan")
		fmt.Println("15. Exit")
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
			fmt.Println("\n[Starting] TPM and Secure Boot Check...")
			runCommandWithFallbacks("TPM Status", Command{
				primary: []string{"powershell", "-Command", "Get-WmiObject -Namespace ROOT\\CIMV2\\Security\\MicrosoftTpm -Class Win32_Tpm | Select-Object IsActivated_InitialValue, IsEnabled_InitialValue, IsOwned_InitialValue, ManufacturerVersion, PhysicalPresenceVersionInfo, SpecVersion"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-Tpm"},
					{"powershell", "-Command", "Get-CimInstance -Namespace ROOT\\CIMV2\\Security\\MicrosoftTpm -ClassName Win32_Tpm"},
				},
			})
			fmt.Println("\n[Checking] Secure Boot Status...")
			runCommandWithFallbacks("Secure Boot", Command{
				primary: []string{"powershell", "-Command", "Confirm-SecureBootUEFI"},
				fallbacks: [][]string{
					{"powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State' -Name UEFISecureBootEnabled | Select-Object -ExpandProperty UEFISecureBootEnabled"},
				},
			})
			fmt.Println("[Complete] TPM and Secure Boot Check finished")
		case "12":
			saveAllToFile(false)
		case "13":
			saveAllToFile(true)
		case "14":
			compareScans(reader)
		case "15":
			fmt.Println("\nExiting HWID Checker...")
			return
		default:
			fmt.Printf("Invalid choice '%s'. Please enter a number between 1-15.\n", choiceStr)
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

	if isPrivilegeError(outputStr) {
		return CommandResult{
			success: false,
			error:   "Administrator privileges required",
			output:  outputStr,
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

	if isPrivilegeError(outputStr) {
		return CommandResult{
			success: false,
			error:   "Administrator privileges required",
			output:  outputStr,
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

func generateTimestampedFilename(cleanList bool) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	if cleanList {
		return fmt.Sprintf("hwid_clean_%s.txt", timestamp)
	}
	return fmt.Sprintf("hwid_info_%s.txt", timestamp)
}

func saveAllToFile(cleanList bool) {
	filename := generateTimestampedFilename(cleanList)

	if cleanList {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Generating clean HWID list...\n")
		fmt.Printf("Output file: %s\n", filename)
		fmt.Printf("========================================\n\n")
	} else {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Starting full system scan...\n")
		fmt.Printf("Output file: %s\n", filename)
		fmt.Printf("========================================\n\n")
	}

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

	if err := writeFileHeader(file, cleanList); err != nil {
		logError(fmt.Sprintf("[Error] Failed to write file header: %s", err))
		return
	}

	commands := buildCommandList()
	totalCommands := len(commands)
	successCount := 0
	failureCount := 0
	startTime := time.Now()
	var hwidList []HWIDData

	for idx, cmdEntry := range commands {
		progress := fmt.Sprintf("[%d/%d]", idx+1, totalCommands)
		percentComplete := float64(idx) / float64(totalCommands) * 100

		fmt.Printf("%s (%.1f%%) Processing: %s\n", progress, percentComplete, cmdEntry.description)

		var success bool
		var cleanData string

		if cleanList {
			success, cleanData = processCommandForCleanList(cmdEntry)
			if success && cleanData != "" {
				hwidList = append(hwidList, HWIDData{
					name:  cmdEntry.description,
					value: cleanData,
				})
			}
		} else {
			success = processCommandForFile(file, cmdEntry, progress)
		}

		if success {
			successCount++
			fmt.Printf("%s SUCCESS - %s\n", progress, cmdEntry.description)
		} else {
			failureCount++
			fmt.Printf("%s FAILED - %s\n", progress, cmdEntry.description)
		}
	}

	elapsed := time.Since(startTime)

	if cleanList {
		if err := writeCleanList(file, hwidList); err != nil {
			logError(fmt.Sprintf("[Error] Failed to write clean list: %s", err))
		}
	}

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

func processCommandForCleanList(cmdEntry FileCommandEntry) (bool, string) {
	result := executeCommandWithResult(cmdEntry.command.primary)

	if result.success {
		return true, extractCleanValue(result.output)
	}

	for _, fallback := range cmdEntry.command.fallbacks {
		if len(fallback) == 0 {
			continue
		}

		result = executeCommandWithResult(fallback)
		if result.success {
			return true, extractCleanValue(result.output)
		}
	}

	return false, ""
}

func extractCleanValue(output string) string {
	lines := strings.Split(output, "\n")
	var values []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.Contains(strings.ToLower(line), "serialnumber") ||
			strings.Contains(strings.ToLower(line), "uuid") ||
			strings.Contains(strings.ToLower(line), "macaddress") ||
			strings.Contains(strings.ToLower(line), "name") ||
			strings.Contains(strings.ToLower(line), "model") ||
			strings.Contains(strings.ToLower(line), "deviceid") ||
			strings.Contains(strings.ToLower(line), "description") ||
			strings.Contains(strings.ToLower(line), "friendlyname") ||
			strings.Contains(strings.ToLower(line), "devicelocator") ||
			strings.Contains(strings.ToLower(line), "status") ||
			strings.Contains(strings.ToLower(line), "volumename") ||
			strings.Contains(strings.ToLower(line), "volumeserialnumber") {
			continue
		}

		if len(line) > 0 && !strings.HasPrefix(line, "-") {
			values = append(values, line)
		}
	}

	if len(values) == 0 {
		return "Not Available"
	}

	return strings.Join(values, " | ")
}

func writeCleanList(file *os.File, hwidList []HWIDData) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	for _, hwid := range hwidList {
		line := fmt.Sprintf("%s: %s\n", hwid.name, hwid.value)
		if _, err := fmt.Fprint(file, line); err != nil {
			return err
		}
	}

	return nil
}

func compareScans(reader *bufio.Reader) {
	fmt.Println("\n========================================")
	fmt.Println("      Compare Previous Scans")
	fmt.Println("========================================")

	files, err := filepath.Glob("hwid_*.txt")
	if err != nil {
		logError(fmt.Sprintf("Error finding scan files: %s", err))
		return
	}

	if len(files) == 0 {
		fmt.Println("No previous scan files found.")
		fmt.Println("Please run option 12 or 13 first to generate scan files.")
		return
	}

	fmt.Println("\nAvailable scan files:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	fmt.Print("\nEnter first file number: ")
	file1Str, err := reader.ReadString('\n')
	if err != nil {
		logError(fmt.Sprintf("Error reading input: %s", err))
		return
	}
	file1Str = strings.TrimSpace(file1Str)

	var file1Idx int
	if _, err := fmt.Sscanf(file1Str, "%d", &file1Idx); err != nil || file1Idx < 1 || file1Idx > len(files) {
		fmt.Println("Invalid file number.")
		return
	}

	fmt.Print("Enter second file number: ")
	file2Str, err := reader.ReadString('\n')
	if err != nil {
		logError(fmt.Sprintf("Error reading input: %s", err))
		return
	}
	file2Str = strings.TrimSpace(file2Str)

	var file2Idx int
	if _, err := fmt.Sscanf(file2Str, "%d", &file2Idx); err != nil || file2Idx < 1 || file2Idx > len(files) {
		fmt.Println("Invalid file number.")
		return
	}

	file1Path := files[file1Idx-1]
	file2Path := files[file2Idx-1]

	fmt.Printf("\nComparing:\n  File 1: %s\n  File 2: %s\n\n", file1Path, file2Path)

	data1, err := parseHWIDFile(file1Path)
	if err != nil {
		logError(fmt.Sprintf("Error reading file 1: %s", err))
		return
	}

	data2, err := parseHWIDFile(file2Path)
	if err != nil {
		logError(fmt.Sprintf("Error reading file 2: %s", err))
		return
	}

	compareResults := compareHWIDData(data1, data2)

	outputFilename := fmt.Sprintf("hwid_comparison_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		logError(fmt.Sprintf("Error creating comparison file: %s", err))
		return
	}
	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			logError(fmt.Sprintf("Error closing comparison file: %s", closeErr))
		}
	}()

	header := fmt.Sprintf(
		"========================================\n"+
			"     HWID Comparison Report\n"+
			"========================================\n"+
			"File 1: %s\n"+
			"File 2: %s\n"+
			"Comparison Date: %s\n"+
			"========================================\n\n",
		file1Path,
		file2Path,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if _, err := fmt.Fprint(outputFile, header); err != nil {
		logError(fmt.Sprintf("Error writing header: %s", err))
	}

	fmt.Println("Comparison Results:")
	fmt.Println("========================================")

	unchangedCount := 0
	changedCount := 0
	addedCount := 0
	removedCount := 0

	for _, result := range compareResults {
		var statusLine string
		switch result.status {
		case "UNCHANGED":
			statusLine = fmt.Sprintf("[=] %s: %s\n", result.name, result.value1)
			unchangedCount++
		case "CHANGED":
			statusLine = fmt.Sprintf("[!] %s:\n    File 1: %s\n    File 2: %s\n", result.name, result.value1, result.value2)
			changedCount++
			fmt.Printf("[CHANGED] %s\n", result.name)
		case "ADDED":
			statusLine = fmt.Sprintf("[+] %s: %s (only in File 2)\n", result.name, result.value2)
			addedCount++
			fmt.Printf("[ADDED] %s\n", result.name)
		case "REMOVED":
			statusLine = fmt.Sprintf("[-] %s: %s (only in File 1)\n", result.name, result.value1)
			removedCount++
			fmt.Printf("[REMOVED] %s\n", result.name)
		}

		if _, err := fmt.Fprint(outputFile, statusLine); err != nil {
			logError(fmt.Sprintf("Error writing comparison line: %s", err))
		}
	}

	summary := fmt.Sprintf(
		"\n========================================\n"+
			"        Comparison Summary\n"+
			"========================================\n"+
			"Total Items: %d\n"+
			"Unchanged: %d\n"+
			"Changed: %d\n"+
			"Added: %d\n"+
			"Removed: %d\n"+
			"========================================\n",
		len(compareResults),
		unchangedCount,
		changedCount,
		addedCount,
		removedCount,
	)

	if _, err := fmt.Fprint(outputFile, summary); err != nil {
		logError(fmt.Sprintf("Error writing summary: %s", err))
	}

	fmt.Println("========================================")
	fmt.Printf("Unchanged: %d\n", unchangedCount)
	fmt.Printf("Changed: %d\n", changedCount)
	fmt.Printf("Added: %d\n", addedCount)
	fmt.Printf("Removed: %d\n", removedCount)
	fmt.Printf("\nDetailed comparison saved to: %s\n", outputFilename)
}

type ComparisonResult struct {
	name   string
	status string
	value1 string
	value2 string
}

func parseHWIDFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logError(fmt.Sprintf("Error closing file: %s", closeErr))
		}
	}()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	currentKey := ""
	var currentValue strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			if currentKey != "" && currentValue.Len() > 0 {
				data[currentKey] = strings.TrimSpace(currentValue.String())
			}

			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				currentKey = strings.TrimSpace(parts[1])
				currentValue.Reset()
			}
		} else if strings.Contains(line, ":") && !strings.Contains(line, "Output:") && !strings.Contains(line, "Status:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				data[key] = value
			}
		} else if currentKey != "" && strings.TrimSpace(line) != "" &&
			!strings.HasPrefix(line, "=") &&
			!strings.HasPrefix(line, "Command:") &&
			!strings.HasPrefix(line, "Status:") &&
			!strings.HasPrefix(line, "Primary") &&
			!strings.Contains(line, "Report") {
			if currentValue.Len() > 0 {
				currentValue.WriteString(" | ")
			}
			currentValue.WriteString(strings.TrimSpace(line))
		}
	}

	if currentKey != "" && currentValue.Len() > 0 {
		data[currentKey] = strings.TrimSpace(currentValue.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func compareHWIDData(data1, data2 map[string]string) []ComparisonResult {
	var results []ComparisonResult
	allKeys := make(map[string]bool)

	for key := range data1 {
		allKeys[key] = true
	}
	for key := range data2 {
		allKeys[key] = true
	}

	for key := range allKeys {
		val1, exists1 := data1[key]
		val2, exists2 := data2[key]

		if exists1 && exists2 {
			if val1 == val2 {
				results = append(results, ComparisonResult{
					name:   key,
					status: "UNCHANGED",
					value1: val1,
					value2: val2,
				})
			} else {
				results = append(results, ComparisonResult{
					name:   key,
					status: "CHANGED",
					value1: val1,
					value2: val2,
				})
			}
		} else if exists1 {
			results = append(results, ComparisonResult{
				name:   key,
				status: "REMOVED",
				value1: val1,
				value2: "",
			})
		} else {
			results = append(results, ComparisonResult{
				name:   key,
				status: "ADDED",
				value1: "",
				value2: val2,
			})
		}
	}

	return results
}

func writeFileHeader(file *os.File, cleanList bool) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	adminStatus := "No (run as Administrator for TPM/Secure Boot/full results)"
	if isRunningAsAdmin() {
		adminStatus = "Yes"
	}

	var header string
	if cleanList {
		header = fmt.Sprintf(
			"========================================\n"+
				"    Clean HWID List\n"+
				"========================================\n"+
				"Generated: %s\n"+
				"System: Windows\n"+
				"Administrator: %s\n"+
				"========================================\n\n",
			time.Now().Format("2006-01-02 15:04:05"),
			adminStatus,
		)
	} else {
		header = fmt.Sprintf(
			"========================================\n"+
				"  Hardware ID Information Report\n"+
				"========================================\n"+
				"Generated: %s\n"+
				"System: Windows\n"+
				"Administrator: %s\n"+
				"========================================\n\n",
			time.Now().Format("2006-01-02 15:04:05"),
			adminStatus,
		)
	}

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
		{"TPM Status", Command{
			primary: []string{"powershell", "-Command", "Get-WmiObject -Namespace ROOT\\CIMV2\\Security\\MicrosoftTpm -Class Win32_Tpm | Select-Object IsActivated_InitialValue, IsEnabled_InitialValue, IsOwned_InitialValue, ManufacturerVersion, PhysicalPresenceVersionInfo, SpecVersion"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-Tpm"},
				{"powershell", "-Command", "Get-CimInstance -Namespace ROOT\\CIMV2\\Security\\MicrosoftTpm -ClassName Win32_Tpm"},
			},
		}},
		{"Secure Boot", Command{
			primary: []string{"powershell", "-Command", "Confirm-SecureBootUEFI"},
			fallbacks: [][]string{
				{"powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State' -Name UEFISecureBootEnabled | Select-Object -ExpandProperty UEFISecureBootEnabled"},
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
