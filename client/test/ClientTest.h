#ifndef CLIENT_TEST_H
#define CLIENT_TEST_H

#include "C2Profile.h"
#include "Compressor.h"
#include "Encoder.h"
#include "Encryptor.h"
#include "Serializer.h"

#include "ClientSubsystem.h"
#include "MessageTypes.h"

#include <string>

class ClientTest
{
public:
	ClientTest();
	~ClientTest();

	void testTransportLayer();
	void testFactoryDirectly();
	bool testSerializer(InternalMessage& inMsg, std::vector<uint8_t>& outMsg, BinarySerializer& serializer);
	bool testEncoder(std::vector<uint8_t>& inMsg, std::vector<uint8_t>& outMsg, Base64Encoder& encoder);
	bool testEncrypter(std::vector<uint8_t>& inMsg, std::vector<uint8_t>& outMsg, XorEncryptor& encryptor);
	bool testAll();


protected:
	static ClientSubsystem* client_;
private:
	std::vector<uint8_t> string2byte(const std::string inMsg);
	std::string byte2string(const std::vector<uint8_t> inMsg);
};

#endif
