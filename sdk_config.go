/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-go/common/crypto/asym"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"chainmaker.org/chainmaker-go/common/log"
	"go.uber.org/zap"
)

const (
	// 单ChainMaker节点最大连接数
	MaxConnCnt = 1024
	// 查询交易超时时间
	GetTxTimeout = 10
	// 发送交易超时时间
	SendTxTimeout = 10
	// 默认grpc客户端接受最大值 4M
	DefaultRpcClientMaxReceiveMessageSize = 4
)

// 节点配置
type NodeConfig struct {
	// 必填项
	// 节点地址
	addr string
	// 节点连接数
	connCnt int
	// 选填项
	// 是否启用TLS认证
	useTLS bool
	// CA ROOT证书路径
	caPaths []string
	// CA ROOT证书内容（同时配置caPaths和caCerts以caCerts为准）
	caCerts []string
	// TLS hostname
	tlsHostName string
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

// 添加CA证书内容
func WithNodeCACerts(caCerts []string) NodeOption {
	return func(config *NodeConfig) {
		config.caCerts = caCerts
	}
}

func WithNodeTLSHostName(tlsHostName string) NodeOption {
	return func(config *NodeConfig) {
		config.tlsHostName = tlsHostName
	}
}

// Archive配置
type ArchiveConfig struct {
	// 非必填
	// secret key
	secretKey string
}

type ArchiveOption func(config *ArchiveConfig)

// 设置Archive的secret key
func WithSecretKey(key string) ArchiveOption {
	return func(config *ArchiveConfig) {
		config.secretKey = key
	}
}

// RPC Client 链接配置
type RPCClientConfig struct {

	//pc客户端最大接受大小 (MB)
	rpcClientMaxReceiveMessageSize int
}

type RPCClientOption func(config *RPCClientConfig)

// 设置RPC Client的Max Receive Message Size
func WithRPCClientMaxReceiveMessageSize(size int) RPCClientOption {
	return func(config *RPCClientConfig) {
		config.rpcClientMaxReceiveMessageSize = size
	}
}

type ChainClientConfig struct {
	// logger若不设置，将采用默认日志文件输出日志，建议设置，以便采用集成系统的统一日志输出
	logger Logger

	// 链客户端相关配置
	// 方式1：配置文件指定（方式1与方式2可以同时使用，参数指定的值会覆盖配置文件中的配置）
	confPath string

	// 方式2：参数指定（方式1与方式2可以同时使用，参数指定的值会覆盖配置文件中的配置）
	orgId    string
	chainId  string
	nodeList []*NodeConfig

	// 以下xxxPath和xxxBytes同时指定的话，优先使用Bytes
	userKeyFilePath     string
	userCrtFilePath     string
	userSignKeyFilePath string
	userSignCrtFilePath string

	userKeyBytes     []byte
	userCrtBytes     []byte
	userSignKeyBytes []byte
	userSignCrtBytes []byte

	// 以下字段为经过处理后的参数
	privateKey crypto.PrivateKey
	userCrt    *bcx509.Certificate

	// 归档特性的配置
	archiveConfig *ArchiveConfig

	// rpc客户端设置
	rpcClientConfig *RPCClientConfig
}

type ChainClientOption func(*ChainClientConfig)

// 设置配置文件路径
func WithConfPath(confPath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.confPath = confPath
	}
}

// 添加ChainMaker节点地址及连接数配置
func AddChainClientNodeConfig(nodeConfig *NodeConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.nodeList = append(config.nodeList, nodeConfig)
	}
}

// 添加用户私钥文件路径配置
func WithUserKeyFilePath(userKeyFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userKeyFilePath = userKeyFilePath
	}
}

// 添加用户证书文件路径配置
func WithUserCrtFilePath(userCrtFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userCrtFilePath = userCrtFilePath
	}
}

// 添加用户签名私钥文件路径配置
func WithUserSignKeyFilePath(userSignKeyFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignKeyFilePath = userSignKeyFilePath
	}
}

// 添加用户签名证书文件路径配置
func WithUserSingCrtFilePath(userSignCrtFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignCrtFilePath = userSignCrtFilePath
	}
}

// 添加用户私钥文件内容配置
func WithUserKeyBytes(userKeyBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userKeyBytes = userKeyBytes
	}
}

// 添加用户证书文件内容配置
func WithUserCrtBytes(userCrtBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userCrtBytes = userCrtBytes
	}
}

// 添加用户签名私钥文件内容配置
func WithUserSignKeyBytes(userSignKeyBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignKeyBytes = userSignKeyBytes
	}
}

// 添加用户签名证书文件内容配置
func WithUserSignCrtBytes(userSignCrtBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignCrtBytes = userSignCrtBytes
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

// 设置Archive配置
func WithArchiveConfig(conf *ArchiveConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.archiveConfig = conf
	}
}

//设置grpc客户端配置
func WithRPCClientConfig(conf *RPCClientConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.rpcClientConfig = conf
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

func setChainConfig(config *ChainClientConfig) {
	if Config.ChainClientConfig.ChainId != "" && config.chainId == "" {
		config.chainId = Config.ChainClientConfig.ChainId
	}

	if Config.ChainClientConfig.OrgId != "" && config.orgId == "" {
		config.orgId = Config.ChainClientConfig.OrgId
	}
}

// 如果参数没有设置，便使用配置文件的配置
func setUserConfig(config *ChainClientConfig) {
	if Config.ChainClientConfig.UserKeyFilePath != "" && config.userKeyFilePath == "" && config.userKeyBytes == nil {
		config.userKeyFilePath = Config.ChainClientConfig.UserKeyFilePath
	}

	if Config.ChainClientConfig.UserCrtFilePath != "" && config.userCrtFilePath == "" && config.userCrtBytes == nil {
		config.userCrtFilePath = Config.ChainClientConfig.UserCrtFilePath
	}

	if Config.ChainClientConfig.UserSignKeyFilePath != "" && config.userSignKeyFilePath == "" && config.userSignKeyBytes == nil {
		config.userSignKeyFilePath = Config.ChainClientConfig.UserSignKeyFilePath
	}

	if Config.ChainClientConfig.UserSignCrtFilePath != "" && config.userSignCrtFilePath == "" && config.userSignCrtBytes == nil {
		config.userSignCrtFilePath = Config.ChainClientConfig.UserSignCrtFilePath
	}
}

func setNodeList(config *ChainClientConfig) {
	if len(Config.ChainClientConfig.NodesConfig) > 0 && len(config.nodeList) == 0 {
		for _, conf := range Config.ChainClientConfig.NodesConfig {
			node := NewNodeConfig(
				// 节点地址，格式：127.0.0.1:12301
				WithNodeAddr(conf.NodeAddr),
				// 节点连接数
				WithNodeConnCnt(conf.ConnCnt),
				// 节点是否启用TLS认证
				WithNodeUseTLS(conf.EnableTLS),
				// 根证书路径，支持多个
				WithNodeCAPaths(conf.TrustRootPaths),
				// TLS Hostname
				WithNodeTLSHostName(conf.TLSHostName),
			)

			config.nodeList = append(config.nodeList, node)
		}
	}
}

func setArchiveConfig(config *ChainClientConfig) {
	if Config.ChainClientConfig.ArchiveConfig != nil && config.archiveConfig == nil {
		archive := NewArchiveConfig(
			// secret key
			WithSecretKey(Config.ChainClientConfig.ArchiveConfig.SecretKey),
		)

		config.archiveConfig = archive
	}
}

func setRPCClientConfig(config *ChainClientConfig) {
	if Config.ChainClientConfig.RPCClientConfig != nil && config.rpcClientConfig == nil {
		rpcClient := NewRPCClientConfig(
			WithRPCClientMaxReceiveMessageSize(Config.ChainClientConfig.RPCClientConfig.MaxRecvMsgSize),
		)
		config.rpcClientConfig = rpcClient
	}
}

func readConfigFile(config *ChainClientConfig) error {
	// 若没有配置配置文件
	if config.confPath == "" {
		return nil
	}

	if err := InitConfig(config.confPath); err != nil {
		return fmt.Errorf("init config failed, %s", err.Error())
	}

	setChainConfig(config)

	setUserConfig(config)

	setNodeList(config)

	setArchiveConfig(config)

	setRPCClientConfig(config)

	return nil
}

// SDK配置校验
func checkConfig(config *ChainClientConfig) error {

	var (
		err error
	)

	if err = readConfigFile(config); err != nil {
		return fmt.Errorf("read sdk config file failed, %s", err.Error())
	}

	// 如果logger未指定，使用默认zap logger
	if config.logger == nil {
		config.logger = getDefaultLogger()
	}

	if err = checkNodeListConfig(config); err != nil {
		return err
	}

	if err = checkUserConfig(config); err != nil {
		return err
	}

	if err = checkChainConfig(config); err != nil {
		return err
	}

	if err = checkArchiveConfig(config); err != nil {
		return err
	}

	if err = checkRPCClientConfig(config); err != nil {
		return err
	}
	return nil
}

func checkNodeListConfig(config *ChainClientConfig) error {
	// 连接的节点地址不可为空
	if len(config.nodeList) == 0 {
		return fmt.Errorf("connect chainmaker node address is empty")
	}

	// 已配置的节点地址连接数，需要在合理区间
	for _, node := range config.nodeList {
		if node.connCnt <= 0 || node.connCnt > MaxConnCnt {
			return fmt.Errorf("node connection count should >0 && <=%d",
				MaxConnCnt)
		}

		if node.useTLS {
			// 如果开启了TLS认证，CA路径必填
			if len(node.caPaths) == 0 && len(node.caCerts) == 0 {
				return fmt.Errorf("if node useTLS is open, should set caPaths or caCerts")
			}

			// 如果开启了TLS认证，需配置TLS HostName
			if node.tlsHostName == "" {
				return fmt.Errorf("if node useTLS is open, should set tls hostname")
			}
		}
	}

	return nil
}

func checkUserConfig(config *ChainClientConfig) error {
	// 用户私钥不可为空
	if config.userKeyFilePath == "" && config.userKeyBytes == nil {
		return fmt.Errorf("user key cannot be empty")
	}

	// 用户证书不可为空
	if config.userCrtFilePath == "" && config.userCrtBytes == nil {
		return fmt.Errorf("user crt cannot be empty")
	}

	return nil
}

func checkChainConfig(config *ChainClientConfig) error {
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

func checkArchiveConfig(config *ChainClientConfig) error {
	return nil
}

func checkRPCClientConfig(config *ChainClientConfig) error {
	if config.rpcClientConfig == nil {
		rpcClient := NewRPCClientConfig(
			WithRPCClientMaxReceiveMessageSize(DefaultRpcClientMaxReceiveMessageSize),
		)
		config.rpcClientConfig = rpcClient
	} else {
		if config.rpcClientConfig.rpcClientMaxReceiveMessageSize <= 0 || config.rpcClientConfig.rpcClientMaxReceiveMessageSize > 100 {
			config.rpcClientConfig.rpcClientMaxReceiveMessageSize = DefaultRpcClientMaxReceiveMessageSize
		}
	}
	return nil
}

func dealConfig(config *ChainClientConfig) error {
	var err error

	if err = dealUserCrtConfig(config); err != nil {
		return err
	}

	if err = dealUserKeyConfig(config); err != nil {
		return err
	}

	if err = dealUserSignCrtConfig(config); err != nil {
		return err
	}

	if err = dealUserSignKeyConfig(config); err != nil {
		return err
	}

	return nil
}

func dealUserCrtConfig(config *ChainClientConfig) (err error) {

	if config.userCrtBytes == nil {
		// 读取用户证书
		config.userCrtBytes, err = ioutil.ReadFile(config.userCrtFilePath)
		if err != nil {
			return fmt.Errorf("read user crt file failed, %s", err.Error())
		}
	}

	// 将证书转换为证书对象
	if config.userCrt, err = ParseCert(config.userCrtBytes); err != nil {
		return fmt.Errorf("ParseCert failed, %s", err.Error())
	}

	return nil
}

func dealUserKeyConfig(config *ChainClientConfig) (err error) {

	if config.userKeyBytes == nil {
		// 从私钥文件读取用户私钥，转换为privateKey对象
		config.userKeyBytes, err = ioutil.ReadFile(config.userKeyFilePath)
		if err != nil {
			return fmt.Errorf("read user key file failed, %s", err)
		}
	}

	config.privateKey, err = asym.PrivateKeyFromPEM(config.userKeyBytes, nil)
	if err != nil {
		return fmt.Errorf("parse user key file to privateKey obj failed, %s", err)
	}

	return nil
}

func dealUserSignCrtConfig(config *ChainClientConfig) (err error) {

	if config.userSignCrtBytes == nil {
		if config.userSignCrtFilePath == "" {
			config.userSignCrtBytes = config.userCrtBytes
			return nil
		}

		config.userSignCrtBytes, err = ioutil.ReadFile(config.userSignCrtFilePath)
		if err != nil {
			return fmt.Errorf("read user sign crt file failed, %s", err.Error())
		}

	}

	if config.userCrt, err = ParseCert(config.userSignCrtBytes); err != nil {
		return fmt.Errorf("ParseSignCert failed, %s", err.Error())
	}

	return nil
}

func dealUserSignKeyConfig(config *ChainClientConfig) (err error) {

	if config.userSignKeyBytes == nil {
		if config.userSignKeyFilePath == "" {
			config.userSignKeyBytes = config.userKeyBytes
			return nil
		}

		config.userSignKeyBytes, err = ioutil.ReadFile(config.userSignKeyFilePath)
		if err != nil {
			return fmt.Errorf("read user sign key file failed, %s", err.Error())
		}
	}

	config.privateKey, err = asym.PrivateKeyFromPEM(config.userSignKeyBytes, nil)
	if err != nil {
		return fmt.Errorf("parse user key file to privateKey obj failed, %s", err)
	}

	return nil
}

func getDefaultLogger() *zap.SugaredLogger {
	config := log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.LEVEL_DEBUG,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	}

	logger, _ := log.InitSugarLogger(&config)
	return logger
}
