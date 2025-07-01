package types

type NodeInfoResponse struct {
	DefaultNodeInfo    DefaultNodeInfo    `json:"default_node_info"`
	ApplicationVersion ApplicationVersion `json:"application_version"`
}

type DefaultNodeInfo struct {
	ProtocolVersion struct {
		P2P   string `json:"p2p"`
		Block string `json:"block"`
		App   string `json:"app"`
	} `json:"protocol_version"`
	DefaultNodeID string `json:"default_node_id"`
	ListenAddr    string `json:"listen_addr"`
	Network       string `json:"network"`
	Version       string `json:"version"`
	Channels      string `json:"channels"`
	Moniker       string `json:"moniker"`
	Other         struct {
		TxIndex    string `json:"tx_index"`
		RPCAddress string `json:"rpc_address"`
	} `json:"other"`
}

type ApplicationVersion struct {
	Name             string `json:"name"`
	AppName          string `json:"app_name"`
	Version          string `json:"version"`
	GitCommit        string `json:"git_commit"`
	BuildTags        string `json:"build_tags"`
	GoVersion        string `json:"go_version"`
	CosmosSDKVersion string `json:"cosmos_sdk_version"`
}

type SyncingResponse struct {
	Syncing bool `json:"syncing"`
}
