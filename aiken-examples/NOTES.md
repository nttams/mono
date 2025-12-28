```sh
aiken blueprint convert > gift.script
cardano-cli address build --testnet-magic 2 --payment-script-file gift.script | tee gift.addr
```