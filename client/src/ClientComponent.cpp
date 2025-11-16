#include "ClientComponent.h"

ClientComponent::ClientComponent()
{
	createClientSubsystem();
}

ClientComponent::~ClientComponent()
{
	if (clientSubsystem_)
	{
		delete clientSubsystem_;
		clientSubsystem_ = 0;
	}
}

ClientComponent::createClientSubsystem()
{
	clientSubsystem_ = new ClientSubsystem();
}
