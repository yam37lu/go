package global

import (
	"template/config"
	"template/utils/license/client"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	SYS_DB           *gorm.DB
	SYS_CONFIG       config.Server
	SYS_REDIS        *redis.Client
	SYS_RedisCluster *redis.ClusterClient
	SYS_LICENSE      *client.License
	SYS_VP           *viper.Viper
	SYS_LOG          *zap.Logger
)

const (
	SystemName = "template"
	MeterDeg   = 8.983152841195214e-6
)
