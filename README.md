![GitHub Release](https://img.shields.io/github/v/release/catalogfi/merry)

<h1 style="display:flex;align-items:center;gap:0.5rem;" > <img height="32" src="./logo.png">  <span>Merry</span> </h1>

Streamline your multi-chain testing with `merry`!

This CLI tool leverages Docker to effortlessly set up a multi-chain testing environment in a single command. Merry includes Bitcoin regtest node, Ethereum localnet node, and essential Garden components, providing a self-contained space to test your applications independently of external services.

It supports a variety of features, including a faucet, electrum services and an orderbook with COBI.

## Prerequisites

Before using merry, ensure you have Docker installed and running on your system. If not, you'll need to download and install Docker from the official [website](https://www.docker.com).

## Installation

You can install merry using the following command.

```bash
curl https://get.merry.dev | bash
```

merry stores its configuration files and other data in a directory on your system. This directory is typically named `.merry/volumes/`.

See the [Install from scratch](#install-from-scratch) section to install merry from scratch.

## Commands

merry provides a variety of commands to manage your testing environment:

### Starting merry

To start merry, run the following command:

```bash
merry go
```

### Stopping merry

```bash
merry stop
```

Stops all running services

### Getting logs

```bash
merry logs -s <service>

# getting logs of evm service
merry logs -s evm
```

Replace <service> with the specific service (e.g., cobi, orderbook, evm) to view its logs.

### Replacing a service with a local one

```bash
merry replace <service>
```

This command allows you to replace a service with your local development version. Make sure you're in the directory containing the local service's Dockerfile. Supported services include cobi, orderbook, and evm.

### Calling bitcoin rpc methods

```bash
merry rpc <method> <params>

# example: get blockchain info
merry rpc getblockchaininfo
```

Interact with the Bitcoin regtest node directly using RPC methods.

### Updating docker images

```bash
merry update
```

Keep your testing environment up-to-date by updating all docker images.

### Fund accounts

```bash
merry faucet <address>
```

Fund any EVM or Bitcoin address for testing purposes. Replace <address> with the address you want to fund. It could be a Bitcoin or Ethereum address.

## Testing with merry

Once your environment is set up:

- Connect to the orderbook using its provided URL within your client application.
- Leverage the built-in Bitcoin regtest and Ethereum testnet nodes to test your multi-chain functionalities.

Contributing

We welcome contributions to merry! No special requirements are needed. Simply fork the repository, make your changes, and submit a pull request.

Let merry simplify your multi-chain testing journey!

## Install from scratch

- Clone the repository

```bash
git clone https://github.com/catalogfi/merry.git
```

- Building and installing

```bash
cd cmd/merry
# build and install the binary
go install
```
