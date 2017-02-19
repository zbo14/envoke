## SPEC

### Transaction 

A BigchainDB transaction wraps every data model sent to or queried in the database. A transaction includes a 64-character hexadecimal ID, the public key of the sender encoded as a base58 string, other cryptographic information and specification details. To learn more about the transaction model, please refer to the BigchainDB [doc](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html) or check out the `bigchain` module in the repo.  

```javascript
// Tx skeleton
{
	"id": "148ad0dbac5296c541e999be82a342064acf91b8269edd1d5f4cb0a6adb67b88",
    "asset": {
    	"data": {
        	"<envoke data model>"
        }
    },
    "outputs": [
    	{
        	"details": {
            	"public_key": "EMecDWWH1dLfyz1eEJL35iWibfY5JBnDjR3QybVSMgGL",
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
    "id": "148ad0dbac5296c541e999be82a342064acf91b8269edd1d5f4cb0a6adb67b88",
    {
      "email": "composer@email.com",
      "instance": {
          "time": "1487380379",
          "type": "agent"
       },
       "name": "composer",
       "socialMedia": "www.composer.com" 
     },
    "publicKey": "EMecDWWH1dLfyz1eEJL35iWibfY5JBnDjR3QybVSMgGL"
}
```

### Right 
The right model indicates ownership of a composition or sound recording. A right is defined by the data model that wraps it; a composition right is embedded in a composition and a recording right is embedded in a recording. Since the right is an embedded model, it does not reference the composition or recording by id. However, it does reference the right holder by id. Valid from/to specifies the timeframe in which the right is valid. Percentage shares indicate the amount of ownership or royalties the right holder is entitled to.

Example:
```javascript
{
    ...
    {
    	"instance": {
        	"time": "1487538908",
            "type": "right"
        },
    	"percentageShares": "20",
    	"rightHolderId": "148ad0dbac5296c541e999be82a342064acf91b8269edd1d5f4cb0a6adb67b88",
    	"validFrom": "2020-01-01",
    	"validTo": "3030-01-01"
    },
    ...
}
```

### Composition

A composition covers the written music and lyrics of a musical work. The model references the composer and publisher by id, contains the title of and the rights to the work. In the example, we have a composition written by our composer, which the composer and a publisher own the rights to.

Example: 
```javascript
{
	"id": "af6db95ed3a9f0827fe1512073a65d564b0789c7a64a64b9f9f963009be0a075",
    {
        "composerId": "148ad0dbac5296c541e999be82a342064acf91b8269edd1d5f4cb0a6adb67b88",
        "instance": {
            "time": "1487538908",
            "type": "composition"
        },
        "publisherId": "8c2155f7b541de9213f6308648fa6cdb036aad7bb21c69a55e27cc968c100bc4",
        "rights": [
            {
                "instance": {
                    "time": "1487538908",
                    "type": "right"
                },
                "percentageShares": "20",
                "rightHolderId": "148ad0dbac5296c541e999be82a342064acf91b8269edd1d5f4cb0a6adb67b88",
                "validFrom": "2020-01-01",
                "validTo": "3030-01-01"
            }, {
                "instance": {
                    "time": "1487538908",
                    "type": "right"
                },
                "percentageShares": "80",
                "rightHolderId": "8c2155f7b541de9213f6308648fa6cdb036aad7bb21c69a55e27cc968c100bc4",
                "validFrom": "2020-01-01",
                "validTo": "2030-01-01"
                }
            ],
        "title": "untitled"
    },
    "publicKey": "EMecDWWH1dLfyz1eEJL35iWibfY5JBnDjR3QybVSMgGL"
}
```

### Recording

A recording represents a digitally recorded performance of a composition. It references the composition, performer, producer, and record label by id. The sigining agent must reference a publishing license as well if it does not hold a right to the composition. In the example, we have a performer's cover of the previous composition. The rights are divided between the performer, producer, and record label accordindly. Note: the record label is the signing agent because it has a publishing license for the composition (see next section). Of course, this may not always be the case. If the composer performed and produced a master of an original composition, the publishing license id field would not be necessary.

Examples:
```javascript
{
    "id": "25603cc253fb3243864ea03a3a20dcbf06608233c5755849a89c90dcff7f5b67",
    {
       "compositionId": "af6db95ed3a9f0827fe1512073a65d564b0789c7a64a64b9f9f963009be0a075",
       "instance":{
          "time": "1487538908",
          "type": "recording"
       },
       "labelId": "df35075601cf0231aa9b373c6c7efc70354a26c45fa7e6cee3c3b52fa7518a25",
       "performerId": "96d3ec0d8dbf8966878b71679f01048977bdfbaac476dcdc057de15e7e1f5028",
       "producerId": "fdd987c53c1a71aa974c1758bf3a1c3fa4919edf10e88361a0de4f309e2e520a",
       "publishingLicenseId": "7bab873a54806186b3130263d59d8756c380044321cac858339b09cbe767d3a2",
       "rights":[
          {
             "instance": {
                "time": "1487538908",
                "type": "right"
             },
             "percentageShares": "50",
             "rightHolderId": "df35075601cf0231aa9b373c6c7efc70354a26c45fa7e6cee3c3b52fa7518a25",
             "validFrom": "2020-01-01",
             "validTo": "2030-01-01"
          },
          {
             "instance":{
                "time": "1487538908",
                "type": "right"
             },
             "percentageShares": "30",
             "rightHolderId": "96d3ec0d8dbf8966878b71679f01048977bdfbaac476dcdc057de15e7e1f5028",
             "validFrom": "2020-01-01",
             "validTo": "3030-01-01"
          },
          {
             "instance":{
                "time": "1487538908",
                "type": "right"
             },
             "percentageShares": "20",
             "rightHolderId": "fdd987c53c1a71aa974c1758bf3a1c3fa4919edf10e88361a0de4f309e2e520a",
             "validFrom": "2020-01-01",
             "validTo": "3030-01-01"
          }
       ]
    },
    "publicKey": "63L5uWVi8aMFvEhzvcngHDfF13J5oannMhD3bbJTVcf"
}
```

### License

A license is issued by a right-holder to a licensee. The licensee is then permitted to use the underlying composition or recording according to the terms of the license. The license references the licensee, licenser, and the composition or recording by id. Note: the license does not reference a right; the rights are embedded in the composition or recording. A license type (e.g. master, mechanical, synchronization) is included and we specify the timeframe in which the license is valid. The two examples here are (1) a publishing license issued by the publisher to the record label and (2) a recording license issued by the record label to a radio station.

Examples:

```javascript
// Publishing license
{
	"id": "7bab873a54806186b3130263d59d8756c380044321cac858339b09cbe767d3a2",
    {
       "compositionId": "af6db95ed3a9f0827fe1512073a65d564b0789c7a64a64b9f9f963009be0a075",
       "instance":{
          "time": "1487538908",
          "type": "publishing_license"
       },
       "licenseType": "mechanical",
       "licenseeId": "df35075601cf0231aa9b373c6c7efc70354a26c45fa7e6cee3c3b52fa7518a25",
       "licenserId": "8c2155f7b541de9213f6308648fa6cdb036aad7bb21c69a55e27cc968c100bc4",
       "validFrom": "2020-01-01",
       "validTo": "2025-01-01"
    },
    "publicKey": "HSC4y72BF615jx38JahJjpXqm6d9DGaEuWpei8mKheqk"
}

// Recording license
{
    "id": "1e41c036632e65dc210423036d349346faa8001079ac432851a3d1a6e65b94f5",
    {
       "instance":{
          "time": "1487538908",
          "type": "recording_license"
       },
       "licenseType": "master",
       "licenseeId": "42d7ad6eae12422426ce9d3504477220a63209553c703011d392481f3bd380c9",
       "licenserId": "df35075601cf0231aa9b373c6c7efc70354a26c45fa7e6cee3c3b52fa7518a25",
       "recordingId": "25603cc253fb3243864ea03a3a20dcbf06608233c5755849a89c90dcff7f5b67",
       "validFrom": "2020-01-01",
       "validTo": "2022-01-01"
    },
    "publicKey": "63L5uWVi8aMFvEhzvcngHDfF13J5oannMhD3bbJTVcf"
}
```

### Validation

The spec implements numeric validation (e.g. making sure total percentage shares do not exceed 100), text validation with regular-expressions (e.g. checking that ids are 64-character hexidecimal), and time validation (e.g. checking that rights and licenses have valid timeframes). The `linked_data` module implements a more intensive validation process, traversing the graph of data models and checking whether referenced entities exist in the database and satisfy the definitions in the spec. 