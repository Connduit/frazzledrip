#include "HTTPSTransportLayer.h"

HTTPSTransportLayer::HTTPSTransportLayer(
    ApiManager* apiManager)
    :
    apiManager_(apiManager),
    TransportLayer()
{
    // TODO: initialize... WinHttpOpen? 
}

HTTPSTransportLayer::HTTPSTransportLayer(
    ApiManager* apiManager, 
    std::string& server, 
    std::string& port)
    : 
    apiManager_(apiManager),
    server_(server),
    port_(port)
{
    // TODO: initialize... WinHttpOpen?
}

HTTPSTransportLayer::~HTTPSTransportLayer()
{
    // TODO: WinHttpCloseHandle(TODO_MY_HANDLE);
}

bool HTTPSTransportLayer::connect()
{
    hSession_ = WinHttpOpen(L"AgentName/1.0", WINHTTP_ACCESS_TYPE_AUTOMATIC_PROXY, WINHTTP_NO_PROXY_NAME, WINHTTP_NO_PROXY_BYPASS, 0); // use WINHTTP_ACCESS_TYPE_AUTOMATIC_PROXY on win8.1 and newer; otherwise, WINHTTP_ACCESS_TYPE_DEFAULT_PROXY
    if (!hSession_)
    {
        return false;
    }
    
    //HINTERNET hConnect = WinHttpConnect(hSession, L"example.com", INTERNET_DEfAULT_HTTPS_PORT, 0);
    hConnect_ = WinHttpConnect(hSession_, L"example.com", INTERNET_DEFAULT_PORT, 0); // INTERNET_DEFAULT_PORT is 80 for http and 443 for https

    if (!hConnect_)
    {
        return false;
    }

    connected_ = true;
    return true;
}

bool HTTPSTransportLayer::send(const RawByteBuffer& msg)
{
    if (!connected_)
    {
        std::cout << "Cannot send data over HTTPS... not connected..." << std::endl;
        return false;
    }

    hRequest_ = WinHttpOpenRequest(
        hConnect_,
        L"POST",
        path_.c_str(),
        NULL, // HTTP/1.1
        WINHTTP_NO_REFERER,
        WINHTTP_DEFAULT_ACCEPT_TYPES,
        WINHTTP_FLAG_SECURE); // use HTTPS... SSL/TLS

    if (!hRequest_)
    {
        return false;
    }

    BOOL result = WinHttpSendRequest(
        hRequest_,
        WINHTTP_NO_ADDITIONAL_HEADERS,
        0,
        (LPVOID)msg.data(), // data to send using POST or PUT operations?
        (DWORD)msg.size(),
        (DWORD)msg.size(),
        0);

    if (!result)
    {
        return false;
    }

    return WinHttpReceiveResponse(hRequest_, NULL);
    
}

void HTTPSTransportLayer::receive()
{
    if (!hRequest_ || !receiveCallback_)
    {
        return;
    }

    DWORD bytesAvailable = 0;
    
    while (WinHttpQueryDataAvailable(hRequest_, &bytesAvailable) && bytesAvailable > 0)
    {
        RawByteBuffer buffer(bytesAvailable);
        DWORD bytesRead = 0;

        if (!WinHttpReadData(hRequest_, buffer.data(), bytesAvailable, &bytesRead))
        {
            break;
        }

        buffer.resize(bytesRead);

        // TODO: need a way to handle muliple chunks and putting them back together if the enitre message isn't sent in a single chunk
        // possible solutions: length-prefixed framing? 
        receiveCallback_(buffer); 
    }
}
