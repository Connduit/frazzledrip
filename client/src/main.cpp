#include "ClientComponent.h"

int main()
{
	//Config config = loadConfig();
    //Client client(config);


	ClientComponent clientComponent;
	clientComponent.createClientSubsystem();
	clientComponent.startClientSubsystem();
	



	/* MainLoop: 
	 * - setup socket
	 * - connect to server
	 * - client main loop
	 * - close connection
	 * */
    return 0;
}
