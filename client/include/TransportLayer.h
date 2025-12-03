/*
 * TransportLayer.h
 *
 * Purpose: Responsible for handling external messages
 *
 *
 *
 */

/*
├── TransportLayer (OWNS connection + processing)
│   ├── ConnectionManager (OWNS socket/network)
│   └── XOREncryptor
*/

#ifndef TRANSPORT_LAYER_H
#define TRANSPORT_LAYER_H

#include "MessageHandler.h"
#include "Serializer.h"
#include "Encryptor.h"
#include "Encoder.h"
#include "ApiManager.h"

#include <vector>
#include <memory>


#include <string>
#include <winsock2.h>

#include <functional>

//#include <ws2tcpip.h>

typedef enum
{
	//NONE,
	TCP,
	//UDP,
	//HTTPS,
	//DNS
} TransportLayerType; // TODO: rename to TransportType?

class TransportLayer
{
public:
	using ReceiveCallback = std::function<void(const RawByteBuffer&)>;
	//void setReceiveCallback(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }
	void setOnMessage(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }

	// constructor
	//TransportLayer();
	//TransportLayer(ApiManager* apiManager);
	//TransportLayer(Encryptor* encryptor);

	// deconstructor
	virtual ~TransportLayer() = default;

	// TODO: might have to change parameters after adding other transportlayers besides tcp
	virtual bool connect() = 0;
	virtual bool send(const RawByteBuffer& msg) = 0;
	virtual void receive() = 0;

	// Getter to see if we've connected to the server
	bool isConnected() { return connected_; }
	//virtual void disconnect() = 0;

	//bool sendMessage(const RawByteBuffer& msg); // TODO: conversion shouldn't be happening here: 

	//void beacon();

protected:
	// default subsystems
	//Encryptor* encryptor_;
	
	//ApiManager* apiManager_;
	ReceiveCallback receiveCallback_;
	bool connected_ = false;

private:
	/*
	std::string server;
	uint16_t port;
	SOCKET socket = INVALID_SOCKET;
	bool connected = false;
	*/


};



/*
// Handles DNS queries/responses, tunneling
class DNSTransportLayer : public TransportLayer
{

public:
private:
};*/

typedef std::unique_ptr<TransportLayer> TransportLayerUniquePtr;

#endif
