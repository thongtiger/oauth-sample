package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thongtiger/oauth-rfc6749/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PaymentStatusNew      = "new"
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusCanceled = "canceled"

	CollectOnlineBanking = "onlineBanking"
	CollectPayment       = "payment"
	CollectUsers         = "users"
	TypeDeposit          = "deposit" // +
	TypePayout           = "payout"  // -
)

var (
	// ErrDatabase error db type
	ErrDatabase = errors.New("error from database")

	// ErrInvalidConnect error connecting
	ErrInvalidConnect = errors.New("error invalid connect")

	// ErrDatabasePing error ping host
	ErrDatabasePing = errors.New("error ping")

	// ErrInvalidAgrs error of invalid agrument
	ErrInvalidAgrs = errors.New("error invalid agrument")

	// ErrDuplicated error duplicate record on database
	ErrDuplicated = errors.New("error duplicated record")

	// ErrInsert error duplicate record on database
	ErrInsert = errors.New("error add new record or duplicate key")

	// ErrDocNotFound error document not found
	ErrDocNotFound = errors.New("error document not found")
)

// Store interface
type Store interface {
	NewUser(username, password, role string, scope []string) (*auth.User, error)
	ValidateUser(username, password string) (bool, auth.User)
	GetUser(username string) (result *auth.User, err error)
	FindUser() (results []*auth.User)
}

// NewMongoStore : context
type mongoContext struct {
	hostname, port, username, password, database string
}

// NewMongoStore : Constructor function
func NewMongoStore(hostname, port, username, password, database string) Store {
	// return interface
	return &mongoContext{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}

func (c *mongoContext) newClient() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", c.hostname, c.port)).SetAuth(options.Credential{
		Username:      c.username,
		Password:      c.password,
		AuthSource:    c.database,
		AuthMechanism: "SCRAM-SHA-1",
	}).SetMaxPoolSize(10).SetMaxConnIdleTime(time.Minute * 5)

	// Connect to MongoDB
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // timeout 10 second
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, ErrInvalidConnect
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, ErrDatabasePing
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}
