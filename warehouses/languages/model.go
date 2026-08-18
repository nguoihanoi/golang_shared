package languages

import (
	"context"
	"sync"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var dbMongo *libDb.DatabaseClass
var languageCollection *libDb.CollectionClass
var languageCodeCollection *libDb.CollectionClass
var groupCodeCollection *libDb.CollectionClass
var languageCache *libCache.Cache

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache) {
	dbMongo = inDb
	languageCache = inRedisClient
	languageCollection = dbMongo.NewCollection("languages")
	languageCodeCollection = dbMongo.NewCollection("language_codes")
	groupCodeCollection = dbMongo.NewCollection("group_codes")
	languageCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_status_-1_order_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}, {Key: "order", Value: 1}},
		},
	})
}

func Create(insertData Language) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := languageCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func Update(inId string, inUpdateData bSon.M) bool {
	output := languageCollection.UpdateId(inId, inUpdateData)
	if output == true {
		languageCache.Del("language:" + inId)
	}
	return output
}

func Delete(inId string) bool {
	output := languageCollection.DeleteId(inId)
	if output == true {
		languageCache.Del("language:" + inId)
	}
	return output
}

func GetById(inId string, isCache bool) (output Language) {
	if isCache == true {
		cacheData := languageCache.Get("language:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(Language)
			return output
		}
	}
	err := languageCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		languageCache.Set("language:"+inId, output, 0)
	}
	return output
}

func Gets() (results []Language) {
	cursor, err := languageCollection.Find(bSon.M{"delete": 0, "status": 1}, bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}, {Key: "order", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}

func GetCodes() (results []string) {
	languages := Gets()
	for i := range languages {
		results = append(results, languages[i].Code)
	}
	return results
}

func Search(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []Language, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := languageCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = languageCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []Language{}
		total = 0
	}
	return results, total
}
