package permissions

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func CreateType(insertData PermissionType) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := permissionTypeCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateType(inId string, inUpdateData bSon.M) bool {
	output := permissionTypeCollection.UpdateId(inId, inUpdateData)
	if output == true {
		permissionCache.Del("permission_type:" + inId)
	}
	return output
}

func DeleteType(inId string) bool {
	output := permissionTypeCollection.DeleteId(inId)
	if output == true {
		permissionCache.Del("permission_type:" + inId)
	}
	return output
}

func GetTypeById(inId string, isCache bool) (output PermissionType) {
	if isCache == true {
		cacheData := permissionCache.Get("permission_type:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(PermissionType)
			return output
		}
	}
	err := permissionTypeCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		permissionCache.Set("permission_type:"+inId, output, 0)
	}
	return output
}

func GetTypes() (results []MiniPermissionType) {
	cursor, err := permissionTypeCollection.Find(bSon.M{"delete": 0}, bSon.D{{Key: "delete", Value: 1}, {Key: "order", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}

func SearchTypes(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []PermissionType, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := permissionTypeCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = permissionTypeCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []PermissionType{}
		total = 0
	}
	return results, total
}
