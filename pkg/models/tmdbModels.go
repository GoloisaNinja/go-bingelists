package models

type TMDBBaseStruct struct {
	Movie interface{} `json:"movie"`
	Tv    interface{} `json:"tv"`
}

func (t *TMDBBaseStruct) Build(i interface{}, mediaType string) {
	if mediaType == "movie" {
		t.Movie = i
	} else {
		t.Tv = i
	}
}

type MediaWithAttributes struct {
	Media     interface{} `json:"media"`
	Credits   interface{} `json:"credits"`
	Providers interface{} `json:"providers"`
	Similars  interface{} `json:"similars"`
}

func (m *MediaWithAttributes) Build(media, credits, providers, similars interface{}) {
	m.Media = media
	m.Credits = credits
	m.Providers = providers
	m.Similars = similars
}
