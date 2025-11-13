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
│   ├── BinarySerializer
│   ├── Base64Encoder
│   └── XOREncryptor
*/

#ifndef TRANSPORT_LAYER_H
#define TRANSPORT_LAYER_H

#include "MessageHandler.h"
#include "Serializer.h"
#include "Encryptor.h"
#include "Encoder.h"

#include <vector>
#include <memory>


#include <string>
#include <winsock2.h>

//#include <ws2tcpip.h>

class MessageHandler;

typedef enum
{
	//NONE,
	TCP,
	//UDP,
	//HTTP,
	//HTTPS,
	//DNS
} TransportLayerType; // TODO: rename to TransportType?

// TODO: rename to TransportLayer
class TransportLayer
{
public:
	// constructor
	//explicit TransportLayer(MessageHandler* handler) : messageHandler_(handler) {}
	explicit TransportLayer(
				MessageHandler& handler,
				SerializerUniquePtr serializer,
				EncoderUniquePtr encoder,
				EncryptorUniquePtr encryptor,
				SerializerType serializerType,
				EncoderType encoderType,
				EncryptorType encryptorType);
	// deconstructor
	virtual ~TransportLayer() = default; // TODO: what does default do?
	virtual bool connect() = 0;
	virtual bool send(const std::vector<uint8_t>& data) = 0;
	virtual std::vector<uint8_t> receive() = 0;
	virtual bool isConnected() = 0;
	//virtual void disconnect() = 0;

	// Process InternalMessage and convert it into raw bytes beforing sending it to the server
	bool sendMessage(const InternalMessage& msg);

	// Process raw bytes from server and convert them into an InternalMessage
	InternalMessage receiveMessage();

	void beacon();

	void testMessage();
protected:
	// default subsystems
	MessageHandler& messageHandler_; // TODO: change to const?
	SerializerUniquePtr serializer_;
	EncoderUniquePtr encoder_;
	EncryptorUniquePtr encryptor_;
private:
	/*
	std::string server;
	uint16_t port;
	SOCKET socket = INVALID_SOCKET;
	bool connected = false;
	*/
	InternalMessage createHeartbeat();

	// TODO: delete?
	uint32_t generateId();

};



class TCPTransportLayer : public TransportLayer
{
public:
	//TCPTransporter(const std::string& server, uint16_t port) : server_(server), port_(port) {}
	//TCPTransporter(MessageHandler* messageHandler, const std::string& server, uint16_t port);
	explicit TCPTransportLayer(
				MessageHandler& messageHandler, 
				const std::string& server, 
				const std::string& port,
				SerializerUniquePtr serializer = nullptr,
				EncoderUniquePtr encoder = nullptr,
				EncryptorUniquePtr encryptor = nullptr,
				SerializerType serializerType = SerializerType::BINARY,
				EncoderType encoderType = EncoderType::BASE64,
				EncryptorType encryptorType = EncryptorType::XOR);

	~TCPTransportLayer();

	// Attemps to send a std::vector<uint8_t> as a raw buffer to server
	bool send(const std::vector<uint8_t>& data);

	// Creates socket and attempts to connect to it
	// Returns true on success
	bool connect();

	// Receives raw buffer from server
	// Returns raw buffer as a std::vector<uint8_t>
	std::vector<uint8_t> receive();

	// Getter to see if we've connected to the server
	bool isConnected() { return connected_; }
	
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
