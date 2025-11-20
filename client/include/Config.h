// config.h 

typedef struct Config
{
        Config() :
transportLayerType_(TransportLayerType::TCP),
                serializerType_(SerializerType::BINARY),
                encoderType_(EncoderType::BASE64),
encryptorType_(EncryptorType::XOR)
        {}

        TransportLayerType transportLayerType_;
SerializerType serializerType_;
EncoderType encoderType_;
EncryptorType encryptorType_;
       

} Config;