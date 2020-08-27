package utils

import "golang.org/x/crypto/bcrypt"

// HashAndSalt The below function uses GenerateFromPassword to generate a salted hash which is returned as a byte slice. We then return the byte slice as a string so that we can store the salted hash in our database as the users password.
func HashAndSalt(pwd []byte) (hash string, err error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	var hashBytes []byte
	hashBytes, err = bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	hash = string(hashBytes)
	return
}

// ComparePasswords Using CompareHashAndPassword we can create a function that returns a boolean to let us know whether or not the passwords match.
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}
