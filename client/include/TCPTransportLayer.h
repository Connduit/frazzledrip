#ifndef TCP_TRANSPORT_LAYER_H
#define TCP_TRANSPORT_LAYER_H

#include "TransportLayer.h"

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

#endif