#include "ComponentFactory.h"
#include "TransportLayer.h"
#include "TCPTransportLayer.h"

#include <iostream>


TransportLayer* TransportLayerFactory::create(TransportLayerType type, Encryptor* encryptor)
{
	switch (type)
	{
	case TransportLayerType::TCP:
		return new TCPTransportLayer(encryptor);
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
/*
TransportLayerUniquePtr TransportLayerFactory::create(
		MessageHandler& messageHandler, 
		const std::string& host, 
		const std::string& port, 
		TransportLayerType type,
		SerializerType serializerType,
		EncoderType encoderType,
		EncryptorType encryptorType)
{

	auto serializer = ComponentFactory::create(serializerType);
	auto encoder = ComponentFactory::create(encoderType);
	auto encryptor = ComponentFactory::create(encryptorType);

	
	switch (type)
	{
		//case NONE:
		//break;
		case TransportLayerType::TCP:
			return std::make_unique<TCPTransportLayer>(
					messageHandler, 
					host, 
					port,
					std::move(serializer),
					std::move(encoder),
					std::move(encryptor));
			   // case UDP:
			   // return std::make_unique<UDPTransportLayer>(messageHandler, host, port);
			   // case HTTP:
			   // return std::make_unique<HTTPTransportLayer>(messageHandler, host, port);
			   // case HTTPS:
			   // return std::make_unique<HTTPSTransportLayer>(messageHandler, host, port);
			   // case DNS:
			   // return std::make_unique<DNSTransportLayer>(messageHandler, host, port);
		default:
			// throw error
			//break;
			std::cout << "throw error" << std::endl;
			return nullptr;

	}
}*/
