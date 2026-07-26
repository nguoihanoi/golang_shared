package libs

import (
	"log"
	"time"

	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libTime "github.com/nguoihanoi/golang_shared/libs/time"
)

func Crypto() *libCrypto.CryptoClass {
	return &libCrypto.CryptoClass{}
}
func Time() *libTime.TimeClass {
	return &libTime.TimeClass{}
}

func main() {
	log.Println(Time().FormatDate(time.Now()))
}
