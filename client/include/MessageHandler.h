/* 
 * MessageHandler.h 
 *
 * Purpose: Responsible for handling internal messages (InternalMessage)
 *
 *
 *
 */


/*  (OWNS execution components) */

// TODO: rename class/file to CommandDispatcher?

#ifndef MESSAGE_HANDLER_H
#define MESSAGE_HANDLER_H

#include "MessageTypes.h"

#include "TransportLayer.h"
#include "MessageTransformer.h"
#include "Packer.h"



#include <string>
#include <vector>
#include <cstdint>
#include <functional>
#include <memory>

class TransportLayer;

class MessageHandler
{
public:
	
	using ReceiveCallback = std::function<void(const InternalMessage&)>;
	void setOnMessage(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }


	MessageHandler(TransportLayer* transportLayer, MessageTransformer* transformer, Packer* packer);

	~MessageHandler();


	void handle(const RawByteBuffer& data); // TODO: change rawbytebuffer to const? 

    bool sendMessage(InternalMessage& msg); // overload this function?


private:

	TransportLayer* transportLayer_; // NOTE: pointer is being used instead of a reference because references require immediate initialization
	MessageTransformer* transformer_;
	Packer* packer_;


	ReceiveCallback receiveCallback_;

};


typedef std::unique_ptr<MessageHandler> MessageHandlerPtr;

#endif
