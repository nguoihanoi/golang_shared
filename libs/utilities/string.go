package utilities

import (
	"log"
	"math/rand"
	"strings"
	"time"

	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
)

type stringClass struct {
	//radius float64 // Private field
}

func (s *stringClass) String() *stringClass {
	return &stringClass{}
}

func (s *stringClass) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (s *stringClass) GetHashPassWord(inPassWord string, inHash string, inEncrypted bool) (string, string) {
	if inHash == "" {
		inHash = s.GenerateRandomString(6)
	}
	if inEncrypted {
		if inPassWord == "" {
			inPassWord = libCrypto.Sha3(s.GenerateRandomString(9))
		}
	} else {
		if inPassWord == "" {
			inPassWord = s.GenerateRandomString(9)
		}
		inPassWord = libCrypto.Sha3(inPassWord)
	}
	return libCrypto.Sha3(inHash + inPassWord), inHash
}

func (s *stringClass) GetFromFullName(inFullName string) (string, string) {
	inFirstName := ""
	inLastName := ""
	libProcess.Try(func() {
		nameArr := strings.Split(inFullName, " ")
		if len(nameArr) > 0 {
			inFirstName = nameArr[0]
			inLastName = strings.Join(nameArr[1:len(nameArr)], " ")
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return inFirstName, inLastName
}
