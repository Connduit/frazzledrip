#ifndef CLIENT_H
#define CLIENT_H


#include "TransportLayerTypes.h"
#include "C2Profile.h"
#include "MessageHandler.h"
#include "TransportLayer.h"
#include "../common/logging.h"

#include <memory>


class Client
{
public:
    Client();
    //Client(std::string server, std::string port); 
    //Client(TCPTransportLayer transportLayer, std::string server, std::string port); 
    //Client(TransportLayerType transportType, const std::string& server, const std::string& port); 
    //Client(TransportLayerUniquePtr transportLayer);

	Client(const std::string& server,
		   const std::string& port,
		   TransportType transportType,
		   SerializerType serializerType,
		   EncoderType encoderType,
		   EncryptorType encryptorType);

    bool run();
private:

    TransportLayerUniquePtr transportLayer_; 
    //TCPTransportLayer transportLayer_; // NOTE: hardcode transport type?, TODO: eventually turn this into a "ConnectionManager" that can own multiple transporters
    MessageHandler messageHandler_;
	// ApiResolver apiResolver_;


    //C2Profile& config_; // TODO: eventually have 

    // TODO: Eventually have multiple transporters if i need to connect to multiple servers.
    //       This allows us to act as both cnc and p2p?
    // std::vector<std::unique_ptr<Transporter>> transporters_;

};


#endif // CLIENT_H


