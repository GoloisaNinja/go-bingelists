package handlers

import (
	"encoding/json"
	"fmt"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/util"
	"net/http"
	"strings"
)

const TMDB_BASE_URL = "https://api.themoviedb.org/3"

var APIKEY = util.GetDotEnv("TMDB_APIKEY")

func trendingByTypeAndPage(mediaType, page string) (*http.Response, error) {
	resp, err := http.Get(TMDB_BASE_URL + "/trending/" + mediaType + "/week?page=" + page + "&api_key=" + APIKEY)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetTrendingLanding(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	fmt.Println(userId)
	var resp responses.Response
	mtResp, err := trendingByTypeAndPage("movie", "1")
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	ttResp, err := trendingByTypeAndPage("tv", "1")
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	defer mtResp.Body.Close()
	defer ttResp.Body.Close()
	var t models.TMDBBaseStruct
	err = json.NewDecoder(mtResp.Body).Decode(&t.Movie)
	err = json.NewDecoder(ttResp.Body).Decode(&t.Tv)
	resp.Build(200, "success", t)
	resp.Respond(w)
}

func GetTrending(w http.ResponseWriter, r *http.Request) {
	mediaType := r.URL.Query().Get("type")
	page := r.URL.Query().Get("page")
	var resp responses.Response
	var target interface{}
	trendingResp, err := trendingByTypeAndPage(mediaType, page)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	err = json.NewDecoder(trendingResp.Body).Decode(&target)
	resp.Build(200, "success", target)
	resp.Respond(w)
}

func GetMediaWithAllAttributes(w http.ResponseWriter, r *http.Request) {
	mediaType := r.URL.Query().Get("media_type")
	mediaId := r.URL.Query().Get("media_id")
	base := "/" + mediaType + "/" + mediaId
	var resp responses.Response
	var mA models.MediaWithAttributes
	// base media response
	media, err := http.Get(TMDB_BASE_URL + base + "?append_to_response=videos&api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	credits, err := http.Get(TMDB_BASE_URL + base + "/credits?api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	providers, err := http.Get(TMDB_BASE_URL + base + "/watch/providers?language=en-US&api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	similars, err := http.Get(TMDB_BASE_URL + base + "/similar?page=1&api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	defer media.Body.Close()
	defer credits.Body.Close()
	defer providers.Body.Close()
	defer similars.Body.Close()
	err = json.NewDecoder(media.Body).Decode(&mA.Media)
	if err != nil {
		resp.Build(500, "json encoding error on media", nil)
		resp.Respond(w)
		return
	}
	err = json.NewDecoder(credits.Body).Decode(&mA.Credits)
	if err != nil {
		resp.Build(500, "json encoding error on credits", nil)
		resp.Respond(w)
		return
	}
	err = json.NewDecoder(providers.Body).Decode(&mA.Providers)
	if err != nil {
		resp.Build(500, "json encoding error on providers", nil)
		resp.Respond(w)
		return
	}
	err = json.NewDecoder(similars.Body).Decode(&mA.Similars)
	if err != nil {
		resp.Build(500, "json encoding error on similars", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", mA)
	resp.Respond(w)
}

func GetMediaCategories(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	movieGenres, err := http.Get(TMDB_BASE_URL + "/genre/movie/list?api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	tvGenres, err := http.Get(TMDB_BASE_URL + "/genre/tv/list?api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	var genres models.TMDBBaseStruct
	defer movieGenres.Body.Close()
	defer tvGenres.Body.Close()
	err = json.NewDecoder(movieGenres.Body).Decode(&genres.Movie)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	err = json.NewDecoder(tvGenres.Body).Decode(&genres.Tv)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", genres)
	resp.Respond(w)
}
func GetCategoryResultsByTypeAndPage(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	mediaType := queryParams.Get("media_type")
	genre := queryParams.Get("genre")
	page := queryParams.Get("page")
	sort := "popularity.desc"
	if queryParams.Has("sort_by") {
		qp := queryParams.Get("sort_by")
		qpArr := strings.Split(qp, ".")
		if qpArr[0] == "rating" {
			sort = "vote_average." + qpArr[1] + "&vote_count.gte=50"
		}
	}
	var resp responses.Response
	catResp, err := http.Get(TMDB_BASE_URL + "/discover/" +
		mediaType + "?language=en-US&include_adult=false&include_video=false&with_genres=" +
		genre + "&sort_by=" + sort + "&page=" + page + "&api_key=" + APIKEY)
	if err != nil {
		resp.Build(500, err.Error(), nil)
		resp.Respond(w)
		return
	}
	defer catResp.Body.Close()
	var c interface{}
	err = json.NewDecoder(catResp.Body).Decode(&c)
	if err != nil {
		resp.Build(500, "json decoding error", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", c)
	resp.Respond(w)
}
