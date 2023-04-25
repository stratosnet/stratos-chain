#!/bin/sh
chown -R $RUN_AS_USER $WORK_DIR

if [ ! -d "$WORK_DIR/.stchaind" ]
then
  echo "[entrypoint] Init stratos node..."
  su -s /bin/sh - $RUN_AS_USER -c "stchaind init stratos-node"

  if [ "$CHAIN_ID" == "DEFAULT" ]
  then
    genesis_file_URL=https://raw.githubusercontent.com/stratosnet/stratos-chain-testnet/main/genesis.json
    config_file_URL=https://raw.githubusercontent.com/stratosnet/stratos-chain-testnet/main/config.toml
  else
    genesis_file_URL=https://raw.githubusercontent.com/stratosnet/stratos-chain-testnet/main/$CHAIN_ID/genesis.json
    config_file_URL=https://raw.githubusercontent.com/stratosnet/stratos-chain-testnet/main/$CHAIN_ID/config.toml
  fi

  echo "[entrypoint] Download genesis.json from $genesis_file_URL"
  su -s /bin/sh - $RUN_AS_USER -c "wget $genesis_file_URL -O $WORK_DIR/.stchaind/config/genesis.json"

  echo "[entrypoint] Download config.toml from $config_file_URL"
  su -s /bin/sh - $RUN_AS_USER -c "wget $config_file_URL -O $WORK_DIR/.stchaind/config/config.toml"

  echo "[entrypoint] Set the node moniker to '$MONIKER'"
  su -s /bin/sh - $RUN_AS_USER -c "sed -i 's/moniker = \".*\"/moniker = \"'$MONIKER'\"/g' $WORK_DIR/.stchaind/config/config.toml"
fi

su -s /bin/sh - $RUN_AS_USER -c "$@"
