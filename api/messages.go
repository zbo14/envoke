package api

type AlbumInfo struct {
	AlbumId  string   `json:"album_id"`
	TrackIds []string `json:"track_ids"`
}

func NewAlbumInfo(albumId string, songIds []string) *AlbumInfo {
	return &AlbumInfo{
		AlbumId:  albumId,
		TrackIds: songIds,
	}
}

type PartnerInfo struct {
	PartnerId string `json:"partner_id"`
	PrivPEM   []byte `json:"private_key"`
	PubPEM    []byte `json:"public_key"`
}

func NewPartnerInfo(partnerId string, privPEM, pubPEM []byte) *PartnerInfo {
	return &PartnerInfo{
		PartnerId: partnerId,
		PrivPEM:   privPEM,
		PubPEM:    pubPEM,
	}
}

type Stream struct {
	AlbumTitle string `json:"album_title"`
	Artist     string `json:"artist"`
	TrackTitle string `json:"track_title"`
	URL        string `json:"url"`
}

func NewStream(albumTitle, artist, trackTitle, url string) *Stream {
	return &Stream{
		AlbumTitle: albumTitle,
		Artist:     artist,
		TrackTitle: trackTitle,
		URL:        url,
	}
}
