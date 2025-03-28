package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func getTopPlaytimePipeline(serverId uint32, maxResult uint8) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$sort", bson.D{{"playTime", -1}}}}, // Sort by playtime in descending order
		{{"$limit", maxResult}},               // Limit the results to maxResult
		{{"$project", bson.D{
			{"_id", 0},                          // Exclude the _id field from the output
			{"characterName", "$characterName"}, // Include the characterName field
			{"count", "$playTime"},              // Include the playtime field
		}}},
	}

	return pipeline
}

func getConnectionsLastMonthPerServerPipeline(serverId uint32) mongo.Pipeline {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"serverId", serverId},
			{"_id", bson.D{
				{"$gte", primitive.NewObjectIDFromTimestamp(firstOfLastMonth)},
				{"$lt", primitive.NewObjectIDFromTimestamp(firstOfThisMonth)},
			}},
		}}},
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m-%d"},
					{"date", bson.D{{"$toDate", "$_id"}}},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getConnectionsPerServerPipeline(serverId uint32) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m"},
					{"date", bson.D{{"$toDate", "$_id"}}},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getAllConnectionsLastMonthPipeline() mongo.Pipeline {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	lt := now
	gte := now.AddDate(0, -1, 0) // Move one month back

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"_id", bson.D{
				{"$gte", primitive.NewObjectIDFromTimestamp(gte)},
				{"$lt", primitive.NewObjectIDFromTimestamp(lt)},
			}},
		}}},
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m-%d"},
					{"date", bson.D{{"$toDate", "$_id"}}},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getAllConnectionsLastYearPipeline() mongo.Pipeline {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	lt := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC) // First day of the current month
	gte := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)        // First day of last year

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"_id", bson.D{
				{"$gte", primitive.NewObjectIDFromTimestamp(gte)},
				{"$lt", primitive.NewObjectIDFromTimestamp(lt)},
			}},
		}}},
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m"},
					{"date", bson.D{{"$toDate", "$_id"}}},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	return pipeline
}

func getKillsPerServerPipeline(serverId uint32, entityType string) mongo.Pipeline {
	now := time.Now()
	lt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"serverId", serverId}}}},
		{{"$match", bson.D{{"type", entityType}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$match", bson.D{
			{"creationDate", bson.D{
				{"$lt", lt},
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

func getAllKillsPipeline(entityType string) mongo.Pipeline {
	now := time.Now()
	lt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"type", entityType}}}},
		{{"$addFields", bson.D{{"creationDate", bson.D{{"$toDate", "$_id"}}}}}},
		{{"$match", bson.D{
			{"creationDate", bson.D{
				{"$lt", lt},
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
