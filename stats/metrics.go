package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vitwit/matic-telemetry/config"
	"github.com/vitwit/matic-telemetry/types"
)

// GetStatus will returns the the node status
func GetStatus(cfg *config.Config) (*types.NodeStatusResponse, error) {
	var status types.NodeStatusResponse
	url := cfg.Endpoints.HeimdallRPCEndpoint + "/status"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return nil, err
		}
		err = json.Unmarshal(body, &status)
		if err != nil {
			return nil, err
		}
	}

	return &status, nil
}

// GetStatus will returns the latest block info
func GetLatestBlock(cfg *config.Config) (*types.BlockResponse, error) {
	var result types.BlockResponse
	url := cfg.Endpoints.HeimdallRPCEndpoint + "/block"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return nil, err
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// GetNetInfo will returns the network information and error if any
func GetNetInfo(cfg *config.Config) (*types.NetInfoResponse, error) {
	var info types.NetInfoResponse
	url := cfg.Endpoints.HeimdallRPCEndpoint + "/net_info"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return nil, err
		}
		err = json.Unmarshal(body, &info)
		if err != nil {
			log.Printf("Error while unmarshelling net info")
			return nil, err
		}
	}
	return &info, nil
}

// SyncStatus will returns the node syncing status and error if any
func SyncStatus(cfg *config.Config) (*types.SyncingResponse, error) {
	var sync types.SyncingResponse
	url := cfg.Endpoints.HeimdallLCDEndpoint + "/cosmos/base/tendermint/v1beta1/syncing"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting sync info: %v", err)
		return nil, err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return nil, err
		}
		err = json.Unmarshal(body, &sync)
		if err != nil {
			log.Printf("Error while unmarshelling sync res : %v", err)
			return nil, err
		}
	}
	return &sync, nil
}

// GetHeimdallVersion will returns the software version of heimdall
func GetHeimdallVersion(cfg *config.Config) (string, error) {
	var nodeInfo types.NodeInfoResponse
	url := cfg.Endpoints.HeimdallLCDEndpoint + "/cosmos/base/tendermint/v1beta1/node_info"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting heimdall version: %v", err)
		return "", err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return "", err
		}
		err = json.Unmarshal(body, &nodeInfo)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return "", err
		}
	}
	log.Printf("Heimdall Verison : %s", nodeInfo.ApplicationVersion.Version)

	return nodeInfo.ApplicationVersion.Version, nil
}

// GetVersion will returns the software version of heimdall
func GetVersion(cfg *config.Config) (string, error) {
	var versionInfo types.VersionResponse
	url := cfg.Endpoints.HeimdallLCDEndpoint + "/version"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting heimdall version: %v", err)
		return "", err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return "", err
		}
		err = json.Unmarshal(body, &versionInfo)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return "", err
		}
	}
	return versionInfo.Version, nil
}

// GetHeimdallVersion will returns the software version of heimdall
func GetTransactions(cfg *config.Config, height int) (int, error) {
	var block types.BlockResponse
	url := cfg.Endpoints.HeimdallRPCEndpoint + fmt.Sprintf("/block?height=%d", height)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting block info: %v", err)
		return 0, err
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while getting block info : %v", err)
			return 0, err
		}
		err = json.Unmarshal(body, &block)
		if err != nil {
			log.Printf("Error while getting block info : %v", err)
			return 0, err
		}
	}

	return len(block.Result.Block.Data.Txs), nil
}
