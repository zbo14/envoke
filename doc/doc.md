## Overview

Note: Envoke is in the early (early) stages of development.

### Identity 
We are looking at identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack), though user registration/login will initially happen through partner organizations (e.g. labels, publishers). The demo uses ED25519 public-key cryptography for partner verification and BigchainDB transactions.

TODO: implement RSA crypto for linked-data signatures.

### Blockchain
This is not a blockchain application (yet). We are considering platforms and consensus engines such as [Tendermint](https://github.com/tendermint), which would be used for the ordering of transactions, state replication across a network, and payments. This infrastructure would be in addition to BigchainDB/IPDB.

### Metadata and Rights 
Metadata and rights will be persisted to BigchainDB/IPDB via user transactions. When an album/track is uploaded to the file system, a transaction containing the title and metadata will be sent to the BigchainDB/IPDB network. Metadata and rights assignments/transfers are defined as Coala IP data models, with additional specifications if needed (e.g. adding music-specific fields to the schema).

TODO: assess whether Music Ontology is better suited for core data models (not including rights).

### File storage 
We are looking at several distributed and decentralized file systems, including S3, minio (compatible with S3), IPFS, and storj. At the moment, we are using minio locally for file uploads and queries.

## Directories

### API
The API communicates with BigchainDB/IPDB and the file system over http. Currently in flux, but we are working on the following functionality:
- Register a partner
- Artist/partner login
- Persist a work to the db
- Upload an album to the file system
- Stream a track from the file system

Note: the API is purely for demo purposes.

### Bigchain
[Handcrafting transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) and sending them to BigchainDB/IPDB.

### Common 

### Crypto
ED25519 and RSA public-key cryptography and a minimal ILP [crypto-conditions](https://tools.ietf.org/html/draft-thomas-crypto-conditions-00) library.

### Spec
[Coala IP](https://github.com/COALAIP/specs/tree/master/data-structure) and [Music Ontology](http://musicontology.com/specification/).

## Walkthrough
TODO:

### Registering a partner

### Logging in

### Uploading an album

### Streaming a track