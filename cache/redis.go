package cache

import (
	"douyin/config"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

// RedisPool 数据库连接池
var RedisPool *redis.Pool

// RedisPoolInit 初始化数据库redis连接池
func RedisPoolInit() error {
	RedisPool = &redis.Pool{
		MaxIdle:     config.C.Redis.Maxidle,   //最大空闲数
		MaxActive:   config.C.Redis.Maxactive, //最大连接数，0不设上
		Wait:        true,
		IdleTimeout: time.Duration(1) * time.Second, //空闲等待时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.C.Redis.Ipaddress+":"+config.C.Redis.Port) //redis IP地址
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			//密码认证
			if _, err = c.Do("AUTH", config.C.Redis.Authpassword); err != nil {
				err := c.Close()
				if err != nil {
					return nil, err
				}
				return nil, err
			}
			redis.DialDatabase(0)
			return c, err
		},
	}
	return nil
}

// getUserCountKey 关于user计数的key
func getUserCountKey(userid int64) string {
	return fmt.Sprintf("%s_%d", "UserCountKey", userid)
}

func getUserRelationKey(userid int64) string {
	return fmt.Sprintf("%s_%d", "UserRelationKey", userid)
}

func getUserVideoRelation(userid int64) string {
	return fmt.Sprintf("%s_%d", "UserVideoRelationKey", userid)
}

// SetUserRelation 建立用户和用户的关系集合
func SetUserRelation(userid, touserId int64) error {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	key := getUserRelationKey(userid)
	//往集合中加关注的人
	_, err := conn.Do("SADD", key, touserId)
	if err != nil {
		return err
	}
	return nil
}

// SetUserVideoRelation 建立用户和视频的关系集合
func SetUserVideoRelation(userid, videoId int64) error {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	key := getUserVideoRelation(userid)
	//往集合中加关注的人
	_, err := conn.Do("SADD", key, videoId)
	if err != nil {
		return err
	}
	return nil
}

// IsUserRelation 判断是否在集合中
func IsUserRelation(userid, touserId int64) bool {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	key := getUserRelationKey(userid)
	res, err := redis.Int64(conn.Do("SISMEMBER", key, touserId))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if res == 0 {
		fmt.Printf("%#v", res)
		return false
	}
	return true
}

// IsUserVideoRelation 判断是否在集合里面
func IsUserVideoRelation(userid, videoId int64) bool {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	key := getUserVideoRelation(userid)

	res, err := redis.Int64(conn.Do("SISMEMBER", key, videoId))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if res == 0 {
		return false
	}
	return true
}
func DeleteUserRelation(userid, touserId int64) error {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	key := getUserRelationKey(userid)
	//往集合中加关注的人
	_, err := conn.Do("SREM", key, touserId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserVideoRelation 删除关系
func DeleteUserVideoRelation(userid, videoId int64) error {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	key := getUserVideoRelation(userid)
	//往集合中加关注的人
	_, err := conn.Do("SREM", key, videoId)
	if err != nil {
		return err
	}
	return nil
}

// SetUserCount 设置user计数
func SetUserCount(userid int64) error {
	conn := RedisPool.Get()

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	key := getUserCountKey(userid)
	_, err := conn.Do("hmset", redis.Args{key}.AddFlat(map[string]int64{
		"followCount":   0,
		"followerCount": 0,
		"workCount":     0,
		"favoriteCount": 0,
		"totalFavorite": 0,
	})...)
	if err != nil {
		return err
	}
	return nil
}

func GetUserFollowCount(userID int64) (int64, error) {
	res, err := get(userID, "followCount")
	return res, err
}
func GetUserFollowerCount(userID int64) (int64, error) {
	res, err := get(userID, "followerCount")
	return res, err
}
func GetUserWorkCount(userID int64) (int64, error) {
	res, err := get(userID, "workCount")
	return res, err
}
func GetUserFavoriteCount(userID int64) (int64, error) {
	res, err := get(userID, "favoriteCount")
	return res, err
}
func GetUserTotalFavoriteCount(userID int64) (int64, error) {
	res, err := get(userID, "totalFavorite")
	return res, err
}

func DecrByUserFollowCount(userID int64) error {
	err := change(userID, "followCount", -1)
	return err
}
func IncrByUserFollowCount(userID int64) error {
	err := change(userID, "followCount", 1)
	return err
}

func DecrByUserFollowerCount(userID int64) error {
	err := change(userID, "followerCount", -1)
	return err
}
func IncrByUserFollowerCount(userID int64) error {
	err := change(userID, "followerCount", 1)
	return err
}
func DecrByUserWorkCount(userID int64) error {
	err := change(userID, "workCount", -1)
	return err
}
func IncrByUserWorkCount(userID int64) error {
	err := change(userID, "workCount", 1)
	return err
}

func DecrByUserTotalFavorite(userID int64) error {
	err := change(userID, "totalFavorite", -1)
	return err
}
func IncrByUserTotalFavorite(userID int64) error {
	err := change(userID, "totalFavorite", 1)
	return err
}

func DecrByUserFavoriteCount(userID int64) error {
	err := change(userID, "favoriteCount", -1)
	return err
}
func IncrByUserFavoriteCount(userID int64) error {
	err := change(userID, "favoriteCount", 1)
	return err
}

func change(userid int64, field string, incr int64) error {

	key := getUserCountKey(userid)
	conn := RedisPool.Get()

	//释放资源
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	isExists, _ := redis.Int64(conn.Do("exists", key))
	//判断键值是否存在
	if isExists == 0 {
		return errors.New("is not exists")
	}

	before, err := redis.Int64(conn.Do("hget", key, field))
	if err != nil {
		return err
	}
	if before+incr < 0 {
		//此时已经小于0了
		fmt.Println(errors.New("already <0"))
		log.Println(errors.New("already <0"))
		return nil
	}
	_, err = conn.Do("HIncrBy", key, field, incr)
	if err != nil {
		return err
	}
	return nil
}

func get(userid int64, field string) (int64, error) {
	key := getUserCountKey(userid)

	conn := RedisPool.Get()
	//释放资源
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	isExists, _ := redis.Int64(conn.Do("exists", key))

	//判断键值是否存在
	if isExists == 0 {
		return 0, errors.New("is not exists")
	}

	res, err := redis.Int64(conn.Do("hget", key, field))

	if err != nil {
		return 0, err
	}
	return res, nil
}

// UserIsExists 提供给中间件验证的方法
func UserIsExists(userid int64) error {
	conn := RedisPool.Get()
	//释放资源
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	key := getUserCountKey(userid)
	isExists, _ := redis.Int64(conn.Do("exists", key))
	if isExists == 0 {
		return errors.New("is not exists")
	}
	return nil
}
