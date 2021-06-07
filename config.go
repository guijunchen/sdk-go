/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config *ChainClientConfigModel

type nodesConfigModel struct {
	// 节点地址
	NodeAddr string `mapstructure:"node_addr"`
	// 节点连接数
	ConnCnt int `mapstructure:"conn_cnt"`
	// RPC连接是否启用双向TLS认证
	EnableTLS bool `mapstructure:"enable_tls"`
	// 信任证书池路径
	TrustRootPaths []string `mapstructure:"trust_root_paths"`
	// TLS hostname
	TLSHostName string `mapstructure:"tls_host_name"`
}

type archiveConfigModel struct {
	// 链外存储类型，已支持：mysql
	Type 		string `mapstructure:"type"`
	// 链外存储目标，格式：
	// 	- mysql: user:pwd:host:port
	Dest 		string `mapstructure:"dest"`
	SecretKey	string `mapstructure:"secret_key"`
}

type chainClientConfigModel struct {
	// 链ID
	ChainId string `mapstructure:"chain_id"`
	// 组织ID
	OrgId string `mapstructure:"org_id"`
	// 客户端用户私钥路径
	UserKeyFilePath string `mapstructure:"user_key_file_path"`
	// 客户端用户证书路径
	UserCrtFilePath string `mapstructure:"user_crt_file_path"`
	// 客户端用户交易签名私钥路径(若未设置，将使用user_key_file_path)
	UserSignKeyFilePath string `mapstructure:"user_sign_key_file_path"`
	// 客户端用户交易签名证书路径(若未设置，将使用user_crt_file_path)
	UserSignCrtFilePath string `mapstructure:"user_sign_crt_file_path"`
	// 节点配置
	NodesConfig []nodesConfigModel `mapstructure:"nodes"`
	// 归档特性的配置
	ArchiveConfig *archiveConfigModel `mapstructure:"archive,omitempty"`
}

type ChainClientConfigModel struct {
	ChainClientConfig chainClientConfigModel `mapstructure:"chain_client"`
}

func InitConfig(confPath string) error {
	var (
		err       error
		confViper *viper.Viper
	)

	if confViper, err = initViper(confPath); err != nil {
		return fmt.Errorf("Load sdk config failed, %s", err)
	}

	Config = &ChainClientConfigModel{}
	if err = confViper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("Unmarshal config file failed, %s", err)
	}

	return nil
}

func initViper(confPath string) (*viper.Viper, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(confPath)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}

	return cmViper, nil
}
