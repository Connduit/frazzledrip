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
    HINTERNET hConnect = WinHttpConnect(hSession, L"example.com", INTERNET_DEfAULT_HTTPS_PORT, 0);
    return true;
}
