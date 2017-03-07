## SPEC

### Transaction 

A JSON-encoded transaction wraps every data model sent to/queried in the database. A transaction includes a 64-character hexadecimal ID, the public keys of the signer/holder encoded as base58 strings, and other information. To learn more about the transaction model, please refer to the BigchainDB [doc](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) or check out the `bigchain` module in the repo. In the following examples, the transaction id and public keys are included. Note: output amount and holder key are only included in the right models. The remaining transaction information is omitted.

### Instance

Every data model in the envoke spec wraps an instance, which contains the unix timestamp when the model was created and the type of model. 

Example:
```javascript
{
	...
    "instance": {
        "time": "1487034702", 
        "type": "agent"
    }
    ...
}
```

### Agent

The agent model represents an envoke user, e.g. a composer, performer, producer, publisher, or record label. 

* **Required** 
  
  `email=[email_address]` | `name=[alphanumeric]` | `socialMedia=[url]` 

Example:

```javascript
{
    "data": {
        "email": "composer@email.com"
        "instance": {
          "type": "agent",
          "time": "1488486588"
        },
        "name": "composer_name",
        "socialMedia": "www.composer.com"
    },
    "id": "<composer_id>",
    "signer_key": "<composer_key>"
}
```


### Composition

A composition represents the written music and lyrics of a musical work. 

* **Required**

`composerId,publisherId=[hexadecimal]` | `title=[alphanumeric]`

* **Optional**

`hfa=[HFA_code]` | `ipi=[IPI_number]` | `iswc=[ISWC_code]` | `pro=[PRO_name]`

Example: 
```javascript
{
    "data": {
        "composerId": "<composer_id>",
        "hfa": "B3107S",
        "instance": {
            "time": "1488768059",
            "type": "composition"
        },
        "ipi": "678175698",
        "iswc": "T-034.524.680-1",
        "pro": "ASCAP",
        "publisherId": "<publisher_id>",
        "title": "untitled"
    },
    "id": "<composition_id>",
    "signer_key": "<composer_key>"
}
```


### Composition Right 
A composition right indicates ownership of a composition. The amount in the transaction output specifies the percentage shares.

* **Required**
	
	`compositionId=[hexadecimal]` | `validFrom,validTo=[yyyy-mm-dd]`
    

* **Optional**

	`territory=[country_codes]`

Example:
```javascript
{
    "amount": 80,
    "data": {
        "compositionId": "<composition_id>",
        "instance": {
            "time": "1488768060",
            "type": "right"
        },
        "territory": [
            "GB",
            "US"
        ],
        "validFrom": "2020-01-01",
        "validTo": "3000-01-01"
    },
    "holder_key": "<publisher_key>",
    "id": "<publisher_right_id>",
    "signer_key": "<composer_key>"
}
```

### Composition Right Assignment
A composition right assignment links a signer and holder agent to a composition right.

* **Required**

	`holderId,rightId,signerId=[hexadecimal]`

Example:
```javascript
{
    "data": {
        "holderId": "<publisher_id>",
        "instance": {
            "time": "1488768060",
            "type": "assignment"
        },
        "rightId": "<publisher_right_id>",
        "signerId": "<composer_id>"
    },
    "id": "<publisher_right_assignment_id>",
    "signer_key": "<composer_key>"
}
```


### Publication

The publication model describes the relationship of composition right assignments to a composition.  

* **Required**
	
    `assignmentIds,compositionId=[hexadecimal]`

Example:
```javascript
{
    "data": {
        "assignmentIds": [
            "<composer_right_assignment_id>",
            "<publisher_right_assignment_id>"
        ],
        "compositionId": "<composition_id>",
        "instance": {
            "time": "1488768061",
            "type": "publication"
        }
    },
    "id": "<publication_id>",
    "signer_key": "<publisher_key>"
}
```

### Mechanical License

A mechanical license, issued by a composition right-holder to a licensee, permits use of a publication.

* **Required**

	`licenseeId,licenserId,publicationId=[hexadecimal]` | `validFrom,validTo=[yyyy-mm-dd]`
    
* **Optional**

	`assignmentId=[hexadecimal]` | `territory=[country_codes]` | `transferId=[hexadecimal]`

* **Notes**
    
    Mechanical license must contain assignmentId or transferId

Example:

```javascript
{
    "data": {
        "assignmentId": "<publisher_right_assignment_id>",
        "instance": {
            "time": "1488768063",
            "type": "mechanical_license"
        },
        "licenseeId": "<label_id>",
        "licenserId": "<publisher_id>",
        "publicationId": "<publication_id>",
        "territory": [
            "US"
        ],
        "validFrom": "2020-01-01",
        "validTo": "2025-01-01"
    },
    "id": "<mechanical_license_id>",
    "signer_key": "<publisher_key>"
}
```

### Composition Right Transfer 

A composition right transfer indicates a transfer of ownership shares from a sender to a recipient.

* **Required**

	`publicationId,recipientId,senderId,txId=[hexadecimal]`

Example:
```javascript 
{
    "data": {
        "instance": {
            "time": "1488768090",
            "type": "transfer"
        },
        "publicationId": "<publication_id>",
        "recipientId": "<publisher_id>",
        "senderId": "<composer_id>",
        "txId": "<TRANSFER_tx_id>"
    },
    "id": "<publisher_transfer_id>",
    "signer_key": "<composer_key>"
}
```

### Recording

A recording represents a digitally recorded performance of a composition.

* **Required**

	`labelId,performerId,producerId,publicationId=[hexadecimal`
    
* **Optional**
	
    `assignmentId=[hexadecimal]` | 	`isrc=[ISRC_code]` 

Example:
```javascript
{
    "data": {
        "instance": {
            "time": "1488768066",
            "type": "recording"
        },
        "isrc": "US-S1Z-99-00001",
        "labelId": "<label_id>",
        "performerId": "<performer_id>",
        "producerId": "<producer_id>",
        "publicationId": "<publication_id>"
    },
    "id": "<recording_id>",
    "signer_key": "<performer_key>" // or producer_key
}
```

### Recording Right 
A recording right indicates ownership of a recording.

* **Required**
	
	`recordingId=[hexadecimal]` | `validFrom,validTo=[yyyy-mm-dd]`
    
* **Optional**

	`territory=[country_codes]`

Example:
```javascript
{
    "amount": 70,
    "data": {
        "instance": {
            "time": "1488768070",
            "type": "right"
        },
        "recordingId": "<recording_id>",
        "territory": [
            "GB",
            "US"
        ],
        "validFrom": "2020-01-01",
        "validTo": "2080-01-01"
    },
    "holder_key": "<label_key>",
    "id": "<label_right_id>",
    "signer_key": "<performer_key>"
}
```

### Recording Right Assignment

A recording right assignment links a signer and holder agent to a recording right.

* **Required**
	
    `holderId,rightId,senderId=[hexadecimal]`

Example:
```javascript
{
    "data": {
        "holderId": "<label_id>",
        "instance": {
            "time": "1488768070",
            "type": "assignment"
        },
        "rightId": "<label_right_id>",
        "signerId": "<performer_id>"
    },
    "id": "<label_right_assignment_id>",
    "signer_key": "<performer_key>"
}
```

### Release

A release describes the relationship of recording right assignments to a recording.

* **Required**
	
    `assignmentIds,recordingId=[hexadecimal]`
   
    
* **Optional**

	`licenseId=[hexadecimal]`

Example:

```javascript
{
    "data": {
        "assignmentIds": [
            "<label_right_assignment_id>",
            "<performer_right_assignment_id>",
            "<producer_right_assignment_id>"
         ],
        "instance": {
            "time": "1488768071",
            "type": "release"
         },
        "licenseId": "<mechanical_license_id>",
        "recordingId": "<recording_id>"
    },
    "id": "<release_id>",
    "signer_key": "<label_key>"
}
```


### Master License

A master license, issued by a recording right-holder to a licensee, permits use of a release.

* **Required**

	`licenseeId,licenserId,releaseId=[hexadecimal]` | `validFrom,validTo=[yyyy-mm-dd]`
    
* **Optional**

	`assignmentId=[hexadecimal]` | `territory=[country_codes]` | `transferId=[hexadecimal]`

* **Notes**
    
    Master license must contain assignmentId or transferId

Example:
```javascript
{
    "data": {
        "assignmentId": "<label_right_assignment_id>",
        "instance": {
            "time": "1488768079",
            "type": "master_license"
        },
        "licenseeId": "<radio_station_id>",
        "licenserId": "<label_id>",
        "releaseId": "<release_id>",
        "territory": [
            "US"
        ],
        "validFrom": "2020-01-01",
        "validTo": "2022-01-01"
    },
    "id": "<master_license_id>",
    "signer_key": "<label_key>"
}
```


### Recording Right Transfer

A recording right transfer indicates a transfer of ownership shares from a sender to a recipient.

* **Required**

	`recipientId,releaseId,senderId,txId=[hexadecimal]`

Example:
```javascript
{
    "data": {
        "instance": {
            "time": "1488768102",
            "type": "transfer"
        },
        "recipientId": "<label_id>",
        "releaseId": "<release_id>",
        "senderId": "<performer_id>",
        "txId": "<TRANSFER_tx_id>"
    },
    "id": "<label_transfer_id>",
    "signer_key": "<performer_key>"
}
``` 

### Validation

The spec implements field validation with regular-expressions, e.g. checking that ids are 64-character hexidecimal. The `linked_data` module implements a more intensive validation process, traversing the graph of data models and checking whether referenced models exist in the database and satisfy the definitions in the spec. 