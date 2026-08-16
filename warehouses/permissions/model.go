package permissions

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
var permissionCollection *libDb.CollectionClass
var permissionTypeCollection *libDb.CollectionClass
var accountTypeCollection *libDb.CollectionClass
var permissionCache *libCache.Cache

func InitModel(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache) {
	dbMongo = inDb
	permissionCache = inRedisClient
	permissionCollection = dbMongo.NewCollection("permissions")
	permissionCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_type_id_1_order_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "type_id", Value: 1}, {Key: "order", Value: 1}},
		},
	})
	permissionTypeCollection = dbMongo.NewCollection("permission_types")
	permissionTypeCollection.UpdateIndex([]libDb.CreateIndex{
		{
			Name: "delete_1",
			Key:  bSon.D{{Key: "delete", Value: 1}},
		},
		{
			Name: "delete_1_order_1",
			Key:  bSon.D{{Key: "delete", Value: 1}, {Key: "order", Value: 1}},
		},
	})
	accountTypeCollection = dbMongo.NewCollection("account_types")
	accountTypeCollection.UpdateIndex([]libDb.CreateIndex{
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

func CreatePermission(insertData Permission) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := permissionCollection.Create(insertData)
	if result == true {
		permissionCache.Del("permission_by_type:" + insertData.PermissionTypeID)
		output = insertData.ID
	}
	return output
}

func UpdatePermission(inId string, inUpdateData bSon.M, inType string) bool {
	output := permissionCollection.UpdateId(inId, inUpdateData)
	if output == true {
		permissionCache.Dels("permission:"+inId, "permission_by_type:"+inId)
	}
	return output
}

func GetPermissionById(inId string, isCache bool) (output Permission) {
	if isCache == true {
		cacheData := permissionCache.Get("permission:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(Permission)
			return output
		}
	}
	err := permissionCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		permissionCache.Set("permission:"+inId, output, 0)
	}
	return output
}

func GetPermissionByCode(inCode string) (output Permission) {
	err := permissionCollection.FindOne(bSon.M{"code": inCode, "delete": 0}).Decode(&output)
	if err != nil {
		output.ID = ""
	}
	return output
}

func CheckPermissionCode(inCode string, inId string) (output Permission) {
	filter := bSon.M{"code": inCode, "delete": 0}
	if inId != "" {
		filter["_id"] = bSon.M{"$ne": inId}
	}
	err := permissionCollection.FindOne(filter).Decode(&output)
	if err != nil {
		output.ID = ""
	}
	return output
}

func GetPermissionsByType(inType string) (results []MiniPermission) {
	cacheData := permissionCache.Get("permission_by_type:" + inType)
	if cacheData != nil {
		output, _ := cacheData.([]MiniPermission)
		return output
	}

	cursor, err := permissionTypeCollection.Find(bSon.M{"delete": 0, "type": inType}, bSon.D{{Key: "delete", Value: 1}, {Key: "type_id", Value: 1}, {Key: "order", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			results = []MiniPermission{}
		} else {
			permissionCache.Set("permission_by_type:"+inType, results, 0)
		}
	}
	return results
}

func GetPermissions() []MiniPermissionType {
	types := GetTypes()
	n := len(types)
	if n == 0 {
		return types
	}

	// Channel đóng vai trò Semaphore giới hạn tối đa 8 goroutines chạy song song
	maxConcurrency := 4
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup

	for i := range types {
		wg.Add(1)
		sem <- struct{}{} // Bơm 1 slot, nếu đủ 8 slot sẽ block đợi slot trống

		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }() // Giải phóng 1 slot khi hoàn thành

			// Gán trực tiếp vào index của slice gốc (Go slice safe cho ghi vào index khác nhau)
			types[idx].Permissions = GetPermissionsByType(types[idx].ID)
		}(i)
	}

	wg.Wait()
	return types
}

func Search(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []Permission, total int64) {
	var wg sync.WaitGroup
	wg.Add(2) // 2 tác vụ song song

	var findErr error
	// Tác vụ 1: Lấy danh sách kết quả
	go func() {
		defer wg.Done()
		cursor, err := permissionCollection.Find(filter, inSortOrder, inPage, inLimit)
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
		total = permissionCollection.Count(filter)
	}()

	// Chờ cả 2 goroutine hoàn thành
	wg.Wait()
	if findErr != nil {
		results = []Permission{}
		total = 0
	}
	return results, total
}
