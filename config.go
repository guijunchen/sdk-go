/**
 * @Author: jasonruan
 * @Date:   2020-11-30 14:41:07
 */
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-go/common/crypto/asym"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"chainmaker.org/chainmaker-go/common/log"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
)

const (
	// 单ChainMaker节点最大连接数
	MaxConnCnt      = 5
	// 查询交易超时时间
	GetTxTimeout    = 10
	// 发送交易超时时间
	SendTxTimeout   = 10
)

// 节点配置
type NodeConfig struct {
	// 必填项
	// 节点地址
	addr                string
	// 节点连接数
	connCnt             int
	// 选填项
	// 是否启用TLS认证
	useTLS              bool
	// CA ROOT证书路径
	caPaths             []string
	// TLS hostname
	tlsHostName         string
}

type NodeOption func(config *NodeConfig)

// 设置节点地址
func WithNodeAddr(addr string) NodeOption {
	return func(config *NodeConfig) {
		config.addr = addr
	}
}

// 设置节点连接数
func WithNodeConnCnt(connCnt int) NodeOption {
	return func(config *NodeConfig) {
		config.connCnt = connCnt
	}
}

// 设置是否启动TLS开关
func WithNodeUseTLS(useTLS bool) NodeOption {
	return func(config *NodeConfig) {
		config.useTLS = useTLS
	}
}

// 添加CA证书路径
func WithNodeCAPaths(caPaths []string) NodeOption {
	return func(config *NodeConfig) {
		config.caPaths = caPaths
	}
}

func WithNodeTLSHostName(tlsHostName string) NodeOption {
	return func(config *NodeConfig) {
		config.tlsHostName = tlsHostName
	}
}

type UserConfig struct {
	userKeyFilePath     string
	userCrtFilePath     string
	// 以下字段为经过处理后的参数
	privateKey          crypto.PrivateKey
	userCrtPEM          []byte
	userCrt             *bcx509.Certificate
}

type UserOption func(config *UserConfig)

// 添加用户私钥文件路径配置
func WithUserKeyFilePath(userKeyFilePath string) UserOption {
	return func(config *UserConfig) {
		config.userKeyFilePath = userKeyFilePath
	}
}

// 添加用户证书文件路径配置
func WithUserCrtFilePath(userCrtFilePath string) UserOption {
	return func(config *UserConfig) {
		config.userCrtFilePath = userCrtFilePath
	}
}

type ChainClientConfig struct {
	orgId               string
	chainId             string
	nodeList            []*NodeConfig
	userConfig          *UserConfig
	// logger若不设置，将采用默认日志文件输出日志，建议设置，以便采用集成系统的统一日志输出
	logger              Logger
}

type ChainClientOption func(*ChainClientConfig)

// 添加ChainMaker节点地址及连接数配置
func AddChainClientNodeConfig(nodeConfig *NodeConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.nodeList = append(config.nodeList, nodeConfig)
	}
}

// 添加User配置
func WithChainClientUserConfig(userConfig *UserConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userConfig = userConfig
	}
}

// 添加OrgId
func WithChainClientOrgId(orgId string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.orgId = orgId
	}
}

// 添加ChainId
func WithChainClientChainId(chainId string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.chainId = chainId
	}
}

// 设置Logger对象，便于日志打印
func WithChainClientLogger(logger Logger) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.logger = logger
	}
}

// 生成SDK配置并校验合法性
func generateConfig(opts ...ChainClientOption) (*ChainClientConfig, error) {
	config := &ChainClientConfig{}
	for _, opt := range opts {
		opt(config)
	}

	// 校验config参数合法性
	if err := checkConfig(config); err != nil {
		return nil, err
	}

	// 进一步处理config参数
	if err := dealConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SDK配置校验
func checkConfig(config *ChainClientConfig) error {

	// 如果logger未指定，使用默认zap logger
	if config.logger == nil {
		config.logger = getDefaultLogger()
	}

	// 连接的节点地址不可为空
	if len(config.nodeList) == 0 {
		return fmt.Errorf("connect chianmaker node address is empty")
	}

	// 已配置的节点地址连接数，需要在合理区间
	for _, node := range config.nodeList {
		if node.connCnt <= 0 || node.connCnt > MaxConnCnt {
			return fmt.Errorf("node connection count should >0 && <=%d",
				MaxConnCnt)
		}

		if node.useTLS {
			// 如果开启了TLS认证，CA路径必填
			if len(node.caPaths) == 0 {
				return fmt.Errorf("if node useTLS is open, should set caPath")
			}

			// 如果开启了TLS认证，需配置TLS HostName
			if node.tlsHostName == "" {
				return fmt.Errorf("if node useTLS is open, should set tls hostname")
			}
		}
	}


	// 用户私钥不可为空
	if config.userConfig.userKeyFilePath == "" {
		return fmt.Errorf("user key file path cannot be empty")
	}

	// 用户证书不可为空
	if config.userConfig.userCrtFilePath == "" {
		return fmt.Errorf("user crt file path cannot be empty")
	}

	// OrgId不可为空
	if config.orgId == "" {
		return fmt.Errorf("orgId cannot be empty")
	}

	// ChainId不可为空
	if config.chainId == "" {
		return fmt.Errorf("chainId cannot be empty")
	}

	return nil
}

func dealConfig(config *ChainClientConfig) error {
	var err error

	// 读取用户证书
	config.userConfig.userCrtPEM, err = ioutil.ReadFile(config.userConfig.userCrtFilePath)
	if err != nil {
		return fmt.Errorf("read user crt file failed, %s", err.Error())
	}

	// 从私钥文件读取用户私钥，转换为privateKey对象
	skBytes, err := ioutil.ReadFile(config.userConfig.userKeyFilePath)
	if err != nil {
		return fmt.Errorf("read user key file failed, %s", err)
	}
	config.userConfig.privateKey, err = asym.PrivateKeyFromPEM(skBytes, nil)
	if err != nil {
		return fmt.Errorf("parse user key file to privateKey obj failed, %s", err)
	}

	// 将证书转换为证书对象
	if config.userConfig.userCrt, err = ParseCert(config.userConfig.userCrtPEM); err != nil {
		return fmt.Errorf("ParseCert failed, %s", err.Error())
	}

	return nil
}

func getDefaultLogger() *zap.SugaredLogger {
	config := log.LogConfig{
		Module:   "[SDK]",
		LogPath:  "./sdk.log",
		LogLevel: log.LEVEL_DEBUG,
		MaxAge: 30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	}

	logger, _ := log.InitSugarLogger(&config)
	return logger
}