package organizations

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func CreateMember(insertData Member) (output string) {
	output = ""
	insertData.ID = primitive.NewObjectID().Hex()
	curDate := time.Now()
	insertData.CreatedAt = curDate
	insertData.UpdatedAt = curDate
	result := memberCollection.Create(insertData)
	if result == true {
		output = insertData.ID
	}
	return output
}

func UpdateMember(inId string, inUpdateData bSon.M, inType string) bool {
	output := memberCollection.UpdateId(inId, inUpdateData)
	if output == true {
		organizationCache.Del("member:" + inId)
	}
	return output
}
func DeleteMember(inId string) bool {
	output := memberCollection.DeleteId(inId)
	if output == true {
		organizationCache.Del("member:" + inId)
	}
	return output
}

func GetMemberById(inId string, isCache bool) (output Member) {
	if isCache == true {
		cacheData := organizationCache.Get("member:" + inId)
		if cacheData != nil {
			output, _ = cacheData.(Member)
			return output
		}
	}
	err := memberCollection.FindById(inId).Decode(&output)
	if err != nil || output.Delete != 0 {
		output.ID = ""
	} else {
		organizationCache.Set("member:"+inId, output, 0)
	}
	return output
}
func GetMembers(inOrganizationId string) (results []Member) {
	cursor, err := memberCollection.Find(bSon.M{"delete": 0, "organization_id": inOrganizationId}, bSon.D{{Key: "delete", Value: 1}, {Key: "organization_id", Value: 1}, {Key: "host", Value: -1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}
func GetOrganizationsByCustomerId(inCustomerId string) (results []Member) {
	cursor, err := memberCollection.Find(bSon.M{"delete": 0, "customer_id": inCustomerId}, bSon.D{{Key: "delete", Value: 1}, {Key: "customer_id", Value: 1}}, 0, 0)
	if err != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
	}
	return results
}
func SearchMembers(filter bSon.M, inSortOrder bSon.D, inPage int64, inLimit int64) (results []Member, total int64) {
	var wg sync.WaitGroup
	wg.Add(2)
	var findErr error
	go func() {
		defer wg.Done()
		cursor, err := memberCollection.Find(filter, inSortOrder, inPage, inLimit)
		if err != nil {
			findErr = err
			return
		}
		if err := cursor.All(context.TODO(), &results); err != nil {
			findErr = err
			return
		}
	}()
	go func() {
		defer wg.Done()
		total = memberCollection.Count(filter)
	}()
	wg.Wait()
	if findErr != nil {
		results = []Member{}
		total = 0
	}
	return results, total
}
