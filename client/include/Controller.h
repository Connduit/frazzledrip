#ifndef CONTROLLER_H
#define CONTROLLER_H


#include "MessageTypes.h"



class Controller
{
public:
	void handleDefault(const InternalMessage& msg);
	void handleNone(const InternalMessage& msg);
	void handleExecuteCommand(const InternalMessage& msg);
private:
};

#endif
