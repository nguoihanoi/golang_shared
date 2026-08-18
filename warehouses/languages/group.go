package languages

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func CreateGroupCode(insertData GroupCode) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := groupCodeCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateGroupCode(inId string, inUpdateData bSon.M) bool {
	output := groupCodeCollection.UpdateId(inId, inUpdateData)
	if output == true {
		languageCache.Del("group_code:" + inId)
	}
	return output
}

func DeleteGroupCode(inId string) bool {
	output := groupCodeCollection.DeleteId(inId)
	if output == true {
		languageCache.Del("group_code:" + inId)
	}
	return output
}
func CheckGroupCode(inCode string, inId string) (output GroupCode) {
	filter := bSon.M{"code": inCode}
	if inId != "" {
		filter["_id"] = bSon.M{"$ne": inId}
	}
	err := groupCodeCollection.FindOne(filter).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	}
	return output
}
func GetGroupCodeById(inId string, isCache bool) (output GroupCode) {
	if isCache == true {
		cacheData := languageCache.Get("group_code:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(GroupCode)
			return output
		}
	}
	err := groupCodeCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		languageCache.Set("group_code:"+inId, output, 0)
	}
	return output
}

func SearchGroupCode(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []GroupCode, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := groupCodeCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = groupCodeCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []GroupCode{}
		total = 0
	}
	return results, total
}
