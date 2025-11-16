#include "ClientSubsystem.h"
#include "ComponentFactory.h"

ClientSubsystem::ClientSubsystem()
{
	// loadConfig?
	// enable debugging messages (depending on config?)
	setupMessaging();
	// setupEvents();
	// setupTasks();
}

ClientSubsystem::~ClientSubsystem()
{
	if (messageHandler_)
	{
		delete messageHandler_;
		messageHandler_ = 0;
	}

	if (transportLayer_)
	{
		delete transportLayer_;
		transportLayer_ = 0;
	}
}


void ClientSubsystem::setupMessaging()
{

	messageHandler_ = new MessageHandler();
	transportLayer_ = new TCPTransportLayer();
	//transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP);
}
