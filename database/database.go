package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

type MongoDB struct {
	Host        string
	Port        int
	User        string
	Pass        string
	DBName      string
	AuthSource  string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
}

type Option func(*MongoDB)

func WithHost(host string) Option {
	return func(db *MongoDB) {
		db.Host = host
	}
}

func WithPort(port int) Option {
	return func(db *MongoDB) {
		db.Port = port
	}
}

func WithUser(user string) Option {
	return func(db *MongoDB) {
		db.User = user
	}
}

func WithPass(pass string) Option {
	return func(db *MongoDB) {
		db.Pass = pass
	}
}

func WithDBName(dbName string) Option {
	return func(db *MongoDB) {
		db.DBName = dbName
	}
}

func WithMaxPoolSize(maxPoolSize uint64) Option {
	return func(db *MongoDB) {
		db.MaxPoolSize = maxPoolSize
	}
}

func WithMinPoolSize(minPoolSize uint64) Option {
	return func(db *MongoDB) {
		db.MinPoolSize = minPoolSize
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(db *MongoDB) {
		db.Timeout = timeout
	}
}

func WithAuthSource(authSource string) Option {
	return func(db *MongoDB) {
		db.AuthSource = authSource
	}
}

func (m *MongoDB) URI() string {
	if m.User != "" && m.Pass != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			m.User, m.Pass, m.Host, m.Port, m.DBName, m.AuthSource)
	}
	return fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
}

func (m *MongoDB) Connect() (*mongo.Client, *mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(m.URI())

	if m.MaxPoolSize > 0 {
		clientOpts.SetMaxPoolSize(m.MaxPoolSize)
	}

	if m.MinPoolSize > 0 {
		clientOpts.SetMinPoolSize(m.MinPoolSize)
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)

	defer cancel()

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, client.Database(m.DBName), nil
}

func (m *MongoDB) Disconnect(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}

func NewMongoDB(opts ...Option) *MongoDB {
	m := &MongoDB{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
