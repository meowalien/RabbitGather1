package token

import (
	"context"
	"core/src/conf"
	"core/src/lib/errs"
	"core/src/lib/uuid"
	"core/src/module/db/redisdb"
	"core/src/module/log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const PublicKeyPrefix = "public_key"

func savePublicKey(tokenID string, key string) error {
	_, err := redisdb.Conn.SetStruct(context.TODO(), redisdb.FormatKey(PublicKeyPrefix, tokenID), key, time.Duration(conf.GlobalConfig.JWT.TokenExpiresAt))
	return err
}

// return stringtoken , publickey , err
func NewTokenWithClaims(cls ...map[string]interface{}) (tk string, err error) {
	if cls == nil {
		return "", errs.WithLine(fmt.Errorf("empty claim"))
	}
	var claim =map[string]interface{}{}
	if len(cls) == 1 {
		claim = cls[0]
	} else {
		for _, v := range cls {
			for k, vv := range v {
				claim[k] = vv
			}
		}
	}

	tokenID, exist := GetTokenIDInClaim(claim)
	if !exist {
		err = fmt.Errorf("token id notfound")
		return
	}
	publicKey := ""
	tk, publicKey, err = newTokenWithClaim(claim)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	err = savePublicKey(tokenID, publicKey)
	if err != nil {
		return "", errs.WithLine(err)
	}
	return tk, nil
}

func GetTokenIDInClaim(claim map[string]interface{}) (string, bool) {
	idkey, ok := claim[TokenIDKey]
	if !ok {
		return "", false
	}

	sk, ok := idkey.(string)
	if !ok {
		log.Logger.Errorf("the token id is not string")
		return "", false
	}
	return sk, true
}

var SignMethod = jwt.SigningMethodRS256

const (
	LowStrengthBitSize      = 512
	MediumStrengthBitSize   = 1024
	HighStrengthBitSize     = 2048
	VeryHighStrengthBitSize = 4096
)

func PublicKeyToPEMBytes(key *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return pubBytes, nil
}

func newTokenWithClaim(claims map[string]interface{}) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, HighStrengthBitSize)
	if err != nil {
		return "", "", fmt.Errorf("error when create PrivateKey, err:%w", err)
	}
	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.Public()

	signedToken, err := jwt.NewWithClaims(SignMethod, jwt.MapClaims(claims)).SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("error when Sign token, err:%w", err)
	}

	pub, err := PublicKeyToPEMBytes(publicKey.(*rsa.PublicKey))
	if err != nil {
		return "", "", errs.WithLine(err)
	}
	return signedToken, string(pub), nil
}

const TokenIDKey = "jti"
const Subject = "sub"
const APIClaim = "api_token"
const RefreshClaim = "refresh_token"

func StandardAPIClaims() map[string]interface{} {
	nowTime := time.Now()
	return map[string]interface{}{
		//"aud": audience,
		Subject:    APIClaim,
		"exp":      nowTime.Add(time.Duration(conf.GlobalConfig.JWT.TokenExpiresAt)).Unix(),
		TokenIDKey: uuid.NewUUID("T"),
		"iat":      nowTime.Unix(),
		"iss":      conf.GlobalConfig.JWT.Issuer,
		"nbf":      nowTime.Add(time.Duration(conf.GlobalConfig.JWT.TokenNotBefore)).Unix(),
	}
}

func StandardRefreshClaims() map[string]interface{} {
	nowTime := time.Now()
	return map[string]interface{}{
		//"aud": audience,
		Subject:    RefreshClaim,
		"exp":      nowTime.Add(time.Duration(conf.GlobalConfig.JWT.TokenExpiresAt)).Unix(),
		TokenIDKey: uuid.NewUUID("T"),
		"iat":      nowTime.Unix(),
		"iss":      conf.GlobalConfig.JWT.Issuer,
		"nbf":      nowTime.Add(time.Duration(conf.GlobalConfig.JWT.TokenNotBefore)).Unix(),
	}
}
