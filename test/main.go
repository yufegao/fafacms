package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func main() {
	fmt.Println(Sha2256Days())
}

func Sha2256(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Sha2256Days() string {
	temp := Sha2256([]byte(TodayString(3)))
	//fmt.Printf("%x\n", temp)
	return Base64E(temp[:])
}

func Base64E(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func TodayString(level int) string {
	formats := "20060102150405"
	switch level {
	case 1:
		formats = "2006"
	case 2:
		formats = "200601"
	case 3:
		formats = "20060102"
	case 4:
		formats = "2006010215"
	case 5:
		formats = "200601021504"
	default:

	}
	return time.Now().Format(formats)
}
