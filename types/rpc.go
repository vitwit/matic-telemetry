package types

type NodeStatusResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string `json:"id"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
			Channels   string `json:"channels"`
			Moniker    string `json:"moniker"`
			Other      struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string `json:"latest_block_hash"`
			LatestAppHash       string `json:"latest_app_hash"`
			LatestBlockHeight   string `json:"latest_block_height"`
			LatestBlockTime     string `json:"latest_block_time"`
			EarliestBlockHash   string `json:"earliest_block_hash"`
			EarliestAppHash     string `json:"earliest_app_hash"`
			EarliestBlockHeight string `json:"earliest_block_height"`
			EarliestBlockTime   string `json:"earliest_block_time"`
			CatchingUp          bool   `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

// type BlockResponse struct {
// 	JSONRPC string `json:"jsonrpc"`
// 	ID      int    `json:"id"`
// 	Result  struct {
// 		BlockID struct {
// 			Hash  string `json:"hash"`
// 			Parts struct {
// 				Total int    `json:"total"`
// 				Hash  string `json:"hash"`
// 			} `json:"parts"`
// 		} `json:"block_id"`
// 		Block struct {
// 			Header struct {
// 				Version struct {
// 					Block string `json:"block"`
// 				} `json:"version"`
// 				ChainID     string `json:"chain_id"`
// 				Height      string `json:"height"`
// 				Time        string `json:"time"`
// 				LastBlockID struct {
// 					Hash  string `json:"hash"`
// 					Parts struct {
// 						Total int    `json:"total"`
// 						Hash  string `json:"hash"`
// 					} `json:"parts"`
// 				} `json:"last_block_id"`
// 				LastCommitHash     string `json:"last_commit_hash"`
// 				DataHash           string `json:"data_hash"`
// 				ValidatorsHash     string `json:"validators_hash"`
// 				NextValidatorsHash string `json:"next_validators_hash"`
// 				ConsensusHash      string `json:"consensus_hash"`
// 				AppHash            string `json:"app_hash"`
// 				LastResultsHash    string `json:"last_results_hash"`
// 				EvidenceHash       string `json:"evidence_hash"`
// 				ProposerAddress    string `json:"proposer_address"`
// 			} `json:"header"`
// 			Data struct {
// 				Txs []string `json:"txs"`
// 			} `json:"data"`
// 			Evidence struct {
// 				Evidence []interface{} `json:"evidence"`
// 			} `json:"evidence"`
// 			LastCommit struct {
// 				Height  string `json:"height"`
// 				Round   int    `json:"round"`
// 				BlockID struct {
// 					Hash  string `json:"hash"`
// 					Parts struct {
// 						Total int    `json:"total"`
// 						Hash  string `json:"hash"`
// 					} `json:"parts"`
// 				} `json:"block_id"`
// 				Signatures []struct {
// 					BlockIDFlag      int    `json:"block_id_flag"`
// 					ValidatorAddress string `json:"validator_address"`
// 					Timestamp        string `json:"timestamp"`
// 					Signature        string `json:"signature"`
// 				} `json:"signatures"`
// 			} `json:"last_commit"`
// 		} `json:"block"`
// 	} `json:"result"`
// }

type NetInfoResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Listening bool     `json:"listening"`
		Listeners []string `json:"listeners"`
		NPeers    string   `json:"n_peers"`
		Peers     []Peer   `json:"peers"`
	} `json:"result"`
}

type Peer struct {
	NodeInfo         NodeInfo         `json:"node_info"`
	IsOutbound       bool             `json:"is_outbound"`
	ConnectionStatus ConnectionStatus `json:"connection_status"`
	RemoteIP         string           `json:"remote_ip"`
}

type NodeInfo struct {
	ProtocolVersion struct {
		P2P   string `json:"p2p"`
		Block string `json:"block"`
		App   string `json:"app"`
	} `json:"protocol_version"`
	ID         string `json:"id"`
	ListenAddr string `json:"listen_addr"`
	Network    string `json:"network"`
	Version    string `json:"version"`
	Channels   string `json:"channels"`
	Moniker    string `json:"moniker"`
	Other      struct {
		TxIndex    string `json:"tx_index"`
		RPCAddress string `json:"rpc_address"`
	} `json:"other"`
}

type ConnectionStatus struct {
	Duration    string        `json:"Duration"`
	SendMonitor MonitorStatus `json:"SendMonitor"`
	RecvMonitor MonitorStatus `json:"RecvMonitor"`
	Channels    []Channel     `json:"Channels"`
}

type MonitorStatus struct {
	Start    string `json:"Start"`
	Bytes    string `json:"Bytes"`
	Samples  string `json:"Samples"`
	InstRate string `json:"InstRate"`
	CurRate  string `json:"CurRate"`
	AvgRate  string `json:"AvgRate"`
	PeakRate string `json:"PeakRate"`
	BytesRem string `json:"BytesRem"`
	Duration string `json:"Duration"`
	Idle     string `json:"Idle"`
	TimeRem  string `json:"TimeRem"`
	Progress int    `json:"Progress"`
	Active   bool   `json:"Active"`
}

type Channel struct {
	ID                int    `json:"ID"`
	SendQueueCapacity string `json:"SendQueueCapacity"`
	SendQueueSize     string `json:"SendQueueSize"`
	Priority          string `json:"Priority"`
	RecentlySent      string `json:"RecentlySent"`
}

type BlockResponse struct {
	Result struct {
		Block struct {
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	} `json:"result"`
}
