Metamask - storage tokens, safety, send & sign transacation in Blockchain 
through Infura

Infure infastructure - with communicate Blockhain; no need create
own blockchain node

go run main.go - run at localhost

curl queries:

/did?address="" - create did account
/eth/balance/"address" - get balance by address
/eth/send - send amount - fromAddress - toAddress

curl -X GET "http://localhost:8080/eth/balance/0x000000000000000000000000000000000000dead"

curl -X POST "http://localhost:8080/eth/send" \
-H "Content-Type: application/json" \
-d '{
"private_key": "privateKeyValue",
"recipient_address": "0x000000000000000000000000000000000000dead",
"amount": 10000000000000000
}'

