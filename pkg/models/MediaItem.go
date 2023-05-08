package models

type MediaItem struct {
	MediaId        string `bson:"mediaId" json:"mediaId"`
	Title          string `bson:"title" json:"title"`
	Type           string `bson:"type" json:"type"`
	PosterPath     string `bson:"posterPath" json:"posterPath"`
	PrimaryGenreId int    `bson:"primaryGenreId" json:"primaryGenreId"`
}

func (m *MediaItem) Build(mi, t, ty, pp string, gid int) {
	m.MediaId = mi
	m.Title = t
	m.Type = ty
	m.PosterPath = pp
	m.PrimaryGenreId = gid
}
