package files

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
var fileCollection *libDb.CollectionClass
var fileCache *libCache.Cache

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache) {
	dbMongo = inDb
	fileCache = inRedisClient
	fileCollection = dbMongo.NewCollection("files")
	//Todo: update file index
	fileCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_-1",
			Key:  bSon.D{{Key: "delete", Value: -1}},
		},
		{
			Name: "delete_-1_created_at_-1",
			Key:  bSon.D{{Key: "delete", Value: -1}, {Key: "created_at", Value: -1}},
		},
		{
			Name: "delete_-1_s3_bucket_1",
			Key:  bSon.D{{Key: "delete", Value: -1}, {Key: "s3_bucket", Value: 1}},
		},
	})
}
func Create(insertData File) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := fileCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}
func Update(inId string, inUpdateData bSon.M) bool {
	output := fileCollection.UpdateId(inId, inUpdateData)
	if output == true {
		fileCache.Del("language:" + inId)
	}
	return output
}
func Delete(inId string) bool {
	output := fileCollection.DeleteId(inId)
	if output == true {
		fileCache.Del("language:" + inId)
	}
	return output
}

func GetById(inId string, isCache bool) (output File) {
	if isCache == true {
		cacheData := fileCache.Get("language:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(File)
			return output
		}
	}
	err := fileCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		fileCache.Set("language:"+inId, output, 0)
	}
	return output
}
func GetByIds(inIds []string, isCache bool) []File {
	if len(inIds) == 0 {
		return []File{}
	}
	output := make([]File, len(inIds))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup

	for i, id := range inIds {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, fileID string) {
			defer wg.Done()
			defer func() { <-sem }()
			output[idx] = GetById(fileID, isCache)
		}(i, id)
	}
	wg.Wait()
	return output
}
func Search(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []File, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := fileCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = fileCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []File{}
		total = 0
	}
	return results, total
}
