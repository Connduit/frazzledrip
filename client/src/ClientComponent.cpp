#include "ClientComponent.h"

#include <iostream>

ClientComponent::ClientComponent() : config_()
{
}

ClientComponent::~ClientComponent()
{
	if (clientSubsystem_)
	{
		delete clientSubsystem_;
		//clientSubsystem_ = 0;
	}
}

void ClientComponent::loadConfig()
{
}

void ClientComponent::createClientSubsystem()
{
	std::cout << "ClientComponent::createClientSubsystem()" << std::endl;
	clientSubsystem_ = new ClientSubsystem(config_);
}

void ClientComponent::startClientSubsystem()
{
	std::cout << "ClientComponent::startClientSubsystem()" << std::endl;
	clientSubsystem_->run();
}

