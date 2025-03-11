// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DIDStorage {
    // Переменная для хранения строки
    string private storedData;

    // Событие для логирования изменений
    event DataChanged(string newData);

    // Функция для изменения строки
    function setData(string memory _data) public {
        storedData = _data;
        emit DataChanged(_data);
    }

    // Функция для получения строки
    function getData() public view returns (string memory) {
        return storedData;
    }
}