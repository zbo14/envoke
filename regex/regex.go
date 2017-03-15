package regex

const (
	DATE      = `^[12][09][0-9]{2}-[01][0-9]-[0-3][0-9]$`
	EMAIL     = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+.[a-zA-Z0-9-.]+$)`
	HFA       = `^[A-Z0-9]{6}$`
	ID        = `^[A-Fa-f0-9]{64}$` // hex
	IPI       = `^[0-9]{9}$`
	ISNI      = `^[0-9X]{16}$` //..
	ISRC      = `^[A-Z]{2}-[A-Z0-9]{3}-[7890][0-9]-[0-9]{5}$`
	ISWC      = `^T-[0-9]{3}.[0-9]{3}.[0-9]{3}-[0-9]$`
	LANGUAGE  = `^[A-Z]{2}$`
	PRO       = `^ASCAP|BMI|SESAC$`
	PUBKEY    = `^[1-9A-HJ-NP-Za-km-z]{43,44}$` // base58
	SIGNATURE = `^[1-9A-HJ-NP-Za-km-z]{87,88}$` // base58
	TERRITORY = `^[A-Z]{2}$`

	// FINGERPRINT_STD = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	// FINGERPRINT_URL = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
)
