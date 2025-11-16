/* Encoder.h */
#ifndef ENCODER_H
#define ENCODER_H

#include "MessageTypes.h"

#include <string>
#include <cstdint>
#include <memory>
#include <vector>



typedef enum
{
	BASE64,
	HEX
} EncoderType;


class Encoder
{
public:
	//TODO: match function signature
	//virtual bool encode() = 0;
	//virtual bool decode() = 0;
	virtual ~Encoder() = default;
	virtual RawByteBuffer encode(RawByteBuffer& data) = 0;
	virtual RawByteBuffer decode(RawByteBuffer& data) = 0;
private:
};


class Base64Encoder : public Encoder
{
public:
	RawByteBuffer encode(RawByteBuffer& msg); // override
	RawByteBuffer decode(RawByteBuffer& msg);
private:
	static constexpr uint8_t base64_chars[] =
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		"abcdefghijklmnopqrstuvwxyz"
		"0123456789+/";
};

//class HexEncoder : public Encoder

typedef std::unique_ptr<Encoder> EncoderUniquePtr;


#endif
