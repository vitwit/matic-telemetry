# heimdall-stats-exporter
This will fetch data from rpc/lcd and exports the metrics, to netstas dashboard by establishing a connection between both (client and server) by using web sockets. Data will be transfered through the socket and displayed over the netstas dashboard. 

- Here is an overview of the netstas dashboard. In matic, will be having bor and heimdall. So for bor(geth based nodes) we can just enable the command --ethstats then metrics will be exported to netstats dashboard. But for heimdall we don't have such option, for that here we have written script, which reads data from rcd and lcd endpoints and exports the data to netstats dashboard.
So that in netstats dashboard you can observe two nodes data.

- You may get doubt, where to run these servers. Actually, `ethstats` dashboard can run on any of your server, it can be validator, sentry or any fresh instance. In the same way, you can run `heimdall-stats-exporter` as well. but for bor you have to enable the `--ethstats` command and retstart, wherever it was running. Note: Those three can run on same server as well.

- Follow below instructions to run `ethstats` and export the data of `bor` and `heimdall`.

## Instructions to start ethstats dashboard

 - https://github.com/cubedro/eth-netstats follow instrctions from here and start eth-netstats
 - While starting the server, if you get `WS_SECRET NOT SET!!!` , then do `export WS_SECRET="key_name"` . (ex: export WS_SECRET=hello)
 - Note this same secret key will be used, while running bor and heimdall nodes to be connected with eth-nestats

## Run your bor node by giving flag of --ethstats

 - If you are running any geth related node then you can give this flag while running the binary of it.
 - If you are running `bor` as a system service, 
   - then you need to edit bor system service file, with editor of your choice. 
   - Add `--ethstats` falg to your `ExecStart` and restart the bor service. (ex: ExecStart=/bin/bash /home/ubuntu/node/bor-start.sh --ethstats bor:hello@localhost:3000 )
   - Then do system demon reload and start your service.

 - If you are running just binary of bor as a non system service file then pass the `--ethstats` falg to it. ex: bor `--ethstats node:secretKey@host:port` (bor:hello@localhost:3000)
 - Note:  you have to give the host and port of netstats listening server.

## Running heimdall stats exporter

 - Get the code
 ``` bash
    git clone https://github.com/PrathyushaLakkireddy/heimdall-node-stats.git
    cd heimdall-node-stats
    git fetch && git checkout main
    cp example.config.toml config.toml
 ```

 - Before running the sevrer make sure to configure the config.toml, by providing the mentioned fields.

 - Edit config.toml
   
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

 - and then start the server, by running `go run main.go`
- After all these steps you can check the stats at http://ip:3000 (ex:http://localhost:3000)

```bash 
Note :: Make sure to give same secret key for other nodes(bor and heimdall) which you have mentioned for `netstats`, then only your nodes can be connected with netstats and details will be displayed there.
```

