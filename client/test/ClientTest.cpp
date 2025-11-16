#include <iomanip>
#include <sstream> // NOTE: must be included before private/protected public defines
#include <iostream>
#include <cassert>

#define private public
#define protected public


#include "test.h" // TODO: not getting the includes from here for some reason?
#include "ClientTest.h"
#include "Client.h"
#include "ComponentFactory.h"

#include "MessageHandler.h"
#include "TransportLayer.h"
#include "TCPTransportLayer.h"



Client* ClientTest::client_ = 0;

ClientTest::ClientTest()
{
	client_ = new Client();
}

ClientTest::~ClientTest()
{
	delete client_;
}

bool ClientTest::testAll()
{

	//testFactoryDirectly();
	testTransportLayer();
	std::cout << "================================" << std::endl;

	XorEncryptor encryptor("testKey");
	Base64Encoder encoder;
	BinarySerializer serializer;
	C2Profile config;

	//Client client(compressor, encryptor, encoder, serializer, config);
	//Client client(encryptor, encoder, serializer, config);
	Client client;

	// InternalMessage
	InternalMessage inMsg;
	inMsg.data = internalMessageData;
	MessageHeader header;
	header.messageType = MessageType::EXECUTE_COMMAND;
	header.dataSize = inMsg.data.size();
	header.messageId = 69;
	inMsg.header = header;


	/*
	printVector(inMsg.data);
	std::cout << inMsg.header.messageId << std::endl;
	std::cout << inMsg.header.messageType << std::endl;
	std::cout << inMsg.header.dataSize << std::endl;
	*/

	// TODO: Validate?

	// Serialize
	std::vector<uint8_t> serializerOutMsg;
	testSerializer(inMsg, serializerOutMsg, serializer);
	std::cout << "Serializer: " << byte2string(serializerOutMsg) << std::endl;
	//printVector(serializerOutMsg);
	//std::cout << std::endl;

	/*
	InternalMessage outMsg = serializer.deserialize(serializerOutMsg);
	printVector(outMsg.data);
	std::cout << outMsg.header.messageId << std::endl;
	std::cout << outMsg.header.messageType << std::endl;
	std::cout << outMsg.header.dataSize << std::endl;
	std::cout << std::endl;
	*/

	// Encode
	std::vector<uint8_t> encoderOutMsg;
	testEncoder(serializerOutMsg, encoderOutMsg, encoder);
	/*
	std::cout << "Encoder" << std::endl;
	printVector(encoderOutMsg);
	std::cout << std::endl;
	*/

	// TODO: Compress? 

	std::vector<uint8_t> encryptOutMsg;
	testEncrypter(encoderOutMsg, encryptOutMsg, encryptor);
	/*
	std::cout << "Encryptor" << std::endl;
	printVector(encryptOutMsg);
	*/
	std::cout << std::endl;

	// TODO: Frame (if needed)

	// TODO: Raw Bytes

	// TODO: Send

	return false;
}

void ClientTest::testTransportLayer()
{
	std::cout << "=== Testing TransportLayer ===" << std::endl;

	// Create a MessageHandler that definitely stays alive
	MessageHandler handler;

	// Create transport layer through factory
	auto transport = TransportLayerFactory::create(
		handler, "localhost", "4444",
		TransportLayerType::TCP,
		SerializerType::BINARY,
		EncoderType::BASE64,
		EncryptorType::XOR
	);

	if (transport)
	{
		std::cout << "Transport layer created successfully" << std::endl;

		InternalMessage testMsg;
		testMsg.header.messageType = MessageType::NONE;
		testMsg.header.dataSize = 0;

		std::cout << "About to call sendMessage..." << std::endl;
		transport->sendMessage(testMsg);  // Does it crash here?
		std::cout << "sendMessage completed" << std::endl;
	}
}

// Simple test - does the factory work in isolation?
void ClientTest::testFactoryDirectly()
{
	std::cout << "Testing factory directly..." << std::endl;

	// Test 1: Create serializer directly
	std::cout << "Creating serializer directly..." << std::endl;
	BinarySerializer directSerializer;
	InternalMessage testMsg;
	testMsg.header.messageType = MessageType::NONE;
	testMsg.header.dataSize = 0;
	auto result1 = directSerializer.serialize(testMsg);
	std::cout << "Direct creation: SUCCESS" << std::endl;

	// Test 2: Create through factory  
	std::cout << "Creating serializer through factory..." << std::endl;
	auto factorySerializer = ComponentFactory::create(SerializerType::BINARY);
	std::cout << "Factory returned: " << (factorySerializer != nullptr) << std::endl;
	if (factorySerializer)
	{
		std::cout << "About to call serialize on factory serializer..." << std::endl;
		auto result2 = factorySerializer->serialize(testMsg);  // Does it crash here?
		std::cout << "Factory creation: SUCCESS" << std::endl;
	}
}


bool ClientTest::testSerializer(InternalMessage& inMsg, std::vector<uint8_t>& outMsg, BinarySerializer& serializer)
{
	outMsg = serializer.serialize(inMsg);
	return false;
}

bool ClientTest::testEncoder(std::vector<uint8_t>& inMsg, std::vector<uint8_t>& outMsg, Base64Encoder& encoder)
//bool ClientTest::testEncoder(std::vector<uint8_t>& inMsg, std::string& outMsg , B64Encoder& encoder)
{
	outMsg = encoder.encode(inMsg);
	return false;
}

bool ClientTest::testEncrypter(std::vector<uint8_t>& inMsg, std::vector<uint8_t>& outMsg, XorEncryptor& encryptor)
//bool ClientTest::testEncrypter(std::string& inMsg, std::vector<uint8_t>& outMsg, XorEncryptor& encryptor)
{
	outMsg = encryptor.encrypt(inMsg);
	return false;

}

std::vector<uint8_t> ClientTest::string2byte(const std::string inMsg)
{
	return std::vector<uint8_t>(inMsg.begin(), inMsg.end());
}

std::string ClientTest::byte2string(const std::vector<uint8_t> inMsg)
{
	//return std::string(inMsg.begin(), inMsg.end());
	std::ostringstream oss;
	oss << std::uppercase;

	for (uint8_t b : inMsg)
	{
		oss << std::hex << std::setw(2) << std::setfill('0') << (int)b << " ";		
	}
	return oss.str();
}
