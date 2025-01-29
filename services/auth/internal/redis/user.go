package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thejixer/jixifood/shared/models"
)

func (rc *RedisStore) CacheUser(u *models.UserEntity) error {
	key := fmt.Sprintf("jf-u-%v", u.ID)

	output, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if err := rc.rdb.Set(rc.ctx, key, string(output), time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetUser(id uint64) *models.UserEntity {
	key := fmt.Sprintf("jf-u-%v", id)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return nil
	}

	user := new(models.UserEntity)

	err = json.Unmarshal([]byte(val), &user)

	if err != nil {
		return nil
	}

	return user
}

func (rc *RedisStore) DelUser(id int) error {
	key := fmt.Sprintf("jf-u-%v", id)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
