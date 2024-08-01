package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CountPerServerResult struct {
	Count uint32 `bson:"count"`
}

type Kills struct {
	ServerId uint32 `bson:"serverId" json:"serverId"`
}
type ConnectionsPerMonth struct {
	Id    string `bson:"_id" json:"_id"`
	Count uint32 `bson:"count" json:"count"`
}

type TopKiller struct {
	CharacterName string `bson:"characterName"`
	Count         uint32 `bson:"count"`
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
