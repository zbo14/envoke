## SPEC

### Transaction 

A JSON-encoded transaction wraps every data model sent to/queried in the database. A transaction includes a 64-character hexadecimal ID, the public keys of the sender/recipient encoded as base58 strings, and other information. To learn more about the transaction model, please refer to the BigchainDB [doc](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) or check out the `bigchain` module in the repo. In the following examples, the transaction id and public keys are included. Note: output amount and recipient key are included in the right models. The remaining transaction information is omitted.

### Party

The party model represents a right-holder of musical content, e.g. composer, performer, producer, publisher, or record label. 

* **Required** 
  
  `email=[email_address]` | `name=[alphanumeric]` | `sameAs=[url]` 

Example:

```javascript
{
    "data": {
        "email": "composer@email.com"
        "name": "composer_name",
        "sameAs": "www.composer.com"
        "type": "party",
    },
    "id": "<composer_id>",
    "sender_key": "<composer_key>"
}
```


### Composition

A composition represents the written music and lyrics of a musical work. 

* **Required**

`composerId=[hexadecimal]` | `title=[alphanumeric]`

* **Optional**

`hfa=[HFA_code]` | `ipi=[IPI_number]` | `iswc=[ISWC_code]` | `pro=[PRO_name]`

Example: 
```javascript
{
    "data": {
        "composerId": "<composer_id>",
        "hfa": "B3107S",
        "ipi": "678175698",
        "iswc": "T-034.524.680-1",
        "pro": "ASCAP",
        "title": "untitled",
        "type": "composition"
    },
    "id": "<composition_id>",
    "sender_key": "<composer_key>"
}
```


### Composition Right 
A composition right indicates a recipient's ownership of a composition. 

* **Required**
    
    `compositionId,recipientId,senderId=[hexadecimal]` | `territory=[country_codes]` | `validFrom,validTo=[yyyy-mm-dd]`

* **Notes**

    The amount in the transaction output specifies the percentage shares.

Example:
```javascript
{
    "amount": 80,
    "data": {
        "compositionId": "<composition_id>",
        "recipientId": "<publisher_id>",
        "senderId": "<composer_id>"
        "territory": [
            "GB",
            "US"
        ],
        "type": "right",
        // usage
        "validFrom": "2020-01-01",
        "validTo": "3000-01-01"
    },
    "recipient_key": "<publisher_key>",
    "id": "<publisher_right_id>",
    "sender_key": "<composer_key>"
}
```

### Publication

The publication model describes the relationship of composition rights to a composition.  

* **Required**
    
    `compositionId,compositionRightIds,publisherId=[hexadecimal]`

Example:
```javascript
{
    "data": {
        "compositionId": "<composition_id>",
        "compositionRightIds": [
            "<composer_right_id>",
            "<publisher_right_id>"
        ],
        "publisherId": "<publisher_id>",
        "type": "publication"
    },
    "id": "<publication_id>",
    "sender_key": "<publisher_key>"
}
```

### Mechanical License

A mechanical license, sent by a composition right-holder to a recipient, permits use of a publication.

* **Required**

    `publicationId,recipientId,senderId=[hexadecimal]` | `territory=[country_codes]` | `validFrom,validTo=[yyyy-mm-dd]`
    
* **Optional**

    `compositionRightId,compositionRightTransferId=[hexadecimal]`

* **Notes**
    
    Mechanical license should link to composition right or composition right transfer.

Example:

```javascript
{
    "data": {
        "compositionRightId": "<publisher_right_id>",
        "publicationId": "<publication_id>",
        "recipientId": "<label_id>",
        "senderId": "<publisher_id>",
        "territory": [
            "US"
        ],
        "type": "license",
        // usage
        "validFrom": "2020-01-01",
        "validTo": "2025-01-01"
    },
    "id": "<mechanical_license_id>",
    "sender_key": "<publisher_key>"
}
```

### Composition Right Transfer 

A composition right transfer indicates a transfer of composition right shares from a sender to a recipient.

* **Required**

    `compositionRightId,publicationId,recipientId,senderId,txId=[hexadecimal]`

* **Optional**

    `compositionRightTransferId=[hexadecimal]`

* Notes 

    The output amount(s) in the linked TRANSFER tx should specify recipient/sender shares.

Example:
```javascript 
{
    "data": {
        "compositionRightId": "<composer_right_id>"
        "publicationId": "<publication_id>",
        "recipientId": "<publisher_id>",
        "senderId": "<composer_id>",
        "txId": "<TRANSFER_tx_id>",
        "type": "transfer"
    },
    "id": "<publisher_transfer_id>",
    "sender_key": "<composer_key>"
}
```

### Recording

A recording represents a digitally recorded performance of a composition.

* **Required**

    `performerId,producerId,publicationId=[hexadecimal]`
    
* **Optional**
    
    `compositionRightId=[hexadecimal]` | `isrc=[ISRC_code]` 

Example:
```javascript
{
    "data": {
        "isrc": "US-S1Z-99-00001",
        "performerId": "<performer_id>",
        "producerId": "<producer_id>",
        "publicationId": "<publication_id>",
        "type": "recording"
    },
    "id": "<recording_id>",
    "sender_key": "<performer_key>" // or producer_key
}
```

### Recording Right 
A recording right indicates a recipient's ownership of a recording.

* **Required**
    
    `recordingId=[hexadecimal]` | `territory=[country_codes]` | `validFrom,validTo=[yyyy-mm-dd]`

* **Notes**

    The amount in the transaction output specifies the percentage shares.

Example:
```javascript
{
    "amount": 70,
    "data": {
        "recipientId": "<label_id>",
        "recordingId": "<recording_id>",
        "senderId": "<performer_id>",
        "territory": [
            "GB",
            "US"
        ],
        "type": "right",
        // usage
        "validFrom": "2020-01-01",
        "validTo": "2080-01-01"
    },
    "recipient_key": "<label_key>",
    "id": "<label_right_id>",
    "sender_key": "<performer_key>"
}
```

### Release

A release describes the relationship of recording rights to a recording.

* **Required**
    
    `recordingId,recordingRightIds,recordLabelId=[hexadecimal]`
    
* **Optional**

    `mechanicalLicenseId=[hexadecimal]`

* **Notes**

    Release should link to mechanical license if underlying recording does not link to composition right.

Example:

```javascript
{
    "data": {
        "mechanicalLicenseId": "<mechanical_license_id>",
        "recordingId": "<recording_id>",
        "recordingRightIds": [
            "<label_right_id>",
            "<performer_right_id>",
            "<producer_right_id>"
        ],
        "recordLabelId": "<label_id>",
        "type": "release"
    },
    "id": "<release_id>",
    "sender_key": "<label_key>"
}
```


### Master License

A master license, sent by a recording right-holder to a recipient, permits use of a release.

* **Required**

    `recipientId,releaseId,senderId=[hexadecimal]` | `territory=[country_codes]` | `validFrom,validTo=[yyyy-mm-dd]`
    
* **Optional**

    `recordingRightId,recordingRightTransferId=[hexadecimal]`

* **Notes**
    
    Master license should link to recording right or recording right transfer.

Example:
```javascript
{
    "data": {
        "instance": {
            "time": "1488768079",
            "type": "master_license"
        },
        "recipientId": "<radio_station_id>",
        "recordingRightId": "<label_right_id>",
        "releaseId": "<release_id>",
        "senderId": "<label_id>",
        "territory": [
            "US"
        ],
        "type": "license"
        // usage
        "validFrom": "2020-01-01",
        "validTo": "2022-01-01"
    },
    "id": "<master_license_id>",
    "sender_key": "<label_key>"
}
```


### Recording Right Transfer

A recording right transfer indicates a transfer of recording right shares from a sender to a recipient.

* **Required**

    `recipientId,recordingRightId,releaseId,senderId,txId=[hexadecimal]`

* **Optional**

    `recordingRightTransferId=[hexadecimal]`

* Notes 

    The output amount(s) in the linked TRANSFER tx should specify recipient/sender shares.

Example:
```javascript
{
    "data": {
        "recipientId": "<label_id>",
        "recordingRightId": "<performer_right_id>",
        "releaseId": "<release_id>",
        "senderId": "<performer_id>",
        "txId": "<TRANSFER_tx_id>",
        "type": "transfer"
    },
    "id": "<label_transfer_id>",
    "sender_key": "<performer_key>"
}
``` 

### Validation

The spec implements field validation with regular-expressions, e.g. checking that ids are 64-character hexidecimal. The `linked_data` module implements a more intensive validation process, traversing the graph of data models and checking whether referenced models exist in the database and satisfy the definitions in the spec. 
