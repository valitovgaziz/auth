package utils

import "golang.org/x/crypto/bcrypt"


/*CompareHashPassword: This function takes the user input password 
and the hashed password stored in the database and compares them. 
If the hashes match, it returns true. This function is used to verify 
if the user has entered the correct password during login.*/
func CompareHashPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}