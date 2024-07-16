package controllers

import (
    "go-auth/models"
    "time"

    "go-auth/utils"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)

/*The string "my_secret_key" is just an example 
and should be replaced with a secret key of sufficient 
length and complexity in a real-world scenario.*/
var jwtKey = []byte("my_secret_key")


/*Login: Handles user login requests by receiving user credentials, 
verifying if the user exists, and checking if the provided password is correct. 
If everything is okay, it generates a JSON Web Token (JWT) and sets it as a cookie 
in the client's browser.*/
func Login(c *gin.Context) {

    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existingUser models.User

    models.DB.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID == 0 {
        c.JSON(400, gin.H{"error": "user does not exist"})
        return
    }

    errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

    if !errHash {
        c.JSON(400, gin.H{"error": "invalid password"})
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)

    claims := &models.Claims{
        Role: existingUser.Role,
        StandardClaims: jwt.StandardClaims{
            Subject:   existingUser.Email,
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        c.JSON(500, gin.H{"error": "could not generate token"})
        return
    }

    c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
    c.JSON(200, gin.H{"success": "user logged in"})
}


/*
Signup: Handles user signup requests by receiving user information,
checking if the user already exists, hashing the password, and creating
a new user in the database.
*/
func Signup(c *gin.Context) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existingUser models.User

    models.DB.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID != 0 {
        c.JSON(400, gin.H{"error": "user already exists"})
        return
    }

    var errHash error
    user.Password, errHash = utils.GenerateHashPassword(user.Password)

    if errHash != nil {
        c.JSON(500, gin.H{"error": "could not generate password hash"})
        return
    }

    models.DB.Create(&user)

    c.JSON(200, gin.H{"success": "user created"})
}

/*
Home: Handles requests to the home page by checking if the user is authorized to access 
it by looking for the JWT in the client's cookie and verifying it. If the user is authorized, 
the function returns the "home page" message along with the user's role.
*/
func Home(c *gin.Context) {

    cookie, err := c.Cookie("token")

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    claims, err := utils.ParseToken(cookie)

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    if claims.Role != "user" && claims.Role != "admin" {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    c.JSON(200, gin.H{"success": "home page", "role": claims.Role})
}


// Premium: Similar to the Home function but only returns the "premium page" message
// and the user's role if the user has an "admin" role.
func Premium(c *gin.Context) {

    cookie, err := c.Cookie("token")

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    claims, err := utils.ParseToken(cookie)

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    if claims.Role != "admin" {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    c.JSON(200, gin.H{"success": "premium page", "role": claims.Role})
}

// Logout: Handles user logout requests by deleting the JWT cookie in the client's browser.
func Logout(c *gin.Context) {
    c.SetCookie("token", "", -1, "/", "localhost", false, true)
    c.JSON(200, gin.H{"success": "user logged out"})
}