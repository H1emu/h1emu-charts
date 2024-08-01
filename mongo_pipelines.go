package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCountPerServerPipeline(serverId uint32) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$count", "count"}},
	}

	return pipeline
}

func getTopKillerPipeline(serverId uint32, entityType string, maxResult uint8) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$match", bson.D{{"type", entityType}}}},
		{{"$group", bson.D{
			{"_id", "$characterName"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"count", -1}}}},
		{{"$limit", maxResult}},
		{{"$project", bson.D{
			{"_id", 0},
			{"characterName", "$_id"},
			{"count", 1},
		}}},
	}

	return pipeline
}

func getConnectionsPerServerPipeline(serverId uint32) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$addFields", bson.D{{"yearMonth", bson.D{
			{"$dateToString", bson.D{
				{"format", "%Y-%m"},
				{"date", "$creationDate"},
			}},
		}}}}},
		{{"$group", bson.D{
			{"_id", "$yearMonth"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getAllConnectionsPipeline() mongo.Pipeline {
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

	return pipeline
}
