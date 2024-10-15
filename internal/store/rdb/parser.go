package rdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// common RDB section indicators
var (
	RDBHeader = []byte("REDIS")
	MetadataHeader = byte(0xFA)
	DatabaseHeader = byte(0xFE)
	KeyExpiryHeaderSec = byte(0xFD)
	KeyExpiryHeaderMS = byte(0xFC)
	EOFHeader = byte(0xFF)
)

/*
	ParseRdbFile is responsible for reading an RDB file dump, validating it, and finally returning a valid
	value of the type *RDBType.

	Function Signature:
		func ParseRdbFile(dirPath, fileName string) (*RDBType, error)

	Parameters:
		- dirPath: Directory where the RDB dump is present. (string)
		- fileName: Name of the RDB dump. (string)

	Returns:
		- *RDBType - A reference to a RDBType which is the parsed out form of the RDB dump.
		- error - Error, if any, else nil.

	Example Usage:
	  rdbData, err := ParseRdbFile("xyz/dumps", "dump.rdb")
	  // Output
	  rdbData = {
	  		Version: 7.2.0
			Databases: [
				{
					DatabaseNumber: 0,
					HashTableSize: 5,
					ExpiryHashTableSize: 5,
					KVMap: {
						"foo": {
							Value: "bar",
							ExpireAt: 1729006766003
						},
						"appple": {
							Value: "banana"
						}
					}
				}
			]
		},
		error = nil
*/
func ParseRdbFile(dirPath, fileName string) (*RDBType, error) {
	parsedRdb := new(RDBType)
	parsedRdb.Databases = []RDBDatabase{}

	data, err := readRdbFile(dirPath, fileName)
	if err != nil {
		return nil, err
	}

	isValidRdb, err := validateRdbFile(data)
	if err != nil {
		return nil, err
	} else if !isValidRdb {
		return nil, fmt.Errorf("malformed rdb file")
	}

	parsedRdb.Version = string(data[5:9])

	databases, err := parseDatabaseSection(data)
	if err != nil {
		return nil, err
	}
	parsedRdb.Databases = databases

	return parsedRdb, nil
}

/*
	readRDBFile reads a rdb dump from a given path and returns it in the form of a byte array

	Function Signature:
		func readRdbFile(dirPath, fileName string) ([]byte, error)

	Parameters:
		- dirPath: Directory where the RDB dump is present. (string)
		- fileName: Name of the RDB dump. (string)
	
	Returns:
		- []byte - A byte array containing the rdb data.
		- error - Error, if any, else nil.

	Example Usage:
	  rdbData, err := readRdbFile("xyz/dumps", "dump.rdb")
	  // Output
	  rdbData = [5 103 114 97 112 101 0 5 97 112 112 108 101 5 97 112 112 108 101 0 6 98 97 110 97 110 97 9 112 105 110 101 97 112 112 108 101 0 6 111 114 97 110 103 101 10 115 116 114 97 119 98 101 114 114 121 0 9 114 97 115 112 98 101 114 114 121 9 98 108 117 101 98 101 114 114 121 255 227 231 5 135 233 60 120 138 10]
	  err = nil
*/
func readRdbFile(dirPath, fileName string) ([]byte, error) {
	path := fileName
	if dirPath != "" {
		path = dirPath + "/" + path
	}

	rdbData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in reading rdb file: %s", err.Error())
	}

	return rdbData, nil
}

/*
	validateRdbFile checks if a given rdb byte array is valid or not by performing some basic checks.

	Function Signature:
		func validateRdbFile(data []byte) (bool, error)

	Parameters:
		- data: A byte array representation of the rdb dump. ([]byte])
	
	Returns:
		- bool - Whether it is a valid rdb or not.
		- error - Error, if any, else nil.

	Example Usage:
		isValid, err := validateRdbFile([5 103 114 97 112 101 0 5 97 112 112 108 101 5 97])
		// Output isValid = false, err = nil
*/
func validateRdbFile(data []byte) (bool, error) {
	// First 5 bytes must have the 'REDIS' magic string
	if len(data) < 5 || !bytes.Equal(data[:5], RDBHeader) {
		return false, fmt.Errorf("malformed rdb file: must start with REDIS header")
	}
	// Next 4 bytes should be present which store the version number
	if len(data) < 9 {
		return false, fmt.Errorf("malformed rdb file: version is missing")
	}
	// EOF indicator should be present
	if !bytes.Contains(data, []byte{EOFHeader}) {
		return false, fmt.Errorf("malformed rdb file: EOF missing")
	}
	return true, nil
}

/*
	parseDatabaseSection takes in a rdb byte array and parses the database sections out of it.

	Function Signature:
		func parseDatabaseSection(data []byte) ([]RDBDatabase, error)

	Parameters:
		- data: A byte array representation of the rdb dump. ([]byte])

	Returns:
		- []RDBDatabase - The parsed database sections from the dump.
		- error - Error, if any, else nil.

	Example Usage:
	  rdbDatabases, err := parseDatabaseSection([5 103 114 97 112 101 0 5 97 112 112 108 101 5 97])
	  // Output
	  rdbDatabases = [
				{
					DatabaseNumber: 0,
					HashTableSize: 5,
					ExpiryHashTableSize: 5,
					KVMap: {
						"foo": {
							Value: "bar",
							ExpireAt: 1729006766003
						},
						"appple": {
							Value: "banana"
						}
					}
				}
			],
		error = nil
*/
func parseDatabaseSection(data []byte) ([]RDBDatabase, error) {
	if !bytes.Contains(data, []byte{DatabaseHeader}) {
		return []RDBDatabase{}, nil
	}
	startIndex := bytes.Index(data, []byte{DatabaseHeader}) + 1

	databases := make([]RDBDatabase, 0)
	for {
		database := new(RDBDatabase)
		kvMap := make(map[string]KVValue)

		databaseNumber, bytesConsumed, err := parseSizeEncoding(data[startIndex:])
		if err != nil {
			return nil, err
		}
		startIndex += bytesConsumed + 1 // + 1 because databseNumber follows 1 byte which has 0xFB indicating hash table size information follows

		hashTableSize, bytesConsumed, _ := parseSizeEncoding(data[startIndex:])
		if err != nil {
			return nil, err
		}
		startIndex += bytesConsumed

		expiryHashTableSize, bytesConsumed, _ := parseSizeEncoding(data[startIndex:])
		if err != nil {
			return nil, err
		}
		startIndex += bytesConsumed
		
		database.DatabaseNumber = databaseNumber
		database.HashTableSize = hashTableSize
		database.ExpiryHashTableSize = expiryHashTableSize

		for {
			if data[startIndex] == KeyExpiryHeaderMS || data[startIndex] == KeyExpiryHeaderSec {
				expiry := uint64(0)
				if data[startIndex] == KeyExpiryHeaderMS {
					expiry = binary.LittleEndian.Uint64(data[startIndex + 1:startIndex + 9])
					startIndex += 9
					if data[startIndex] != 0 {
						return databases, fmt.Errorf("only string type is supported")
					}
					startIndex++
				} else {
					expiry = binary.LittleEndian.Uint64(data[startIndex + 1:startIndex + 5]) * 1000
					startIndex += 5

					if data[startIndex] != 0 {
						return databases, fmt.Errorf("only string type is supported")
					}
					startIndex++
				}

				key, bytesConsumed, err := stringEncoding(data[startIndex:])
				if err != nil {
					return nil, err
				}
				startIndex += bytesConsumed

				val, bytesConsumed, err := stringEncoding(data[startIndex:])
				if err != nil {
					return nil, err
				}
				startIndex += bytesConsumed
				
				if uint64(time.Now().UnixMilli()) < expiry {
					kvMap[key] = KVValue {
						Value: val,
						ExpireAt: expiry,
					}
				}
				
			} else {
				if data[startIndex] != 0 {
					return databases, fmt.Errorf("only string type is supported")
				}
				startIndex = startIndex + 1

				key, bytesConsumed, err := stringEncoding(data[startIndex:])
				if err != nil {
					return nil, err
				}
				startIndex += bytesConsumed

				val, bytesConsumed, err := stringEncoding(data[startIndex:])
				if err != nil {
					return nil, err
				}
				startIndex += bytesConsumed
				
				kvMap[key] = KVValue {
					Value: val,
				}
			}

			if startIndex >= len(data) || data[startIndex] == DatabaseHeader || data[startIndex] == EOFHeader {
				break
			}
		}

		database.KVMap = kvMap
		databases = append(databases, *database)

		if startIndex >= len(data) || data[startIndex] == DatabaseHeader || data[startIndex] == EOFHeader {
			break
		}
	}

	return databases, nil
}

/*
	stringEncoding takes in a rdb byte array and returns a string based on rdb string encoding specification

	Function Signature:
		func stringEncoding(data []byte) (string, int, error)

	Parameters:
		- data: A byte array representation of the rdb dump. ([]byte])
	
	Returns:
		- string - The decoded string.
		- int - The number of bytes consumed.
		- error - Error, if any, else nil.

	Example Usage:
		size, bytesConsumed, err := stringEncoding([5 103 114 97 112 101 0 5 97 112 112 108 101 5 97])
		// Output size = "abc", bytesConsumed = 1, err = nil
*/
func stringEncoding(data []byte) (string, int, error) {
    strSize, bytesConsumed, err := parseSizeEncoding(data)
    if err != nil {
        return "", 0, fmt.Errorf("error parsing size encoding: %v", err)
    }

    if bytesConsumed + strSize > len(data) {
        return "", 0, fmt.Errorf("string size exceeds data length: %d > %d", bytesConsumed + strSize, len(data))
    }

    // Extract the string based on the parsed size
    result := string(data[bytesConsumed : bytesConsumed+strSize])

    return result, bytesConsumed + strSize, nil
}

/*
	parseSizeEncoding takes in a rdb byte array and finds the size based on rdb length encoding specification

	Function Signature:
		func parseSizeEncoding(data []byte) (int, int, error)

	Parameters:
		- data: A byte array representation of the rdb dump. ([]byte])
	
	Returns:
		- int - The length / size.
		- int - The number of bytes consumed.
		- error - Error, if any, else nil.

	Example Usage:
		size, bytesConsumed, err := parseSizeEncoding([5 103 114 97 112 101 0 5 97 112 112 108 101 5 97])
		// Output size = 10, bytesConsumed = 1, err = nil
*/
func parseSizeEncoding(data []byte) (int, int, error) {
	firstByte, firstTwoSignificantBits := data[0],  (data[0] & 0b11000000)

	switch firstTwoSignificantBits {
		case byte(0b00000000): {
			// size is in the next 6 bits of first byte
			return int(firstByte & 0b00111111), 1, nil
		}
		case byte(0b01000000): {
			// size is in the next 14 bits, last 6 of first byte and all bits from next byte
			if len(data) < 2 {
        		return -1, 0, fmt.Errorf("insufficient data for 14-bit size")
    		}

			secondByte := data[1]
			return int(((firstByte & 0b00111111) << 8) | secondByte), 2, nil
		}
		case byte(0b10000000): {
			// Size is in the next 32 bits or 64 bits (ignore remaining 6 bits of the first byte)
			if len(data) < 5 {
				return -1, 0, fmt.Errorf("insufficient data for 32-bit size")
			}

			if byte(firstByte) == byte(0x80) {
				size := int(binary.BigEndian.Uint32(data[1:5]))
				return size, 5, nil
			} else {
				size := int(binary.BigEndian.Uint64(data[1:9]))
				return size, 9, nil
			}
			
		}
		case byte(0b11000000): {
			// Handle special string encodings based on remaining 6 bits
			switch int(firstByte & 0b00111111) {
				case 0: {
					// 8-bit integer encoding
					return int(data[1]), 2, nil
				}
				case 1: {
					// 16-bit integer encoding (little-endian)
					if len(data) < 3 {
						return -1, 0, fmt.Errorf("insufficient data for 16-bit integer")
					}
					return int(binary.LittleEndian.Uint16(data[1:3])), 3, nil
				}
				case 2: {
					// 32-bit integer encoding (little-endian)
					if len(data) < 5 {
						return -1, 0, fmt.Errorf("insufficient data for 32-bit integer")
					}
					return int(binary.LittleEndian.Uint32(data[1:5])), 5, nil
				}
				case 3: {
					// LZF compression not supported
					return -1, 0, fmt.Errorf("lzf compression is not supported")
				}
				default: {
					return -1, 0, fmt.Errorf("unknown string encoding: 0x%X", firstByte)
				}
			}

		}
	}
	return -1, 0, fmt.Errorf("error occurred while parsing size encoding")
}
