package main

import (
	"time"

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

func getConnectionsLastMonthPerServerPipeline(serverId uint32) mongo.Pipeline {
	now := time.Now()
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},

		{{"$match", bson.D{
			{"creationDate", bson.D{
				{"$gte", firstOfLastMonth},
				{"$lt", now},
			}},
		}}},

		{{"$addFields", bson.D{{"day", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m-%d"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$day"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
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

func getAllConnectionsLastMonthPipeline() mongo.Pipeline {
	now := time.Now()
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)

	pipeline := mongo.Pipeline{
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$match", bson.D{
			{"creationDate", bson.D{
				{"$gte", firstOfLastMonth},
				{"$lt", now},
			}},
		}}},
		{{"$addFields", bson.D{{"day", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m-%d"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$day"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getKillsPerServerPipeline(serverId uint32) mongo.Pipeline {
	// TODO: could be great if we could merge getAllKillsPipeline in it
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$match", bson.D{{"type", "player"}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$addFields", bson.D{{"day", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m-%d"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$day"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getAllKillsPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"type", "player"}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$addFields", bson.D{{"day", bson.D{{"$dateToString", bson.D{
			{"format", "%Y-%m-%d"},
			{"date", "$creationDate"},
		}}}}}}},
		{{"$group", bson.D{
			{"_id", "$day"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getAllConnectionsPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{
		// {{"$match", bson.D{
		// 	{"_id", bson.D{{"$gte", primitive.NewObjectIDFromTimestamp(time.Date(2023, 2, 01, 0, 0, 0, 0, time.UTC))}}},
		// }}},
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
