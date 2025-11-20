#include "Dispatcher.h"
#include <iostream>

void Dispatcher::registerHandler(MessageType type, Handler handler)
{
	handlers_[type] = handler;
}

void Dispatcher::dispatch(const InternalMessage& msg)
{
	std::cout << "Dispatcher::dispatch()" << std::endl;
	/*
	 * TODO: use iterator instead
	 * for (int i = 0; i < handlers_.size(); ++i)
	 * {
	 * 		if (handlers_[msg.header_.messageType])
	 * 		{
	 * 			handlers[i]->handle(); ?
	 * 		}
	 *
	 * }
	 * */
}
