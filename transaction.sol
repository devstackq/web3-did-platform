// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TransactionManager {
    // struct for save info transaction
    struct Transaction {
        address sender;
        address receiver;
        uint256 amount;
        uint256 timestamp;
    }

    // list of all - transactions
    Transaction[] public transactions;

    //event for logging trx
    event TransactionSent(address indexed sender, address indexed receiver, uint256 amount, uint256 timestamp);


    function sendEth(address payable _receiver) public payable{
        require(msg.value > 0, "amount must be greater than 0");
        require(_receiver != address(0), "invalid receiver address");

        //send ETH to receiver
        (bool success,) = _receiver.call{value: msg.value}("");
        require(success, "transfer failed");

        //save trx
        transactions.push(Transaction({
            sender: msg.sender,
            receiver: _receiver,
            amount: msg.value,
            timestamp: block.timestamp
        }));

        //log event
        emit TransactionSent(msg.sender, _receiver, msg.value, block.timestamp);
    }

    // get history trx
    function getTxHistory() public view returns (Transaction[] memory) {
        return transactions;
    }

    //get history by address
    function getTransactionHistoryByAddress(address _address) public view returns (Transaction[] memory) {
        uint256 count = 0;

        for (uint256 i = 0; i < transactions.length; i++) {
            if (transactions[i].sender == _address || transactions[i].receiver == _address ) {
                count++;
            }
        }

        Transaction[] memory result = new Transaction[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < transactions.length; i++) {
            if (transactions[i].sender == _address || transactions[i].receiver == _address ) {
                result[index] = transactions[i];
                index++;
            }
    }
        
    return result;
    }

}