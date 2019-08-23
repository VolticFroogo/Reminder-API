package jwt

import (
	"context"
	"crypto/rsa"
	"os"
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
func LoadKeys() (err error) {
	err = LoadPublicKey()
	if err != nil {
		return
	}

	err = LoadPrivateKey()
	return
}

// LoadPublicKey loads the public key.
func LoadPublicKey() (err error) {
	key := os.Getenv("PUBLIC_KEY")

	bytes := []byte(key)

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	return
}

// LoadPrivateKey loads the private key.
func LoadPrivateKey() (err error) {
	key := os.Getenv("PRIVATE_KEY")
	password := os.Getenv("PRIVATE_KEY_PASSWORD")

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
	expiry := now.Add(model.AuthDuration).Unix()

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
