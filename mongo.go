package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_DEFAULT_URL           = "mongodb://localhost:27017"
	DATABASE_NAME               = "h1server"
	KILLS_COLLECTION_NAME       = "kills"
	CONNECTIONS_COLLECTION_NAME = "connections"
	SERVERS_COLLECTION_NAME     = "servers"
)

type Kills struct {
	ServerId int32 `bson:"serverId" json:"serverId"`
}
type ConnectionsPerMonth struct {
	Id    string `bson:"_id" json:"_id"`
	Count int32  `bson:"count" json:"count"`
}

func getMongoCtx() (context.Context, context.CancelFunc) {
	mongoCtx, cancel := context.WithCancel(context.Background())
	return mongoCtx, cancel
}

func getDb(mongoCtx context.Context) *mongo.Database {
	mongo_url := os.Getenv("MONGO_URL")
	if mongo_url == "" {
		mongo_url = MONGO_DEFAULT_URL
	}
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongo_url))
	if err != nil {
		panic("Failed to connect to mongo")
	}
	return mongoClient.Database(DATABASE_NAME)
}

func getServers(db *mongo.Database, mongoCtx context.Context) []Server {
	serversCollection := db.Collection(SERVERS_COLLECTION_NAME)
	cursor, error := serversCollection.Find(mongoCtx, bson.M{})
	if error != nil {
		panic(error)
	}
	var results []Server
	cursor.All(mongoCtx, &results)
	return results
}

type Server struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Locked           bool               `bson:"locked" json:"locked"`
	AllowedAccess    bool               `bson:"allowedAccess" json:"allowedAccess"`
	IsOfficial       bool               `bson:IsOfficial json:"IsOfficial"`
	ServerId         uint32             `bson:"serverId" json:"serverId"`
	ServerAddress    string             `bson:"serverAddress" json:"serverAddress"`
	Name             string             `bson:"name" json:"name"`
	NameId           uint32             `bson:"nameId" json:"nameId"`
	MaxPop           uint32             `bson:"maxPopulationNumber" json:"maxPopulationNumber"`
	PopulationNumber uint32             `bson:"populationNumber" json:"populationNumber"`
	Region           string             `bson:"region" json:"region"`
	H1emuVersion     string             `bson:"h1emuVersion" json:"h1emuVersion"`
	GameVersion      uint32             `bson:"gameVersion" json:"gameVersion"`
}

func getKills(db *mongo.Database, mongoCtx context.Context, serverId uint32) uint32 {
	killsCollection := db.Collection(KILLS_COLLECTION_NAME)
	cursor, error := killsCollection.Find(mongoCtx, bson.M{"serverId": serverId})
	if error != nil {
		panic(error)
	}
	var results []Kills
	cursor.All(mongoCtx, &results)
	return uint32(len(results))
}

func getConnectionsToServer(db *mongo.Database, mongoCtx context.Context, serverId uint32) []ConnectionsPerMonth {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$addFields", bson.D{{"yearMonth", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$yearMonth"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []ConnectionsPerMonth
	cursor.All(mongoCtx, &results)
	return results
}

func getAllConnections(db *mongo.Database, mongoCtx context.Context) []ConnectionsPerMonth {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := mongo.Pipeline{
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$addFields", bson.D{{"yearMonth", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$yearMonth"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []ConnectionsPerMonth
	cursor.All(mongoCtx, &results)
	return results
}
