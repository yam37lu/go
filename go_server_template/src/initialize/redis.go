package initialize

import (
	"template/global"
	"strings"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func Redis() error {
	redisCfg := global.SYS_CONFIG.Redis
	if redisCfg.Mode == 1 {
		addrs := strings.Split(redisCfg.Addr, ",")
		client := redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: addrs,
			MasterName:    redisCfg.Name,
			Password:      redisCfg.Password, // no password set
			DB:            redisCfg.DB,       // use default DB
		})
		pong, err := client.Ping().Result()
		if err != nil {
			global.SYS_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
			return err
		} else {
			global.SYS_LOG.Info("redis connect ping response:", zap.String("pong", pong))
			global.SYS_REDIS = client
			return nil
		}
	} else if redisCfg.Mode == 2 {
		addrs := strings.Split(redisCfg.Addr, ",")
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: redisCfg.Password, // no password set
		})
		pong, err := client.Ping().Result()
		if err != nil {
			global.SYS_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
			return err
		} else {
			global.SYS_LOG.Info("redis connect ping response:", zap.String("pong", pong))
			global.SYS_RedisCluster = client
			return nil
		}
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password, // no password set
			DB:       redisCfg.DB,       // use default DB
		})
		pong, err := client.Ping().Result()
		if err != nil {
			global.SYS_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
			return err
		} else {
			global.SYS_LOG.Info("redis connect ping response:", zap.String("pong", pong))
			global.SYS_REDIS = client
			return nil
		}
	}
}
