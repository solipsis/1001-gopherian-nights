# Creating and Deriving extended keys
## 

```go
//********************************************************
// Basic structure of an Extended Key
//********************************************************
type ExtendedKey struct {
	key         []byte (32-33 bytes)
	chaincode   []byte (32 bytes)
	depth       uint8 
	fingerprint []byte (4 bytes) checksum of parent
	index       uint32 (4 bytes)
	version     []byte (4 bytes)
	isPrivate   bool
}

//********************************************************
// Deserializing an XPUB or XPRV into an extended key
//********************************************************
{	
	// Decode the serialized key from base58
	decoded := base58.Decode(key)
	
	// serialized format
	// version (4) depth (1) fingerprint (4) Index (4) Chaincode (32) keydata (33) checksum (4)
	payload := decoded[:len(decoded)-4]
	
	// the checksum is the last 4 bytes of the payload
	// hash the payload twice with sha256 and compare the first 4 bytes to the checksum
	chksum := decoded[len(decoded-4):]
	if DoubleHashB(payload)[:4] != checksum {
		ERROR		
	}

	// Fetch all the data from the payload
	version := payload[:4] /
	depth := payload[4:5][0]
	fingerprint := payload[5:9]
	index := binary.BigEndian.Uint32(payload[9:13])
	chaincode := payload[13:45]
	keyData := payload[45:78]

	// We know the key is private if it starts with 0x00 
	isPrivate := keyData[0] == 0x00
	// drop the 0x00 prefix if private
	if isPrivate {
		keyData = keyData[1:]
	}
	
	// Verify the new key is on the ecliptic curve and you are done
	return NewExtendedKey(keyData, chaincode, depth, fingerprint, index, version, isPrivate), nil
}

//********************************************************
// Converting x and y points on the curve into a compressed public key
//********************************************************
func pointToCompressedPubKey(x, y big.Int) []byte {
	
	// compressed public key format is either 0x02 or 0x03 depending on the
	// 'Arity' of the y coordinate
	var format byte = 0x2
	if y.Bit(0) == 1 {
		format |= 0x1
	}

	// serialize key to 33 byte compressed format
	// the first byte is the format byte which tells you which side of the curve the coordinate is on
	// followed by the bytes of the x coordinate
	// using these 2 bits of data you can reconstruct the y coordinate later
	b := make([]byte, 0, 33)
	b = append(b, format)
	b = append(b, x.Bytes()...)
	return b
}

//********************************************************
// Getting and address from an extended key
//********************************************************
func (k *ExtendedKey) Address(versionByte byte) string {
	// get the pkhash with ripemd160(sha256(b))   Commonly called Hash160
	pkHash := hashBytes(hashBytes(k.PubKeyBytes(), sha256.New()), ripemd160.New())

	// prepend version byte and append checksum
	// The version byte determines the prefix of the address once encoded
	// 0x00 for bitcoin p2pkh 0x30 for lite coin p2pkh etc.
	b := make([]byte, 0, 1+len(pkHash)+4) // version + hash + checksum

	b = append(b, versionByte)
	b = append(b, pkHash...)

	// compute the checksum and append
	chksum := checksum(b)
	b = append(b, chksum[:]...)

	// encode the buffer and you have your address
	return base58.Encode(b)
}

// the checksum is the first 4 bytes of the input hashed twice
func checksum(input []byte) (checksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(checksum[:], h2[:4])
	return
}



//********************************************************
// Serializing an extended key into an xpub or xprv (or whatever prefix is in the version field of the key)
//********************************************************
func (k *ExtendedKey) serializedString() string {

	// serialized format
	// version (4) depth (1) fingerprint (4) Index (4) Chaincode (32) keydata (33) checksum (4)
	serializedKeyLength := 4 + 1 + 4 + 4 + 32 + 33
	serializedBytes := make([]byte, 0, serializedKeyLength+4) // 4 bytes for checksum

	// append all fields to the buffer
	serializedBytes = append(serializedBytes, k.version...)
	serializedBytes = append(serializedBytes, k.depth)
	serializedBytes = append(serializedBytes, k.fingerprint...)
	
	// Pay attention to the ‘Endian’ness off the child index
	var indexBytes [4]byte
	binary.BigEndian.PutUint32(indexBytes[:], k.index)

	serializedBytes = append(serializedBytes, indexBytes[:]...)
	serializedBytes = append(serializedBytes, k.chaincode...)
	
	// if the key is private we need to prepend the byte 0x00 to make it 33 bytes long
	if k.isPrivate {
		serializedBytes = append(serializedBytes, 0x00)
		serializedBytes = append(serializedBytes, k.key...)
	} else {
		serializedBytes = append(serializedBytes, k.PubKeyBytes()...)
	}

	// Compute the checksum and append it to the end
	chksum := checksum(serializedBytes)
	serializedBytes = append(serializedBytes, chksum[:]...)

	// Encode it and you are done
	return base58.Encode(serializedBytes)
}





//********************************************************
// Deriving child key from a public or private extended key
//********************************************************
func (k *ExtendedKey) Derive(i uint32) (*ExtendedKey, error) {
	
	// verify we haven't exceeded the depth limit. 2^8???
	if i > 2^8 {
		Error
	}	

	// 4 possible derivation scenarios
	// 1. private EK -> hardened child private EK
	// 2. private EK -> non-hardened child private EK
	// 3. public EK -> non-hardened child public EK
	// 4. public EK -> hardened child public EK. (NOT ALLOWED)

	// check for invalid case where public key trying to make hardened child
	// hardedened keys start at 2^31 (0x80000000)
	isHardened := i >= 0x80000000
	if !k.isPrivate && isHardened {
		return nil, errors.New("not allowed to create hardened child from public key")
	}

	// as per BIP32 create a 33 byte key
	// hardened child key = 0x00 || ser256(parentKey) || ser32(index)
	// normal child key = serP(parentPublicKey) || ser32(index)
	keyData := make([]byte, 33+4) // keylength + index(4 bytes)
	if isHardened {
		// hardedend child of private EK
		// copy starting at position one so we end up with 33 bytes and the first byte 0x00
		copy(keyData[1:], k.key)
	} else {
		// normal child of either pub or private EK
		copy(keyData, k.PubKeyBytes())
	}
	// pay attention to the endianness of the child index
	binary.BigEndian.PutUint32(keyData[33:], i)

	// get intermediate key with HMAC(Sha512)
	// with chain code as the inner key
	hmac512 := hmac.New(sha512.New, k.chaincode)
	hmac512.Write(keyData)
	iKey := hmac512.Sum(nil)

	// split into 2 halves. left half is key. right half is new chain code
	iL := iKey[:len(iKey)/2]
	childChainCode := iKey[len(iKey)/2:]

	iLNum := new(big.Int).SetBytes(iL)
	// Check that the above point is on the curve

	// Derive child key
	// if private child:
	// 	childKey = parse256(Il) + parentKey
	// if public child:
	// 	childKey = serP(point(parse256(iL)) + parentKey)
	if k.isPrivate {
		
		// interpret key as large int and add intermediate private key to parent private key
		keyNum := new(big.Int).SetBytes(k.key)
		iLNum.Add(iLNum, keyNum)
		iLNum.Mod(iLNum, S256().N) // The N value of the ecliptic curve
		childKey := iLNum.Bytes()

		// Fingerprint with Hash160
		parentFP := hashBytes(hashBytes(k.PubKeyBytes(), sha256.New()), ripemd160.New())[:4]
		return NewExtendedKey(childKey, childChainCode, k.depth+1, parentFP, i, k.version, true), nil
	} else {
		// get the intermediate public key from the intermediate private key
		iX, iY := S256().ScalarBaseMult(iL) // find new point on the curve for intermediate public key
		

		// convert serialized and compressed parent key to x and y coordinates
		pubkey, err := ParsePubKey(k.key, S256())
		if err != nil {
			return nil, err
		}

		// Add intermediate pubkey to parent pubkey and serialize in compressed format
		childX, childY := btcec.S256().Add(iX, iY, pubkey.X, pubkey.Y)
		var format byte = 0x2
		if childY.Bit(0) == 1 { // check if the y coordinate is odd
			format |= 0x1
		}

		// serialize the key to 33 byte compressed format
		childKey := make([]byte, 0, 33)
		childKey = append(childKey, format)
		childKey = append(childKey, childX.Bytes()...)

		// compute fingerprint with Hash160
		parentFP := hashKey(hashKey(k.PubKeyBytes(), sha256.New()), ripemd160.New())[:4]

		// We have a shiny new extended key
		return NewExtendedKey(childKey, childChainCode, k.depth+1, parentFP, i, k.version, false), nil

	}
}


//********************************************************
// Encode an extended key into a hex encoded ethereum address   
//********************************************************                                         
  func encodeEthAddress(ek *hdkeychain.ExtendedKey) string {                                                          
          pubKey, err := ek.ECPubKey()                                                                                
          if err != nil {                                                                                             
                  log.Fatal("Unable to derive public key", err)                                                       
          }                                                                                                           
          pBytes := pubKey.SerializeUncompressed()                                                                    
                                                                                                                      
          // Eth address is the hex encoding of the last 20 bytes of the kekkak256 of the 
	  // uncompressed public key without the first byte                 
          hasher := sha3.NewKeccak256()                                                                               
          hasher.Write(pBytes[1:])                                                                                    
                                                                                                                      
          hash := hasher.Sum(nil)                                                                                     
          return hex.EncodeToString(hash[len(hash)-20:])                                                              
  }                                        
```


