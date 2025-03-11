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

curl -X GET "http://localhost:8080/eth/history"

success trx hash - 0x9c7cc213ffafa0f5f236d6f48c80d5f5aada0218ac7a36b33f8671189f475204

Links:
https://sepolia.etherscan.io/tx/0x9c7cc213ffafa0f5f236d6f48c80d5f5aada0218ac7a36b33f8671189f475204
https://remix.ethereum.org/
https://developer.metamask.io/

Flow:
    MetaMask: Пользователь подключает MetaMask к вашему приложению и выбирает аккаунт.
    Frontend: Приложение отправляет данные (адрес получателя, сумму) на бэкенд.
    Backend: Бэкенд создает транзакцию, подписывает её приватным ключом и отправляет в сеть через ноду (Infura/Alchemy).
    Smart Contract: Смарт-контракт на Solidity обрабатывает транзакцию (например, отправляет ETH и сохраняет данные).
    Ethereum Node: Нода передает транзакцию майнерам для включения в блок.
    Mining: Майнеры включают транзакцию в блок и подтверждают её.
    Transaction Hash: Бэкенд возвращает хэш транзакции пользователю.
    User: Пользователь проверяет статус транзакции на блокчейн-эксплорере с помощью хэша.