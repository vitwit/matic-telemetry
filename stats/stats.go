package stats

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/PrathyushaLakkireddy/heimdall-node-stats/config"
)

type connWrapper struct {
	conn *websocket.Conn

	rlock sync.Mutex
	wlock sync.Mutex
}

// nodeInfo is the collection of meta information about a node that is displayed
// on the monitoring page.
type nodeInfo struct {
	Name     string `json:"name"`
	Node     string `json:"node"`
	Port     int    `json:"port"`
	Network  string `json:"net"`
	Protocol string `json:"protocol"`
	API      string `json:"api"`
	Os       string `json:"os"`
	OsVer    string `json:"os_v"`
	Client   string `json:"client"`
	History  bool   `json:"canUpdateHistory"`
}

// nodeStats is the information to report about the local node.
type nodeStats struct {
	Active          bool   `json:"active"`
	Syncing         bool   `json:"syncing"`
	Mining          bool   `json:"mining"`
	Hashrate        int    `json:"hashrate"`
	Peers           int    `json:"peers"`
	GasPrice        int    `json:"gasPrice"`
	Uptime          int    `json:"uptime"`
	HeimdallVersion string `json:"hversion"`
}

// // blockStats is the information to report about individual blocks.
type blockStats struct {
	Number *big.Int `json:"number"`
	Hash   string   `json:"hash"`
	// ParentHash common.Hash    `json:"parentHash"`
	Timestamp *big.Int `json:"timestamp"`
	// Miner      common.Address `json:"miner"`
	GasUsed   uint64    `json:"gasUsed"`
	GasLimit  uint64    `json:"gasLimit"`
	Diff      string    `json:"difficulty"`
	TotalDiff string    `json:"totalDifficulty"`
	Txs       []txStats `json:"transactions"`
	TxHash    string    `json:"transactionsRoot"`
	// Root   common.Hash `json:"stateRoot"`
	Uncles          []string `json:"uncles"`
	HeimdallVersion string   `json:"heimdallVersion"`
	BorVersion      string   `json:"borVersion"`
}

// type uncleStats []string

type txStats struct {
	Hash string `json:"hash"`
}

// authMsg is the authentication infos needed to login to a monitoring server.
type authMsg struct {
	ID     string   `json:"id"`
	Info   nodeInfo `json:"info"`
	Secret string   `json:"secret"`
}

func Dailer(cfg *config.Config) error {
	// Resolve the URL, defaulting to TLS, but falling back to none too
	path := fmt.Sprintf("ws://%s/api", cfg.StatsDetails.NetStatsIPAddress)
	urls := []string{path}
	var (
		conn *connWrapper
		err  error
	)

	dialer := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	header := make(http.Header)

	header.Set("origin", "http://localhost")
	for _, url := range urls {
		c, _, e := dialer.Dial(url, header)
		err = e
		if err == nil {
			conn = newConnectionWrapper(c)
			break
		}
	}

	if err != nil {
		log.Printf("Stats server unreachable", err)
		return err
	}

	if err = login(conn, cfg); err != nil {
		log.Printf("Stats login failed : %v", err)
		conn.conn.Close()
		return err
	}

	// Send the initial stats so our node looks decent from the get go
	if err = report(conn, cfg); err != nil {
		log.Printf("Initial stats report failed", "err", err)
	}
	return nil
}

// WriteJSON wraps corresponding method on the websocket but is safe for concurrent calling
func (w *connWrapper) WriteJSON(v interface{}) error {
	w.wlock.Lock()
	defer w.wlock.Unlock()

	return w.conn.WriteJSON(v)

}

func newConnectionWrapper(conn *websocket.Conn) *connWrapper {
	return &connWrapper{conn: conn}
}

// ReadJSON wraps corresponding method on the websocket but is safe for concurrent calling
func (w *connWrapper) ReadJSON(v interface{}) error {
	w.rlock.Lock()
	defer w.rlock.Unlock()

	return w.conn.ReadJSON(v)
}

// login tries to authorize the client at the remote server.
func login(conn *connWrapper, cfg *config.Config) error {
	status, err := GetLatestBlock(cfg)
	if err != nil {
		log.Printf("Error while getting network details : %v", err)
		return err
	}
	node := cfg.StatsDetails.Node
	str := strings.Split(cfg.StatsDetails.NetStatsIPAddress, ":")
	port, _ := strconv.Atoi(str[1])
	auth := &authMsg{
		ID: node,
		Info: nodeInfo{
			Name:    node,
			Node:    node,
			Port:    port,
			Network: status.Result.Network,
			// Protocol: strings.Join(protocols, ", "),
			API:     "No",
			Os:      runtime.GOOS,
			OsVer:   runtime.GOARCH,
			Client:  "0.1.1",
			History: true,
		},
		Secret: cfg.StatsDetails.SecretKey,
	}
	login := map[string][]interface{}{
		"emit": {"hello", auth},
	}
	if err := conn.WriteJSON(login); err != nil {
		return err
	}
	// Retrieve the remote ack or connection termination
	var ack map[string][]string
	if err := conn.ReadJSON(&ack); err != nil || len(ack["emit"]) != 1 || ack["emit"][0] != "ready" {
		return errors.New("unauthorized")
	}
	return nil
}

func report(conn *connWrapper, cfg *config.Config) error {
	for {
		err := ReportBlock(conn, cfg)
		if err != nil {
			log.Printf("Error while reporting block details : %v", err)
			return err
		}

		if err = reportStats(conn, cfg); err != nil {
			log.Printf("Error while reporting node stats : %v", err)
			return err
		}
		time.Sleep(4 * time.Second)
	}

	return nil
}

// ReportBlock retrieves the current block details and reports it to the stats server.
func ReportBlock(conn *connWrapper, cfg *config.Config) error {
	block, err := GetLatestBlock(cfg)
	if err != nil {
		log.Printf("Error while getting block details : %v", err)
		return err
	}
	if block.Result.SyncInfo.LatestBlockHeight == "" {
		log.Printf("Got an empty block result ")
		return err
	}

	number := new(big.Int)
	number, ok := number.SetString(block.Result.SyncInfo.LatestBlockHeight, 10)
	if !ok {
		log.Println("SetString: error")
		// return
	}
	log.Printf("Block height : %v", number)

	thetime, err := time.Parse(time.RFC3339, block.Result.SyncInfo.LatestBlockTime)
	if err != nil {
		panic("Can't parse time format")
	}
	epoch := thetime.Unix()
	s := strconv.FormatInt(epoch, 10)

	blockTime := new(big.Int)
	blockTime, ok = number.SetString(s, 10)
	if !ok {
		log.Println("SetString: error")
		// return
	}
	log.Printf("Block Time : %v", blockTime)

	details := blockStats{
		Number:    number,
		Hash:      block.Result.SyncInfo.LatestBlockHash,
		Timestamp: blockTime,
		TxHash:    "---", // dummy data
		Txs: []txStats{
			{
				Hash: "---", // dummy data
			},
		},
		Uncles: []string{
			"---", // dummy data as frontend is not accepting empty response
		},
	}

	// Assemble the block report and send it to the server
	log.Printf("Sending new block to ethstats", "number", details.Number)

	stats := map[string]interface{}{
		"id":    cfg.StatsDetails.Node,
		"block": details,
	}
	report := map[string][]interface{}{
		"emit": {"block", stats},
	}

	return conn.WriteJSON(report)
}

// reportStats retrieves various stats about the node and
// reports it to the stats server.
func reportStats(conn *connWrapper, cfg *config.Config) error {
	netInfo, err := GetNetInfo(cfg)
	if err != nil {
		log.Printf("Error while getting net info : %v", err)
		return err
	}

	sync, err := SyncStatus(cfg)
	if err != nil {
		log.Printf("Error while getting sync info : %v", err)
		return err
	}

	heimdallVersion, err := GetHeimdallVersion(cfg)
	if err != nil {
		log.Printf("Error while getting heimdall version : %v", err)
	}

	peers, _ := strconv.Atoi(netInfo.Result.NPeers)
	stats := map[string]interface{}{
		"id": cfg.StatsDetails.Node,
		"stats": &nodeStats{
			Active: netInfo.Result.Listening,
			Mining: true,
			// Hashrate: 1,
			Peers:           peers,
			GasPrice:        1000,
			Syncing:         sync.Syncing,
			HeimdallVersion: heimdallVersion,
			// Uptime:   100,
		},
	}
	report := map[string][]interface{}{
		"emit": {"stats", stats},
	}
	log.Printf("reporting node stats..", report)
	return conn.WriteJSON(report)
}
