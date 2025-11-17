package graph

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"platform-go-challenge/models"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

type contextKey string

const audienceLoaderKey = contextKey("audienceLoader")

// AudienceLoader wraps the dataloader
type AudienceLoader struct {
	loader *dataloader.Loader[string, *models.Audience]
}

// NewAudienceLoader creates a new audience dataloader
func NewAudienceLoader(db *gorm.DB) *AudienceLoader {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*models.Audience] {
		// Convert string keys to uint IDs
		ids := make([]uint, len(keys))
		keyMap := make(map[uint]int) // maps ID to position in keys
		for i, key := range keys {
			id, err := strconv.ParseUint(key, 10, 32)
			if err != nil {
				continue
			}
			ids[i] = uint(id)
			keyMap[uint(id)] = i
		}

		// Fetch all audiences in a single query
		var audiences []*models.Audience
		if err := db.Find(&audiences, ids).Error; err != nil {
			// Return errors for all keys
			results := make([]*dataloader.Result[*models.Audience], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*models.Audience]{Error: err}
			}
			return results
		}

		// Map audiences back to their keys
		audienceMap := make(map[uint]*models.Audience)
		for _, audience := range audiences {
			audienceMap[audience.ID] = audience
		}

		// Build results in the same order as keys
		results := make([]*dataloader.Result[*models.Audience], len(keys))
		for i, id := range ids {
			if audience, ok := audienceMap[id]; ok {
				results[i] = &dataloader.Result[*models.Audience]{Data: audience}
			} else {
				results[i] = &dataloader.Result[*models.Audience]{
					Error: fmt.Errorf("audience not found"),
				}
			}
		}

		return results
	}

	loader := dataloader.NewBatchedLoader(
		batchFn,
		dataloader.WithCache[string, *models.Audience](&dataloader.NoCache[string, *models.Audience]{}),
		dataloader.WithBatchCapacity[string, *models.Audience](100),
		dataloader.WithWait[string, *models.Audience](1*time.Millisecond),
	)

	return &AudienceLoader{loader: loader}
}

// Load a single audience by ID
func (l *AudienceLoader) Load(ctx context.Context, id string) (*models.Audience, error) {
	return l.loader.Load(ctx, id)()
}

// LoadMany audiences by IDs
func (l *AudienceLoader) LoadMany(ctx context.Context, ids []string) ([]*models.Audience, []error) {
	thunk := l.loader.LoadMany(ctx, ids)
	return thunk()
}

// DataLoaderMiddleware injects the dataloader into the context
func DataLoaderMiddleware(db *gorm.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		audienceLoader := NewAudienceLoader(db)
		return context.WithValue(ctx, audienceLoaderKey, audienceLoader)
	}
}

// GetAudienceLoader retrieves the audience loader from context
func GetAudienceLoader(ctx context.Context) *AudienceLoader {
	return ctx.Value(audienceLoaderKey).(*AudienceLoader)
}