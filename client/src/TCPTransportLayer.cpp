#include "TCPTransportLayer.h"

#include <winsock2.h>
#include <ws2tcpip.h>
#include <iomanip>  // for std::hex, std::setw, std::setfill
#include <iostream>

//TCPTransportLayer::TCPTransportLayer(MessageHandler* hdlr, const std::string& server, uint16_t port)
// TCPTransportLayer::TCPTransportLayer(
//     MessageHandler& messageHandler,
//     const std::string& server,
//     const std::string& port,
//     SerializerUniquePtr serializer,
//     EncoderUniquePtr encoder,
//     EncryptorUniquePtr encryptor,
//     SerializerType serializerType,
//     EncoderType encoderType,
//     EncryptorType encryptorType)
//     :
//     TransportLayer(
//         messageHandler,
//         std::move(serializer),
//         std::move(encoder),
//         std::move(encryptor),
//         serializerType,
//         encoderType,
//         encryptorType),
//     server_(server),
//     port_(port)
// {
//     initializeWinsock();
// }


/*
TCPTransportLayer::TCPTransportLayer() :
    TransportLayer(messageHandler_)
{
    initializeWinsock();
}

TCPTransportLayer::TCPTransportLayer(MessageHandler* messageHandler) :
    TransportLayer(messageHandler_)
{
    initializeWinsock();
}

TCPTransportLayer::TCPTransportLayer(
    MessageHandler* messageHandler,
    Encryptor* encryptor) 
    :
    TransportLayer(messageHandler_)
{
    initializeWinsock();
}*/

// TODO: this constructor shouldn't ever be hit? i also think it's just wrong
TCPTransportLayer::TCPTransportLayer() :
    TransportLayer(encryptor_)
{
    //connected_ = false;
    initializeWinsock();
}

TCPTransportLayer::TCPTransportLayer(Encryptor* encryptor) :
    TransportLayer(encryptor)
{
    //connected_ = false;
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



bool TCPTransportLayer::send(const RawByteBuffer& data)
{
    /*
    if (!connected_ && !connect()) return false;

    return ::send(socket_, (const char*)data.data(), data.size(), 0) > 0;
    */

    if (!connected_ && !connect())
    {
        std::cout << "Not connected and connect failed" << std::endl;
        return false;
    }

    int result = ::send(socket_, (const char*)data.data(), data.size(), 0);

    return result > 0;

}

bool TCPTransportLayer::connect()
{
    socket_ = ::socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (socket_ == INVALID_SOCKET)
    {
        std::cout << "socket invalid" << std::endl;
        return false;
    }



    ADDRINFOA hints, * result = nullptr;
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

RawByteBuffer TCPTransportLayer::receive()
{
    if (!connected_) return {};

    RawByteBuffer buffer(4096);
    int received = recv(socket_, (char*)buffer.data(), buffer.size(), 0);
    if (received <= 0)
    {
        connected_ = false;
        return {};
    }
    buffer.resize(received);
    return buffer;
}

void TCPTransportLayer::run()
{
    while (connected_)
    {
        RawByteBuffer buffer(4096);
        int received = recv(socket_, (char*)buffer.data(), buffer.size(), 0);

        if (received <= 0)
        {
            connected_ = false;
            break;
        }

        buffer.resize(received);

        if (receiveCallback_)
        {
            receiveCallback_(buffer);   // 🔥 Fire callback!
        }
    }
}

void TCPTransportLayer::update()
{
	if (!connected_)
		return;

	u_long bytesAvailable = 0;
	ioctlsocket(socket_, FIONREAD, &bytesAvailable);

	if (bytesAvailable == 0)
		return;

	RawByteBuffer buffer(bytesAvailable);
	int received = recv(socket_, (char*)buffer.data(), buffer.size(), 0);

	if (received <= 0)
	{
		connected_ = false;
		return;
	}

	buffer.resize(received);

	if (receiveCallback_)
		receiveCallback_(buffer);
}
