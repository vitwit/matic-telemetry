### Prerequisites

- node
- npm

Clone the repository and install the dependencies

```
git clone https://github.com/PrathyushaLakkireddy/eth-netstats.git
cd eth-netstats
git fetch && git checkout master
npm install
sudo npm install -g grunt-cli
```
While starting that service, if you get `WS_SECRET NOT SET!!!` , then do `export WS_SECRET="key_name"` ex: `export WS_SECRET=hello`
You can check the dashboard on [localhost:3000](). Note that data will be displayed only after exporting bor stats.
##### Restart your bor node with the ethstats flag

  
   - Add `--ethstats` flag to your bor bash script which will be present at `~/node/bor/start.sh`. After adding the flag to the bash file it should look like this:
   ```#!/usr/bin/env sh

set -x #echo on

BOR_DIR=${BOR_DIR:-~/.bor}
DATA_DIR=$BOR_DIR/data

bor --datadir $DATA_DIR \
  --ethstats <node-name>:<key>@localhost:3000 \
  --port 30303 \
  --http --http.addr '0.0.0.0' \
  --http.vhosts '*' \
  --http.corsdomain '*' \
  ......
  ......
```
- Restart your bor service `sudo systemctl restart bor`
