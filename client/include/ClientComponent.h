#ifndef CLIENT_COMPONENT_H
#define CLIENT_COMPONENT_H

#include "ClientSubsystem.h"


class ClientComponent
{
public:
    ClientComponent();
    ~ClientComponent();

private:
	void loadConfig();
	void createClientSubsystem();

	ClientSubsystem clientSubsystem_;

};




#endif
