package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"strconv"
)

type NewListRequest struct {
	Name string `json:"name"`
}

var bingelistCollection = db.GetCollection(db.DB, "bingelists")

func parseTitlesBySelectors(titles []*models.MediaItem, typeFilter string, genreFilter int) []*models.MediaItem {
	var result []*models.MediaItem = make([]*models.MediaItem, 0)
	evalType := false
	evalGenre := false
	if typeFilter != "" {
		evalType = true
	}
	if genreFilter != -1 {
		evalGenre = true
	}
	for _, title := range titles {
		if evalType && evalGenre {
			if title.Type == typeFilter && title.PrimaryGenreId == genreFilter {
				result = append(result, title)
			}
		} else if evalType {
			if title.Type == typeFilter {
				result = append(result, title)
			}
		} else if evalGenre {
			if title.PrimaryGenreId == genreFilter {
				result = append(result, title)
			}
		}
	}
	return result
}

func CreateNewBingeList(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	userId := r.Context().Value("userId").(string)
	rBody, err := io.ReadAll(r.Body)
	var listReq NewListRequest
	err = json.Unmarshal(rBody, &listReq)
	if err != nil {
		resp.Build(400, "binge list name error", nil)
		resp.Respond(w)
		return
	}
	listName := listReq.Name
	if err != nil || len(listName) < 1 {
		resp.Build(400, "binge list name error", nil)
		resp.Respond(w)
		return
	}
	blId := primitive.NewObjectID()
	var newList models.BingeList
	newList.Build(blId, listName, userId)
	_, err = bingelistCollection.InsertOne(context.TODO(), newList)
	if err != nil {
		resp.Build(500, "internal server error - list could not be created", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "list created successfully", nil)
	resp.Respond(w)
}
func DeleteBingeList(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	lIdStr := r.URL.Query().Get("id")
	listId, err := primitive.ObjectIDFromHex(lIdStr)
	if err != nil {
		resp.Build(400, "bad request - invalid list id", nil)
		resp.Respond(w)
		return
	}
	filter := bson.M{"_id": listId, "owner": owner}
	result := bson.M{}
	err = bingelistCollection.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		resp.Build(400, "bad request - list could not be deleted", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success - list deleted", nil)
	resp.Respond(w)
}
func GetMinifiedBingeLists(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	miniLists, err := services.BuildMinifiedBingeSlicesByOwner(owner)
	if err != nil {
		resp.Build(500, "internal server error - minified db contexts failed", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", miniLists)
	resp.Respond(w)
}
func GetBingeList(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	lStr := r.URL.Query().Get("id")
	typeFilter := r.URL.Query().Get("type")
	genreFilter := r.URL.Query().Get("genre")
	var gfToUse int
	lObj, err := primitive.ObjectIDFromHex(lStr)
	if err != nil {
		resp.Build(400, "there was a problem with the list id", nil)
		resp.Respond(w)
		return
	}
	var list models.BingeList
	filter := bson.M{"_id": lObj, "$or": bson.A{bson.M{"owner": owner}, bson.M{"users": owner}}}
	err = bingelistCollection.FindOne(context.TODO(), filter).Decode(&list)
	if err != nil {
		resp.Build(400, "bad request - problem with user or list", nil)
		resp.Respond(w)
		return
	}
	var titles models.BingeTitles
	if typeFilter == "" && genreFilter == "" {
		titles.Build(list.Name, list.Titles)
		resp.Build(200, "success", titles)
		resp.Respond(w)
		return
	}
	if genreFilter == "" {
		gfToUse = -1
	} else {
		if s, sErr := strconv.ParseInt(genreFilter, 10, 0); sErr == nil {
			gfToUse = int(s)
		} else {
			gfToUse = -1
			fmt.Println("couldn't parse genre id")
		}
	}
	filteredTitles := parseTitlesBySelectors(list.Titles, typeFilter, gfToUse)
	titles.Build(list.Name, filteredTitles)
	resp.Build(200, "success", titles)
	resp.Respond(w)
}
func GetBingeLists(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	var lists []*models.BingeList
	filter := bson.M{"$or": bson.A{bson.M{"owner": owner}, bson.M{"users": owner}}}
	cursor, err := bingelistCollection.Find(context.TODO(), filter)
	if cursor.All(context.TODO(), &lists); err != nil {
		resp.Build(500, "internal server error - db list fetch failed", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", lists)
	resp.Respond(w)
}
func AddToBingeList(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	lIdStr := r.URL.Query().Get("id")
	listId, err := primitive.ObjectIDFromHex(lIdStr)
	if err != nil {
		resp.Build(400, "bad request - invalid list id", nil)
		resp.Respond(w)
		return
	}
	var mediaItem models.MediaItem
	err = json.NewDecoder(r.Body).Decode(&mediaItem)
	if err != nil {
		resp.Build(400, "bad request - invalid media item", nil)
		resp.Respond(w)
		return
	}
	alreadyExists := services.AlreadyExists("bingelists", owner, mediaItem.MediaId, mediaItem.Type, listId)
	if alreadyExists {
		resp.Build(400, "bad request - item already in list", nil)
		resp.Respond(w)
		return
	}
	var genreName string
	genreName, err = GetGenreNameFromId(mediaItem.Type, mediaItem.PrimaryGenreId)
	if err != nil {
		resp.Build(500, "internal server error - problem getting genre from TMDB", nil)
		resp.Respond(w)
		return
	}
	mediaItem.PrimaryGenreName = genreName
	filter := bson.M{"_id": listId, "$or": bson.A{bson.M{"owner": owner}, bson.M{"users": owner}}}
	update := bson.M{"$push": bson.M{"titles": mediaItem}, "$inc": bson.M{"mediaCount": 1}}
	var result *mongo.UpdateResult
	result, err = bingelistCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		resp.Build(400, "bad request - error with list or media", nil)
		resp.Respond(w)
		return
	}
	if result.ModifiedCount == 0 {
		resp.Build(400, "bad request - invalid action", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success - item added to list", nil)
	resp.Respond(w)
}
func RemoveFromBingeList(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	owner := r.Context().Value("userId").(string)
	lIdStr := r.URL.Query().Get("id")
	mediaId := r.URL.Query().Get("mediaId")
	mediaType := r.URL.Query().Get("type")
	listId, err := primitive.ObjectIDFromHex(lIdStr)
	if err != nil {
		resp.Build(400, "bad request - invalid list id", nil)
		resp.Respond(w)
		return
	}
	alreadyExists := services.AlreadyExists("bingelists", owner, mediaId, mediaType, listId)
	if !alreadyExists {
		resp.Build(400, "bad request - item not in list", nil)
		resp.Respond(w)
		return
	}
	filter := bson.M{"_id": listId, "$or": bson.A{bson.M{"owner": owner}, bson.M{"users": owner}}}
	update := bson.M{"$pull": bson.M{"titles": bson.M{"mediaId": mediaId, "type": mediaType}}, "$inc": bson.M{"mediaCount": -1}}
	var result *mongo.UpdateResult
	result, err = bingelistCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		resp.Build(400, "bad request - invalid list or item", nil)
		resp.Respond(w)
		return
	}
	if result.MatchedCount == 0 {
		resp.Build(400, "bad request - invalid action", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success - removed item from list", nil)
	resp.Respond(w)
}
