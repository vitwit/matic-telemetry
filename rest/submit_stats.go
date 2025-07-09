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
	"time"

	"github.com/vitwit/matic-telemetry/client"
	"github.com/vitwit/matic-telemetry/config"
	"github.com/vitwit/matic-telemetry/stats"
)

type NodeStatsRequest struct {
	Secret string `json:"secret"`
	Stats  Stats  `json:"stats"`
}

type Stats struct {
	NodeID              string   `json:"nodeID"`
	IsSyncing           bool     `json:"isSyncing"`
	Height              uint64   `json:"height"`
	BlockTime           int64    `json:"blockTime"`
	VotingPower         int      `json:"votingPower"`
	Peers               []string `json:"peers"`
	EarliestBlockHeight uint64   `json:"earliestBlockHeight"`
	EarliestAppHash     string   `json:"earliestAppHash"`
	LatestBlockHeight   uint64   `json:"latestBlockHeight"`
	LatestAppHash       string   `json:"latestAppHash"`
	Address             string   `json:"address"`
	Moniker             string   `json:"moniker"`
	Version             string   `json:"version"`
	OS                  string   `json:"os"`
	GoVersion           string   `json:"goVersion"`
	Network             string   `json:"network"`
	Transactions        int      `json:"transactions"`
	Latency             int64    `json:"latency"`
	Latitude            float64  `json:"latitude"`
	Longitude           float64  `json:"longitude"`
	Country             string   `json:"country"`
}

func SubmitStats(ctx *client.AppContext, cfg *config.Config, lat, lon float64, country string) error {

	start := time.Now()
	status, err := stats.GetStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	latency := time.Since(start)
	version, err := stats.GetVersion(cfg)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	netInfo, err := stats.GetNetInfo(cfg)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	h, _ := strconv.Atoi(status.Result.SyncInfo.LatestBlockHeight)
	eh, _ := strconv.Atoi(status.Result.SyncInfo.EarliestBlockHeight)

	peers := []string{}

	for _, peer := range netInfo.Result.Peers {
		peers = append(peers, peer.NodeInfo.ListenAddr)
	}

	blockTime, err := time.Parse(time.RFC3339Nano, status.Result.SyncInfo.LatestBlockTime)
	if err != nil {
		return err
	}

	transactionsLen, err := stats.GetTransactions(cfg, h)
	if err != nil {
		return err
	}

	vp, _ := strconv.Atoi(status.Result.ValidatorInfo.VotingPower)
	payload := NodeStatsRequest{
		Secret: cfg.StatsDetails.SecretKey,
		Stats: Stats{
			NodeID:              status.Result.NodeInfo.ID,
			Version:             version,
			Network:             status.Result.NodeInfo.Network,
			OS:                  fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH),
			GoVersion:           runtime.Version(),
			Address:             status.Result.ValidatorInfo.Address,
			IsSyncing:           status.Result.SyncInfo.CatchingUp,
			Height:              uint64(h),
			BlockTime:           blockTime.Unix(),
			VotingPower:         vp,
			Moniker:             status.Result.NodeInfo.Moniker,
			LatestAppHash:       status.Result.SyncInfo.LatestAppHash,
			LatestBlockHeight:   uint64(h),
			EarliestAppHash:     status.Result.SyncInfo.EarliestAppHash,
			EarliestBlockHeight: uint64(eh),
			Peers:               peers,
			Transactions:        transactionsLen,
			Latency:             latency.Milliseconds(),
			Latitude:            lat,
			Longitude:           lon,
			Country:             country,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.StatsDetails.StatsServiceURL+"/api/submit-stats", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("404 received but failed to parse body: %s", string(body))
		}
		return errors.New(errResp.Message)
	}

	return nil
}
