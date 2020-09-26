package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// User struct
	User struct {
		ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Username string             `json:"username,omitempty"`
		Fullame  string             `json:"fullname"`
		Password string             `json:"password,omitempty"`
	}
)
