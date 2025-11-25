#ifndef HTTPS_TRANSPORT_LAYER_H
#define HTTPS_TRANSPORT_LAYER_H

#include "TransportLayer.h"

class HTTPSTransportLayer : public TransportLayer
{
public:
private:
	// encryptor object for https tls encryption (not for encrypting the InternalMessage)
	Encryptor* encryptor_;
};

#endif
