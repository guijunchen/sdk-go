package utils

const (
	// System Block Contract keys
	KeyBlockContractWithRWSet   = "withRWSet"
	KeyBlockContractBlockHash   = "blockHash"
	KeyBlockContractBlockHeight = "blockHeight"
	KeyBlockContractTxId        = "txId"

	// System Chain Config Contract keys
	KeyChainConfigContractRoot        = "root"
	KeyChainConfigContractOrgId       = "org_id"
	KeyChainConfigContractNodeId      = "node_id"
	KeyChainConfigContractNewNodeId   = "new_node_id"
	KeyChainConfigContractNodeIds     = "node_ids"
	KeyChainConfigContractBlockHeight = "block_height"

	// CoreConfig keys
	KeyTxSchedulerTimeout         = "tx_scheduler_timeout"
	KeyTxSchedulerValidateTimeout = "tx_scheduler_validate_timeout"

	// BlockConfig keys
	KeyTxTimeOut       = "tx_timeout"
	KeyBlockTxCapacity = "block_tx_capacity"
	KeyBlockSize       = "block_size"
	KeyBlockInterval   = "block_interval"

	// CertManage keys
	KeyCertHashes = "cert_hashes"
	KeyCerts      = "certs"
	KeyCertCrl    = "cert_crl"

	// ArchiveConfig consts
	MysqlDBNamePrefix     = "cm_archived_chain"
	MysqlTableNamePrefix  = "t_block_info"
	RowsPerBlockInfoTable = 100000
)
