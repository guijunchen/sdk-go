package utils

const (
	KeyWithRWSet   = "withRWSet"
	KeyBlockHash   = "blockHash"
	KeyBlockHeight = "blockHeight"
	KeyTxId        = "txId"

	KeyOrgId     = "org_id"
	KeyNodeId    = "node_id"
	KeyNewNodeId = "new_node_id"
	KeyNodeIds   = "node_ids"

	// CoreConfig keys
	KeyTxSchedulerTimeout         = "tx_scheduler_timeout"
	KeyTxSchedulerValidateTimeout = "tx_scheduler_validate_timeout"

	// BlockConfig keys
	KeyTxTimeOut       = "tx_timeout"
	KeyBlockTxCapacity = "block_tx_capacity"
	KeyBlockSize       = "block_size"
	KeyBlockInterval   = "block_interval"

	// ArchiveConfig consts
	MysqlDBNamePrefix     = "cm_archived_chain"
	MysqlTableNamePrefix  = "t_block_info"
	RowsPerBlockInfoTable = 100000
)
