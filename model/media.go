package model

const (
	MediaFile = iota
	MediaImage
	MediaAudio
	MediaVideo
)

// 图片，音频，和其它附件类型
type Media struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"-" meddler:"created_at"`

	PostId    int64  `json:"post_id"   meddler:"post_id"`
	AuthorId  int64  `json:"author_id" meddler:"author_id"`

	// media type: file/image/audio/video
	Type      int64  `json:"type"    meddler:"_type"`

	// media path (json array): ["http://...jpg", ..]
	Path      string `json:"path"    meddler:"path"`

	// meta-data of media
	Meta      string `json:"meta"    meddler:"meta"`
}

