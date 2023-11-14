package repositories

import (
	"context"
	"time"
)

func DefaultDbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
