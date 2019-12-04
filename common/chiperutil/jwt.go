package chiperutil

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// token := jwt.New(jwt.SigningMethodHS256) 和 token.Claims = claims
// 等同于 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
func CreateJWT(secretKey string, claimsMap map[string]string, expireSecond int64) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	// alg属性表示签名使用的算法,默认为HMAC SHA256(写为HS256);typ属性表示令牌的类型;JWT令牌统一写为JWT
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	// Claims
	claims := make(jwt.MapClaims)
	for key, value := range claimsMap {
		claims[key] = value
	}
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expireSecond)).Unix() // 过期时间
	token.Claims = claims
	// Signature
	// 使用自定义字符串加密,并将完整的编码令牌作为字符串
	tokenString, err = token.SignedString([]byte(secretKey))
	return
}

func reverseJWT(secretKey, tokenString string) (t *jwt.Token, errType string, err error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				return nil, "errInputData", err
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, "errExpired", err
			} else {
				// Couldn't handle this token
				return nil, "errInputData", err
			}
		} else {
			// Couldn't handle this token
			return nil, "errInputData", err
		}
	}
	if !token.Valid {
		return nil, "errInputData", err
	}
	return token, "", err
}

func ParseJWT(secretKey, tokenString string) (map[string]interface{}, error) {
	token, _, err := reverseJWT(secretKey, tokenString)
	if err == nil {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return claims, nil
		}
	}
	return nil, err
}
