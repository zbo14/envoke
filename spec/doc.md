## spec


### Instance

Every data model in the envoke spec wraps an instance, which indicates the time of creation with a unix timestamp and the type of model. 

Example:
```javascript
{
    "time": 1487034702, 
    "type": "artist"
}
```

### Agent

The agent model represents an envoke user, e.g. an artist, publisher, record label, or other organization. The model contains the user's email, name, and Ed25519 public key encoded as a base58 string. The public key is generated during the user registration process.

Example:

```javascript
{ 
    "id": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
    {
        "email": "zbo@gmail.com"
        "instance": {
            "time": 1487034702,
            "type": "artist"
        },
        "name": "zbo",
        "public_key": "7YJokhq2kDcdpoFz38F54nf1GtPGNuX3jC4d3tXUHa1y"
    },
    ...
}
```

### Album

The album model contains metadata pertaining to a musical project and the ID of the artist, publisher, and/or record label involved in the undertaking of the project. The IDs are hexadecimal strings corresponding to transactions, e.g. we can query the transaction containing the artist with the artist id. Lastly, we have the title of the project.

Example:
```javascript
{
    "id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
    {
        "artist_id": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
        "instance": {
            "time": 1487034739,
            "type": "album"
        },
        "label_id": "145b9d9c51998ee08420a8031f4f9f4b8ce9a37598099034a178ddda6bccb9cc",
        "publisher_id": "94b28c606e229e83728ea9b4fd07895e0ef17742bd1ee21b83ff5d5bf19432d4",
        "title": "Splash"
    },
    ...
}
```

### Track

The track model contains metadata pertaining to an audio file and the ID of the album, artist, publisher, and/or record label involved in the creation of digital content. The IDs are hexidecimal strings as before. Note: the example track does not contain an artist, label, or publisher id; it points to the album in the previous example, which points to an artist, label, and publisher. The fingerprint is a url-safe base64 encoded string of an acoustic fingerprint computed on the audio file. We also have the title of track and the track number, since the track is part of an album.

Example: 
```javascript
{
    "id": "67a5851dedc1bc14e947092ceaac4c506e5ada048d34e4052efeca47877a707b",
    {
        "album_id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
        "fingerprint": "13_1d9d9dXfXfXV3131199d_dffX_1HX1__R1dX_0VXV-fF11fn1d9Xp9XfVafV31Wn1d9Vr93fVe9dX1293Vddt93X3_fV99_3xfff9cV33__Fd13_xXdd_9X_Vf_V_1X9V_9V9Vf_VfdXf1X3139V99f_V_fX_1f3V_9X99f_V_fV_1X_1X9V_9V_Xf3VX139VV9d_VTfX_3V31_9xd9ftdXfX7XV31f1Vd91_V_c",
        "instance": {
            "time": 1487034741,
            "type": "track"
        },
        "title": "Luv Tub",
        "track_number": 1
    },
    ...
}
```

### Signature

The signature model represents an agent's signature of another data model. One case where a signature is used is the issuance/acceptance of a right to a track or album. The model ID points to the model that was signed and the signer ID points to the agent who generated the signature. The value is a base58 encoded string of an Ed25519 signature, which can be verified with the signer's public key. In the following example, our artist generates a signature for the album created earlier.

Example:
```javascript
{
    "instance": {
        "time": 1487134839,
        "type": "signature"
    },
    "model_id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
    "signer_id": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
    "value": "3h9SpyQ855cCk7RMqd9FycVaxZ4z5ronxhrKpDcWhDyHTM7jBvKkLMGS419rn1TkKFt8bHXAcL31FcF4QjHVSZJ5"
}
```

### Right

The right model specifies terms of usage regarding a track or album, distribution of content, and collection of royalty payments within a context and timeframe. Entities involved in the creation of digital content have the ability to issue rights to others (or among themselves). We include the issuer type, in this case artist. Note: the or label, and the ID of the recipient. The music ID points to a track or album while the percentage shares indicate the portion of royalty payments for the recipient. The signature field contains a valid signature of the track or album by the issuer. Note: I copy-and-pasted the previous example, changing the timestamp. I could have assigned an ID to the previous example and had the right point to the signature. Usually, a signature is generated with the right, so we do not have two separate transactions. In the following example, our artist issues an album right to the publisher.

Example:
```javascript 
{
    "id": "bcd5e4b38f67a7529adc085725f2ccd9b3a116e5af7e54db865aa91847612721",
    {
        "context": "commercial_use",
        // TODO: "copyright"
        "instance": {
            "time": 14924142831,
            "type": "right"
        },
        "issuer_type": "artist"
        // TODO: "license"
        "percentage_shares": 70,
        "recipient_id": "94b28c606e229e83728ea9b4fd07895e0ef17742bd1ee21b83ff5d5bf19432d4",
        "signature": {
            "instance": {
                "time": 14924142831,
                "type": "signature"
            },
            "model_id": "14bd4afc291413a63b4ab8c0754e875a9f6abb7560f1e828d2fdd280e0505578",
            "signer_id": "772d206f7f5c541f8e9153262601b5229fc81b512667f1d7c35f5101ba347188",
            "value": "3h9SpyQ855cCk7RMqd9FycVaxZ4z5ronxhrKpDcWhDyHTM7jBvKkLMGS419rn1TkKFt8bHXAcL31FcF4QjHVSZJ5"
        },
        "usage": ["copy", "play"],
        "valid_from": "2018-01-01",
        "valid_to": "2020-01-01",
    },
    ...
}
```

### Right Signature

A right signature is a recipient's signature of a right. This model represents the recipient's acknowledgement of the right and agreement to its terms. In the example, the publisher signs the right issued by the artist.

Example:
```javascript
{
    "id": "a1a5d1ebe0242fefc521d4e15a8e5abcc28bb06b1386318798da1e30fdfd8365",
    {
        "instance": {
            "time": 1492424868,
            "type": "signature"
        },
        "model_id": "bcd5e4b38f67a7529adc085725f2ccd9b3a116e5af7e54db865aa91847612721",
        "signer_id": "94b28c606e229e83728ea9b4fd07895e0ef17742bd1ee21b83ff5d5bf19432d4",
        "value": "24Jciw3jzP87WivHLo8FdcBhP5qPJhWtfvK1WJv3ePBR2RADx6hMhEfQEdNhwi8dLHU9ATjz33aUJwB5FQ7W4EiV"
    },
    ...
}
```

### Validation

The spec implements numeric validation, e.g. making sure percentage shares do not exceed 100, and text validation with regular-expressions, e.g. checking that IDs are 32-byte hexidecimal, signatures are 64-byte base58, emails have appropriate formats, etc. The linked-data module implements a more intensive validation process, traversing the graph of data models and checking whether referenced entities exist in the database and satisfy the definitions in the spec.