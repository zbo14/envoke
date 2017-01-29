## Overview

Note: go-resonate is in the early (early) stages of development.

### Identity
The demo currently uses Ed25519 public-key cryptography for user identification and verification. We are looking at identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack).

### Blockchain
This is not a blockchain application (yet). We are considering various platforms and consensus engines such as [Tendermint](https://github.com/tendermint), which would be used for the ordering of transactions (e.g. uploads, plays), state replication across a network, and payments. This infrastructure would be in addition to BigchainDB/IPDB.

### Works, Metadata, and Rights 
This information will be persistend to BigchainDB/IPDB via user txs. For example, when an artist uploads a track to the file system, a `create` tx containing the work's title and metadata will be sent to the BigchainDB/IPDB network. Works, metadata, and rights assignments/transfers are defined as Coala IP data models, with additional specifications if needed (e.g. adding music-specific fields to the schema).

## File storage 
We are looking at several distributed and decentralized file systems, including S3, minio, IPFS, and storj. Currently implementing minio for uploads and queries.

## Directories

### API
A client-side API. It's currently designed to run locally and communicate with BigchainDB/IPDB over http, though this may change. Working on functionality for the following actions:
- Login/Register (eventually via parter org)
- Persist/Query a work
- Upload/Stream a track
- other TBD

Note: the API is purely for demo purposes.

### Bigchain
[Handcrafting transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) and sending them to BigchainDB/IPDB.

### Coala
A json-ld implementation of the Coala IP [spec](https://github.com/COALAIP/specs/tree/master/data-structure).

### Crypto
Ed25519 public-key cryptography and a minimal ILP crypto-conditions library.

### Types

### Util

## Walkthrough
TODO

### Registering a new user

### Logging in

### Creating a project

### Streaming a song

More to come!