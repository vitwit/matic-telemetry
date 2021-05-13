package stats

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vitwit/matic-telemetry/config"
)

// Status is a struct which holds the parameter of status response
type Status struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		NodeInfo interface{} `json:"node_info"`
		Network  string      `json:"network"`
		SyncInfo struct {
			LatestBlockHash   string `json:"latest_block_hash"`
			LatestBlockHeight string `json:"latest_block_height"`
			LatestBlockTime   string `json:"latest_block_time"`
			CatchingUp        bool   `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo interface{} `json:"validator_info"`
	} `json:"result"`
}

// NetInfo is a structre which holds the parameters of net info
type NetInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Listening bool     `json:"listening"`
		Listeners []string `json:"listeners"`
		NPeers    string   `json:"n_peers"`
		Peers     []struct {
			NodeInfo struct {
				ProtocolVersion interface{} `json:"protocol_version"`
				ID              string      `json:"id"`
				ListenAddr      string      `json:"listen_addr"`
				Network         string      `json:"network"`
				Version         string      `json:"version"`
				Moniker         string      `json:"moniker"`
			} `json:"node_info"`
			RemoteIP string `json:"remote_ip"`
		} `json:"peers"`
	} `json:"result"`
}

// Caughtup is a struct which holds the fields of syncing
type Caughtup struct {
	Syncing bool `json:"syncing"`
}

// ApplicationInfo is a struct which holds the details app
type HeimdallVersion struct {
	NodeInfo           interface{} `json:"node_info"`
	ApplicationVersion struct {
		Name       string `json:"name"`
		ServerName string `json:"server_name"`
		ClientName string `json:"client_name"`
		Version    string `json:"version"`
	} `json:"application_version"`
}

type HeimdallStats struct {
	// Success bool `json:"success"`
	Data struct {
		AverageBlockTime float64 `json:"averageBlockTime"`
		TxCount          int     `json:"txCount"`
		TotalSupply      int     `json:"totalSupply"`
	} `json:"data"`
}

// GetLatestBlock will returns the latest block info and error if any
func GetLatestBlock(cfg *config.Config) (Status, error) {
	var block Status
	url := cfg.Endpoints.HeimdallRPCEndpoint + "/status?"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return block, err
	}

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return block, err
		}
		err = json.Unmarshal(body, &block)
		if err != nil {
			return block, err
		}
	}

	return block, nil
}

// GetNetInfo will returns the network information and error if any
func GetNetInfo(cfg *config.Config) (NetInfo, error) {
	var info NetInfo
	url := cfg.Endpoints.HeimdallRPCEndpoint + "/net_info?"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return info, err
	}

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return info, err
		}
		err = json.Unmarshal(body, &info)
		if err != nil {
			log.Printf("Error while unmarshelling net info")
			return info, err
		}
	}
	return info, nil
}

// SyncStatus will returns the node syncing status and error if any
func SyncStatus(cfg *config.Config) (Caughtup, error) {
	var sync Caughtup
	url := cfg.Endpoints.HeimdallLCDEndpoint + "/syncing"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting sync info: %v", err)
		return sync, err
	}

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading resp body : %v", err)
			return sync, err
		}
		err = json.Unmarshal(body, &sync)
		if err != nil {
			log.Printf("Error while unmarshelling sync res : %v", err)
			return sync, err
		}
	}
	return sync, nil
}

// GetHeimdallVersion will returns the software version of heimdall
func GetHeimdallVersion(cfg *config.Config) (string, error) {
	var v string
	var version HeimdallVersion
	url := cfg.Endpoints.HeimdallLCDEndpoint + "/node_info"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting heimdall version: %v", err)
		return v, err
	}

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return v, err
		}
		err = json.Unmarshal(body, &version)
		if err != nil {
			log.Printf("Error while getting heimdall version : %v", err)
			return v, err
		}
	}
	v = version.ApplicationVersion.Version
	log.Printf("Heimdall Verison : %s", v)

	return v, nil
}

// GetHeimdallNodeStats which returns avg block time and tx count etc
func GetHeimdallNodeStats(cfg *config.Config) (HeimdallStats, error) {
	var h HeimdallStats
	url := "http://159.65.184.68:3456" + "/stats"
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error while getting heimdall stats: %v", err)
		return h, err
	}

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error while reading heimdall stats resp body : %v", err)
			return h, err
		}
		err = json.Unmarshal(body, &h)
		if err != nil {
			log.Printf("Error while unmarshelling heimdall stats res : %v", err)
			return h, err
		}
	}
	return h, nil
}
