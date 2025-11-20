#ifndef CLIENT_SUBSYSTEM_H
#define CLIENT_SUBSYSTEM_H

#include "TransportLayer.h"
#include "MessageParser.h"
#include "Dispatcher.h"
#include "Controller.h"

// TODO: rename to ClientController? 
class ClientSubsystem
{
public:
    ClientSubsystem();
    // ClientSubsystem(ClientConfig config); // TODO: 
    ~ClientSubsystem();

	void run();

private:
	void setupMessaging();
	void setupEvents();
	void setupTasks();

	void setupSubcomponents();
	//void setupModules();

	TransportLayer* transportLayer_;
	MessageParser* messageParser_;

	Dispatcher* dispatcher_;
	Controller* controller_;

	Serializer* serializer_;
	Encoder* encoder_;
	Encryptor* encryptor_;

};
#endif
