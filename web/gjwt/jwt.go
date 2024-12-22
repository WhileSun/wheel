package gjwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtConf
// Version 版本号用于
type GjwtConf struct {
	Secret  string  `default:""`
	Exp     int64   `default:"1800"`
	LastExp int64   `default:"43200"` // 失效后多久内自动续期
	Version float64 `default:"1.0"`
}

func New(gjwtConf GjwtConf) *GjwtConf {
	if gjwtConf.Secret == "" {
		fmt.Println("gjwt secret is empty, please check your config")
		return nil
	}
	return &gjwtConf
}

func (gjwtConf *GjwtConf) CreateToken(values jwt.MapClaims) (string, error) {
	values["exp"] = time.Now().Unix() + gjwtConf.Exp
	values["lastexp"] = time.Now().Unix() + gjwtConf.LastExp
	values["iat"] = time.Now().Unix()
	values["version"] = gjwtConf.Version
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, values)
	token, err := at.SignedString([]byte(gjwtConf.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (gjwtConf *GjwtConf) ParseToken(tokenString string) (jwt.MapClaims, error) {
	//检测加密方式是否一致
	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		err, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return err, nil
		}
		return []byte(gjwtConf.Secret), nil
	})
	if token == nil {
		return nil, tokenErr
	}
	if !token.Valid {
		if token.Claims == nil {
			return nil, tokenErr
		}
		claims := token.Claims.(jwt.MapClaims)
		return claims, tokenErr
	} else {
		claims := token.Claims.(jwt.MapClaims)
		return claims, nil
	}
}

// 刷新token
func (gjwtConf *GjwtConf) RefreshToken(tokenString string) (string, error) {
	claims, err := gjwtConf.ParseToken(tokenString)
	if err == nil {
		return tokenString, nil
	}
	// 检测是否过期
	if strings.Contains(err.Error(), "expired") {
		nowTime := time.Now().Unix()
		if int64(claims["lastexp"].(float64)) >= nowTime {
			// 生成新token
			token, err := gjwtConf.CreateToken(claims)
			if err == nil {
				return token, nil
			}
		}
	}
	return "", err
}
