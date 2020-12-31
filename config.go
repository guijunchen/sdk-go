/**
 * @Author: jasonruan
 * @Date:   2020-12-29 11:05:48
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config *ChainClientConfigModel

type nodesConfigModel struct {
	// 节点地址
	NodeAddr                string              `mapstructure:"node_addr"`
	// 节点连接数
	ConnCnt                 int                 `mapstructure:"conn_cnt"`
	// RPC连接是否启用双向TLS认证
	EnableTLS               bool                `mapstructure:"enable_tls"`
	// 信任证书池路径
	TrustRootPaths          []string            `mapstructure:"trust_root_paths"`
	// TLS hostname
	TLSHostName             string              `mapstructure:"tls_host_name"`
}

type chainClientConfigModel struct {
	// 链ID
	ChainId                 string              `mapstructure:"chain_id"`
	// 组织ID
	OrgId                   string              `mapstructure:"org_id"`
	// 客户端用户私钥路径
	UserKeyFilePath         string              `mapstructure:"user_key_file_path"`
	// 客户端用户证书路径
	UserCrtFilePath         string              `mapstructure:"user_crt_file_path"`
	// 客户端管理员私钥路径
	AdminKeyFilePath        string              `mapstructure:"admin_key_file_path"`
	// 客户端管理员证书路径
	AdminCrtFilePath        string              `mapstructure:"admin_crt_file_path"`
	// 节点配置
	NodesConfig             []nodesConfigModel  `mapstructure:"nodes"`
}

type ChainClientConfigModel struct {
	ChainClientConfig       chainClientConfigModel           `mapstructure:"chain_client"`
}

func InitConfig(confPath string) error {
	var (
		err error
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
