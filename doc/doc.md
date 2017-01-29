## Overview

Note: go-resonate is in the early (early) stages of development.

### Identity
The demo currently uses Ed25519 public-key cryptography for user identification and verification. We are looking at other identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack).

TODO: Test compatibility of BigchainDB crypto-conditions with Ed25519 library 

### Blockchain
This is not a blockchain application (yet). We are considering various platforms and consensus engines such as [Tendermint](https://github.com/tendermint), which would be used for the ordering of transactions (e.g. uploads, plays), state replication across a network, and payments. This infrastructure would be in addition to BigchainDB/IPDB.

## Directories

### API
A client-side API. It's currently designed to run locally and communicate with BigchainDB/IPDB over http, though this may change. Working on functionality for the following actions:
- Login/Register via parter org
- Persist a track to the db
- Stream a track
- other TBD

Note: the API is purely for demo purposes. We hope to eventually integrate go-resonate into an existing API.

### Bigchain
[Handcrafting transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) and sending them to BigchainDB/IPDB.

### Coala
A json-ld implementation of the Coala IP [spec](https://github.com/COALAIP/specs/tree/master/data-structure).

### Crypto
Ed25519 public-key cryptography and a minimal ILP crypto-conditions library.

### Types

### Util

## Walkthrough

### Registering a new user 
TODO: through partner orgaizations

### Logging in

### Creating a project

### Streaming a song

More to come!