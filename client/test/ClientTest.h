#ifndef CLIENT_TEST_H
#define CLIENT_TEST_H

#include "Compressor.h"
#include "Encoder.h"
#include "Encryptor.h"
#include "Serializer.h"

#include "Config.h"

#include "ClientSubsystem.h"
#include "MessageTypes.h"

#include <string>

class ClientTest
{
public:
	ClientTest();
	~ClientTest();

	bool testSerializer(InternalMessage& msg);
	bool testEncoder(RawByteBuffer& msg);
	bool testEncrypter(RawByteBuffer& msg);
	bool testAll();


protected:
	static ClientSubsystem* client_;
private:
	Config config_;

};

#endif

