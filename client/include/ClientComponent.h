#ifndef CLIENT_COMPONENT_H
#define CLIENT_COMPONENT_H

#include "ClientSubsystem.h"


class ClientComponent
{
public:
    ClientComponent();
    ~ClientComponent();

	void loadConfig();
	void createClientSubsystem();
	void startClientSubsystem();
private:

	ClientSubsystem* clientSubsystem_;
	// ClientConfig* config_;

};




#endif
