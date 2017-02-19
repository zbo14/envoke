## spec

### Transaction 

A BigchainDB transaction wraps every data model sent to or queried in the database. A transaction includes a 64-character hexadecimal ID, the public key of the sender encoded as a base58 string, other cryptographic information and specification details. To learn more about the transaction model, please refer to the BigchainDB [doc](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) or check out the `bigchain` module in the repo.  

```javascript
// BigchainDB Tx skeleton
{
	"id": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
    "asset": {
    	"data": {
        	"<envoke data model>"
        }
    },
    "outputs": [
    	{
        	"details": {
            	"public_key": "5CYiuPjtftwXS59XBmJnakUt3CyQLVoWbg3mGPPgNWeg",
                ...
            },
    		...
    	}
    ],
    ...
}
```

### Instance

Every data model in the envoke spec wraps an instance, which includes the unix timestamp when the model was created and the type of model. 

Example:
```javascript
{
    "time": "1487034702", 
    "type": "agent"
}
```

### Agent

The agent model represents an envoke user, e.g. a composer, performer, producer, publisher, or record label. The model contains the agent's name, email, and URL to website/social media profile. The agent's public key is generated during a registration process and included in the create transaction. In the following examples, we have the data model, id and public key fields. The other transaction fields are omitted for brevity. 

Example:

```javascript
{ 
    "id": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
    {
      "email": "artist@gmail.com",
      "instance": {
          "time": "1487380379",
          "type": "agent"
       },
       "name": "artist",
       "socialMedia": "<soundcloud-page-url>" 
     },
    "publicKey": "676zY8CJ8tK9VcA1ZcoLdKtRBisGR9vviy4zkbLvTtSc"
}
```

### Right 
The right model indicates ownership of a composition or sound recording. A right is defined by the data model that wraps it; a composition right is embedded in a composition and a recording right is embedded in a recording. Since the right is an embedded model, it does not reference the composition or recording by id. However, it does reference the right holder by id. Valid from/to specifies the timeframe in which the right is valid. Percentage shares indicates the amount of ownership or royalties the right holder is entitled to.

Example:
```javascript
{
    "id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
    {
    	"instance": {
        	"time": "1487380380",
            "type": "right"
        },
    	"percentageShares": "30",
    	"rightHolderId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
    	"validFrom": "2018-01-01",
    	"validTo": "2020-01-01"
    },
    "publicKey": "676zY8CJ8tK9VcA1ZcoLdKtRBisGR9vviy4zkbLvTtSc"
}
```

### Composition

A composition covers the written music and lyrics of a musical work. The model references the composer and publisher by id, contains the title of and the rights to the work. In the example, we have a composition written by our artist, which the artist and a publisher own the rights to. 

Example: 
```javascript
{
	"id": "553ff534a508796050a6026768ad6c55a6277cf509985a07b0e63d84de6ae39f",
    {
        "composerId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
        "id": "553ff534a508796050a6026768ad6c55a6277cf509985a07b0e63d84de6ae39f",
        "instance": {
            "time": "1487380380",
            "type": "composition"
        },
        "publisherId": "80a59a5dd91bb7a3d4eb9138a6e3ab17fab66e274222df76199408df68c5d838",
        "rights": [
            {
                "instance": {
                    "time": "1487380380",
                    "type": "right"
                },
                "percentageShares": "30",
                "rightHolderId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
                "validFrom": "2018-01-01",
                "validTo": "2020-01-01"
            },
            {
                "instance": {
                    "time": "1487380380",
                    "type": "right"
                },
                "percentageShares": "70",
                "rightHolderId": "80a59a5dd91bb7a3d4eb9138a6e3ab17fab66e274222df76199408df68c5d838",
                "validFrom": "2018-01-01",
                "validTo": "2020-01-01"
            }
        ],
        "title": "untitled"
    },
    "publicKey": "676zY8CJ8tK9VcA1ZcoLdKtRBisGR9vviy4zkbLvTtSc"
}
```

### Recording

A recording represents a digitally recorded performance of a composition. It references the composition, performer, producer, and record label by id. The agent must reference a publishing license as well if it does not hold a right to the composition. The fingerprint is a url-safe base64 encoded string of an acoustic fingerprint computed on the audio file. Once again, we have rights embedded in the data model. In the first example, the artist who performs and produces the master retains 40% ownership of the recording while the record label has 60% ownership. As the composition rights holder, the artist does not need a publishing license. In the second example, another artist produces a cover of the composition with the same label, sam division of rights, and references a publishing license accordingly.

Examples:
```javascript
// No publishing license
{	
	"id": "e4b63b94a743602a4bb471811368e2d7be70261727c6a5e407a985087e30683e",
    {
    	"compositionId": "553ff534a508796050a6026768ad6c55a6277cf509985a07b0e63d84de6ae39f",
    	"fingerprint": "V0VHa09XR0xXb2VnbGt3ZW93ZWZ3ZUZ3ZWZ3Z3dlZ2VnZ2VyZ2U",
        "instance": {
            "time": "1487380380",
            "type": "recording"
        },
        "labelId": "f043be23912d9da24ab647dec83e5b74b729a771110585acb323fc02372650d5",
        "performerId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
        "producerId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
    "	rights": [
    		{
          		"instance": {
            		"time": "1487380380",
            		"type": "right"
          		},
          		"percentageShares": "60",
          		"rightHolderId": "f043be23912d9da24ab647dec83e5b74b729a771110585acb323fc02372650d5",
          		"validFrom": "2018-01-01",
          		"validTo": "2022-01-01"
      		},
      		{
       	 		"instance": {
          			"time": "1487380380",
          			"type": "right"
       	 		},
        		"percentageShares": "40",
        		"rightHolderId": "80b37eece2ef36b2fdabe3511f7409b0e69b97c48090688700dd6536dc5c835d",
        		"validFrom": "2018-01-01",
        		"validTo": "2023-01-01"
      		}
    	]
  	},
    "publicKey": "676zY8CJ8tK9VcA1ZcoLdKtRBisGR9vviy4zkbLvTtSc"
}

// With publishing license
{	
	"id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
    {
    	"compositionId": "553ff534a508796050a6026768ad6c55a6277cf509985a07b0e63d84de6ae39f",
    	"fingerprint": "V0VHa09XR0xXb2VnbGt3ZW93ZWZ3ZUZ3ZWZ3Z3dlZ2VnZ2VyZ2U",
        "instance": {
            "time": "1487380380",
            "type": "recording"
        },
        "labelId": "f043be23912d9da24ab647dec83e5b74b729a771110585acb323fc02372650d5",
        "performerId": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
        "producerId": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
        "publishingLicenseId": "ef6ca05ef80c5662ea9f34353b2f94260258ba4dff65b35324028121ca5c6f7e",
    "	rights": [
    		{
          		"instance": {
            		"time": "1487380380",
            		"type": "right"
          		},
          		"percentageShares": "60",
          		"rightHolderId": "f043be23912d9da24ab647dec83e5b74b729a771110585acb323fc02372650d5",
          		"validFrom": "2018-01-01",
          		"validTo": "2022-01-01"
      		},
      		{
       	 		"instance": {
          			"time": "1487380380",
          			"type": "right"
       	 		},
        		"percentageShares": "40",
        		"rightHolderId": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
        		"validFrom": "2018-01-01",
        		"validTo": "2023-01-01"
      		}
    	]
  	},
    "publicKey": "676zY8CJ8tK9VcA1ZcoLdKtRBisGR9vviy4zkbLvTtSc"
}
```

### License

A license is issued by a right-holder to a licensee. The licensee is then permitted to use the underlying composition or recording according to the terms of the license. The license references the licensee, licenser, and the composition or recording by id. Note: the license does not reference a right; the rights are embedded in the composition or recording. A license type (e.g. master, mechanical, synchronization) is included and we specify the timeframe in which the license is valid. The two examples here are (1) a publishing license issued by the publisher to a licensee and (2) a recording license issued by the record label to a licensee.

Examples:

```javascript
// Publishing license
{
	"id": "ef6ca05ef80c5662ea9f34353b2f94260258ba4dff65b35324028121ca5c6f7e",
    {
        "compositionId": "553ff534a508796050a6026768ad6c55a6277cf509985a07b0e63d84de6ae39f",
        "instance": {
            "time": "1487380387",
            "type": "publishing_license"
        },	
        "licenseType": "mechanical_license",
        "licenseeId": "6e34ce797516318663c12a893475a4327259726317aa21c1ea8c5a4bcfd9eb93",
        "licenserId": "80a59a5dd91bb7a3d4eb9138a6e3ab17fab66e274222df76199408df68c5d838",
        "validFrom": "2018-01-01",
        "validTo": "2025-01-01"
    },
    "publicKey": "7YJokhq2kDcdpoFz38F54nf1GtPGNuX3jC4d3tXUHa1y"
}

// Recording license
{
	"id": "218d3b5372786f4916ab5cb2afaaf7c3a78d0024d79d708cdf77987e40b3789d",
    {
    	"instance": {
        	"time": "1487380382",
        	"type": "recording_license"
      	},
      	"licenseType": "master_license",
      	"licenseeId": "d9fe6f37dbf26274dce15dc7d5e51c544e70bade17140964d28f711aa0771706",
      	"licenserId": "f043be23912d9da24ab647dec83e5b74b729a771110585acb323fc02372650d5",
      	"recordingId": "e4b63b94a743602a4bb471811368e2d7be70261727c6a5e407a985087e30683e",
      	"validFrom": "2018-01-01",
      	"validTo": "2019-01-01"
    },
    "publicKey":
}
```

### Validation

The spec implements numeric validation (e.g. making sure total percentage shares do not exceed 100), text validation with regular-expressions (e.g. checking that ids are 64-character hexidecimal), and time validation (e.g. checking that rights and licenses have valid timeframes). The `linked_data` module implements a more intensive validation process, traversing the graph of data models and checking whether referenced entities exist in the database and satisfy the definitions in the spec. 