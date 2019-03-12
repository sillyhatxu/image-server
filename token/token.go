package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserToken struct {
	UserId    string
	UserName  string
	LoginTime time.Time
	jwt.StandardClaims
}
