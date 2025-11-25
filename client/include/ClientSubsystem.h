#ifndef CLIENT_SUBSYSTEM_H
#define CLIENT_SUBSYSTEM_H

#include "TransportLayer.h"
#include "MessageHandler.h"
#include "Dispatcher.h"
#include "Controller.h"
#include "Config.h"
#include "Packer.h"
#include "MessageTransformer.h"

// TODO: rename to ClientController? 
class ClientSubsystem
{
public:
    // ClientSubsystem();
    ClientSubsystem(Config& config);
    ~ClientSubsystem();

	// run/start listening?
	void run();

	// void sendResult(); // TODO: it feels like the wrong place to put this function

	// TODO:
	// void recvMessage(); 
	// void startReceiving(); ???

private:
	void setupMessaging();
	void setupEvents();
	void setupTasks();

	void setupSubcomponents();
	//void setupModules();

	TransportLayer* transportLayer_;
	MessageHandler* messageHandler_;

	Packer* packer_;
	MessageTransformer* messageTransformer_;

	Dispatcher* dispatcher_;
	Controller* controller_;

	Serializer* serializer_;
	Encoder* encoder_;
	Encryptor* encryptor_;

	Config& config_;

};
#endif
