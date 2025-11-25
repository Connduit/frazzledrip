#include "ComponentFactory.h"
#include "TransportLayer.h"
#include "TCPTransportLayer.h"

#include <iostream>

TransportLayer* TransportLayerFactory::create(TransportLayerType type)
{
	switch (type)
	{
	case TransportLayerType::TCP:
		return new TCPTransportLayer();
	default:
		// throw error
		//break;
		std::cout << "throw error" << std::endl;
		return nullptr;
	}
}

TransportLayer* TransportLayerFactory::create(
	TransportLayerType type, 
	std::string& server, 
	std::string& port)
{
	switch (type)
	{
	case TransportLayerType::TCP:
		return new TCPTransportLayer(server, port);
	default:
		// throw error
		//break;
		std::cout << "throw error" << std::endl;
		return nullptr;
	}
}

Serializer* ComponentFactory::create(SerializerType type)
{
	switch (type)
	{
	case SerializerType::BINARY:
		return new BinarySerializer();
	case SerializerType::JSON:
		// return std::make_unique<JSON_TODO>();
	default:
		return nullptr;
	}
}

Encoder* ComponentFactory::create(EncoderType type)
{
	switch (type)
	{
	case EncoderType::BASE64:
		return new Base64Encoder();
	case EncoderType::HEX:
		// return std::make_unique<HEX_TODO>();
	default:
		return nullptr;
	}
}

Encryptor* ComponentFactory::create(EncryptorType type)
{
	switch (type)
	{
	case EncryptorType::XOR:
		return new XorEncryptor();
	case EncryptorType::AES:
		// return std::make_unique<AES_TODO>();
	default:
		return nullptr;
	}
}
