#include "Dispatcher.h"
#include <iostream>

void Dispatcher::registerHandler(MessageType type, Handler handler)
{
	handlers_[type] = std::move(handler);
}

void Dispatcher::dispatch(const InternalMessage& msg)
{
	std::cout << "Dispatcher::dispatch()" << std::endl;
	auto it = handlers_.find(msg.header_.messageType_);
	if (it != handlers_.end())
	{
		it->second(msg);
	}
	else
	{
		std::cout << "No handlers found for the given messageType: " << msg.header_.messageType_ << std::endl;
	}
}
