## SPEC

### Transaction 

A JSON-encoded transaction wraps every data model sent to/queried in the database. A transaction includes a 64-character hexadecimal ID, the public keys of the signer and recipient encoded as base58 strings, and other information. To learn more about the transaction model, please refer to the BigchainDB [doc](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) or check out the `bigchain` module in the repo. In the following examples, the transaction id and public keys are included. Note: output amount and recipient key are only included in the right model. The remaining transaction information is omitted.

### Instance

Every data model in the envoke spec wraps an instance, which contains the unix timestamp when the model was created and the type of model. 

Example:
```javascript
{
    "time": "1487034702", 
    "type": "agent"
}
```

### Agent

The agent model represents an envoke user, e.g. a composer, performer, producer, publisher, or record label. The model contains the agent's name, email, and URL to website/social media profile. The agent's public key is generated during a registration process and included in the create transaction.

Example:

```javascript
{
    "id": "3316...",
    "signer_key": "9vd8...",
    {
        "name": "composer_name",
        "instance": {
          "type": "agent",
          "time": "1488486588"
        },
        "socialMedia": "www.composer.com",
        "email": "composer@email.com"
    }
}
```


### Composition

A composition represents the written music and lyrics of a musical work. The model contains an HFA song code, ISWC code, title, and links to the composer and publisher.

Example: 
```javascript
{
    "id": "ebdb...",
    "signer_key": "9vd8...",
    {
        "composerId": "3316...",
        "hfa": "B3107S",
        "instance": {
            "time": "1488487088",
            "type": "composition"
        },
        "iswc": "T-034.524.680-1",
        "publisherId": "67c3...",
        "title": "Luv Tub"
    }
}
```


### Composition Right 
A composition right indicates ownership of a composition. The amount in the transaction output specifies the percentage shares of ownership the recipient receives. Territory and timeframe are also specified.

Example:
```javascript
{
    "amount": 80,
    "id": "2238...",
    "recipient_key": "7tev...",
    "signer_key": "9vd8...",
    {
        "compositionId": "ebdb...",
        "instance": {
            "time": "1487538908",
            "type": "right"
        },
        "territory": [
            "GB",
            "US"
        ]
        "validFrom": "2020-01-01",
        "validTo": "3000-01-01"
    }
}
```

### Publication

The publication model describes the relationship of composition rights to a composition.  

Example:
```javascript
{
    "id": "fd7b...",
    "signer_key": "7tev...",
    {
        "compositionId": "ebdb..."
        "instance": {
            "type": "publication",
            "time": "1488487537"
        },
        "rightIds": [
            "4481...",
            "2338..."
        ]
    }
}
```

### Recording

A recording represents a digitally recorded performance of a composition. It contains an ISRC code and links to the performer, producer, and record label. If the performer or producer is a right-holder of the composition, the composition right id is included as well.

Example:
```javascript
{
    "id": "868c...",
    "signer_key": "3y8v...",
    {
        "isrc": "US-S1Z-99-00001",
        "instance": {
            "type": "recording",
            "time": "1488488596"
        },
        "labelId": "5ea2..."
        "performerId": "6c77...",
        "producerId": "2e3a...",
        "publicationId": "fd7b..."
    }
}
```

### Recording Right 
A recording right indicates ownership of a recording. The model is identical to the composition right model except that it links to a recording.

Example:
```javascript
{
    "amount": 70,
    "id": "1c52...",
    "recipient_key": "Eup9...",
    "signer_key": "3y8v...",
    {
        "instance": {
            "time": "1488488739",
            "type": "right"
        },
        "recordingId": "868c...",
        "territory": [
            "GB",
            "US"
        ]
        "validFrom": "2020-01-01",
        "validTo": "2080-01-01",
    }
}
```

### Release

A release describes the relationship of recording rights to a recording. The model links to a mechanical license if the recording does not link to a composition right. 

```javascript
{
    "id": "9502..."
    "signer_key": "Eup9..."
    {
        "instance": {
            "type": "release",
            "time": "1488488895"
        },
        "licenseId": "6572...",
        "recordingId": "868c...",
        "rightIds": [
            "4e30...", 
            "a5dd...",
            "1c52..."
        ],
    }
}
```

### Mechanical License

A mechanical license is issued by a composition right-holder to a licensee. The model links to the licensee, licenser, and the licensed publication. Territory and timeframe are also specified.

Example:

```javascript
{
    "id": "6572...",
    "signer_key": "7tev...",
    {
        "instance": {
            "type": "mechanical_license",
            "time": "1488488047"
        },
        "licenseeId": "5ea2...",
        "licenserId": "67c3..."
        "publicationId": "fd7b...",
        "rightId": "2338...",
        "territory": [
            "GB",
            "US"
        ],
        // usage?
        "validFrom": "2020-01-01",
        "validTo": "2024-01-01"
    }
}
```

### Master License

A master license is issued by a recording right-holder to a licensee. The model is identical to the mechanical license except that it links to a licensed release. 

Example:
```javascript
{
    "id": "0d19...",
    "signer_key": "Eup9...",
    {
        "instance": {
            "type": "master_license",
            "time": "1488489015"
        },
        "licenseeId": "1280...",
        "licenserId": "5ea2...",
        "rightId": "1c52..",
        "territory": [
            "GB",
            "US"
        ],
        "validFrom": "2020-01-01",
        "validTo": "2024-01-01"
    }
}
```

### Composition Right Transfer 

A composition right transfer indicates a transfer of ownership shares from a sender to a recipient. The sender holds a composition right or ownership shares from a previous transfer. The model links to a publication, recipient, sender, and [TRANSFER tx](https://docs.bigchaindb.com/en/latest/transaction-concepts.html#transfer-transactions).

Example:
```javascript 
{
    "id": "56ac...",
    "signer_key": "9vd8...",
    {
        "instance": {
            "time": "1488654659",
            "type": "transfer"
        },
        "publicationId": "fd7b...",
        "recipientId": "67c3...",
        "senderId": "3316...", 
        "txId": "1047..."
    }
}
```

### Recording Right Transfer

A recording right transfer indicates a transfer of ownership shares from a sender to a recipient. The sender holds a recording right or ownership shares from a previous transfer. The model links to a recipient, release, sender, and TRANSFER tx.

Example:
```javascript
{
    "id": "3ce9...",
    "signer_key": "3y8v...",
    {
        "instance": {
            "time": "1488654667",
            "type": "transfer"
        },
        "recipientId": "5ea2...",
        "releaseId": "9502...",
        "senderId": "6c77...",
        "txId": "30ec..."
    }
}
``` 

### Validation

The spec implements field validation with regular-expressions, e.g. checking that ids are 64-character hexidecimal. The `linked_data` module implements a more intensive validation process, traversing the graph of data models and checking whether referenced models exist in the database and satisfy the definitions in the spec. 