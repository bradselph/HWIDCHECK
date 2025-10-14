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

### Advanced Features
- **Timestamped File Output:** Save all information to uniquely timestamped text files (format: `hwid_info_YYYY-MM-DD_HH-MM-SS.txt`) to prevent overwriting previous scans
- **Real-time Progress Tracking:** Visual progress indicators showing completion percentage and status during full system scans
- **Comprehensive Error Handling:** Robust error handling with detailed logging and graceful recovery from failures
- **Multiple Fallback Commands:** Automatic command fallback when primary commands fail, ensuring maximum compatibility
- **PowerShell Integration:** Support for both WMI (legacy) and CIM (modern) PowerShell cmdlets
- **Detailed Status Reporting:** Success/failure tracking with execution time and success rate statistics
- **Enhanced User Experience:** Clear visual feedback with [Starting], [Complete], [Success], and [Failed] status tags
- **Interactive Operation:** "Press Enter to continue" prompts between operations for better control

## Download

You can download the pre-compiled executable for Windows from the [Releases](https://github.com/bradselph/HWIDCHECK/releases) section of this repository. Simply download the `HWIDCHECK.exe` file from the latest release.

## Requirements

- Windows operating system
- PowerShell (any version, improved functionality with PowerShell 3.0+)
- Administrator privileges recommended for complete hardware information access

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
11. Print All to File and Save
12. Exit
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

#### Save All Information to File
```
Enter your choice: 11

========================================
Starting full system scan...
Output file: hwid_info_2024-10-14_15-30-45.txt
========================================

[1/14] (0.0%) Processing: SMBIOS (UUID)
[1/14] SUCCESS - SMBIOS (UUID)
[2/14] (7.1%) Processing: BIOS (Serial Number)
[2/14] SUCCESS - BIOS (Serial Number)
[3/14] (14.3%) Processing: Motherboard (Serial Number)
[3/14] SUCCESS - Motherboard (Serial Number)
...
[14/14] (92.9%) Processing: MAC Addresses (IPConfig)
[14/14] SUCCESS - MAC Addresses (IPConfig)

========================================
           Scan Complete
========================================
Total Commands: 14
Successful: 14
Failed: 0
Success Rate: 100.0%
Execution Time: 5.234s
Output saved to: hwid_info_2024-10-14_15-30-45.txt
========================================

Press Enter to continue...
```

### Output File Format

When you select option 11, a timestamped file is created with the following structure:
```
========================================
  Hardware ID Information Report
========================================
Generated: 2024-10-14 15:30:45
System: Windows
========================================

[1/14] SMBIOS (UUID)
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
Total Commands Executed: 14
Successful: 14
Failed: 0
Success Rate: 100.0%
Execution Time: 5.234s
Completion Time: 2024-10-14 15:30:50
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

## License

This project is licensed under the AGPL-3.0 License. See the [LICENSE](LICENSE) file for more details.

## Changelog
- Added timestamped file output to prevent overwriting previous scans
- Implemented real-time progress tracking with percentage indicators
- Improved visual feedback with status tags
- Added interactive prompts between operations
- Refactored code to eliminate duplication and improve maintainability
- All error results are now properly handled