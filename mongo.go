package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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

func getMongoCtx() (context.Context, context.CancelFunc) {
	mongoCtx, cancel := context.WithCancel(context.Background())
	return mongoCtx, cancel
}

func getDb(mongoCtx context.Context) *mongo.Database {
	mongo_url := os.Getenv("MONGO_URL")
	if mongo_url == "" {
		println("Use MONGO_DEFAULT_URL")
		mongo_url = MONGO_DEFAULT_URL
	}
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongo_url))
	if err != nil {
		panic("Failed to connect to mongo")
	}
	return mongoClient.Database(DATABASE_NAME)
}

func getCountPerServer(db *mongo.Database, mongoCtx context.Context, serverId uint32, collectionName string) uint32 {
	coll := db.Collection(collectionName)
	pipeline := getCountPerServerPipeline(serverId)
	cursor, error := coll.Aggregate(mongoCtx, pipeline)

	if error != nil {
		panic(error)
	}
	defer cursor.Close(mongoCtx)
	var result []CountPerServerResult
	cursor.All(mongoCtx, &result)
	return result[0].Count
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

func getConnectionsToServer(db *mongo.Database, mongoCtx context.Context, serverId uint32) []ConnectionsPerMonth {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := getConnectionsPerServerPipeline(serverId)
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
	pipeline := getAllConnectionsPipeline()
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []ConnectionsPerMonth
	cursor.All(mongoCtx, &results)
	return results
}
