Metamask - storage tokens, safety, send & sign transacation in Blockchain
through Infura

Infura infastructure - with communicate Blockhain; no need create
own blockchain node

Metamask: Кошелек и инструмент для взаимодействия с блокчейном
- save keys, address users
- sign trans
- can connect dApps
- can connect - Networks(Eth, Sepolia, etc)

Go Code:
- Interact with Blockchain by RPC(Infura)
- Send trx & call methods Smart Contracts
- fetch data from Blocchain
- API for front end

Solidity code: - Роль: Смарт-контракты, которые выполняют логику на блокчейне
- logic apps(votes, manage trxs, etc)
- save data in Blockchain
- send Eth or call some func


Example:
1.Solidity functional - compile and deploy - to Sepolia network( from Remix)
2.Integrate GO code + Infura;
2.1 use ABI contract for call methods (in Solidity)
3.User - connect through Metamask - in own App
4.Metamask Sign trx -> Go code send -> Blockhain


==================================================================
go run main.go - run at localhost

curl queries:

/did?address="" - create did account
/eth/balance/"address" - get balance by address
/eth/send - send amount - fromAddress - toAddress

curl -X GET "http://localhost:8080/eth/balance/0x000000000000000000000000000000000000dead"

curl -X POST "http://localhost:8080/eth/send" \
-H "Content-Type: application/json" \
-d '{
"private_key": "7bd50df4554e459af04e94576fc0e4b7d21a16b96abb857baa322f1b97f64ea0",
"recipient_address": "0x000000000000000000000000000000000000dead",
"amount": 5000000000000000
}'

