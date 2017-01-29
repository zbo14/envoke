package types

const (
	// User types
	ARTIST       = "artist"
	LISTENER     = "listener"
	ORGANIZATION = "organization"
	LABEL        = "label"
)

// User
type User struct {
	Email    string `json:"email"`
	Region   string `json:"region"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Username string `json:"username"`
	// Other info?
}

func NewUser(email, region, password, _type, username string) *User {
	return &User{
		Email:    email,
		Region:   region,
		Password: password,
		Type:     _type,
		Username: username,
	}
}

func NewArtist(email, region, password, username string) *User {
	return &User{
		Email:    email,
		Region:   region,
		Password: password,
		Type:     ARTIST,
		Username: username,
	}
}

func NewListener(email, region, password, username string) *User {
	return &User{
		Email:    email,
		Region:   region,
		Password: password,
		Type:     LISTENER,
		Username: username,
	}
}
