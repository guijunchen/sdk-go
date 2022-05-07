// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract D {
    uint public x;
    constructor(uint a) payable {
        x = a;
    }

    function get() public view returns(uint){
        return x;
    }
}

contract C {
    function createDSalted(uint arg, string calldata name) public {
	    string memory str = name;

        bytes32 n;
	    assembly{
	        n := mload(add(str, 32))
	    }

        D d = new D{salt: n, value: 5}(arg);
    }
}


