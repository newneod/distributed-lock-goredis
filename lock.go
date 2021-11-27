package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

const (
	DistributedLockPrefix  = "LOCK:" // prefix of lock-key
	DistributedLockTimeout = 5       // timeout for lock, unit:second
)

// @param lockName string 锁名（唯一区分不同锁业务）
// @return strUUID string 锁UUID
// @return err error
func Lock(lockName string) (string, error) {
	if conn == nil {
		err := errors.New("Please init redis connection first.")
		log.Fatal(err)
		return "", err
	}

	ctx := context.TODO()
	strLockName := strings.Join([]string{DistributedLockPrefix, lockName}, "")
	strUUID := uuid.NewV4().String()
	timeout := DistributedLockTimeout * time.Second
	iTimeBegin := time.Now().Unix()
	for {
		replySetNx := conn.SetNX(ctx, strLockName, strUUID, timeout)
		if replySetNx.Err() != nil {
			return "", replySetNx.Err()
		}
		if replySetNx.Val() {
			return strUUID, nil
		}

		replyTTL := conn.TTL(ctx, strLockName)
		if replyTTL.Err() != nil {
			return "", replyTTL.Err()
		}
		if replyTTL.Val() == -1 {
			replyExpire := conn.Expire(ctx, strLockName, timeout)
			if replyExpire.Err() != nil {
				return "", replyExpire.Err()
			}
		}

		if time.Now().Unix()-iTimeBegin > DistributedLockTimeout {
			return "", errors.New("Operation timeout, please try again.")
		}
		time.Sleep(time.Microsecond * 10)
	}
}

// @param lockName string 锁名（唯一区分不同锁业务）
// @param strUUID string 锁UUID
// @return err error
func Unlock(lockName string, strUUID string) error {
	ctx := context.TODO()
	strLockName := strings.Join([]string{DistributedLockPrefix, lockName}, "")
	replyGet := conn.Get(ctx, strLockName)
	if replyGet.Err() != nil {
		return replyGet.Err()
	}
	if replyGet.Val() != strUUID {
		return errors.New("Currently someone else has the lock, you cannot unlock it.")
	}

	tx := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, strLockName)
			return nil
		})
		return err
	}
	if err := conn.Watch(ctx, tx, strLockName); err != nil {
		return errors.New("The lock currently belongs to someone else, you cannot unlock it.")
	}
	return nil
}
