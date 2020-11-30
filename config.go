/**
 * @Author: jasonruan
 * @Date:   2020-11-30 14:41:07
 */
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/common/log"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	// SDK单ChainMaker节点最大连接数
	MaxConnCnt      = 5
	// 查询交易超时时间
	GetTxTimeout    = 10 * time.Second
	// 发送交易超时时间
	SendTxTimeout   = 10 * time.Second
)

type Config struct {
	addrsWithConnCnt    map[string](int)
	logger              Logger
	useTLS              bool
	caPaths             []string
	userKeyFilePath     string
	userCrtFilePath     string
}

type Option func(*Config)

// 添加ChainMaker节点地址及连接数配置
func AddNodeAddrWithConnCnt(nodeAddr string, connCnt int) Option {
	return func(config *Config) {
		if config.addrsWithConnCnt == nil {
			config.addrsWithConnCnt = make(map[string](int))
		}
		config.addrsWithConnCnt[nodeAddr] = connCnt
	}
}

// 设置Logger对象，便于日志打印
func WithLogger(logger Logger) Option {
	return func(config *Config) {
		config.logger = logger
	}
}

// 设置是否启动TLS开关
func WithUseTLS(useTLS bool) Option {
	return func(config *Config) {
		config.useTLS = useTLS
	}
}

// 添加CA证书路径
func WithCAPaths(caPaths []string) Option {
	return func(config *Config) {
		config.caPaths = caPaths
	}
}

// 添加用户私钥文件路径配置
func WithUserKeyFilePath(userKeyFilePath string) Option {
	return func(config *Config) {
		config.userKeyFilePath = userKeyFilePath
	}
}

// 添加用户证书文件路径配置
func WithUserCrtFilePath(userCrtFilePath string) Option {
	return func(config *Config) {
		config.userCrtFilePath = userCrtFilePath
	}
}

// 生成SDK配置并校验合法性
func generateConfig(opts ...Option) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	if err := checkConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SDK配置校验
func checkConfig(config *Config) error {
	// 如果logger未指定，使用默认zap logger
	if config.logger == nil {
		config.logger = getDefaultLogger()
	}

	// 连接的节点地址不可为空
	if config.addrsWithConnCnt == nil {
		return fmt.Errorf("connect chianmaker node address is empty")
	}

	// 已配置的节点地址连接数，需要在合理区间
	for _, cnt := range config.addrsWithConnCnt {
		if cnt <= 0 || cnt > MaxConnCnt {
			return fmt.Errorf("node connection count should >0 && <=%d",
				MaxConnCnt)
		}
	}

	if config.useTLS {
		// 如果开启了TLS认证，CA路径必填
		if len(config.caPaths) == 0 {
			return fmt.Errorf("useTLS is open, should set caPath")
		}

		// 如果开启了TLS认证，需配置用户证书
		if config.userCrtFilePath == "" {
			return fmt.Errorf("useTLS is open, should set user crt file path")
		}
	}

	// 用户私钥不可为空
	if config.userKeyFilePath == "" {
		return fmt.Errorf("user key file path cannot be empty")
	}

	return nil
}

func getDefaultLogger() *zap.SugaredLogger {
	config := log.LogConfig{
		Module:   "[DEFAULT]",
		LogPath:  "./default.log",
		LogLevel: log.LEVEL_DEBUG,
		MaxAge: 30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	}

	logger, _ := log.InitSugarLogger(&config)
	return logger
}