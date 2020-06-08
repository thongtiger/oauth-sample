package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thongtiger/oauth-rfc6749/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *mongoContext) NewUser(username, password, role string, scope []string) (*auth.User, error) {
	client, err := c.newClient()
	defer client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	insertData := auth.User{
		ID:             primitive.NewObjectID(),
		Role:           role,
		Scope:          scope,
		Username:       username,
		Password:       password,
		CreateTime:     time.Now().UTC(),
		LatestLoggedin: time.Now().UTC(),
	}

	insertData.HashingPassword() // hashing

	collection := client.Database(c.database).Collection(CollectUsers)
	result, err := collection.InsertOne(context.TODO(), insertData)
	if err != nil {
		fmt.Println(ErrInsert)
		return nil, err
	}
	fmt.Printf("inserted ID: %v\n", result.InsertedID)

	// m := OnlineBankingMethod(insertData)
	return &insertData, nil
}
func (c *mongoContext) ValidateUser(username, password string) (ok bool, result auth.User) {
	client, err := c.newClient()
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	collection := client.Database(c.database).Collection(CollectUsers)
	// single find
	filter := bson.M{"username": username}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	if err = collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return
	}
	if result.VerifyPassword(password) {
		ok = true
		return
	}
	return
}
func (c *mongoContext) GetUser(username string) (result *auth.User, err error) {
	client, err := c.newClient()
	defer client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err)
		return
	}
	// single find
	filter := bson.M{"username": username}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(c.database).Collection(CollectUsers)
	if err = collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return
	}
	return
}

func (c *mongoContext) FindUser() (results []*auth.User) {
	client, err := c.newClient()
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	collection := client.Database(c.database).Collection(CollectUsers)

	findOptions := options.Find()
	//findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem auth.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	//fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return
}
