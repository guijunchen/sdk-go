package chainmaker_sdk_go

const (
	keyWithRWSet   = "with_rwset"
	keyBlockHash   = "block_hash"
	keyBlockHeight = "block_height"
	keyTxId        = "tx_id"

	keyOrgId     = "org_id"
	keyNodeId    = "node_id"
	keyNewNodeId = "new_node_id"
	keyNodeIds   = "node_ids"

	// archive consts
	mysqlDBNamePrefix     = "cm_archived_chain"
	mysqlTableNamePrefix  = "t_block_info"
	rowsPerBlockInfoTable = 100000
)
