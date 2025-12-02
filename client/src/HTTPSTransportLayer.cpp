#include "HTTPSTransportLayer.h"

bool HTTPSTransportLayer::connect()
{
    HINTERNET hConnect = WinHttpConnect(hSession, L"example.com", INTERNET_DEfAULT_HTTPS_PORT, 0);
    return true;
}
