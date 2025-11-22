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

#include "MessageParser.h"
#include "Serializer.h"
#include "Encryptor.h"
#include "Encoder.h"

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
	//HTTP,
	//HTTPS,
	//DNS
} TransportLayerType; // TODO: rename to TransportType?

class TransportLayer
{
public:
	using ReceiveCallback = std::function<void(const RawByteBuffer&)>;
	//void setReceiveCallback(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }
	void setOnMessage(ReceiveCallback cb) { receiveCallback_ = std::move(cb); }
	//virtual void update() = 0;

	// constructor
	TransportLayer(Encryptor* encryptor);

	// deconstructor
	virtual ~TransportLayer() = default;
	virtual bool connect() = 0; // wrapper for ::connect? 
	virtual bool send(const RawByteBuffer& data) = 0;
	virtual RawByteBuffer receive() = 0;
	virtual bool isConnected() = 0;
	//virtual void disconnect() = 0;

	virtual void run() = 0;

	// Process InternalMessage and convert it into raw bytes beforing sending it to the server
	// wrapper for ::send ? 
	bool sendMessage(const RawByteBuffer& msg); // TODO: conversion shouldn't be happening here: 

	// Process raw bytes from server and convert them into an InternalMessage
	RawByteBuffer receiveMessage(); // wrapper for ::recv? 

	void beacon();

protected:
	// default subsystems
	Encryptor* encryptor_;
	ReceiveCallback receiveCallback_;

private:
	/*
	std::string server;
	uint16_t port;
	SOCKET socket = INVALID_SOCKET;
	bool connected = false;
	*/


};




/*

// Handles HTTPS headers, TLS, cookies, etc.
class HTTPSTransportLayer : public TransportLayer
{

public:
private:
};

// Handles DNS queries/responses, tunneling
class DNSTransportLayer : public TransportLayer
{

public:
private:
};*/

typedef std::unique_ptr<TransportLayer> TransportLayerUniquePtr;

#endif
