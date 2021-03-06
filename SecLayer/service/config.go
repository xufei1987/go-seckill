package service

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	etcd "go.etcd.io/etcd/client/v3"
)

var secLayerContext = &SecLayerContext{}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
	RedisQueueName   string
}

type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecLayerConf struct {
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf
	EtcdConfig       EtcdConf
	LogPath          string
	LogLevel         string

	WriteGoroutineNum       int
	ReadGoroutineNum        int
	HandleUserGoroutineNum  int
	Read2HandleChanSize     int
	Handle2WriteChanSize    int
	MaxRequestWaitTimeout   int64
	SendToWriteChanTimeout  int64
	SendToHandleChanTimeout int64

	SecProductInfoMap map[int]*SecProductInfoConf
	TokenPasswd       string
}

type SecLayerContext struct {
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	etcdClient           *etcd.Client
	RWSecProductLock     sync.RWMutex

	secLayerConf     *SecLayerConf
	waitGroup        sync.WaitGroup
	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResponse

	HistoryMap     map[int]*UserBuyHistory
	HistoryMapLock sync.Mutex

	// 商品的计数
	productCountMgr *ProductCountMgr
}

type SecProductInfoConf struct {
	ProductId         int
	StartTime         int64
	EndTime           int64
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	// 每秒最多能卖多少个
	soldMaxLimit int
	// 限速控制
	secLimit *SecLimit
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientReferer string
	// CloseNotify   <-chan bool
}

type SecResponse struct {
	ProductId int
	UserId    int
	Token     string
	TokenTime int64
	Code      int
}
