package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_DEFAULT_URL             = "mongodb://localhost:27017"
	DATABASE_NAME                 = "h1server"
	KILLS_COLLECTION_NAME         = "kills"
	CONNECTIONS_COLLECTION_NAME   = "connections"
	CONSTRUCTIONS_COLLECTION_NAME = "construction"
	CHARACTERS_COLLECTION_NAME    = "characters"
	CROPS_COLLECTION_NAME         = "crops"
	SERVERS_COLLECTION_NAME       = "servers"
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

func getTopKiller(db *mongo.Database, mongoCtx context.Context, serverId uint32, entityType string) []Top {
	coll := db.Collection(KILLS_COLLECTION_NAME)
	pipeline := getTopKillerPipeline(serverId, entityType, 10)
	cursor, error := coll.Aggregate(mongoCtx, pipeline)

	if error != nil {
		panic(error)
	}
	defer cursor.Close(mongoCtx)
	var result []Top
	cursor.All(mongoCtx, &result)
	return result
}

func getTopPlayTime(db *mongo.Database, mongoCtx context.Context, serverId uint32) []Top {
	coll := db.Collection(CHARACTERS_COLLECTION_NAME)
	pipeline := getTopPlaytimePipeline(serverId, 10)
	cursor, error := coll.Aggregate(mongoCtx, pipeline)

	if error != nil {
		panic(error)
	}
	defer cursor.Close(mongoCtx)
	var result []Top
	cursor.All(mongoCtx, &result)
	return result
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
	if len(result) < 1 {
		return 0
	}
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

func getCharacters(db *mongo.Database, mongoCtx context.Context, serverId uint32) []Character {
	charactersCollection := db.Collection(CHARACTERS_COLLECTION_NAME)
	cursor, error := charactersCollection.Find(mongoCtx, bson.M{"serverId": serverId})
	if error != nil {
		panic(error)
	}
	var results []Character
	cursor.All(mongoCtx, &results)
	return results
}

func getConnectionsToServer(db *mongo.Database, mongoCtx context.Context, serverId uint32) []CountPerTime {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := getConnectionsPerServerPipeline(serverId)
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}

func getConnectionsLastMonthToServer(db *mongo.Database, mongoCtx context.Context, serverId uint32) []CountPerTime {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := getConnectionsLastMonthPerServerPipeline(serverId)
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}

func getKillsToServer(db *mongo.Database, mongoCtx context.Context, serverId uint32) []CountPerTime {
	ConnectionsCollection := db.Collection(KILLS_COLLECTION_NAME)
	pipeline := getKillsPerServerPipeline(serverId)
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}

func getAllKills(db *mongo.Database, mongoCtx context.Context) []CountPerTime {
	ConnectionsCollection := db.Collection(KILLS_COLLECTION_NAME)
	// TODO: change this to only get players kills
	pipeline := getAllKillsPipeline()
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}

func getAllConnections(db *mongo.Database, mongoCtx context.Context) []CountPerTime {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := getAllConnectionsPipeline()
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}

func getAllConnectionsLastMonth(db *mongo.Database, mongoCtx context.Context) []CountPerTime {
	ConnectionsCollection := db.Collection(CONNECTIONS_COLLECTION_NAME)
	pipeline := getAllConnectionsLastMonthPipeline()
	cursor, error := ConnectionsCollection.Aggregate(mongoCtx, pipeline)
	if error != nil {
		panic(error)
	}
	var results []CountPerTime
	cursor.All(mongoCtx, &results)
	return results
}
