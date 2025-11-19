#ifndef CLIENT_SUBSYSTEM_H
#define CLIENT_SUBSYSTEM_H

#include "TransportLayer.h"
#include "MessageHandler.h"

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
	MessageHandler* messageHandler_;

	ClientDispatcher* dispatcher_;
	ClientController* controller_;

	Serializer* serializer_;
	Encoder* encoder_;
	Encryptor* encryptor_;

};
#endif
