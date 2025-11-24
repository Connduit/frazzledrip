#ifndef TCP_TRANSPORT_LAYER_H
#define TCP_TRANSPORT_LAYER_H

#include "TransportLayer.h"


#include <thread>
#include <atomic>

class TCPTransportLayer : public TransportLayer
{
public:
	TCPTransportLayer();
	//TCPTransportLayer(Encryptor* encryptor);
	~TCPTransportLayer();

	// Attemps to send a std::vector<uint8_t> as a raw buffer to server
	bool send(const RawByteBuffer& data);

	// Creates socket and attempts to connect to it
	// Returns true on success
	bool connect();

	// Receives raw buffer from server
	// Returns raw buffer as a std::vector<uint8_t>
	RawByteBuffer receive();

	// Getter to see if we've connected to the server
	bool isConnected() { return connected_; }

	void run();

	// call this function manually in main loop instead of using threads
	// void update();

private:
	// Initializes WSADATA by being called through the constructor 
	bool initializeWinsock();


	// TODO: move this to abstract class?
	std::string server_;
	std::string port_;
	//uint16_t port_;
	//std::string port_;
	SOCKET socket_ = INVALID_SOCKET;
	bool connected_ = false;
	//std::atomic<bool> connected_;


};

#endif
