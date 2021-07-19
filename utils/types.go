package utils

const (
	KeyWithRWSet   = "with_rwset"
	KeyBlockHash   = "block_hash"
	KeyBlockHeight = "block_height"
	KeyTxId        = "tx_id"

	KeyOrgId     = "org_id"
	KeyNodeId    = "node_id"
	KeyNewNodeId = "new_node_id"
	KeyNodeIds   = "node_ids"

	// archive consts
	MysqlDBNamePrefix     = "cm_archived_chain"
	MysqlTableNamePrefix  = "t_block_info"
	RowsPerBlockInfoTable = 100000
)
