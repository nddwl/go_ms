package utils

import (
	"applet/internal/types"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

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

func GenerateToken(secretKey string, accessExpire int64, payloads map[string]interface{}) (types.Token, error) {
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
		return types.Token{}, err
	}
	return types.Token{
		AccessToken:  tokenString,
		AccessExpire: exp,
	}, nil

}
