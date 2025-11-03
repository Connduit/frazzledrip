#define private public
#define protected public


#include "test.h" // TODO: not getting the includes from here for some reason?
#include "ClientTest.h"
#include "Client.h"

#include "MessageHandler.h"
#include "TransportLayer.h"


#include <cassert>
#include <iostream>


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
	XorEncryptor encryptor("testKey");
	B64Encoder encoder;
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
	//std::cout << "Serializer" << std::endl;
	printVector(serializerOutMsg);
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





bool ClientTest::testSerializer(InternalMessage& inMsg, std::vector<uint8_t>& outMsg, BinarySerializer& serializer)
{
	outMsg = serializer.serialize(inMsg);
	return false;
}

bool ClientTest::testEncoder(std::vector<uint8_t>& inMsg, std::vector<uint8_t>& outMsg, B64Encoder& encoder)
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
	return std::string(inMsg.begin(), inMsg.end());
}

void ClientTest::printVector(const std::vector<uint8_t>& vec)
{
	for (auto b : vec)
	{
		std::cout << static_cast<int>(b) << " ";
	}
	std::cout << std::endl;
}
