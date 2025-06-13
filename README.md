![gophish logo](https://raw.github.com/gophish/gophish/master/static/images/gophish_purple.png)

Gophish
=======

![Build Status](https://github.com/ezops-br/gophish/workflows/CI/badge.svg) [![GoDoc](https://godoc.org/github.com/ezops-br/gophish?status.svg)](https://godoc.org/github.com/ezops-br/gophish)

Gophish: Open-Source Phishing Toolkit

[Gophish](https://getgophish.com) is an open-source phishing toolkit designed for businesses and penetration testers. It provides the ability to quickly and easily setup and execute phishing engagements and security awareness training.

---

## Python Requirement

> **Important:**  
> Gophish requires **Python 3.11 or newer** to be installed on your system **before you can run the application**.

### How to Install Python

1. **Download Python:**  
   Go to the [official Python downloads page](https://www.python.org/downloads/) and download the latest stable version (Python 3.x) installer for your operating system (Windows or Linux).

2. **Windows Users:**  
   - Run the installer.
   - **Make sure to check the box that says "Add Python to PATH"** before clicking "Install Now".
   - Complete the installation.

3. **Linux Users:**  
   - Most distributions include Python 3 by default.  
   - To install or upgrade, use your package manager. For example, on Ubuntu/Debian:  
     ```sh
     sudo apt update
     sudo apt install python3 python3-venv python3-pip
     ```

4. **Verify Python Installation:**  
   Open a new terminal or command prompt and run:
   ```sh
   python --version
   ```
   You should see a version number like `Python 3.8.x`, `3.9.x`, `3.10.x`, etc.

   If you see an error or a version lower than 3.8, please check your installation and ensure Python is added to your system's PATH.

---

## LibreOffice Requirement

> **Important:**  
> To generate PDF reports, Gophish requires [LibreOffice](https://www.libreoffice.org/download/download/) to be installed on your system.

### How to Install LibreOffice

**Linux (Debian/Ubuntu):**
```sh
sudo apt update
sudo apt install libreoffice
```

**Windows:**
- Download and install from [https://www.libreoffice.org/download/download/](https://www.libreoffice.org/download/download/)
- During installation, ensure the option to add LibreOffice to your system PATH is enabled, or manually add the LibreOffice `program` directory (e.g., `C:\Program Files\LibreOffice\program`) to your PATH.

**After installation, you should be able to run `soffice --version` from your terminal or command prompt.**

---

## Getting Started

### Install

Installation of Gophish is dead-simple - just download and extract the zip containing the [release for your system](https://github.com/ezops-br/gophish/releases/), and run the binary. Gophish has binary releases for Windows and Linux platforms.

### Building From Source
**If you are building from source, please note that Gophish requires Go v1.10 or above!**

To build Gophish from source, simply run ```git clone https://github.com/ezops-br/gophish.git``` and ```cd``` into the project source directory. Then, run ```go build```. After this, you should have a binary called ```gophish``` in the current directory.

### Database Setup

- The release comes with an **empty `gophish.db`** file (SQLite database) in the project root directory.
- **If you already have a Gophish SQLite database** and want to use it:
  - **Copy your own database file** into the project root directory, replacing the existing `gophish.db` (**name must be `gophish.db`**).

> **Disclaimer:** We strongly recommend you keep a backup of your original database before replacing any files.

### Running Gophish

- Run the Gophish binary for your platform:
  - On Windows: double-click `gophish.exe` or run it from the command prompt.
  - On Linux: run `./gophish` from the terminal.
    - > **Note for Linux users:** If you get a "Permission denied" error, you may need to give execute permission to the binary first:
      > ```sh
      > chmod u+x gophish
      > ```
- The application will start and display log output in the terminal.

### Initial Login

- **Default login:**  
  - **Username:** `admin`  
  - **Password:** The password will appear in the terminal output when you first run Gophish.
- Example log output:
  ```
  time="2020-07-29T01:24:08Z" level=info msg="Please login with the username admin and the password 4304d5255378177d"
  ```
- > **Note:** This applies only if you are using the default (empty) database provided in the release. If you use your own database, your existing users and passwords will be used.

### Accessing the Web Interface

- Open your browser and go to:  
  `http://127.0.0.1:3333` (or the address shown in the terminal)
- Log in with the credentials above.

---

## Using the Report Template Feature

### Template Requirements

When creating or uploading a DOCX template for campaign reports, **your template must include the following tags**:

- `{{sender}}` – Will be replaced with the sender information.
- `{{subject}}` – Will be replaced with the email subject.
- `{{screenshot}}` – Will be replaced with a screenshot of the email.

**Place these tags in your template wherever you want the corresponding information to appear in the generated document.**

---

## Building From Source

**If you are building from source, please note that Gophish requires Go v1.10 or above!**

To build Gophish from source, simply run:

```sh
git clone https://github.com/ezops-br/gophish.git
cd gophish
go build
```

After this, you should have a binary called `gophish` in the current directory.
