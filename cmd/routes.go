package main

import (
	"github.com/gorilla/mux"
	"go-bingelists/pkg/handlers"
	"go-bingelists/pkg/middleware"
	"net/http"
)

func routes() http.Handler {
	r := mux.NewRouter()

	// TMDB Handlers

	getTrendingLanding := http.HandlerFunc(handlers.GetTrendingLanding)
	getTrending := http.HandlerFunc(handlers.GetTrending)
	getMedia := http.HandlerFunc(handlers.GetMediaWithAllAttributes)
	getMediaCategories := http.HandlerFunc(handlers.GetMediaCategories)
	getCategoryResult := http.HandlerFunc(handlers.GetCategoryResultsByTypeAndPage)
	searchMedia := http.HandlerFunc(handlers.SearchMedia)

	// TMDB Endpoints

	r.Handle("/api/v1/trending/landing", middleware.Authenticate(getTrendingLanding)).Methods("GET")
	r.Handle("/api/v1/trending", middleware.Authenticate(getTrending)).Methods("GET")
	r.Handle("/api/v1/media", middleware.Authenticate(getMedia)).Methods("GET")
	r.Handle("/api/v1/categories", middleware.Authenticate(getCategoryResult)).Methods("GET")
	r.Handle("/api/v1/categories/list", middleware.Authenticate(getMediaCategories)).Methods("GET")
	r.Handle("/api/v1/search", middleware.Authenticate(searchMedia)).Methods("GET")

	// User Handlers

	createUser := http.HandlerFunc(handlers.CreateNewUser)
	loginUser := http.HandlerFunc(handlers.LoginUser)
	logoutUser := http.HandlerFunc(handlers.Logout)

	// User Endpoints

	r.Handle("/api/v1/user/register", middleware.Registration(createUser)).Methods("POST")
	r.Handle("/api/v1/user/login", loginUser).Methods("POST")
	r.Handle("/api/v1/user/logout", middleware.Authenticate(logoutUser)).Methods("POST")

	// BingeList Handlers
	createNewBingeList := http.HandlerFunc(handlers.CreateNewBingeList)
	deleteBingeList := http.HandlerFunc(handlers.DeleteBingeList)
	getMinifiedBingeLists := http.HandlerFunc(handlers.GetMinifiedBingeLists)
	getBingeList := http.HandlerFunc(handlers.GetBingeList)
	getBingeLists := http.HandlerFunc(handlers.GetBingeLists)
	addToBingeList := http.HandlerFunc(handlers.AddToBingeList)
	removeFromBingeList := http.HandlerFunc(handlers.RemoveFromBingeList)

	// BingeList Endpoints
	r.Handle("/api/v1/bingelist/create", middleware.Authenticate(createNewBingeList)).Methods("POST")
	r.Handle("/api/v1/bingelist/delete", middleware.Authenticate(deleteBingeList)).Methods("DELETE")
	r.Handle("/api/v1/bingelists/minified", middleware.Authenticate(getMinifiedBingeLists)).Methods("GET")
	r.Handle("/api/v1/bingelist", middleware.Authenticate(getBingeList)).Methods("GET")
	r.Handle("/api/v1/bingelists", middleware.Authenticate(getBingeLists)).Methods("GET")
	r.Handle("/api/v1/bingelist/add", middleware.Authenticate(addToBingeList)).Methods("POST")
	r.Handle("/api/v1/bingelist/remove", middleware.Authenticate(removeFromBingeList)).Methods("POST")

	// Favorite Handlers
	getFavorites := http.HandlerFunc(handlers.GetFavorites)
	getMinifiedFavorites := http.HandlerFunc(handlers.GetMinifiedFavorites)
	addToFavorites := http.HandlerFunc(handlers.AddToFavorites)
	removeFromFavorites := http.HandlerFunc(handlers.RemoveFromFavorites)

	// Favorite Endpoints
	r.Handle("/api/v1/favorites", middleware.Authenticate(getFavorites)).Methods("GET")
	r.Handle("/api/v1/favorites/minified", middleware.Authenticate(getMinifiedFavorites)).Methods("GET")
	r.Handle("/api/v1/favorites/add", middleware.Authenticate(addToFavorites)).Methods("POST")
	r.Handle("/api/v1/favorites/remove", middleware.Authenticate(removeFromFavorites)).Methods("POST")

	return r
}
