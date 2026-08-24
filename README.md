[![Build and Release](https://github.com/bradselph/HWIDCHECK/actions/workflows/main.yml/badge.svg)](https://github.com/bradselph/HWIDCHECK/actions/workflows/main.yml)
# HWID Checker

HWID Checker is a Go application that allows you to easily gather various hardware information about your system, such as SMBIOS UUID, BIOS serial number, motherboard serial number, CPU serial number, and more. The application provides a user-friendly CLI menu for selecting and displaying the information, and it also offers the option to save all information to a timestamped text file with detailed progress tracking.

## Features

- Display SMBIOS UUID
- Display BIOS serial number
- Display motherboard serial number
- Display chassis serial number
- Display CPU serial number
- Display HDD/SSD serial number
- Display volume information
- Display RAM serial number
- Display Windows product ID
- Display MAC addresses
- Display TPM status and TPM Endorsement Key
- Display Secure Boot status

### Advanced Features
- **Timestamped File Output:** Save all information to uniquely timestamped text files (format: `hwid_info_YYYY-MM-DD_HH-MM-SS.txt`) to prevent overwriting previous scans
- **Clean HWID List:** Generate a compact, human-readable summary (`hwid_clean_YYYY-MM-DD_HH-MM-SS.txt`) with one line per identifier instead of raw command output
- **Scan Comparison:** Compare two previous scan files to see what changed between runs
- **Administrator Detection:** Reports whether the process is elevated, both on-screen and in the file header, since TPM/Secure Boot/Endorsement Key checks require Administrator privileges to return real data
- **Real-time Progress Tracking:** Visual progress indicators showing completion percentage and status during full system scans
- **Comprehensive Error Handling:** Robust error handling with detailed logging and graceful recovery from failures, including detection of commands that report a privilege error or return an empty result while still exiting successfully
- **Multiple Fallback Commands:** Automatic command fallback when primary commands fail, ensuring maximum compatibility
- **Locale-Independent:** Checks use CIM/PowerShell property access instead of parsing localized command-line text, so results are consistent regardless of the system's display language
- **PowerShell Integration:** Support for both WMI (legacy) and CIM (modern) PowerShell cmdlets
- **Detailed Status Reporting:** Success/failure tracking with execution time and success rate statistics
- **Enhanced User Experience:** Clear visual feedback with [Starting], [Complete], [Success], and [Failed] status tags
- **Interactive Operation:** "Press Enter to continue" prompts between operations for better control

## Download

You can download the pre-compiled executable for Windows from the [Releases](https://github.com/bradselph/HWIDCHECK/releases) section of this repository. Simply download the `HWIDCHECK.exe` file from the latest release.

## Requirements

- Windows operating system
- PowerShell (any version, improved functionality with PowerShell 3.0+)
- Administrator privileges required for complete results — TPM Status, TPM Endorsement Key, and Secure Boot all fail or return incomplete data without elevation. The application warns on startup and records `Administrator: Yes/No` in every report if not run elevated.

**Note on OS Support:** While HWID Checker is primarily developed and optimized for Windows, the Go programming language allows for cross-platform compilation. However, this tool relies heavily on Windows-specific commands (WMIC, PowerShell, cmd.exe) and will not function correctly on other operating systems without significant modifications. If you're interested in using HWID Checker on another operating system, you would need to modify the source code to use platform-appropriate system commands.

## Installation

### Option 1: Using the pre-compiled executable (Windows)

1. Download `HWIDCHECK.exe` from the [Releases](https://github.com/bradselph/HWIDCHECK/releases) section.
2. Run the downloaded `HWIDCHECK.exe` file.
3. Follow the on-screen menu to select the information you want to display.

### Option 2: Building from source

If you prefer to build the application yourself:

1. Ensure you have Go 1.16 or higher installed on your system.

2. Clone the repository:
```bash
   git clone https://github.com/bradselph/HWIDCHECK.git
   cd HWIDCHECK
```

3. Build the application:
```bash
   go build -o HWIDCHECK.exe
```

4. Run the application:
```bash
   .\HWIDCHECK.exe
```

## Usage

When you run the application, you will be presented with a menu of options to choose from. Simply enter the number corresponding to the information you want to display.

### Menu Options
```
========================================
           HWID Checker
========================================
Select an option:
1. SMBIOS (UUID)
2. BIOS (Serial Number)
3. Motherboard (Serial Number)
4. Chassis (Serial Number)
5. CPU (Serial Number)
6. HDD/SSD (Serial Number)
7. Volume Information
8. RAM (Serial Number)
9. Windows Product ID
10. MAC Addresses
11. TPM and Secure Boot Status
12. Print All to File (Detailed)
13. Print Clean HWID List
14. Compare with Previous Scan
15. Exit
========================================
```

### Example Usage

#### Individual Hardware Check
```
Enter your choice: 1

[Starting] SMBIOS UUID Check...
[Attempt 1] Executing primary command...
Command: wmic csproduct get uuid
[Success] Command completed successfully
----------------------------------------
UUID
XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
----------------------------------------
[Complete] SMBIOS UUID Check finished

Press Enter to continue...
```

#### Save All Information to File (Detailed)
```
Enter your choice: 12

========================================
Starting full system scan...
Output file: hwid_info_2026-08-23_21-08-45.txt
========================================

[1/17] (0.0%) Processing: SMBIOS (UUID)
[1/17] SUCCESS - SMBIOS (UUID)
[2/17] (5.9%) Processing: BIOS (Serial Number)
[2/17] SUCCESS - BIOS (Serial Number)
[3/17] (11.8%) Processing: Motherboard (Serial Number)
[3/17] SUCCESS - Motherboard (Serial Number)
...
[17/17] (94.1%) Processing: Secure Boot
[17/17] SUCCESS - Secure Boot

========================================
           Scan Complete
========================================
Total Commands: 17
Successful: 17
Failed: 0
Success Rate: 100.0%
Execution Time: 6.419s
Output saved to: hwid_info_2026-08-23_21-08-45.txt
========================================

Press Enter to continue...
```

#### Print Clean HWID List
```
Enter your choice: 13

========================================
Generating clean HWID list...
Output file: hwid_clean_2026-08-23_21-00-53.txt
========================================
...
```
Produces a compact summary file, one line per identifier:
```
SMBIOS (UUID): XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
Motherboard (Serial Number): XXXXXXXXXXXXXXXX
Volume Information: C: 470.88 GB free of 952.98 GB, G: 14.83 GB free of 953.85 GB, E: 205.34 GB free of 447.11 GB
MAC Addresses (GetMac): XX:XX:XX:XX:XX:XX, XX:XX:XX:XX:XX:XX, ...
TPM Status: Activated=True, Enabled=True, Owned=True, SpecVersion=2.0, 0, 1.38
Secure Boot: True
```

### Output File Format

When you select option 12, a timestamped file is created with the following structure:
```
========================================
  Hardware ID Information Report
========================================
Generated: 2026-08-23 21:08:45
System: Windows
Administrator: Yes
========================================

[1/17] SMBIOS (UUID)
========================================
Primary Command: wmic csproduct get uuid
Status: SUCCESS
Output:
UUID
XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX

...

========================================
           Report Summary
========================================
Total Commands Executed: 17
Successful: 17
Failed: 0
Success Rate: 100.0%
Execution Time: 6.419s
Completion Time: 2026-08-23 21:08:52
========================================
```

### Fallback Mechanism

The application features an multi-tier fallback mechanism:

1. **Primary Command:** Traditional Windows Management Instrumentation Command-line (WMIC)
2. **First Fallback:** PowerShell with Get-WmiObject cmdlet (Windows 7/8/10 compatible)
3. **Second Fallback:** PowerShell with Get-CimInstance cmdlet (Windows 8+ preferred method)

When a command fails, the application automatically:
- Attempts the next available fallback command
- Logs the error with detailed information
- Continues with the remaining checks
- Reports which method succeeded in the output file

**Example Fallback Sequence:**
```
[Attempt 1] Executing primary command...
Command: wmic csproduct get uuid
[Failed] Primary command error: Command 'wmic' not found in system PATH
[Attempt 2] Trying fallback command...
Command: powershell -Command Get-WmiObject Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID
[Success] Fallback 1 completed successfully
```

### Progress Tracking

Real-time progress indication includes:
- Sequential counter: `[X/Y]`
- Percentage completion: `(XX.X%)`
- Current operation description
- Status updates: SUCCESS or FAILED
- Total execution time
- Overall success rate

## Troubleshooting

### "Command not found" errors
- Run the application as Administrator
- Ensure PowerShell is installed and available in system PATH
- Try running individual commands manually to verify system configuration

### Empty or missing information
- Some hardware doesn't provide serial numbers or may report them as "To Be Filled By O.E.M."
- Virtual machines may have limited hardware information
- Run as Administrator for complete access

### File creation errors
- Ensure you have write permissions in the current directory
- Check available disk space
- Verify antivirus isn't blocking file creation

### TPM Status, TPM Endorsement Key, or Secure Boot report FAILED
- These checks require Administrator privileges; the report header's `Administrator:` line confirms whether the process was elevated
- Even when elevated, some hardware/firmware combinations don't expose a TPM Endorsement Key — the report will show the actual PowerShell error message rather than a blank result

## License

This project is licensed under the AGPL-3.0 License. See the [LICENSE](LICENSE) file for more details.

## Changelog
- Added TPM Endorsement Key check as a real per-chip TPM identifier
- Added Administrator-privilege detection, with a startup warning and an `Administrator:` field in every report
- Fixed commands (e.g. `Get-Tpm`, `Confirm-SecureBootUEFI`, `Get-TpmEndorsementKeyInfo`) being logged as SUCCESS when they actually printed a privilege error or returned an empty result
- Fixed a Windows command-line quoting bug that broke `findstr`-based piped commands
- Replaced locale-dependent `findstr` label parsing (English-only `systeminfo`/`ipconfig` labels) with locale-independent CIM/PowerShell equivalents
- Fixed the Clean HWID List producing garbled, unreadable output for MAC address, volume, and TPM checks
- Added TPM and Secure Boot status checks
- Implemented Clean HWID List generation (option 13)
- Added scan comparison functionality (option 14)
- Added timestamped file output to prevent overwriting previous scans
- Implemented real-time progress tracking with percentage indicators
- Improved visual feedback with status tags
- Added interactive prompts between operations
- Refactored code to eliminate duplication and improve maintainability
- All error results are now properly handled