package organizations

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
var organizationCollection *libDb.CollectionClass
var memberCollection *libDb.CollectionClass
var organizationCache *libCache.Cache

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache) {
	dbMongo = inDb
	organizationCache = inRedisClient
	organizationCollection = dbMongo.NewCollection("organizations")
	organizationCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_status_-1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "status", Value: -1}},
		},
	})
	memberCollection = dbMongo.NewCollection("members")
	memberCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_customer_id_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "customer_id", Value: 1}},
		},
		{
			Name: "delete_1_organization_id_1_host_-1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "organization_id", Value: 1}, {Key: "host", Value: -1}},
		},
	})
}

func Create(insertData Organization) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := organizationCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func Update(inId string, inUpdateData bSon.M, inType string) bool {
	output := organizationCollection.UpdateId(inId, inUpdateData)
	if output == true {
		organizationCache.Del("organization:" + inId)
	}
	return output
}
func Delete(inId string) bool {
	output := organizationCollection.DeleteId(inId)
	if output == true {
		organizationCache.Del("organization:" + inId)
	}
	return output
}

func GetById(inId string, isCache bool) (output Organization) {
	if isCache == true {
		cacheData := organizationCache.Get("organization:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(Organization)
			return output
		}
	}
	err := organizationCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		organizationCache.Set("organization:"+inId, output, 0)
	}
	return output
}

func Search(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []Organization, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := organizationCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = organizationCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []Organization{}
		total = 0
	}
	return results, total
}
