#include "ClientComponent.h"

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
	clientSubsystem_ = new ClientSubsystem(config_);
}

void ClientComponent::startClientSubsystem()
{
	clientSubsystem_->run();
}

