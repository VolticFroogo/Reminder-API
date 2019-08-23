package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/model"
	"github.com/dgrijalva/jwt-go"
)

// Variables
var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

// LoadKeys loads both keys.
func LoadKeys(public, private, password string) (err error) {
	err = LoadPublic(public)
	if err != nil {
		return
	}

	err = LoadPrivate(private, password)
	return
}

// LoadPublic loads the public key.
func LoadPublic(key string) (err error) {
	bytes := []byte(key)

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	return
}

// LoadPrivate loads the private key.
func LoadPrivate(key, password string) (err error) {
	bytes := []byte(key)

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEMWithPassword(bytes, password)
	return
}

// NewTokens generates new tokens given a user key.
func NewTokens(user *datastore.Key, client *datastore.Client) (auth, refresh string, err error) {
	auth, err = newAuth(user)
	if err != nil {
		return
	}

	refresh, err = newRefresh(user, client)
	return
}

func newAuth(user *datastore.Key) (token string, err error) {
	// Get all of the time information.
	now := time.Now()
	issued := now.Unix()
	expiry := now.Add(model.AuthDuration).Unix()

	// Create the claims for the token.
	claims := model.Token{
		jwt.StandardClaims{
			Subject:   user.Encode(),
			IssuedAt:  issued,
			ExpiresAt: expiry,
		},
	}

	// Create a new unsigned token.
	unsigned := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token.
	token, err = unsigned.SignedString(privateKey)

	return
}

func newRefresh(user *datastore.Key, client *datastore.Client) (token string, err error) {
	// Get all of the time information.
	now := time.Now()
	issued := now.Unix()
	expiry := now.Add(model.RefreshDuration).Unix()

	// Create the JTI.
	jti := model.JTI{
		Expiry: expiry,
	}

	// Get context.
	ctx := context.Background()

	// Create a key.
	key := datastore.IncompleteKey(model.KindJTI, user)

	// Save the new entity.
	key, err = client.Put(ctx, key, &jti)
	if err != nil {
		return
	}

	// Create the claims for the token.
	claims := model.Token{
		jwt.StandardClaims{
			Id:        key.Encode(),
			Subject:   user.Encode(),
			IssuedAt:  issued,
			ExpiresAt: expiry,
		},
	}

	// Create a new unsigned token.
	unsigned := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token.
	token, err = unsigned.SignedString(privateKey)

	return
}

// CheckAuth checks the validity of an auth token.
func CheckAuth(tokenString string) (valid bool, user string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, model.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return
	}

	if !token.Valid {
		return
	}

	claims, ok := token.Claims.(*model.Token)
	if !ok {
		return
	}

	if claims.Id != "" {
		return false, "", fmt.Errorf("got refresh token not auth")
	}

	return token.Valid, claims.Subject, nil
}

// CheckRefresh checks the validity of a refresh token.
func CheckRefresh(tokenString string, client *datastore.Client) (valid bool, user string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, model.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return
	}

	if !token.Valid {
		return
	}

	claims, ok := token.Claims.(*model.Token)
	if !ok {
		return
	}

	ctx := context.Background()

	key, err := datastore.DecodeKey(claims.Id)
	if err != nil {
		return
	}

	var jti model.JTI

	err = client.Get(ctx, key, &jti)
	if err != nil {
		return
	}

	if time.Now().Unix() > jti.Expiry {
		return
	}

	return token.Valid, claims.Subject, nil
}
