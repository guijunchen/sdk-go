// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.0;

contract Factory {
    function create(int rtType, string calldata name, bytes calldata code) public returns (address addr){
        assert(0 < rtType && rtType < 8);

	    string memory str = name;
	    bytes memory bytesCode = code;

        bytes32 n;
        assembly {
	        n := mload(add(str, 32))
            addr := create2(rtType, add(bytesCode,0x20), mload(bytesCode), n)
        }
    }
}

