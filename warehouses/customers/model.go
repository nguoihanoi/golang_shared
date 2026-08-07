package customers

import (
	"context"
	"sync"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var dbMongo *libDb.DatabaseClass
var customerCollection *libDb.CollectionClass
var groupCollection *libDb.CollectionClass
var customerCache *libCache.Cache
var libJwt *libCrypto.JwtClass

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	dbMongo = inDb
	customerCache = inRedisClient
	customerCollection = dbMongo.NewCollection("customers")
	customerCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_first_name_1_last_name_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "first_name", Value: 1}, {Key: "last_name", Value: 1}},
		},
		{
			Name: "delete_1_customer_group_id_1_first_name_1_last_name_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "customer_group_id", Value: 1}, {Key: "first_name", Value: 1}, {Key: "last_name", Value: 1}},
		},
	})
	groupCollection = dbMongo.NewCollection("customer_groups")
	groupCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_status_-1_order_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}, {Key: "order", Value: 1}},
		},
	})
	libJwt = libCrypto.JWT(inJwtToken)
}

func CreateGroup(insertData CustomerGroup) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := groupCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateGroup(inId string, inUpdateData bSon.M) bool {
	output := groupCollection.UpdateId(inId, inUpdateData)
	if output == true {
		customerCache.Del("customer_group:" + inId)
	}
	return output
}

func DeleteGroup(inId string) bool {
	output := groupCollection.DeleteId(inId)
	if output == true {
		customerCache.Del("customer_group:" + inId)
	}
	return output
}

func GetGroupById(inId string, isCache bool) (output CustomerGroup) {
	if isCache == true {
		cacheData := customerCache.Get("customer_group:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(CustomerGroup)
			return output
		}
	}
	err := groupCollection.FindById(inId).Decode(&output)
	if err != nil {
		output.ID = ""
	}
	return output
}

func GetGroups() (results []CustomerGroup) {
	cursor, err := groupCollection.Find(bSon.M{"delete": 0, "status": 1}, bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}, {Key: "order", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}

func SearchGroups(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []CustomerGroup, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := groupCollection.Find(filter, inSortOrder, inPage, inLimit)
		if err != nil {
			findErr = err
			return
		}
		if err := cursor.All(context.TODO(), &results); err != nil {
			findErr = err
			return
		}
	}()

	// Tác vụ 2: Đếm tổng số lượng bản ghi
	go func() {
		defer wg.Done()
		total = groupCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []CustomerGroup{}
		total = 0
	}
	return results, total
}

func Login(customerDetail Customer, inPassword string, inEncrypted bool) (Customer, bool) {
	resultStatus := false
	newPassword, _ := libUtilities.String().GetHashPassWord(inPassword, customerDetail.PasswordHash, inEncrypted)
	if newPassword == customerDetail.Password {
		newToken, nextTime, err := libJwt.CreateToken(CustomerToken{
			CustomerID:   customerDetail.ID,
			LanguageCode: customerDetail.LanguageCode,
			Time:         time.Now().Unix(),
		})
		if err != nil {
			panic(err)
		} else {
			resultStatus = UpdateCustomer(customerDetail.ID, bSon.M{
				"reset_token":     newToken,
				"reset_token_exp": nextTime,
			})
			customerDetail.ResetToken = newToken
		}
	}
	return customerDetail, resultStatus
}

func CreateCustomer(insertData Customer) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := customerCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateCustomer(inId string, inUpdateData bSon.M) bool {
	output := customerCollection.UpdateId(inId, inUpdateData)
	if output == true {
		customerCache.Del("customer:" + inId)
	}
	return output
}

func GetCustomerById(inId string, isCache bool) (output Customer) {
	if isCache == true {
		cacheData := customerCache.Get("customer:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(Customer)
			return output
		}
	}
	err := customerCollection.FindById(inId).Decode(&output)
	if err != nil {
		output.ID = ""
	}
	return output
}

func GetCustomerByEmail(inEmail string, inId string) (output Customer) {
	filter := bSon.M{"email": inEmail}
	if inId != "" {
		filter["_id"] = bSon.M{"$ne": inId}
	}
	err := customerCollection.FindOne(filter).Decode(&output)
	if err != nil {
		output.ID = ""
	}
	return output
}

func Search(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []Customer, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := customerCollection.Find(filter, inSortOrder, inPage, inLimit)
		if err != nil {
			findErr = err
			return
		}
		if err := cursor.All(context.TODO(), &results); err != nil {
			findErr = err
			return
		}
	}()

	// Tác vụ 2: Đếm tổng số lượng bản ghi
	go func() {
		defer wg.Done()
		total = customerCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []Customer{}
		total = 0
	}
	return results, total
}
