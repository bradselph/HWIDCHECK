# HWID Checker

HWID Checker is a Go application that allows you to easily gather various hardware information about your system, such as SMBIOS UUID, BIOS serial number, motherboard serial number, CPU serial number, and more. The application provides a user-friendly menu for selecting and displaying the information, and it also offers the option to save all information to a text file.

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

## Requirements

- Go 1.16 or higher

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/bradselph/HWIDCHECK.git
   cd HWIDCHECK
   ```

2. **Build the application:**
   ```bash
   go build -o HWIDCHECK
   ```

3. **Run the application:**
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

   ```plaintext
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

   ```plaintext
   SMBIOS (UUID)
   UUID
   XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
   ```

3. **Save all information to a file:** Select option `11` to save all gathered information to a text file (`hwid_info.txt`).

   ```plaintext
   Enter your choice: 11
   All information saved to hwid_info.txt
   ```

## Contributing

If you would like to contribute to the project, please fork the repository and create a pull request with your changes.

## License

This project is licensed under the AGPL-3.0 License. See the [LICENSE](LICENSE) file for more details.