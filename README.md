[![Build and Release](https://github.com/bradselph/HWIDCHECK/actions/workflows/main.yml/badge.svg)](https://github.com/bradselph/HWIDCHECK/actions/workflows/main.yml)
# HWID Checker

HWID Checker is a Go application that allows you to easily gather various hardware information about your system, such as SMBIOS UUID, BIOS serial number, motherboard serial number, CPU serial number, and more. The application provides a user-friendly CLI menu for selecting and displaying the information, and it also offers the option to save all information to a text file.

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
- Save all information to a text file
- Improved error handling and informative error messages
- Enhanced compatibility with different Windows versions
- **NEW:** PowerShell fallback commands for more reliable hardware information retrieval
- **NEW:** Automatic command fallback when primary commands fail
- **NEW:** Support for both WMI (legacy) and CIM (modern) PowerShell cmdlets

## Download

You can download the pre-compiled executable for Windows from the [Releases](https://github.com/bradselph/HWIDCHECK/releases) section of this repository. Simply download the `HWIDCHECK.exe` file from the latest release.

## Requirements

- Windows operating system
- PowerShell (any version, improved functionality with PowerShell 3.0+)
- Improved compatibility with various Windows versions
- Experimental support for other operating systems (see note below)

**Note on OS Support:** While HWID Checker is primarily developed for Windows, the Go programming language allows for cross-platform compilation. Depending on the specific system calls and libraries used, it may be possible to compile and run this application on other operating systems like Linux or macOS. However, some features might be limited or require modifications to work correctly on non-Windows systems. If you're interested in using HWID Checker on another operating system, you may need to modify the source code and compile it yourself.

## Installation

### Option 1: Using the pre-compiled executable (Windows)

1. Download `HWIDCHECK.exe` from the [Releases](https://github.com/bradselph/HWIDCHECK/releases) section.
2. Run the downloaded `HWIDCHECK.exe` file.

### Option 2: Building from source

If you prefer to build the application yourself, or if you want to try running it on a non-Windows system, follow these steps:

1. Ensure you have Go 1.16 or higher installed on your system.

2. Clone the repository:
   ```bash
   git clone https://github.com/bradselph/HWIDCHECK.git
   cd HWIDCHECK
   ```

3. Build the application:
   - For Windows:
     ```bash
     go build -o HWIDCHECK.exe
     ```
   - For Linux/macOS:
     ```bash
     go build -o HWIDCHECK
     ```

4. Run the application:
   - For Windows:
     ```bash
     .\HWIDCHECK.exe
     ```
   - For Linux/macOS:
     ```bash
     ./HWIDCHECK
     ```

## Usage

When you run the application, you will be presented with a menu of options to choose from. Simply enter the number corresponding to the information you want to display.

### Menu Options

1. **SMBIOS (UUID):** Display the SMBIOS UUID.
2. **BIOS (Serial Number):** Display the BIOS serial number.
3. **Motherboard (Serial Number):** Display the motherboard serial number.
4. **Chassis (Serial Number):** Display the chassis serial number.
5. **CPU (Serial Number):** Display the CPU serial number.
6. **HDD/SSD (Serial Number):** Display the HDD/SSD serial number.
7. **Volume Information:** Display volume information.
8. **RAM (Serial Number):** Display the RAM serial number.
9. **Windows Product ID:** Display the Windows product ID.
10. **MAC Addresses:** Display the MAC addresses.
11. **Print All to File and Save:** Save all information to a text file (`hwid_info.txt`).
12. **Exit:** Exit the application.

### Example Usage

1. **Select an option:** Enter the number corresponding to your choice and press Enter.
   ```
   HWID Checker
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
   Enter your choice: 1
   ```

2. **View the output:** The selected information will be displayed on the screen.
   ```
   SMBIOS (UUID)
   UUID
   XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
   ```

3. **Save all information to a file:** Select option `11` to save all gathered information to a text file (`hwid_info.txt`).
   ```
   Enter your choice: 11
   All information saved to hwid_info.txt
   ```

### Fallback Mechanism

The application now features an intelligent fallback mechanism:

1. When a primary command fails (e.g., `wmic` commands), the application automatically tries alternative PowerShell commands
2. Multiple fallback options ensure maximum compatibility across different Windows versions
3. The application will notify you when falling back to alternative commands
4. All fallbacks are also used when saving information to a file

## Technical Details

The application uses multiple methods to retrieve hardware information:

1. **Primary Method:** Traditional Windows Management Instrumentation Command-line (WMIC)
2. **Fallback Methods:**
   - PowerShell with Get-WmiObject cmdlet
   - PowerShell with Get-CimInstance cmdlet (preferred on newer Windows versions)
   - Other specialized PowerShell commands where applicable

This multi-layered approach ensures reliable hardware information retrieval even if some commands are deprecated or unavailable on your specific Windows version.

## Contributing

If you would like to contribute to the project, please fork the repository and create a pull request with your changes. We welcome contributions that improve error handling, cross-platform compatibility, or add new features.

## License

This project is licensed under the AGPL-3.0 License. See the [LICENSE](LICENSE) file for more details.
