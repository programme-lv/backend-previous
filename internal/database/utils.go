package database

import (
	"log"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func randomString(length int) string {
	const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(result)
}

func passwordsMatch(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
