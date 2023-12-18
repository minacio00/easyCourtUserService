package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/minacio00/easyCourtUserService/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName   = "tokens"
	collectionName = "blacklistedTokens"
)

type AuthenticatorService struct {
	mongo *mongo.Client
}

func NewAuthenticatorService(mongo *mongo.Client) *AuthenticatorService {
	return &AuthenticatorService{mongo: mongo}
}

func (as *AuthenticatorService) BlacklistToken(tkStr string) error {
	collection := as.mongo.Database(databaseName).Collection(collectionName)

	_, err := collection.InsertOne(context.Background(), models.Token{
		TokenString: tkStr,
	})
	return err
}
func (as *AuthenticatorService) IsTokenBlacklisted(tokenString string) (bool, error) {
	collection := as.mongo.Database(databaseName).Collection(collectionName)

	var result models.Token
	err := collection.FindOne(context.Background(), bson.M{"tokenString": tokenString}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (as *AuthenticatorService) ValidateToken(tokenString string) (bool, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return false, nil
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token") // Invalid token
	}
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	currentTime := time.Now()

	return expirationTime.Before(currentTime), nil

}

func (as *AuthenticatorService) GetAllBlacklistedTokens() ([]models.Token, error) {
	collection := as.mongo.Database(databaseName).Collection(collectionName)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var blacklistedTokens []models.Token
	for cursor.Next(context.Background()) {
		var token models.Token
		if err := cursor.Decode(&token); err != nil {
			return nil, err
		}
		blacklistedTokens = append(blacklistedTokens, token)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return blacklistedTokens, nil
}
