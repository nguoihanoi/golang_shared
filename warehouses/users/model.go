package users

import (
	"log"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var dbMongo *libDb.DatabaseClass
var userCollection *libDb.CollectionClass
var userCache *libCache.Cache

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache) {
	dbMongo = inDb
	userCache = inRedisClient
	userCollection = dbMongo.NewCollection("users")
}

func Login(userDetail User, inPassword string) (User, bool) {
	resultStatus := false
	newPassword, _ := libUtilities.String().GetHashPassWord(inPassword, userDetail.PasswordHash, true)
	if newPassword == userDetail.Password {
		newToken, nextTime, err := libCrypto.JWT().CreateToken(UserToken{
			UserID:       userDetail.ID,
			LanguageCode: userDetail.LanguageCode,
			Time:         time.Now().Unix(),
		})
		if err != nil {
			panic(err)
		} else {
			resultStatus = UpdateUser(userDetail.ID, bSon.M{
				"reset_token":     newToken,
				"reset_token_exp": nextTime,
			})
			userDetail.ResetToken = newToken
		}
	}
	return userDetail, resultStatus
}

func CreateUser(insertData User) (output string) {
	output = userCollection.Create(insertData)
	return output
}

func UpdateUser(inId string, inUpdateData bSon.M) bool {
	output := userCollection.UpdateId(inId, inUpdateData)
	if output == true {
		userCache.Del("user:" + inId)
	}
	return output
}

func GetUserById(inId string, isCache bool) (output User) {
	if isCache == true {
		cacheData := userCache.Get("user:" + inId)
		if cacheData != nil {
			return cacheData.(User)
		}
	}
	result := userCollection.FindById(inId)
	if result != nil {
		output = result.(User)
	}
	return output
}

func GetUserByEmail(inEmail string, inId string) (output User) {
	filter := bSon.M{"email": inEmail}
	if inId != "" {
		filter["_id"] = bSon.M{"$ne": inId}
	}
	result := userCollection.FindOne(filter)
	log.Println(result)
	if result != nil {
		output = result.(User)
	}
	return output
}
