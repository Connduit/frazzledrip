#ifndef CONFIG_H
#define CONFIG_H

#include "TransportLayer.h"
#include "MessageTypes.h"

// TODO: rename to C2Profile? 
typedef struct Config
{
        Config() :
                transportLayerType_(TransportLayerType::TCP),
                serializerType_(SerializerType::BINARY),
                encoderType_(EncoderType::BASE64),
                encryptorType_(EncryptorType::XOR),
                server_("10.0.0.48"), // NOTE: empty server uses localhost
                port_("4444")
                {}

        TransportLayerType transportLayerType_;
        SerializerType serializerType_;
        EncoderType encoderType_;
        EncryptorType encryptorType_;
        std::string server_;
        std::string port_;
       

} Config;

#endif
