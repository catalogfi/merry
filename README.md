![GitHub Release](https://img.shields.io/github/v/release/catalogfi/merry)

<h1> <img height="36px" src="./logo.png">  <span>Merry</span> </h1>

Streamline your multi-chain testing with `Merry`!

This CLI tool leverages Docker to effortlessly set up a multi-chain testing environment in a single command. Merry includes Bitcoin regtest node, Ethereum localnet node, and essential [Catalog.fi](https://www.catalog.fi/) services, providing a self-contained space to test your applications independently of external services.

It supports various features, including a faucet, electrum services, and an orderbook with COBI.

## Prerequisites

Before using Merry, please ensure you have Docker installed and running on your system. If not, download and install Docker from the official [website](https://www.docker.com).

## Installation

You can install Merry using the following command.

```bash
curl https://get.merry.dev | bash
```

Merry stores its configuration files and other data in a directory on your system, typically named `~/.merry/`.

See the [Install from scratch](#install-from-scratch) section to install Merry from scratch.

## Commands

Merry provides a variety of commands to manage your testing environment:

### Starting Merry

```bash
merry go
```

This starts all services, including the Bitcoin regtest node, Ethereum localnet node, explorers for the nodes, and the Catalog services.

- `--bare` flag: Starts only the multi-chain services (Bitcoin and Ethereum nodes with explorers) and excludes Catalog services. This option is helpful if you don't need the additional functionalities like COBI and Orderbook provided by Catalog.

- `--headless` flag: Starts all services except for frontend interfaces. This can help run Merry in headless environments (e.g., servers) where a graphical user interface is not required.

### Stopping Merry

```bash
merry stop

# Reset data
merry stop -d
```

Stops all running services. Use `--delete` or `-d` to remove data.

### Getting logs

```bash
merry logs -s <service>

# Getting logs of EVM service
merry logs -s evm
```

Replace <service> with the specific service (e.g., cobi, orderbook, evm) to view its logs.

### Replacing a service with a local one

```bash
merry replace <service>
```

This command allows you to replace a service with your local development version. Make sure you're in the directory containing the local service's Dockerfile. Supported services include COBI, Orderbook, and EVM.

### Calling Bitcoin RPC methods

```bash
merry rpc <method> <params>

# example: get blockchain info
merry rpc getblockchaininfo
```

Interact with the Bitcoin regtest node directly using RPC methods.

### Updating Docker images

```bash
merry update
```

Keep your testing environment up-to-date by updating all Docker images.

### Fund accounts

```bash
merry faucet <address>
```

Fund any EVM, Bitcoin or Solana address for testing purposes. Replace `<address>` with the address you want to fund. It could be a Bitcoin, Ethereum or Solana address.

### List active services

```bash
merry status
```

Lists all currently running services with their respective ports and operational status.

### List all commands

```bash
merry --help
```

## Testing with Merry

Once your environment is set up:

- Connect to the Orderbook using its provided URL within your client application.
- Leverage the built-in Bitcoin regtest and Ethereum testnet nodes to test your multi-chain functionalities.

Contributing

We welcome contributions to Merry! There are no special requirements needed. Fork the repository, make your changes, and submit a pull request.

Let Merry simplify your multi-chain testing journey!

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

# Service Endpoints

The following services are available with their respective ports:

## Blockchain Nodes
- **Ethereum**: 8545
- **Arbitrum**: 8546
- **Starknet DevNet**: 8547
- **Bitcoin**: 8080/tcp
- **Solana Validator**: 8899-8900

## Explorers
- **Ethereum Explorer**: 5100
- **Arbitrum Explorer**: 5101
- **Esplora**: 5050

## Relayers & Executors
- **Relayer**: 4426
- **Starknet Relayer**: 4436
- **Solana Relayer**: 5014
- **Starknet Executor**: 3000/tcp

## Blockchain Tools
- **Electrs**: 30000
- **Chopsticks**: 3000
- **Cosigner**: 11818
- **Quote**: 6969
- **Virtual Balance**: 3008

## Databases
- **Postgres**: 5432
- **Garden DB**: 5433
- **Redis**: 6379