# mc-dispute-mon

The `mc-dispute-mon` is an off-chain service to monitor dispute games.

## Quickstart

Clone this repo. Then run:

```shell
just mc-dispute-mon
```

This will build the `mc-dispute-mon` binary which can be run with
`./mc-dispute-mon/bin/mc-dispute-mon`.

## Usage

`mc-dispute-mon` is configurable via command line flags and environment variables. The help menu
shows the available config options and can be accessed by running `./bin/mc-dispute-mon --help`.

```shell

# Start the mc-dispute-mon with predefined network and RPC endpoints
./bin/mc-dispute-mon \
  --network <Predefined-Network> \
  --l1-eth-rpc <L1-Ethereum-RPC-URL> \
  --rollup-rpc <Rollup-RPC-URL>,<Secondary-RPC-URL>,<Tertiary-RPC-URL>

# For networks using super root RPCs:
./bin/mc-dispute-mon \
  --network <Predefined-Network> \
  --l1-eth-rpc <L1-Ethereum-RPC-URL> \
  --superroot-rpc <SuperRoot-RPC-URL>,<Secondary-RPC-URL>,<Tertiary-RPC-URL>

```
