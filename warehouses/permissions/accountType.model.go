package permissions

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func CreateAccountType(insertData AccountType) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := accountTypeCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateAccountType(inId string, inUpdateData bSon.M) bool {
	output := accountTypeCollection.UpdateId(inId, inUpdateData)
	if output == true {
		permissionCache.Del("account_type:" + inId)
	}
	return output
}

func DeleteAccountType(inId string) bool {
	output := accountTypeCollection.DeleteId(inId)
	if output == true {
		permissionCache.Del("account_type:" + inId)
	}
	return output
}

func GetAccountTypeById(inId string, isCache bool) (output AccountType) {
	if isCache == true {
		cacheData := permissionCache.Get("account_type:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(AccountType)
			return output
		}
	}
	err := accountTypeCollection.FindById(inId).Decode(&output)
	if err != nil {
		output.ID = ""
	} else {
		permissionCache.Set("account_type:"+inId, output, 0)
	}
	return output
}

func GetAccountTypes() (results []MiniAccountType) {
	cursor, err := accountTypeCollection.Find(bSon.M{"delete": 0, "status": 1}, bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}, {Key: "order", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}

func SearchAccountTypes(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []AccountType, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := accountTypeCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = accountTypeCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []AccountType{}
		total = 0
	}
	return results, total
}
