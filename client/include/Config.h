#ifndef CONFIG_H
#define CONFIG_H

#include "TransportLayer.h"
#include "MessageTypes.h"

typedef struct Config
{
        Config() :
                transportLayerType_(TransportLayerType::TCP),
                serializerType_(SerializerType::BINARY),
                encoderType_(EncoderType::BASE64),
                encryptorType_(EncryptorType::XOR)
                {}

        // int server_;
        // int port_;
        TransportLayerType transportLayerType_;
        SerializerType serializerType_;
        EncoderType encoderType_;
        EncryptorType encryptorType_;
       

} Config;

#endif
