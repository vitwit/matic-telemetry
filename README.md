# heimdall-stats-exporter
This will fetch data from rpc/lcd and export the metrics, to netstas dashboard by establishing a connection between both (client and server) by using web sockets. Data will be transfered through the socket and displayed over the netstas dashboard.

## Instructions to start bor with ethstats

 - https://github.com/cubedro/eth-netstats follow instrctions from here and start eth-netstats
 - While starting the server, if you get `WS_SECRET NOT SET!!!` , then do `export WS_SECRET="key_name"` . (ex: export WS_SECRET=hello)
 - Note this same secret key will be used, while running bor and heimdall nodes to be connected wth eth-nestats

## Run your bor node by giving flag of --ethstats

 - If you are running any geth related node then you can give this flag while running it.
 - For bor ex: bor `--ethstats node:secretKey@host:port` (bor:hello@localhost:3000)
 - and you have to give the host and port of netstats listening server.

## Running ethstats with heimdall

 - Geth the code
 ``` bash
    git clone https://github.com/PrathyushaLakkireddy/heimdall-node-stats.git
    git fetch && git checkout  heimdall-stats
    cp example.config.toml config.toml
 ```

 - Before running the code make sure to configure the config.toml, by providing the mentioned fields.
 - and then start the server, by running `go run main.go`
- After all these steps you can check the stats at http://ip:3000 (ex:http://localhost:3000)

```bash 
Note :: Make sure to give netstats secret key for other nodes, then only your nodes can be connected with netstats and details will be displayed there.
```

