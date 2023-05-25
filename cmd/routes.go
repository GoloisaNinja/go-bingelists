package main

import (
	"github.com/gorilla/mux"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/handlers"
	"go-bingelists/pkg/middleware"
	"net/http"
)

func routes(config *config.Repository) http.Handler {
	r := mux.NewRouter()

	// TMDB Handlers

	getTrendingLanding := handlers.GetTrendingLanding(config)
	getTrending := handlers.GetTrending(config)
	getMedia := handlers.GetMediaWithAllAttributes(config)
	getMediaCategories := handlers.GetMediaCategories(config)
	getCategoryResult := handlers.GetCategoryResultsByTypeAndPage(config)
	searchMedia := handlers.SearchMedia(config)

	// TMDB Endpoints

	r.Handle("/api/v1/trending/landing", middleware.Authenticate(config, getTrendingLanding)).Methods("GET")
	r.Handle("/api/v1/trending", middleware.Authenticate(config, getTrending)).Methods("GET")
	r.Handle("/api/v1/media", middleware.Authenticate(config, getMedia)).Methods("GET")
	r.Handle("/api/v1/categories", middleware.Authenticate(config, getCategoryResult)).Methods("GET")
	r.Handle("/api/v1/categories/list", middleware.Authenticate(config, getMediaCategories)).Methods("GET")
	r.Handle("/api/v1/search", middleware.Authenticate(config, searchMedia)).Methods("GET")

	// User Handlers

	createUser := handlers.CreateNewUser(config)
	loginUser := handlers.LoginUser(config)
	logoutUser := handlers.Logout(config)
	getPublicUsers := handlers.GetPublicUsers(config)

	// User Endpoints

	r.Handle("/api/v1/user/register", middleware.Registration(config, createUser)).Methods("POST")
	r.Handle("/api/v1/user/login", loginUser).Methods("POST")
	r.Handle("/api/v1/user/logout", middleware.Authenticate(config, logoutUser)).Methods("POST")
	r.Handle("/api/v1/user/users", middleware.Authenticate(config, getPublicUsers)).Methods("GET")

	// BingeList Handlers
	createNewBingeList := handlers.CreateNewBingeList(config)
	deleteBingeList := handlers.DeleteBingeList(config)
	getMinifiedBingeLists := handlers.GetMinifiedBingeLists(config)
	getBingeList := handlers.GetBingeList(config)
	getBingeLists := handlers.GetBingeLists(config)
	addToBingeList := handlers.AddToBingeList(config)
	removeFromBingeList := handlers.RemoveFromBingeList(config)

	// BingeList Endpoints
	r.Handle("/api/v1/bingelist/create", middleware.Authenticate(config, createNewBingeList)).Methods("POST")
	r.Handle("/api/v1/bingelist/delete", middleware.Authenticate(config, deleteBingeList)).Methods("DELETE")
	r.Handle("/api/v1/bingelists/minified", middleware.Authenticate(config, getMinifiedBingeLists)).Methods("GET")
	r.Handle("/api/v1/bingelist", middleware.Authenticate(config, getBingeList)).Methods("GET")
	r.Handle("/api/v1/bingelists", middleware.Authenticate(config, getBingeLists)).Methods("GET")
	r.Handle("/api/v1/bingelist/add", middleware.Authenticate(config, addToBingeList)).Methods("POST")
	r.Handle("/api/v1/bingelist/remove", middleware.Authenticate(config, removeFromBingeList)).Methods("POST")

	// Favorite Handlers
	getFavorites := handlers.GetFavorites(config)
	getMinifiedFavorites := handlers.GetMinifiedFavorites(config)
	addToFavorites := handlers.AddToFavorites(config)
	removeFromFavorites := handlers.RemoveFromFavorites(config)

	// Favorite Endpoints
	r.Handle("/api/v1/favorites", middleware.Authenticate(config, getFavorites)).Methods("GET")
	r.Handle("/api/v1/favorites/minified", middleware.Authenticate(config, getMinifiedFavorites)).Methods("GET")
	r.Handle("/api/v1/favorites/add", middleware.Authenticate(config, addToFavorites)).Methods("POST")
	r.Handle("/api/v1/favorites/remove", middleware.Authenticate(config, removeFromFavorites)).Methods("POST")

	// Invite Handlers
	getPendingInvites := handlers.GetPendingInvites(config)
	createNewInvite := handlers.CreateNewInvite(config)
	processInvite := handlers.ProcessInvite(config)

	// Invite Endpoints
	r.Handle("/api/v1/invites", middleware.Authenticate(config, getPendingInvites)).Methods("GET")
	r.Handle("/api/v1/invites/create", middleware.Authenticate(config, createNewInvite)).Methods("POST")
	r.Handle("/api/v1/invites/process", middleware.Authenticate(config, processInvite)).Methods("POST")

	return r
}
