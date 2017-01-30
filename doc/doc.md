## Overview

Note: Envoke is in the early (early) stages of development.

### Identity 
We are looking at identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack), though user registration/login will eventually happen through partner organizations. The demo currently uses Ed25519 public-key cryptography for user verification and BigchainDB transactions.

### Blockchain
This is not a blockchain application (yet). We are considering various platforms and consensus engines such as [Tendermint](https://github.com/tendermint), which would be used for the ordering of transactions (e.g. uploads, plays), state replication across a network, and payments. This infrastructure would be in addition to BigchainDB/IPDB.

### Works, Metadata, and Rights 
This information will be persisted to BigchainDB/IPDB via user transactions. When an artist uploads a track to the file system, a `create` transaction containing the work's title and metadata will be sent to the BigchainDB/IPDB network. Works, metadata, and rights assignments/transfers are defined as Coala IP data models, with additional specifications if needed (e.g. adding music-specific fields to the schema).

## File storage 
We are looking at several distributed and decentralized file systems, including S3, minio (compatible with S3), IPFS, and storj. At the moment, using minio locally for file uploads and streams.

## Directories

### API
A client-side API that communicates with BigchainDB/IPDB and the file system over http. Working on the following functionality:
- Login/Register
- Persist a work to BigchainDB
- Upload a track to the file system
- Stream a track from the file system

Note: the API is purely for demo purposes. We eventually want to integrate permissioned access and artist payments into the application, but the focus for now is the aforementioned functionality.

### Bigchain
[Handcrafting transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) and sending them to BigchainDB/IPDB.

### Coala
A json-ld implementation of the Coala IP [spec](https://github.com/COALAIP/specs/tree/master/data-structure).

### Crypto
Ed25519 public-key cryptography and a minimal ILP [crypto-conditions](https://tools.ietf.org/html/draft-thomas-crypto-conditions-00) library.

### Types

### Util

## Walkthrough
TODO

### Registering a new user

### Logging in

### Creating a project

### Streaming a song

More to come!