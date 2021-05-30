package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Id                  primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Name                string             `bson:"name,omitempty" json:"name,omitempty"`
	Conference          string             `bson:"conference,omitempty" json:"conference,omitempty"`
	Divisional          string             `bson:"divisional,omitempty" json:"divisional,omitempty"`
	Stadium             string             `bson:"stadium,omitempty" json:"stadium,omitempty"`
	State               string             `bson:"state,omitempty" json:"state,omitempty"`
	Titles              int64              `bson:"titles,omitempty" json:"titles,omitempty"`
	SuperBowlAppearance int64              `bson:"superBowlAppearance,omitempty" json:"superBowlAppearance,omitempty"`
}
