# Rapid Bridge Service Documentation

Rapid Bridge is a cryptographic proxy designed to facilitate secure, authenticated, and verifiable communication between a bank's internal systems (Rapid) and trusted third-party services. It removes the burden of implementing complex security protocols, ensuring sensitive data is protected in transit and accessible only to authorized parties.

## Project Purpose

The primary purpose of Rapid Bridge is to solve the challenge of secure and authenticated communication between third-party applications and bank Rapid. It acts as a intermediary, handling all cryptographic operations to ensure:

- **Data Protection**: Sensitive data remains encrypted while in transit.
- **Authentication**: Only authorized parties can send or receive information.
- **Verification**: Data origin and integrity are guaranteed.

## Core Functionality

- **Secure Message Relay**: Rapid Bridge receives plaintext requests from trusted third-party services, then signs and encrypts them before forwarding them to the bank's Rapid system. It performs the reverse operation for responses, decrypting and verifying them before returning plaintext to the requester.
- **End-to-End Security**: All data passing through the bridge is both encrypted and digitally signed, ensuring both privacy and authenticity.
- **Automated Cryptography**: The bridge transparently handles all cryptographic operations, including encryption, decryption, signing, and verification.

## Use Cases

- **Banking Integrations**: Integrating with new banks becomes effortless, streamlining the process without exposing sensitive infrastructure or cryptographic keys.
- **Regulated Data Exchange**: Ideal for any scenario where regulated or sensitive data must be exchanged between a secure backend and external services, with strong guarantees of privacy, authenticity, and integrity.
- **Simplified Security for Partners**: Third-party developers can interact with the bridge using simple, plaintext JSON APIs, as the bridge manages all underlying security requirements.

## Key Features

- **Encryption**: All data in transit is encrypted, ensuring only the intended recipient can read it.
- **Digital Signatures**: Requests are digitally signed to prove their origin and prevent tampering.
- **Signature Verification**: All responses are verified to ensure they originate from a trusted source.
- **Transparent Security**: Clients interact with the bridge using simple, plaintext JSON, with all cryptographic operations handled internally.
- **Extensible API**: The bridge exposes a set of RESTful endpoints for common banking operations (e.g., balance, statement, payment initiation/approval).

## Service Usage

Third-party services interact with Rapid Bridge by sending HTTP POST requests to specific API endpoints. The bridge handles all security, meaning clients only need to send and receive JSON payloads. The bridge then communicates securely with the bank's Rapid system, returning decrypted and verified responses to the client.

## API Documentation

### Base URL
```
/api/v1/resource
```

### Endpoints
- `POST /api/v1/resource/balance`
- `POST /api/v1/resource/statement`

### Request Body
```json
{
  "Message": "ISO20022 XML MESSAGE PAYLOAD"
}
```

### Required Headers
- `X-Source-Slug`
- `X-Destination-Slug`
- `X-Key-Version`

## Rapid Bridge CLI Documentation

The Rapid Bridge CLI is a command-line tool designed for initializing and managing application and bank cryptographic configurations for the Rapid Bridge backend.

### Building the CLI

To build the CLI tool yourself, simply run:

```bash
go build -o rapid-bridge ./cmd/main.go
```

This command will produce the `rapid-bridge` executable for CLI operations.

### Usage

```bash
rapid-bridge init [app|bank|server] [flags]
```

## Commands

### 1. init app

Initializes an application configuration.

**Usage:**
```bash
rapid-bridge init app --slug <application-slug>
```

**Required Flags:**
- `--slug`: The unique identifier for the application.

**Workflow:**
1. Checks if the application is already registered.
2. If registered, prompts whether to re-initialize.
3. Prompts to either:
    - Generate a new key pair (RSA and Ed25519), or
    - Use your own existing key pair (prompts for file paths).
4. Stores key files and configuration under `_rapid_bridge_data/application/<slug>/<ulid>/`.
5. Updates the CLI configuration and saves it to disk.

**Interactive Prompts:**
- Choice to re-initialize if already registered.
- Choice to generate or provide keys.
- If providing keys, prompts for file paths to RSA and Ed25519 public/private keys.

### 2. init bank

Initializes a bank configuration.

**Usage:**
```bash
rapid-bridge init bank --slug <bank-slug> --rapidUrl <rapid-url>
```

**Required Flags:**
- `--slug`: The unique identifier for the bank.
- `--rapidUrl`: The Rapid Bridge service URL.

**Workflow:**
1. Checks if the bank is already registered.
2. If registered, prompts whether to re-initialize.
3. Prompts to either:
    - Fetch the bank's public keys from the Rapid Bridge service, or
    - Provide your own public key files (prompts for file paths).
4. Stores key files and configuration under `_rapid_bridge_data/bank/<slug>/`.
5. Updates the CLI configuration and saves it to disk.

**Interactive Prompts:**
- Choice to re-initialize if already registered.
- Choice to fetch or provide keys.
- If providing keys, prompts for file paths to RSA and Ed25519 public keys.

### 3. init server

Initializes the backend server configuration.

**Usage:**
```bash
rapid-bridge init server
```

**Workflow:**
This command is registered in the CLI, but the specific flags and interactive prompts depend on the implementation in `cmd/server/server.go`. Typically, it will set up the backend server environment and configuration.

## General Notes

- All commands support the `--help` flag for more information.
- Configuration and key files are stored under the `_rapid_bridge_data` directory.
- All initialization commands are interactive and will prompt for user input as needed.
- Only the flags and options described above are currently supported.

## CLI Environment Requirements

To use the Rapid Bridge CLI, ensure the following environment is set up:

### .env file
A `.env` file must be present in the same directory as the `rapid-bridge` executable. This file must contain the `SERVER_PORT` variable, for example:

```env
SERVER_PORT=8080
```

### _rapid_bridge_data folder
A folder named `_rapid_bridge_data` must exist in the same directory as the `rapid-bridge` executable.

Inside the `_rapid_bridge_data` folder, there must be a `core.json` file.

The `core.json` file must contain the `"rapid_links_url"` key, for example:

```json
{
  "rapid_links_url": "http://localhost:9000/rapid-links"
}
```

For further details, run any command with the `--help` flag or refer to the source code in the `cmd/cli` and `cmd/server` directories.