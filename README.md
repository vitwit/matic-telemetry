# Matic-stats-exporter

- ## Heimdall:-
Telemetry data for Heimdall nodes on Mainnet and Mumbai-testnet can be found here [https://heimdall-mainnet.vitwit.com](https://heimdall-mainnet.vitwit.com) and [https://heimdall-mumbai.vitwit.com](https://heimdall-mumbai.vitwit.com)

To export your nodes telemetry data to these dashboards do the following:-

```
git clone https://github.com/vitwit/matic-telemetry.git
cd matic-telemetry
mkdir -p ~/.telemetry/config
cp example.config.toml ~/.telemetry/config/config.toml
```
Replace default value of `node` with your <node-name> in `~/.telemetry/config/config.toml`.

> **Note:**
> The config file must be present in the `~/.telemetry/config` directory by default. You can specify a different config directory using the `--config` flag when running the telemetry binary. For example:
> 
> ```sh
> ./telemetry --config /path/to/your/configdir
> ```
> If the flag is not provided, the default path is `~/.telemetry/config`.

Use the following secret_key and IP to connect to **Mainnet** dashboard

```
[stats_details]
secret_key = "heimdall_mainnet"  
node = "<node-name>" 
stats_service_url = "heimdall-mainnet.vitwit.com:3000"
```

Use the following secret_key and IP to connect to **Testnet** dashboard

```
[stats_details]
secret_key = "heimdall_testnet"  
node = "<node-name>" 
stats_service_url = "heimdall-mumbai.vitwit.com:3000"
```
Build the binary :-
```
go build -o telemetry
mv telemetry $GOBIN
```
Create systemd file :-
```
echo "[Unit]
Description=Telemtry
After=network-online.target
[Service]
User=$USER
ExecStart=$(which telemetry) --config $HOME/.telemetry/config
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target" | sudo tee "/lib/systemd/system/telemetry.service"
```
Start the telemetry service

```
sudo systemctl enable telemetry.service
sudo systemctl start telemetry.service
```

View the logs using 

`journalctl -u telemetry -f`

