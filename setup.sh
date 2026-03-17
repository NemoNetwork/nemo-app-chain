nemod init validator --chain-id nemo-testnet
echo "pet apart myth reflect stuff force attract taste caught fit exact ice slide sheriff state since unusual gaze practice course mesh magnet ozone purchase" | nemod keys add validator --keyring-backend test
echo "bottom soccer blue sniff use improve rough use amateur senior transfer quarter" | nemod keys add validator1 --keyring-backend test --recover
nemod add-genesis-account $(nemod keys show validator -a --keyring-backend test) 100000000000000unemo,2500000000000000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
nemod add-genesis-account $(nemod keys show validator1 -a --keyring-backend test) 110000000000000unemo,2500000000000000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
nemod gentx validator 90000000000000unemo --keyring-backend test --chain-id nemo-testnet
nemod collect-gentxs
sed -i 's/stake/unemo/g' ~/.nemo/config/genesis.json