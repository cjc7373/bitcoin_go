This project is based on [Jeiwan/blockchain_go](https://github.com/Jeiwan/blockchain_go) and his related blog post.

This project is only a simplified bitcoin implementation not intended to be compatible with current bitcoin protocols.

## Designs
### Roles
There are no roles such as miner, full nodes and SPVs. Each node will store the whole chain, do PoWs and validate blocks. There is a seperate Cli tool to act like a wallet to send/receive coins.

### Tech Stack
- Serialization: protobuf, this is used both for network traffic and data persistence
- Network communication: grpc
- Key-Value Database: bbolt
- Cli framework: cobra

### Implemetation details
- node discovery: one node must manually specify a node to connect, then it can use RPC to get other nodes' info
    - since our nodes' number is relatively small, all nodes are connected directly to each other
    - a node periodically announces its connected nodes to all its neighbours

## References
- [Bitcoin whitepaper](https://bitcoin.org/bitcoin.pdf)
- [Bitcoin Protocol Documentation](https://en.bitcoin.it/wiki/Protocol_documentation)

## Things which are not implemented
- Merkle Tree
- Scripts

## Issues
- TODO: when verifying a block, need to check if an output is spent twice
- when to update uxto set?
- we can't use an UTXO twice in one block