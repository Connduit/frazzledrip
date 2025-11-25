
/* ComponentFactory.h */

// TODO: factory should take in a config file? either from C2Profile
// or maybe even directly from a provided json file? 


#ifndef COMPONENT_FACTORY_H
#define COMPONENT_FACTORY_H


//#include "Compressor.h"
#include "Encoder.h"
#include "Encryptor.h"
#include "Serializer.h"
//#include "TransportLayerTypes.h"
#include "TransportLayer.h"
#include "TCPTransportLayer.h"


#include <memory>

class ComponentFactory
{
public:
	static Serializer* create(SerializerType type); // TODO: mark as explicit
    static Encoder* create(EncoderType type); // TODO: mark as explicit
    static Encryptor* create(EncryptorType type); // TODO: mark as explicit
private:
};

class TransportLayerFactory
{
public:
	static TransportLayer* create(TransportLayerType type);
	static TransportLayer* create(TransportLayerType type, std::string& server, std::string& port);
	//static TransportLayer* create(TransportLayerType type, Encryptor* encryptor);
private:
};



#endif


