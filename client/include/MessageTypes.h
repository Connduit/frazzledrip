/* MessageTypes.h */
// TODO: reorganize order of data types in this file so they're easier to find?
#ifndef MESSAGE_TYPES_H
#define MESSAGE_TYPES_H

#include <vector>
#include <cstdint>


typedef std::vector<uint8_t> RawByteBuffer; // TODO: change to ByteBuffer? 


// TODO: messagetype and their structs should be in their own file called MessageType/s.h?
typedef enum
{
	NONE, 			// default/null/none msg
	DEFAULT, 			// default/null/none msg
	HANDSHAKE,          // Initial connection
	HEARTBEAT,          // Regular check-in
	SYS_INFO,			// Victim system data (RECON)
	COMMAND_RESULT,     // Output from executed command // TODO: is this needed? maybe just for outbound messages
	DATA_EXFIL,         // Stolen data
	UPLOAD_FILE,        // File transfer (client to server)
	EXECUTE_COMMAND,	// CMD command?
	EXECUTE_SHELL_CODE, // shell code
	DOWNLOAD_FILE,		// file transfer (server to client)
	CONFIG_UPDATE,

	ERROR_REPORT        // Error information

} MessageType;

/*
// TODO: does this belong in the MessageHeader? or should these already be established before compile time?
typedef struct {
	bool use_compression;
	SerializationType serialization_type;
	ProtocolType protocol;
	CryptoKey crypto_key;
	Endpoint server_endpoint;

} Config;
*/

#pragma pack(push, 1)
typedef struct MessageHeader
{
	MessageHeader() :
		messageType(MessageType::NONE),
		dataSize(0),
		messageId(0)
	{}

	MessageType messageType;
	uint32_t dataSize;     // Size of following data // size of payload?  // TODO: doesn't have to big as big as a uint32_t
	uint32_t messageId;    // Unique ID for tracking // // TODO: doesn't have to big as big as a uint32_t
	//ULONG messageId; ?
	//ULONG payloadSize; // MessageDataSize
	//ULONG checksum; ?

	//message_id: String = GENERATE_UUID()
	//victim_id: String
	//timestamp: Long
	//compression: Boolean
	//priority: Integer

} MessageHeader;
#pragma pack(pop)

// NOTE: InternalMessage will have to serialize,encrypt,enocde,...etc both the header and payload separately? then combine them together into one array?
// ex: serialize(InternalMessage) wont work, need to do serialize(InternalMessage.header) + serialize(InternalMessage.payload)
// NOTE: this stuct will be easier to detect (struct containing structs), instead of just having the raw fields listed out
#pragma pack(push, 1)
typedef struct InternalMessage
{
	InternalMessage() :
		header(),
		data()
	{}

	MessageHeader header; // the header... always use custom header? no need for tlv... MessageType enum should be defined in header
	//MessageData payload; // the actual payload
	std::vector<uint8_t> data; 
} InternalMessage;
#pragma pack(pop)




#endif
