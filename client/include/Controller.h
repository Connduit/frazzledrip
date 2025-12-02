#ifndef CONTROLLER_H
#define CONTROLLER_H


#include "MessageTypes.h"
#include "MessageHandler.h"


// TODO: 
// https://attack.mitre.org/tactics/enterprise/
class Controller
{
public:
	Controller(MessageHandler* messageHandler, ApiManager* apiManager);
	void handleDefault(const InternalMessage& msg);
	void handleNone(const InternalMessage& msg);

	// TODO: rename to handleExecuteConsoleCommand? 
	void handleExecuteCommand(const InternalMessage& msg);

	void handleExecuteShellcode(const InternalMessage& msg);

	void handleSystemInfo(const InternalMessage& msg);

	// TODO: make a subclass called PersistenceController?
	void handleEstablishPersistence(const InternalMessage& msg);
private:
	MessageHandler* messageHandler_;
	ApiManager* apiManager_;

};

#endif
