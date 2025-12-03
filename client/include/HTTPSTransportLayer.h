#ifndef HTTPS_TRANSPORT_LAYER_H
#define HTTPS_TRANSPORT_LAYER_H

#include "TransportLayer.h"
#include "ApiManager.h"

#include <windows.h>
#include <winhttp.h>

#pragma comment(lib, "winhttp.lib")

class HTTPSTransportLayer : public TransportLayer
{
public:
	HTTPSTransportLayer(ApiManager* apiManager);
	HTTPSTransportLayer(ApiManager* apiManager, std::string& server, std::string& port);
	~HTTPSTransportLayer();

	bool connect();
	bool send(const RawByteBuffer& msg);
	void receive();
private:
	ApiManager* apiManager_;
	HINTERNET hSession_ = nullptr;
	HINTERNET hConnect_ = nullptr;
	HINTERNET hRequest_ = nullptr;

	std::string server_; // TODO: change to wstring
	std::string path_; // TODO: change to wstring
	std::string port_; // TODO: change to INTERNET_PORT? 
	bool connected_ = false;
	//std::atomic<bool> connected_;

};

#endif
