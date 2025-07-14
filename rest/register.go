package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"

	"github.com/vitwit/matic-telemetry/client"
	"github.com/vitwit/matic-telemetry/config"
	"github.com/vitwit/matic-telemetry/stats"
)

type RegisterRequest struct {
	Secret   string   `json:"secret"`
	NodeInfo NodeInfo `json:"nodeInfo"`
}

type NodeInfo struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Network     string `json:"network"`
	OS          string `json:"os"`
	GoVersion   string `json:"goVersion"`
	Address     string `json:"address"`
	PubkeyType  string `json:"pubkeyType"`
	Pubkey      string `json:"pubkey"`
	IsValidator bool   `json:"isValidator"`
}

func RegisterNode(ctx *client.AppContext, cfg *config.Config) error {

	status, err := stats.GetStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	version, err := stats.GetVersion(cfg)
	if err != nil {
		version, err = stats.GetHeimdallVersion(cfg)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}
	}

	vp, _ := strconv.Atoi(status.Result.ValidatorInfo.VotingPower)
	payload := RegisterRequest{
		Secret: cfg.StatsDetails.SecretKey,
		NodeInfo: NodeInfo{
			ID:          status.Result.NodeInfo.ID,
			Version:     version,
			Network:     status.Result.NodeInfo.Network,
			OS:          fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH),
			GoVersion:   runtime.Version(),
			Address:     status.Result.ValidatorInfo.Address,
			PubkeyType:  status.Result.ValidatorInfo.PubKey.Type,
			Pubkey:      status.Result.ValidatorInfo.PubKey.Value,
			IsValidator: vp > 0,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.StatsDetails.StatsServiceURL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("404 received but failed to parse body: %s", string(body))
		}

		if errResp.Message == "Node already registered" {
			return nil
		}
		return errors.New(errResp.Message)
	}

	return nil
}
