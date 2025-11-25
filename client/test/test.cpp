#include <iostream>
#include <typeinfo>

#define private public
#define protected public

#include "test.h"
#include "Client.h"
#include "ClientTest.h"







int main()
{
	std::cout << "inside main" << std::endl;
	//ClientTest clientTest;
	//clientTest.testAll();

	ClientSubsystem client;
	//client.messageHandler_;
	//client.transportLayerPtr_;
	std::cout << "===================" << std::endl;
	std::cout << typeid(client.messageParser_).name() << std::endl;
	//BinarySerializer bs;
	//std::cout << typeid(bs).name() << std::endl;
	//std::cout << typeid(obj).name() << std::endl;
	std::cout << "===================" << std::endl;


        std::cout << "inside main" << std::endl;

        //Client client;

        // Test components are good
        if (client.messageParser_->serializer_)
        {
            std::cout << "Serializer type: "
                << typeid(*client.messageParser_->serializer_).name() << std::endl;

            // Test method call
            InternalMessage testMsg;
            testMsg.header.messageType = MessageType::NONE;
            testMsg.header.dataSize = 0;
            auto data = client.messageParser_->serializer_->serialize(testMsg);
            std::cout << "Pre-run test: GOOD" << std::endl;
        }

        std::cout << "About to call client.run()..." << std::endl;

        // Add canary check right before run
        // (You'd need to add this method to TransportLayer)
        // if (!client.transportLayerPtr_->checkCanary()) {
        //     std::cout << "Memory corrupted before run()!" << std::endl;
        //     return 1;
        // }

        client.run();
        std::cout << "client.run() completed" << std::endl;

}
