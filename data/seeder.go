package data

import (
	"context"
	"fmt"
	"go-bingelists/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

var nameArr = []string{
	"abby",
	"argentine",
	"argento",
	"arragato",
	"arthur",
	"barry",
	"batima",
	"baxter",
	"betsy",
	"betches",
	"balaclava",
	"barrington",
	"bethster",
	"bingeyguy",
	"bingeygirl",
	"bingbong",
	"crepesgirl",
	"crepesguy",
	"crepesthey",
	"crepeswe",
	"cynthia",
	"craig",
	"curtis",
	"charlie",
	"chuggington",
	"chupacabra",
	"generaldevious",
	"distractedjellies",
	"disenfranchisedwaffles",
	"darklydave",
	"davedarkly",
	"dipsanddots",
	"dotsanddips",
	"ditsanddops",
	"eggburton",
	"eggerton",
	"eggsalads",
	"easydoesit",
	"existentialbeans",
	"effortlessmice",
	"erstwhilesausages",
	"franklydumb",
	"ferociousears",
	"fringefries",
	"frinklesink",
	"franklesinks",
	"frunklesanks",
	"fizzlebottoms",
}

func SeedDB() {
	var docs []interface{}

	for _, name := range nameArr {
		newRecord := bson.M{
			"name":      name,
			"email":     name + "@gmail.com",
			"password":  "password",
			"isPrivate": false,
			"token":     "faketokenthatisnotreallyatoken",
			"invites":   bson.A{},
			"createdAt": primitive.NewDateTimeFromTime(time.Now()),
		}
		docs = append(docs, newRecord)
	}

	client := db.DB
	collection := client.Database("bingelist").Collection("users")
	result, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success! Added %d records!", len(result.InsertedIDs))
}
