package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var ErrNotCached = errors.New("channel not cached")

// Returns nil if we cannot create a channel
// Returns ErrNotCached if not cached
func GetDMChannel(userId uint64) (*uint64, error) {
	key := fmt.Sprintf("dmchannel:%d", userId)

	res, err := Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotCached
		}

		return nil, err
	}

	if res == "null" {
		return nil, nil
	}

	parsed, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func StoreNullDMChannel(userId uint64) error {
	key := fmt.Sprintf("dmchannel:%d", userId)
	return Client.Set(key, "null", time.Hour * 24).Err()
}

func StoreDMChannel(userId, channelId uint64) error {
	key := fmt.Sprintf("dmchannel:%d", userId)
	return Client.Set(key, strconv.FormatUint(channelId, 10), 0).Err()
}
