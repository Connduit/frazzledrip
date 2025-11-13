#include "Client.h"

int main()
{
	//Config config = loadConfig();
    //Client client(config);


	Client client;
	client.run();
	


    //const char server_ip[] = "127.0.0.1";
    //client.run(server_ip, "4444"); // TODO: add/replace with C2Profile Config param?

	/* MainLoop: 
	 * - setup socket
	 * - connect to server
	 * - client main loop
	 * - close connection
	 * */
    return 0;
}
