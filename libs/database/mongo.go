package database

import (
	"context"
	"encoding/json"
	"log"
	"time"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type databaseClass struct {
	database *mongo.Database
}

func NewDatabase(mongoClient *mongo.Client, dbName string) *databaseClass {
	myDatabase := mongoClient.Database(dbName)
	return &databaseClass{database: myDatabase}
}

type collectionClass struct {
	collection *mongo.Collection
}

type CollectionBase struct {
	ID        string
	Delete    bool      `bson:"delete" json:"delete"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
type LapTrinhVien struct {
	CollectionBase   // Anonymous field (Struct Embedding)
	NgonNguChuyenMon string
}

func (db *databaseClass) NewCollection(inCollection string) *collectionClass {
	myCollection := db.database.Collection(inCollection)
	return &collectionClass{collection: myCollection}
}

func (col *collectionClass) CreateIndex(indexModel mongo.IndexModel) {
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

func (col *collectionClass) DeleteIndex(inName string) {
	libProcess.Try(func() {
		col.collection.Indexes().DropOne(context.TODO(), inName)
		log.Println("Remove index: ", inName)
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
}

func (col *collectionClass) GetListIndex(nameIndex any) (err error) {
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

func (col *collectionClass) Create(insertData any) (output string) {
	libProcess.Try(func() {
		createData, _ := structToBsonM(insertData)
		id := primitive.NewObjectID().Hex()
		createData["_id"] = id
		curDate := time.Now()
		createData["created_at"] = curDate
		createData["updated_at"] = curDate
		_, err := col.collection.InsertOne(context.TODO(), createData)
		if err != nil {
			log.Println(err)
		} else {
			output = id
		}
	}).Catch(func(e libProcess.E) {
		log.Println(e)
	})
	return output
}

func (col *collectionClass) UpdateId(inId string, inUpdateData bSon.M) bool {
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

func (col *collectionClass) DeleteId(inId string) bool {
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

func (col *collectionClass) RemoveId(inId string) bool {
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
