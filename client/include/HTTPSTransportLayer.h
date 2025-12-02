#ifndef HTTPS_TRANSPORT_LAYER_H
#define HTTPS_TRANSPORT_LAYER_H

#include "TransportLayer.h"

class HTTPSTransportLayer : public TransportLayer
{
	HTTPSTransportLayer(ApiManager* apiManager);
	HTTPSTransportLayer(ApiManager* apiManager, std::string& server, std::string& port);
	~HTTPSTransportLayer();

	// Attemps to send a std::vector<uint8_t> as a raw buffer to server
	bool send(const RawByteBuffer& msg);

	// Creates socket and attempts to connect to it
	// Returns true on success
	bool connect();

	// Receives raw buffer from server
	// Returns raw buffer as a std::vector<uint8_t>
	void receive();

	// Getter to see if we've connected to the server
	bool isConnected() { return connected_; } // TODO: move to parent class? 

private:
	// Initializes WSADATA by being called through the constructor 
	bool initializeWinsock();

	// TODO: move this to abstract class?
	std::string server_;
	std::string port_;
	SOCKET socket_ = INVALID_SOCKET;
	bool connected_ = false;
	//std::atomic<bool> connected_;

	ApiManager* apiManager_;
};

#endif
