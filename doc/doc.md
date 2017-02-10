## Overview

Note: Envoke is in the early (early) stages of development.

### Identity and Cryptography
We are looking at identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack), though user registration/login will initially happen through partner organizations (e.g. labels, publishers). The demo uses ED25519 public-key cryptography for user registration/login and fulfillment of BigchainDB transactions.

### Metadata and Rights 
Metadata and rights are persisted to BigchainDB/IPDB via user transactions containing asset payloads. Metadata assets are defined as Music Ontology data models while rights assignments/transfers are defined as Coala IP data models.

TODO: use own spec for core data models?
TODO: integrate core data models with Coala IP rights

## Directories

### API
The API communicates with BigchainDB/IPDB over http. Currently in flux, but we are working on the following functionality:
- Register/login (artist/partner)
- Persist track/album metadata to the db (artist)
- Persist a signature to the db (partner)
- Verify a signature (artist/partner)

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

### Registration

### Login

### Persisting track/album metadata

### Persisting a signature

### Verifying a signature