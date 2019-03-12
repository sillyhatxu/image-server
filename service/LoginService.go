package service

import (
	"errors"
	"github.com/sillyhatxu/sillyhat-cloud-web/jwt"
	"image-server/api/dto"
	"image-server/token"
	"sillyhat-cloud-utils/cache"
	"time"
)

func Login(login dto.LoginDTO) (string, error) {
	if login.LoginName != "sillyhat" || login.Password != "123" {
		return "", errors.New("get user error")
	}
	userToken := *&token.UserToken{UserId: "1001", UserName: "Cookie", LoginTime: time.Now()}
	cache.Set(userToken.UserId, userToken.UserName, cache.DefaultExpiration)
	tokenSrc, err := jwt.CreateTokenStringHS512(userToken)
	if err != nil {
		return "", err
	}
	return tokenSrc, err
}
