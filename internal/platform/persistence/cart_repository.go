package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/log"
	"github.com/redis/go-redis/v9"
)

type CartRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCartRepository(client *redis.Client, ttl time.Duration) *CartRepository {
	return &CartRepository{client: client, ttl: ttl}
}

func (r *CartRepository) FindByUserID(ctx context.Context, userID string) (*shop.Cart, error) {
	bytes, err := r.client.Get(ctx, userID).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Debugf(ctx, "cart not found for user %s", userID)
			return nil, &ErrEntityNotFound{entity: "cart"}
		}
		return nil, fmt.Errorf("(FindByUserID) failed retrieving cart from redis: %w", err)
	}

	cart := shop.Cart{}
	if err = json.Unmarshal(bytes, &cart); err != nil {
		return nil, fmt.Errorf("(FindByUserID) failed unmarshalling cart: %w", err)
	}

	return &cart, nil
}

func (r *CartRepository) Save(ctx context.Context, cart *shop.Cart) error {
	bytes, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("(Save) failed marshalling cart: %w", err)
	}

	err = r.client.Set(ctx, cart.UserID, bytes, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("(Save) failed saving cart to redis: %w", err)
	}

	return nil
}

func (r *CartRepository) DeleteByUserID(ctx context.Context, userID string) error {
	err := r.client.Del(ctx, userID).Err()
	if err != nil {
		return fmt.Errorf("(DeleteByUserID) failed deleting cart from redis: %w", err)
	}

	return nil
}
