package token

import (
	"context"
	"core/sec/db/redisdb"
	"core/sec/lib/uuid"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"time"
)

const (
	LowStrengthBitSize      = 512
	MediumStrengthBitSize   = 1024
	HighStrengthBitSize     = 2048
	VeryHighStrengthBitSize = 4096
)

type SignMethodString string

const (
	RS256SignMethodString string = "RS256"
)

func SigningMethodStringToInterface(s string) (method jwt.SigningMethod, err error) {
	switch s {
	case RS256SignMethodString:
		return jwt.SigningMethodRS256, err
	default:
		return nil, fmt.Errorf("undefined method: %s", s)
	}
}

func NewTokenWithClaims(SignMethod string, claims map[string]interface{}) (string, *rsa.PublicKey, error) {
	signMethod, err := SigningMethodStringToInterface(SignMethod)
	if err != nil {
		return "", nil, fmt.Errorf("error when SigningMethodStringToInterface, err:%w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, HighStrengthBitSize)
	if err != nil {
		return "", nil, fmt.Errorf("error when create PrivateKey, err:%w", err)
	}
	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return "", nil, err
	}

	publicKey := privateKey.Public()

	signedToken, err := jwt.NewWithClaims(signMethod, jwt.MapClaims(claims)).SignedString(privateKey)
	if err != nil {
		return "", nil, fmt.Errorf("error when Sign token, err:%w", err)
	}

	return signedToken, publicKey.(*rsa.PublicKey), nil
}

func PEMBytesToPublicKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Key type is not RSA")
	}
}

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

const PublicKeyPrefix = "public_key"

func SaveClaimPublicKey(tokenID string, key *rsa.PublicKey) error {
	byteKey, err := PublicKeyToPEMBytes(key)
	if err != nil {
		return fmt.Errorf("error when PublicKeyToPEMBytes: %w", err)
	}
	//fmt.Println("pubBytes: ", string(byteKey))

	_, err = redisdb.Conn.SetStruct(context.TODO(),fmt.Sprintf("%s_%s",PublicKeyPrefix,tokenID) , byteKey, DefaultExpiresAtTimeDuration)
	return err
}

func GetClaimPublicKey(tokenID string) (*rsa.PublicKey, bool, error) {
	var res []byte
	err := redisdb.Conn.GetUnmarshal(context.TODO(), fmt.Sprintf("%s_%s",PublicKeyPrefix,tokenID) , &res)

	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}

	//fmt.Println("err: ",err)
	//fmt.Println("the public key: ",string(res))
	pubKey, err := PEMBytesToPublicKey(res)
	if err != nil {
		return nil, false, err
	}

	return pubKey, true, err
}

func ParseToken(t string, signMethodString string) (*jwt.Token, error) {
	sm, e := SigningMethodStringToInterface(signMethodString)
	if e != nil {
		return nil, fmt.Errorf("error when SigningMethodStringToInterface: %s", e.Error())
	}

	tk, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if token.Method != sm {
			return nil, errors.New("wrong SignMethod")
		}

		if token.Claims == nil {
			return nil, errors.New("claims is empty")
		}
		stdMap := (*token.Claims.(*jwt.MapClaims))["StandardClaims"]
		if stdMap == nil {
			return nil, errors.New("StandardClaims not found")
		}
		var std jwt.StandardClaims
		err := mapstructure.Decode(stdMap, &std)
		if err != nil {
			return nil, fmt.Errorf("error when Decode: %w", err)
		}

		key, exist, err := GetClaimPublicKey(std.Id)
		if err != nil {
			return nil, err
		} else if !exist {
			return nil, errors.New("public pey not exist")
		}
		return key, e
	})
	if tk == nil {
		return nil, ErrFailToParseClaims
	}
	return tk, err
}

var ErrFailToParseClaims = fmt.Errorf("fail to parse claims")

func ParseTokenWithClaim(stringToken string, signMethodString string, u interface{}) (*jwt.Token, error) {
	sm, e := SigningMethodStringToInterface(signMethodString)
	if e != nil {
		return nil, fmt.Errorf("error when SigningMethodStringToInterface: %s", e.Error())
	}

	var mpc jwt.MapClaims

	tk, err := jwt.ParseWithClaims(stringToken, &mpc, func(token *jwt.Token) (interface{}, error) {

		if token.Method != sm {
			return nil, errors.New("wrong SignMethod")
		}

		if token.Claims == nil {
			return nil, errors.New("tk is empty")
		}
		stdMap := (*token.Claims.(*jwt.MapClaims))["StandardClaims"]
		if stdMap == nil {
			return nil, errors.New("StandardClaims not found")
		}
		var std jwt.StandardClaims
		err := mapstructure.Decode(stdMap, &std)
		if err != nil {
			return nil, fmt.Errorf("error when Decode: %w", err)
		}

		key, exist, err := GetClaimPublicKey(std.Id)
		if err != nil {
			return nil, fmt.Errorf("error when GetClaimPublicKey: %w", err)
		} else if !exist {
			return nil, errors.New("public key not exist")
		}
		return key, nil
	})

	if err != nil {
		return tk, fmt.Errorf("error when jwt.ParseWithClaims: %s", err.Error())
	}
	if tk == nil {
		return nil, ErrFailToParseClaims
	}

	switch u.(type) {
	case *jwt.MapClaims:
		*(u.(*jwt.MapClaims)) = mpc
	case map[string]interface{}:
		*(u.(*map[string]interface{})) = mpc
	default:
		err = mapstructure.Decode(mpc, u)
		if err != nil {
			return nil, fmt.Errorf("error when Decode tk, err: %w, raw:%s", err, tk.Raw)
		}
	}
	return tk, nil
}

// DefaultExpiresAtTimeDuration 預設Token過期時長
var DefaultExpiresAtTimeDuration = time.Hour * 12

// DefaultNotBeforeTimeDuration 預設Token開始有效時長
var DefaultNotBeforeTimeDuration = time.Second * 3




func CreateDefaultStandardClaims(audience string) jwt.StandardClaims {
	nowTime := time.Now()
	return jwt.StandardClaims{
		Audience:  audience,
		ExpiresAt: nowTime.Add(DefaultExpiresAtTimeDuration).Unix(),
		Id:        uuid.NewUUIDWithPrefix("TOKEN"),
		IssuedAt:  nowTime.Unix(),
		Issuer:    "joy-games.online",
		NotBefore: nowTime.Add(DefaultNotBeforeTimeDuration).Unix(),
		Subject:   "",
	}
}

const (
	KEYStandardClaims = "StandardClaims"
	KEYId             = "Id"
)

const UserPrefix = "UserPrefix"

const InvalidateTokenPrefix = "InvalidateStringToken"

func InvalidateTokenByID(tokenID string) error {

	ctx := context.TODO()
	_, err := redisdb.Conn.Set(ctx, fmt.Sprintf("%s_%s", InvalidateTokenPrefix, tokenID), nil, DefaultExpiresAtTimeDuration).Result()
	if err != nil {
		return fmt.Errorf("error whem set InvalidateToken to redis - Set: %w", err)
	}

	_, err = redisdb.Conn.Del(ctx,fmt.Sprintf("%s_%s",PublicKeyPrefix,tokenID)).Result()
	if err != nil {
		return fmt.Errorf("error whem set InvalidateToken to redis - Del: %w", err)
	}

	return nil
}

func GetUserOwnToken(userID string) (string, bool, error) {
	res, err := redisdb.Conn.Get(context.TODO(), fmt.Sprintf("%s_%s", UserPrefix, userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return res, true, nil
}

func SetUserOwnToken(userID string, tokenID string) error {
	_, err := redisdb.Conn.Set(context.TODO(), fmt.Sprintf("%s_%s", UserPrefix, userID), tokenID, DefaultExpiresAtTimeDuration).Result()
	if err != nil {
		return fmt.Errorf("error whem set UserOwnToken to redis: %w", err)
	}
	return nil
}
func DeleteUserOwnTokenByUUID(uuid string) error {
	_ , err := redisdb.Conn.Del(context.TODO() , fmt.Sprintf("%s_%s", UserPrefix, uuid)).Result()
	if err != nil {
		return fmt.Errorf("error whem set DeleteUserOwnTokenByUUID on redis: %w", err)
	}
	return nil
}
// 檢查 token 有沒有被作廢
func CheckTokenActiveWithClaim(claim map[string]interface{}) (bool, error) {
	tokenID, err := GetTokenIDWithClaim(claim)
	if err != nil {
		return false, err
	}
	return CheckTokenActiveWithTokenID(tokenID)
}

func CheckTokenActiveWithTokenID(tokenID string) (bool, error) {

	_, err := redisdb.Conn.Get(context.TODO(), fmt.Sprintf("%s_%s", InvalidateTokenPrefix, tokenID)).Result()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		} else {

			return false, nil
		}
	}

	return false, nil
}

func GetTokenIDWithClaim(claim map[string]interface{}) (string, error) {
	std, ok := claim[KEYStandardClaims]
	if !ok {
		return "", fmt.Errorf("the StandardClaims field is not exist")
	}

	standardClaim, ok := std.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("the StandardClaims field is not type of map[string]interface{}")
	}
	i, ok := standardClaim[KEYId]
	if !ok {
		return "", fmt.Errorf("the token Id field is not exist")
	}
	return fmt.Sprint(i), nil
}
