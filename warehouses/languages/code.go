package languages

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func CreateCode(insertData LanguageCode) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := languageCodeCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateCode(inId string, inUpdateData bSon.M) bool {
	output := languageCodeCollection.UpdateId(inId, inUpdateData)
	if output == true {
		languageCache.Del("language_code:" + inId)
	}
	return output
}

func DeleteCode(inId string) bool {
	output := languageCodeCollection.DeleteId(inId)
	if output == true {
		languageCache.Del("language_code:" + inId)
	}
	return output
}

func GetCodeById(inId string, isCache bool) (output LanguageCode) {
	if isCache == true {
		cacheData := languageCache.Get("language_code:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(LanguageCode)
			return output
		}
	}
	err := languageCodeCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		languageCache.Set("language_code:"+inId, output, 0)
	}
	return output
}
func GetCodeByName(inName string) (output LanguageCode) {
	err := languageCodeCollection.FindOne(bSon.M{"name": inName}).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	}
	return output
}
func CheckCodeByName(inName string, inId string) (output LanguageCode) {
	filter := bSon.M{"name": inName}
	if inId != "" {
		filter["_id"] = bSon.M{"$ne": inId}
	}
	err := languageCodeCollection.FindOne(filter).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	}
	return output
}
func GetCodeByKey(inKey, inLanguageCode string) string {
	if val, ok := languageCache.HGet("codes:"+inLanguageCode, inKey).(string); ok && val != "" {
		return val
	}
	var codeDetail LanguageCode
	if err := languageCodeCollection.FindOne(bSon.M{"name": inKey}).Decode(&codeDetail); err == nil {
		for lang, val := range codeDetail.Value {
			languageCache.HSet("codes:"+lang, inKey, val)
		}
		if targetVal, exists := codeDetail.Value[inLanguageCode]; exists {
			return targetVal
		}
	}
	return inKey
}

func GetCodesByType(inType int) (results []LanguageCode) {
	cursor, err := languageCollection.Find(bSon.M{"delete": 0, "type": inType}, bSon.D{{Key: "delete", Value: 1}, {Key: "type", Value: 1}, {Key: "group_id", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}
func SearchCode(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []LanguageCode, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := languageCodeCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = languageCodeCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []LanguageCode{}
		total = 0
	}
	return results, total
}
