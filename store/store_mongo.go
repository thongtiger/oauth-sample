package store

import (
	"context"
	"fmt"
	"time"

	"github.com/thongtiger/oauth-rfc6749/auth"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *mongoContext) NewUser(username, password, role, displayName string, scope []string) (*auth.User, error) {
	client, err := c.newClient()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(c.database).Collection(CollectUsers)

	insertData := auth.User{
		Role:           role,
		Scope:          scope,
		Username:       username,
		Password:       password,
		Name:           displayName,
		CreateTime:     time.Time{}.UTC(),
		LatestLoggedin: time.Time{}.UTC(),
	}

	insertData.HashingPassword() // hashing

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
	ok = false
	client, err := c.newClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(context.TODO())
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
