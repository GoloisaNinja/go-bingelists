package models

type TMDBBaseStruct struct {
	Movie interface{}
	Tv    interface{}
}

func (t *TMDBBaseStruct) Build(i interface{}, mediaType string) {
	if mediaType == "movie" {
		t.Movie = i
	} else {
		t.Tv = i
	}
}

type MediaWithAttributes struct {
	Media     interface{}
	Credits   interface{}
	Providers interface{}
	Similars  interface{}
}

func (m *MediaWithAttributes) Build(media, credits, providers, similars interface{}) {
	m.Media = media
	m.Credits = credits
	m.Providers = providers
	m.Similars = similars
}
