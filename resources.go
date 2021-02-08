package api

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
)

type Resources struct {
	SharedKey []byte
	Logger    *log.Logger
	DB        *pg.DB
	Redis     *redis.Client
}

func NewResources(ctx context.Context, c *Config) (*Resources, error) {
	db, err := newDB(ctx, c)
	if err != nil {
		return nil, err
	}

	conn := db.Conn()
	defer conn.Close()

	if err = createSchema(conn); err != nil {
		return nil, err
	}

	redis := newRedisClient()

	logger := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)

	return &Resources{
		[]byte(c.SharedKey),
		logger,
		db,
		redis}, nil
}

func (r *Resources) Close() error {
	if err := r.DB.Close(); err != nil {
		return err
	}

	if err := r.Redis.Close(); err != nil {
		return err
	}

	return nil
}
