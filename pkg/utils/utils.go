package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
}

// GenerateVerificationCode 生成 100000 到 999999 之间的随机数
func GenerateVerificationCode() string {
	return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(900000) + 100000)
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func GenerateFromPassword(pwd string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b)
}

func CompareHashAndPassword(hashPwd, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(hashPwd)) == nil
}

func GenerateToken(secretKey string, accessExpire int64, payloads map[string]interface{}) (Token, error) {
	iat := time.Now().Add(-time.Minute).Unix()
	exp := iat + accessExpire
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["ita"] = iat
	claims["exp"] = exp
	for k, v := range payloads {
		claims[k] = v
	}
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return Token{}, err
	}
	return Token{
		AccessToken:  tokenString,
		AccessExpire: exp,
	}, nil

}

func GenerateUUID() string {
	return uuid.New().String()
}
