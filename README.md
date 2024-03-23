# Blockchain Activity

## About

It's a service that displays the top 5 addresses on the ETH mainnet with activity metrics for the last 100 blocks, the activity metric increases if the wallet sends or receives any ERC20

Based on GetBlock

## Quickstart

0. Put your `GetBlock` access token in `secrets/rpc-access-token` instead of `your-token`

1. Start server

```
    make start
```

2. Open endpoint in your browser and wait a few seconds

```
http://127.0.0.1:8080/api/v1/top
```

3. Stop server

```
    make down
```

4. Start tests

```
    make test
```

5. Start lint

```
    make lint
```

## Project Structure

**cmd** - contains entry point to start the server

**config** - app's config with env parsing

**internal** - internal app logic

* **service** - contains business logic
* **controllers** - contains logic for working with web protocols, in my case only with http
* **adapters** - contains logic related to working with external systems, in my case with GetBlock RPC
* **utils**
    * **erc20** - tool for work with ERC-20

**pkg** - packages for general use
* **jsonrpc** - client for JSON-RPC

**secrets** - contains sensative data for service