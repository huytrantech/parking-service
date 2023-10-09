package utils

import (
	"crypto/sha512"
	"fmt"
	"time"
)

func GenToken(username string , driverId int) string {
	sum := sha512.Sum512([]byte(fmt.Sprintf("%s-%d-%d-login-parking" , username , driverId , time.Now().Unix())))
	return fmt.Sprintf("%x" , sum)
}

func GenInternalAccountToken(username string , password string) string {
	sum := sha512.Sum512([]byte(fmt.Sprintf("%s-%s-%d-login-internal-parking" , username , password , time.Now().Unix())))
	return fmt.Sprintf("%x" , sum)
}


func GenPassword(password string) string {
	sum := sha512.Sum512([]byte(fmt.Sprintf("%s-%d-password" , password , time.Now().Unix())))
	return fmt.Sprintf("%x" , sum)
}