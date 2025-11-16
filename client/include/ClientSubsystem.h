#ifndef CLIENT_SUBSYSTEM_H
#define CLIENT_SUBSYSTEM_H

class ClientSubsystem
{
public:
    ClientSubsystem();
    // ClientSubsystem(ClientConfig config); // TODO: 
    ~ClientSubsystem();

private:
	void setupMessaging();
	void setupEvents();
	void setupTasks();

	TransportLayer* transportLayer_;
	MessageHandler* messageHandler_;

};
#endif
