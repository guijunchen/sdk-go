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

	// ArchiveConfig consts
	MysqlDBNamePrefix     = "cm_archived_chain"
	MysqlTableNamePrefix  = "t_block_info"
	RowsPerBlockInfoTable = 100000

	// CertManage keys
	KeyCertHashes = "cert_hashes"
	KeyCerts      = "certs"
	KeyCertCrl    = "cert_crl"

	// PrivateCompute keys
	KeyOrderId      = "order_id"
	KeyPrivateDir   = "private_dir"
	KeyContractName = "contract_name"
	KeyCodeHash     = "code_hash"
	KeyResult       = "result"
	KeyCodeHeader   = "code_header"
	KeyVersion      = "version"
	KeyIsDeploy     = "is_deploy"
	KeyRWSet        = "rw_set"
	KeyEvents       = "events"
	KeyReportHash   = "report_hash"
	KeySign         = "sign"
	KeyKey          = "key"
	KeyPayload      = "payload"
	KeyOrgIds       = "org_ids"
	KeySignPairs    = "sign_pairs"
	KeyCaCert       = "ca_cert"
	KeyEnclaveId    = "enclave_id"
	KeyReport       = "report"
	KeyProof        = "proof"
	KeyDeployReq    = "deploy_req"
	KeyPrivateReq   = "private_req"
)
