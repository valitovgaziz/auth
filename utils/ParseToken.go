package utils

import (
    "go-auth/models"

    "github.com/dgrijalva/jwt-go"
)


/*
ParseToken: This function takes a JWT token as input and returns the claims contained in it. 
Claims are a set of key-value pairs that represent the information being transmitted between parties. 
In this case, the claims may include the subject (email), role, and expiration time of the token. 
This function is used to validate if the token is legitimate and to retrieve information contained in it.
*/
func ParseToken(tokenString string) (claims *models.Claims, err error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("my_secret_key"), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*models.Claims)

    if !ok {
        return nil, err
    }

    return claims, nil
}