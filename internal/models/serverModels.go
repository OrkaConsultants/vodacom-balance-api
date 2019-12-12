package models

import "github.com/dgrijalva/jwt-go"

type JWTClaimsStruct struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
