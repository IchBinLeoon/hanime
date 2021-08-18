package utils

type Video struct {
	HentaiVideo    HentaiVideo    `json:"hentai_video"`
	VideosManifest VideosManifest `json:"videos_manifest"`
	StreamIndex    int
	OutputPath     string
}

type HentaiVideo struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	Views          int64  `json:"views"`
	Interests      int64  `json:"interests"`
	Brand          string `json:"brand"`
	Likes          int64  `json:"likes"`
	Dislikes       int64  `json:"dislikes"`
	Downloads      int64  `json:"downloads"`
	MonthlyRank    int64  `json:"monthly_rank"`
	CreatedAtUnix  int64  `json:"created_at_unix"`
	ReleasedAtUnix int64  `json:"released_at_unix"`
}

type VideosManifest struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	ID     int64  `json:"id"`
	Height string `json:"height"`
	Size   int64  `json:"filesize_mbs"`
}
