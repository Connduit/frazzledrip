/* MessageHandler.cpp */

#include "MessageHandler.h"
#include "MessageTypes.h"
#include "Recon.h"
#include "TransportLayer.h"
#include "utils.h"

#include <iostream>
#include <sstream>

#include <iostream> // TODO: delete later?

#include <cstdlib> // NOTE: needed to execute shell cmd



MessageHandler::MessageHandler(
	TransportLayer* transportLayer, 
	MessageTransformer* transformer)
	:
	transportLayer_(transportLayer),
	transformer_(transformer)
{
}

MessageHandler::~MessageHandler()
{
	// debug statement
}

// handle inbound messages
void MessageHandler::handle(const RawByteBuffer& data)
{
	std::cout << "MessageHandler::handle(RawByteBuffer)" << std::endl;


	// transform inbound raw bytes into InternalMessage
	InternalMessage msg = transformer_->transform(data);


	if (receiveCallback_)
	{
		receiveCallback_(msg);
	}


}

bool MessageHandler::sendMessage(InternalMessage& msg)
{
	std::cout << "MessageHandler::sendMessage()" << std::endl;
	auto transformed = transformer_->transform(msg);
	transportLayer_->send(transformed);
	return false;
}




