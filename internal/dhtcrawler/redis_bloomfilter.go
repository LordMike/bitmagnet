package dhtcrawler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/redis/go-redis/v9"
)

const redisBloomFilterKey = "bf:blobs"

type redisBloomFilter struct {
	client *redis.Client
}

func newRedisBloomFilter(redisURL string) (*redisBloomFilter, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &redisBloomFilter{client: redis.NewClient(opt)}, nil
}

func (r *redisBloomFilter) Exists(ctx context.Context, infoHash protocol.ID) (bool, error) {
	key := infoHash[:]
	raw, err := r.client.Do(ctx, "BF.EXISTS", redisBloomFilterKey, key).Result()
	if err != nil {
		return false, err
	}
	return parseRedisBoolish(raw)
}

func (r *redisBloomFilter) Add(ctx context.Context, infoHash protocol.ID) error {
	key := infoHash[:]
	raw, err := r.client.Do(ctx, "BF.ADD", redisBloomFilterKey, key).Result()
	if err != nil {
		return err
	}
	_, err = parseRedisBoolish(raw)
	return err
}

func (r *redisBloomFilter) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}

func parseRedisBoolish(v any) (bool, error) {
	switch t := v.(type) {
	case bool:
		return t, nil
	case int64:
		switch t {
		case 0:
			return false, nil
		case 1:
			return true, nil
		default:
			return false, fmt.Errorf("unexpected integer bool: %d", t)
		}
	case string:
		switch t {
		case "0":
			return false, nil
		case "1":
			return true, nil
		}
		parsed, err := strconv.ParseBool(t)
		if err != nil {
			return false, fmt.Errorf("unexpected string bool: %q", t)
		}
		return parsed, nil
	case []byte:
		return parseRedisBoolish(string(t))
	default:
		return false, fmt.Errorf("unexpected type=%T for bool", v)
	}
}
