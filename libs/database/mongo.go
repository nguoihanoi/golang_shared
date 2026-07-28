package database

import (
	"context"
	"encoding/json"
	"log"
	"time"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	"go.mongodb.org/mongo-driver/v2/bson"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MongoConnect(mongoURI string) *mongo.Client {
	log.Println("Connecting to MongoDB...")
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	/*
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	*/
	return mongoClient
}

type DatabaseClass struct {
	database *mongo.Database
}

func NewDatabase(mongoClient *mongo.Client, dbName string) *DatabaseClass {
	myDatabase := mongoClient.Database(dbName)
	return &DatabaseClass{database: myDatabase}
}

type CollectionClass struct {
	collection *mongo.Collection
}

func (db *DatabaseClass) NewCollection(inCollection string) *CollectionClass {
	myCollection := db.database.Collection(inCollection)
	return &CollectionClass{collection: myCollection}
}

func (col *CollectionClass) CreateIndex(indexModel mongo.IndexModel) {
	libProcess.Try(func() {
		name, err := col.collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			log.Println(err)
		}
		log.Println("Create index: ", name)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
}

func (col *CollectionClass) DeleteIndex(inName string) {
	libProcess.Try(func() {
		col.collection.Indexes().DropOne(context.TODO(), inName)
		log.Println("Remove index: ", inName)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
}

func (col *CollectionClass) GetListIndex(nameIndex any) (err error) {
	var cursor *mongo.Cursor
	libProcess.Try(func() {
		cursor, err = col.collection.Indexes().List(context.TODO())
		if err == nil {
			var results []bSon.D
			if err = cursor.All(context.TODO(), &results); err == nil {
				var indexResult []byte
				indexResult, err = json.Marshal(results)
				if err == nil {
					err = json.Unmarshal([]byte(string(indexResult)), &nameIndex)
				}
			}
		} else {
			return
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return err
}

func structToBsonM(v any) (bSon.M, error) {
	// Step 1: Marshal struct thành bytes BSON
	data, err := bSon.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Step 2: Unmarshal bytes BSON ngược lại vào bson.M
	var doc bSon.M
	err = bSon.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (col *CollectionClass) Create(insertData any) (output bool) {
	output = false
	libProcess.Try(func() {
		_, err := col.collection.InsertOne(context.TODO(), insertData)
		if err != nil {
			log.Println(err)
		} else {
			output = true
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return output
}

func (col *CollectionClass) UpdateId(inId string, inUpdateData bSon.M) bool {
	output := false
	libProcess.Try(func() {
		filter := bSon.M{"_id": inId}
		inUpdateData["updated_at"] = time.Now()
		update := bSon.M{"$set": inUpdateData}
		// Updates the first document that has the specified "_id" value
		result, err := col.collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
		if result.MatchedCount > 0 {
			output = true
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return output
}

func (col *CollectionClass) DeleteId(inId string) bool {
	output := false
	libProcess.Try(func() {
		filter := bSon.M{"_id": inId}
		inUpdateData := bSon.M{
			"delete":     1,
			"updated_at": time.Now(),
		}
		update := bSon.M{"$set": inUpdateData}
		// Updates the first document that has the specified "_id" value
		result, err := col.collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
		if result.MatchedCount > 0 {
			output = true
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return output
}

func (col *CollectionClass) RemoveId(inId string) bool {
	output := false
	libProcess.Try(func() {
		filter := bSon.M{"_id": inId}
		// Updates the first document that has the specified "_id" value
		result, err := col.collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Println(err)
		}
		if result.DeletedCount > 0 {
			output = true
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return output
}

func (col *CollectionClass) FindById(inId string) (output any) {
	var err error
	status := true
	libProcess.Try(func() {
		filter := bSon.M{"_id": inId}
		err = col.collection.FindOne(context.TODO(), filter).Decode(output)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
		status = false
	})
	if err != nil || status == false {
		return nil
	}
	return output
}
func (col *CollectionClass) FindOne(filter bSon.M) any {
	var err error
	var result any
	status := true
	libProcess.Try(func() {
		err = col.collection.FindOne(context.TODO(), filter).Decode(&result)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
		status = false
	})
	if err != nil || status == false {
		return nil
	}
	return result
}

func (col *CollectionClass) Find(filter bSon.M, inSortOrder bSon.M, inPage int64, inLimit int64) (output []any) {
	var err error
	var cursor *mongo.Cursor
	status := true
	libProcess.Try(func() {
		opts := options.Find().SetSort(inSortOrder)
		if inPage > 0 {
			opts = opts.SetSkip((inPage - 1) * inLimit)
		}
		if inLimit > 0 {
			opts = opts.SetLimit(inLimit)
		}
		cursor, err = col.collection.Find(context.TODO(), filter, opts)
		if err == nil {
			err = cursor.All(context.TODO(), &output)
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
		status = false
	})
	if err != nil || status == false {
		return nil
	}
	return output
}

func (col *CollectionClass) Count(filter bSon.M) (output int64) {
	var err error
	status := true
	libProcess.Try(func() {
		output, err = col.collection.CountDocuments(context.TODO(), filter)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
		status = false
	})
	if err != nil || status == false {
		output = 0
	}
	return output
}

func (col *CollectionClass) Pipe(matchFilter bSon.D, groupFilter bSon.D, inSortOrder bSon.D) (output []any) {
	var err error
	var cursor *mongo.Cursor
	status := true
	libProcess.Try(func() {
		matchStage := bSon.D{{Key: "$match", Value: matchFilter}}
		groupStage := bson.D{{Key: "$group", Value: groupFilter}}
		sortStage := bson.D{{Key: "$sort", Value: inSortOrder}}
		cursor, err = col.collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage, sortStage})
		if err == nil {
			err = cursor.All(context.TODO(), &output)
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
		status = false
	})
	if err != nil || status == false {
		return nil
	}
	return output
}
