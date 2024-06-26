<h1 align=center> <img src="./logo.png"/> </h1>
<h1 align=center><code>merry</code></h1>

Streamline Garden component testing with a user-friendly CLI tool that quickly sets up isolated testing environment.

## Features

- **Quick Setup**: Set up a testing environment with a single command.
- **Isolated Environment**: Test your application in an isolated environment without depending on external services.
- **Orderbook**: Fully functional orderbook along with COBI to get you started with cross-chain swapping.
- **Bitcoin regtest**: Get a bitcoin regtest node to test your Bitcoin specific logic.
- **Ethereum testnet**: Get an Ethereum testnet node to test your Ethereum specific logic.

## Installation

You can install `merry` using the following command.

```bash
curl https://get.merry.dev | bash
```

## Quickstart

To get started, run the following command.

```bash
merry go
```

This command will pull docker images (if they do not exist already) of various garden components and launches a testing environment.

When you're finished testing, simply stop the `merry` with this command:

```bash
merry stop
```

## Documentation

The documentation for `merry` is available [here](https://docs.garden.finance/developers/merry).

## Testing

Once you have the environment set up, you can connect to orderbook by providing the orderbook URL in your client.
With `merry`, you also get bitcoin regtest and ethereum testnet nodes to test your application.

## Contributing

If you're interested in contributing to `merry`, there are no special requirements. Just fork the repository, make your changes, and submit a pull request. Happy coding!
