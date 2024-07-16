package utils

import "golang.org/x/crypto/bcrypt"


/*
GenerateHashPassword: This function takes a plain text password as input 
and returns a hash value generated from it using a one-way hashing algorithm. 
The purpose of this function is to store a user's password securely in the database. 
This way, even if the database is compromised, the attacker cannot retrieve 
the original password as it is encrypted.
*/
func GenerateHashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}