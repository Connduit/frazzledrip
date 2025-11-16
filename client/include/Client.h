#ifndef CLIENT_H
#define CLIENT_H


//#include "TransportLayerTypes.h"
// #include "C2Profile.h"
// #include "MessageHandler.h"
// #include "TransportLayer.h"
// #include "ComponentFactory.h"
// //#include "../common/logging.h"
//
// #include <memory>
//
//
// class Client
// {
// public:
//     Client(); // TODO: remove
//
//     [>
//     Client(const std::string& server,
//            const std::string& port,
//            TransportLayerType transportType,
//            SerializerType serializerType,
//            EncoderType encoderType,
//            EncryptorType encryptorType);*/
//
//     ~Client();
//
//     bool run();
//     bool setupComponents(); // change to bool?
//     bool setupTasks(); // change to bool?
//     bool configure(); // change return type to a custom status type?
// private:
//     MessageHandler* messageHandler_;
//     TransportLayer* transportLayer_; // TODO: convert to std::vector<TransportLayerUniquePtr> transporters_;
//     //TransportLayerUniquePtr transportLayerPtr_; // TODO: convert to std::vector<TransportLayerUniquePtr> transporters_;
//
//     //MessageHandler* messageHandlerPtr_;
//     //TransportLayerFactory* transportLayerFactory_;
//     //TransportLayer* transportLayer_;
//
//
//
//     // TODO: eventually turn this into a "ConnectionManager" that can own multiple transporters
//
//     // ApiResolver apiResolver_; // TODO: should be a parameter for MessageHandler constructor?
//
//     //C2Profile config_; // TODO: eventually have
//
//     // TODO: Eventually have multiple transporters if i need to connect to multiple servers.
//     //       This allows us to act as both cnc and p2p?
//     // std::vector<std::unique_ptr<Transporter>> transporters_;
//
// };
//
//
// #endif // CLIENT_H
//
//
