## Overview

Note: go-resonate is in the early (early) stages of development and is not usable in any way for anything yet.

### Identity
The demo currently uses Ed25519 public-key cryptography with BigchainDB for user identification and verification. We are looking at integration of other identity modules such as [uPort](https://github.com/ConsenSys/uport-lib) and [Blockstack](https://github.com/blockstack).

TODO: Test compatibility of BigchainDB crypto-conditions with Ed25519 library 

### Blockchain
This is not a blockchain application (yet). We are considering various platforms and consensus engines such as [Tendermint](https://github.com/tendermint), which would be used for the ordering of transactions (e.g. uploads, plays), state replication across a network, and payments. This infrastructure would be in addition to BigchainDB/IPDB.

### Users
We are currently focusing on two types of users in the demo, `artist` and `listener`, though we may include others (e.g. organization, label). That is, different types of users will have different designated actions and permissions. 

## Directories

### API
The client-side application interface. It's currently designed to run locally and communicate with BigchainDB/IPDB over http, though this may change. Working on functionality for the following user actions:
- Login/Register
- Create a project
- Stream a song
- other TBD

### Bigchain
[Handcrafting transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) and sending them to BigchainDB/IPDB.

### Coala
A json-ld implementation of the Coala IP [spec](https://github.com/COALAIP/specs/tree/master/data-structure).

### Types
Contains the `user` type (below), http/socket streaming services, and others.

```go
type User struct {
	Email    string `json:"email"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	// other TBD
}
```

### Util
TODO

## Walkthrough

### Registering a new user 
What you do:
-  Enter your email, username, password, and account type into the interface

What happens next:
- A `public key` and `private key` are generated for your account 
- Your personal information is signed with your `private key`, producing a `user signature`
- Your `user signature` is sent to BigchainDB/IPDB and stored

What you get:
- Your `public key`, `private key`, and a `user id`, which identifies your `user signature` in BigchainDB/IPDB

Note: none of your personal information is communicated or stored in raw form!

### Logging in
What you do:
- Enter your `private key`, `user id`, email, username, password, and account type into the interface. Woah, that's a lot, maybe we can have a file containing this information that a user uploads to login?

What happens next:
- The application asks BigchainDB/IPDB for the `user signature` corresponding to `user id` 
- It verifies the `user signature` with your `public key` (derived from your `private key`) and personal information
- Your `public key`, `private key`, `user id`, and personal information are kept in memory for the remainder of the session

What you get:
- A friendly welcome message

### Creating a project
TODO

### Streaming a song
TODO

More to come!