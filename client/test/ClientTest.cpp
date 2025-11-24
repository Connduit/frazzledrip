#include <iomanip>
#include <sstream> // NOTE: must be included before private/protected public defines
#include <iostream>
#include <cassert>

#define private public
#define protected public


#include "ClientTest.h"
#include "Client.h"
#include "ComponentFactory.h"

#include "MessageParser.h"
#include "TransportLayer.h"
#include "TCPTransportLayer.h"
#include "utils.h"



ClientSubsystem* ClientTest::client_ = 0;

ClientTest::ClientTest() : config_()
{
	client_ = new ClientSubsystem(config_);
}

ClientTest::~ClientTest()
{
	delete client_;
}

bool ClientTest::testAll()
{


	//Client client(compressor, encryptor, encoder, serializer, config);
	//Client client(encryptor, encoder, serializer, config);

	// InternalMessage
	InternalMessage inMsg;
	inMsg.data_ = {'h', 'e', 'l', 'l', 'o', ',', 'w', 'o', 'r', 'l', 'd'};
	MessageHeader header;
	header.messageType_ = MessageType::EXECUTE_COMMAND;
	header.dataSize_ = inMsg.data_.size();
	header.messageId_ = 69;
	inMsg.header_ = header;




	// Serialize
	testSerializer(inMsg);



	// Encode
	//testEncoder(serializerOutMsg);

	// Encrypt
	//testEncrypter(encoderOutMsg);

	// TODO: Frame (if needed)

	// TODO: Raw Bytes

	// TODO: Send

	return false;
}



bool ClientTest::testSerializer(InternalMessage& msg)
{
	printInternalMessage(msg);
	auto outMsg = client_->serializer_->serialize(msg);
	printRawByteBuffer(outMsg);
	auto orgiMsg = client_->serializer_->deserialize(outMsg);
	printInternalMessage(msg);

	return false;
}

bool ClientTest::testEncoder(RawByteBuffer& msg)
{
	auto outMsg = client_->encoder_->encode(msg);
	return false;
}

bool ClientTest::testEncrypter(RawByteBuffer& msg)
{
	auto outMsg = client_->encryptor_->encrypt(msg);
	return false;
}

