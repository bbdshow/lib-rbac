package mongo

import (
	"context"
	"github.com/bbdshow/bkit/db/mongo"
	"time"
)

type Conf struct {
	URI      string
	Database string
}

func NewDatabase(c *Conf) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return mongo.NewDatabase(ctx, c.URI, c.Database)
}
