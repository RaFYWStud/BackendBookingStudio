package token

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/RaFYWStud/BackendBookingStudio/config"
	"github.com/golang-jwt/jwt/v5"
)

// jwtConfig holds the loaded JWT configuration and RSA keys used for signing and verifying tokens.
// It is initialized by the Load() function and accessed by other token-related functions in this package.
var jwtConfig *jwtStruct

type jwtStruct struct {
	jwtLifeTime        uint
	jwtRefreshLifeTime uint
	privateKey         *rsa.PrivateKey
	publicKey          *rsa.PublicKey
}

// Load reads RSA public and private keys from configured file paths,
// parses them, and sets the global jwtConfig variable with the loaded keys and token lifetimes.
func Load() {
	cfg := config.Get()

	publicKeyPath := cfg.PublicKeyPath
	privateKeyPath := cfg.PrivateKeyPath

	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("Error reading public key file: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key file: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	jwtConfig = &jwtStruct{
		jwtLifeTime:        cfg.AccessTokenLifeTime,
		jwtRefreshLifeTime: cfg.RefreshTokenLifeTime,
		publicKey:          publicKey,
		privateKey:         privateKey,
	}
}
