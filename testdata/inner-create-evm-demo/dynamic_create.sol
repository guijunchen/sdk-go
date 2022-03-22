contract Factory {
    function create(string calldata name, bytes calldata code) public returns (address addr){
	    string memory str = name;
	    bytes memory bytesCode = code;

        bytes32 n;
        assembly {
	        n := mload(add(str, 32))
            addr := create2(0x05, add(bytesCode,0x20), mload(bytesCode), n)
        }
    }
}

