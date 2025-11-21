/* Serializer.h */

#ifndef SERIALIZER_H
#define SERIALIZER_H

#include "MessageTypes.h"

#include <vector>
#include <memory>


typedef enum
{
    BINARY,
    JSON
} SerializerType;


class Serializer
{
public:
	virtual ~Serializer() = default;
	virtual RawByteBuffer serialize(const InternalMessage& msg) = 0;
	virtual InternalMessage deserialize(const RawByteBuffer& data) = 0;
private:
};


class BinarySerializer : public Serializer
{
public:
    RawByteBuffer serialize(const InternalMessage& msg) override;
	InternalMessage deserialize(const RawByteBuffer& msg) override;
private:
	// NOTE: might be useful later? 
    uint16_t calculate_checksum(const std::vector<uint8_t>& data)
    {
        uint16_t checksum = 0;
        for (uint8_t byte : data)
        {
            checksum += byte;
        }
        return checksum;
    }
};

/*
class JsonSerializer : public Serializer
{
public:
	bool serialize() { return true; };
	bool deserialize() { return true; };
private:
};
*/

// TODO: custom Serializer class
typedef std::unique_ptr<Serializer> SerializerUniquePtr;

#endif
