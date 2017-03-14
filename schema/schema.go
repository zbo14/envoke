package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"

	. "github.com/zbo14/envoke/common"
	regex "github.com/zbo14/envoke/regex"
)

const SCHEMA = "http://json-schema.org/draft-04/schema#"

func ValidateModel(model Data, _type string) error {
	var schemaLoader jsonschema.JSONLoader
	modelLoader := jsonschema.NewGoLoader(model)
	switch _type {
	case "party":
		schemaLoader = PartyLoader
	case "composition":
		schemaLoader = CompositionLoader
	case "composition_right_transfer":
		schemaLoader = CompositionRightTransferLoader
	case "master_license":
		schemaLoader = MasterLicenseLoader
	case "mechanical_license":
		schemaLoader = MechanicalLicenseLoader
	case "publication":
		schemaLoader = PublicationLoader
	case "recording":
		schemaLoader = RecordingLoader
	case "recording_right_transfer":
		schemaLoader = RecordingRightTransferLoader
	case "release":
		schemaLoader = ReleaseLoader
	case "right":
		schemaLoader = RightLoader
	default:
		return ErrorAppend(ErrInvalidType, _type)
	}
	result, err := jsonschema.Validate(schemaLoader, modelLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return Error("Validation failed")
	}
	return nil
}

var link = Sprintf(`{
	"title": "Link",
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"pattern": "%s"
		}
	},
	"required": ["id"]
}`, regex.ID)

var itemList = Sprintf(`{
	"title": "ItemList",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"itemListElement": {
			"type": "array",
			"items": {
				"properties": {
					"item": {
						"$ref": "#/definitions/link"
					},
					"position": {
						"type": "number"
					}
				},
				"required": ["item", "position"]
			},
			"minItems": 1,
			"uniqueItems": true
		},
		"numberOfItems": {
			"type": "number"
		}
	},
	"required": ["itemListElement", "numberOfItems"]
}`, link)

var PartyLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "Party",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"email": {
			"type": "string",
			"pattern": "%s"
		},
		"ipiNumber": {
			"type": "string",
			"pattern": "%s"
		},
		"isniNumber": {
			"type": "string",
			"pattern": "%s"
		},
		"member": {
			"type": "array",
			"items": {
				"$ref": "#/definitions/link"
			}
		},
		"name": {
			"type": "string"
		},
		"pro": {
			"type": "string",
			"pattern": "%s"
		},
		"sameAs": {
			"type": "string"
		}
	},
	"required": ["email", "name", "sameAs"]
}`, SCHEMA, link, regex.EMAIL, regex.IPI, regex.ISNI, regex.PRO))

var CompositionLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "MusicComposition",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"composer": {
			"$ref": "#/definitions/link"
		},
		"hfaCode": {
			"type": "string",
			"pattern": "%s"
		},
		"inLanguage": {
			"type": "string",
			"pattern": "%s"
		},
		"iswcCode": {
			"type": "string",
			"pattern": "%s"
		},
		"name": {
			"type": "string"
		},
		"sameAs": {
			"type": "string"
		}
	},
	"required": ["composer", "name", "sameAs"]
}`, SCHEMA, link, regex.HFA, regex.LANGUAGE, regex.ISWC))

var PublicationLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "MusicPublication",
	"type": "object",
	"definitions": {
		"itemList": %s,
		"link": %s
	},
	"properties": {
		"composition": {
			"$ref": "#/definitions/itemList"
		},
		"compositionRight": {
			"$ref": "#/definitions/itemList"
		},
		"name": {
			"type": "string"
		},
		"publisher": {
			"$ref": "#/definitions/link"
		}
	},
	"required": ["composition", "compositionRight", "name"]
}`, SCHEMA, itemList, link))

var RecordingLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema":  "%s",
	"title": "MusicRecording",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"byArtist": {
			"$ref": "#/definitions/link"
		},
		"compositionRight": {
			"$ref": "#/definitions/link"
		},
		"duration": {
			"type": "string"			
		},
		"isrc": {
			"type": "string",
			"pattern": "%s"
		},
		"publication": {
			"$ref": "#/definitions/link"
		},
		"recordingOf": {
			"$ref": "#/definitions/link"
		}
	},
	"required": ["byArtist", "duration", "recordingOf"]
}`, SCHEMA, link, regex.ISRC))

var ReleaseLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema":  "%s",
	"title": "MusicRelease",
	"type": "object",
	"definitions": {
		"itemList": %s,
		"link": %s
	},
	"properties": {
		"name": {
			"type": "string"
		},
		"recording": {
			"$ref": "#/definitions/itemList"
		},
		"recordingRight": {
			"$ref": "#/definitions/itemList"
		},
		"recordLabel": {
			"$ref": "#/definitions/link"
		}
	},
	"required": ["name", "recording", "recordingRight"]
}`, SCHEMA, itemList, link))

var RightLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "Right",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"recipient": {
			"$ref": "#/definitions/link"
		},
		"sender": {
			"$ref": "#/definitions/link"
		},
		"territory": {
			"type": "array",
			"items": {
				"type": "string",
				"pattern": "%s"
			}
		},
		"validFrom": {
			"type": "string",
			"pattern": "%s"
		},
		"validThrough": {
			"type": "string",
			"pattern": "%s"
		}
	},
	"required": ["recipient", "sender", "territory", "validFrom", "validThrough"]
}`, SCHEMA, link, regex.TERRITORY, regex.DATE, regex.DATE))

var CompositionRightTransferLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "CompositionRightTransfer",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"compositionRight": {
			"$ref": "#/definitions/link"
		},
		"publication": {
			"$ref": "#/definitions/link"
		},
		"recipient": {
			"$ref": "#/definitions/link"
		},
		"sender": {
			"$ref": "#/definitions/link"
		},
		"tx": {
			"$ref": "#/definitions/link"
		}
	},
	"required": ["compositionRight", "publication", "recipient", "sender", "tx"]
}`, SCHEMA, link))

var RecordingRightTransferLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "RecordingRightTransfer",
	"type": "object",
	"definitions": {
		"link": %s
	},
	"properties": {
		"recipient": {
			"$ref": "#/definitions/link"
		},
		"recordingRight": {
			"$ref": "#/definitions/link"
		},
		"release": {
			"$ref": "#/definitions/link"
		},
		"sender": {
			"$ref": "#/definitions/link"
		},
		"tx": {
			"$ref": "#/definitions/link"
		}
	},
	"required": ["recipient", "recordingRight", "release", "sender", "tx"]
}`, SCHEMA, link))

var MechanicalLicenseLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "MechanicalLicense",
	"type": "object",
	"definitions": {
		"itemList": %s,
		"link": %s
	},
	"properties": {
		"composition": {
			"$ref": "#/definitions/itemList"
		},
		"compositionRight": {
			"$ref": "#/definitions/link"
		},
		"compositionRightTransfer": {
			"$ref": "#/definitions/link"
		},
		"publication": {
			"$ref": "#/definitions/link"
		},
		"recipient": {
			"$ref": "#/definitions/link"
		},
		"sender": {
			"$ref": "#/definitions/link"
		},
		"territory": {
			"type": "array",
			"items": {
				"type": "string",
				"pattern": "%s"
			}
		},
		"usage": {
			"oneOf": [
				{
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				{
					"type": "null"
				}
			]
		},
		"validFrom": {
			"type": "string",
			"pattern": "%s"
		},
		"validThrough": {
			"type": "string",
			"pattern": "%s"
		}
	},
	"anyOf": [
		{
			"required": ["composition"]
		},
		{
			"required": ["publication"]
		}
	],
	"dependencies": {
		"publication": {
			"oneOf": [
				{
					"required": ["compositionRight"]
				},
				{
					"required": ["compositionRightTransfer"]
				}
			]
		}
	},
	"required": ["recipient", "sender", "territory", "usage", "validFrom", "validThrough"]
}`, SCHEMA, itemList, link, regex.TERRITORY, regex.DATE, regex.DATE))

var MasterLicenseLoader = jsonschema.NewStringLoader(Sprintf(`{
	"$schema": "%s",
	"title": "MasterLicense",
	"type": "object",
	"definitions": {
		"itemList": %s,
		"link": %s
	},
	"properties": {
		"recipient": {
			"$ref": "#/definitions/link"
		},
		"recording": {
			"$ref": "#/definitions/itemList"
		},
		"recordingRight": {
			"$ref": "#/definitions/link"
		},
		"recordingRightTransfer": {
			"$ref": "#/definitions/link"
		},
		"release": {
			"$ref": "#/definitions/link"
		},
		"sender": {
			"$ref": "#/definitions/link"
		},
		"territory": {
			"type": "array",
			"items": {
				"type": "string",
				"pattern": "%s"
			}
		},
		"usage": {
			"oneOf": [
				{
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				{
					"type": "null"
				}
			]
		},
		"validFrom": {
			"type": "string",
			"pattern": "%s"
		},
		"validThrough": {
			"type": "string",
			"pattern": "%s"
		}
	},
	"anyOf": [
		{
			"required": ["recording"]
		},
		{
			"required": ["release"]
		}
	]
	"dependencies": {
		"release": {
			"oneOf": [
				{
					"required": ["recordingRight"]
				},
				{
					"required": ["recordingRightTransfer"]
				}
			]
		}
	},
	"required": ["recipient", "sender", "territory", "usage", "validFrom", "validThrough"]
}`, SCHEMA, itemList, link, regex.TERRITORY, regex.DATE, regex.DATE))
