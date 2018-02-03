package wxmp

// Use to reply music message
type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

// Use to reply news message
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}
