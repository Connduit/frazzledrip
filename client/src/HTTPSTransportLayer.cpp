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
