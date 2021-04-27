# heimdall-stats-exporter
This will fetch data from rpc and lcd endpoints and export the metrics to a netstats dashboard by establishing a connection between the client and server using web sockets. Data will be transferred through the socket and displayed over the netstats dashboard. 

It is recommended to run the ethstats dashboard and heimdall-stats-exporter on a sentry node which has both heimdall and bor already configured and operational. 

## Instructions to start ethstats dashboard

 - Start eth-netstats using these instructions https://github.com/cubedro/eth-netstats 
 - While starting that service, if you get `WS_SECRET NOT SET!!!` , then do `export WS_SECRET="key_name"` ex: `export WS_SECRET=hello`
 - You can check the running dashboard on ip_address:3000. Note that data will be displayed only after exporting bor and heimdall stats.
 

## Restart your bor node with the ethstats flag

  
   - Add `--ethstats` flag to your bor bash script which will be present at `~/node/bor/start.sh`. After adding the flag the bash file it should look like this:
   ```#!/usr/bin/env sh

set -x #echo on

BOR_DIR=${BOR_DIR:-~/.bor}
DATA_DIR=$BOR_DIR/data

bor --datadir $DATA_DIR \
  --ethstats bor:hello@localhost:3000 \
  --port 30303 \
  --http --http.addr '0.0.0.0' \
  --http.vhosts '*' \
  --http.corsdomain '*' \
  ......
  ......
```
   - Restart your bor service `sudo systemctl restart bor`

 
 

## Running heimdall stats exporter

 - Get the code
 ``` bash
    git clone https://github.com/PrathyushaLakkireddy/heimdall-node-stats.git
    cd heimdall-node-stats
    git fetch && git checkout main
    cp example.config.toml config.toml
 ```


 - Edit config.toml (you can use the default values of the `config.toml` if you're deploying this exporter on a sentry node itself)
   
   - *heimdall_rpc_endpoint*
        
        Heimdall rpc end point (RPC of your own validator) is used to gather information about network info, latest block info etc.

   - *heimdall_lcd_endpoint*

        Heimdall lcd end point is used to gather information about node info and syncing status.

   - *secret_key*

      Secret key, which will be used to connect with netstats. Mention the one which you have exported while starting netstats dashbaord.

   - *node*

      Name of your node (ex : heimdall_node). In dahsboard metrics will be displayed on the name of it. So it will be easy to classify multiple networks data based on the node name.

   - *net_stats_ip*

      Provide the ip where the netstats dashboard was running (ex : localhost:3000). Based on this IP only, this script connect to the netstats server and export the metrics to the dashboard.

 - Start the exporter service using `go run main.go`
- You can check the telemetry stats at http://server-ip:3000 

```bash 
Note :: Make sure to use the same secret key for both bor and heimdall .
```
