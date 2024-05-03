# Hex - CLI tool designed to streamline and speed up QGP Deployments

## Introduction

**Hex** is a CLI tool, developed in GoLang, designed to streamline and speed up QGP deployments. Hex effectively simplifies the whole process by wrapping around our existing QGP CLI.

Why Hex? Well, we all know how time-consuming and repetitive deployments can be, especially on ephemeral environments. The usual routine involves hopping between multiple files to configure snapshot tags for each repository, which isn't exactly anyone’s idea of fun.

Here's where Hex steps in. It removes all the unnecessary hassle by providing pre-built configurations. Imagine creating a deploy file for all selected services, including relevant dependencies, with just a few commands. That’s Hex for you—making deployments fast and simple.
[Demo Video](static%2Fdemo_vid.mov)

## Installing Hex

To start using Hex, you have two main options for installation:

### 1. Install from the Latest Release (Recommended)
Follow these steps to download and install Hex directly from the latest release:
- Navigate to the **releases section** of the Hex repository and click on the release tagged as 'latest'.
- Download the binary file named `hex`.
- Move to the folder where `hex` is downloaded and grant it executable permissions with the following command:
    ```bash
    chmod +x ./hex
    ```
- Relocate this executable to your local binary path to make it accessible from anywhere in the terminal:
    ```bash
    mv ./hex /usr/local/bin/
    ```

### 2. Build from Source
If you prefer to build Hex from the source code, follow these instructions:
- Clone the repository to your local machine.
- Ensure you have GoLang version 1.22 installed on your system.
- Navigate to the project directory and install the required packages:
    ```bash
    go mod tidy
    ```
- Build the project for your system architecture (example shown for macOS):
    ```bash
    GOOS=darwin GOARCH=amd64 go build -o hex
    ```
- Move the newly created `hex` executable to your local binary path to execute it from anywhere in the terminal:
    ```bash
    mv ./hex /usr/local/bin/
    ```

By following these steps, you'll have Hex installed and ready to streamline your QGP deployments.

## Getting Started with Hex

Once you have completed the installation steps, using Hex is straightforward. Follow these steps to begin deploying services using Hex from any terminal:

### Starting Hex
- Open your terminal.
- Type `hex` and press Enter. This will start the Hex interface.
- Use `Ctrl+c` to quite Hex interface anytime.

### Selecting Services to Deploy
- Hex currently supports four services with hardcoded configurations:
  - Finance Calcy Service
  - Finance Job Service
  - Finance Dashboard
  - Finance Orchestrator
- Use the Up and Down arrow keys to navigate through the list of services.
- Type the name of a service to search for it.
- Use the Spacebar to select the services you want to deploy.
- Press Enter to confirm your selections and move forward.

### Configuring Deployment
- You will be prompted to enter the snapshot tags for each of the selected services. You can type these manually or copy and paste them directly from GitHub.
- Next, you will be asked to enter the name of the environment instance where you wish to deploy your services.

### Generating Deployment Configuration
- Hex will generate a deployment configuration file based on your inputs. This file will be stored in the `/usr/local/bin/` directory.

### Deploying Services
- Confirm the details, and Hex will initiate the deployment using the GQP command.
- Hex will display a message indicating whether the deployment was successful or if it failed.

By following these steps, you'll be able to deploy services efficiently using Hex.

## Next Steps for Hex Development

Hex is poised for further development with numerous enhancements on the horizon. Here are some planned additions to improve functionality:

### Planned Enhancements
1. **Database Integration:**
  - Connect Hex CLI with a database to manage configurations for different services dynamically, enhancing flexibility and scalability.

2. **Improved Service Display:**
  - Modify the current service listing to show the top 4 relevant results instead of the full list, making the terminal interface more user-friendly.

3. **Status Indicators:**
  - Introduce status indicators for API calls, such as fetching service configurations and deploying to environments. This will provide real-time feedback and improve the user experience during operations.

4. **Graceful Exit:**
  - Implement a feature to quit the CLI gracefully post-deployment instead of using `ctrl+c`. This will enhance the tool’s usability and professionalism.

These updates aim to streamline operations and enhance the overall functionality of Hex, making it an even more robust tool for deployment processes.
