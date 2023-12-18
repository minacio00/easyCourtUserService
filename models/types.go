package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type Tenant struct {
	ID          uint   `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	TrialPeriod bool   `json:"periodo_teste"`
}

type Token struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TokenString    string             `bson:"tokenString"`
	ExpirationTime time.Time          `bson:"expirationTime"`
}
