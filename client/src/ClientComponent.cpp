#include "ClientComponent.h"

ClientComponent::ClientComponent()
{
	createClientSubsystem(); // TODO: should take in a config arguement
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
	clientSubsystem_ = new ClientSubsystem();
}

void ClientComponent::startClientSubsystem()
{
	clientSubsystem_->run();
}

