/* Encoder.h */
#ifndef ENCODER_H
#define ENCODER_H

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
	virtual std::vector<uint8_t> encode(std::vector<uint8_t>& data) = 0;
private:
};


class B64Encoder : public Encoder
{
public:
	std::vector<uint8_t> encode(std::vector<uint8_t>& msg); // override
	std::vector<uint8_t> decode(std::vector<uint8_t>& msg);
private:
	static constexpr uint8_t base64_chars[] =
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		"abcdefghijklmnopqrstuvwxyz"
		"0123456789+/";
};

//class HexEncoder : public Encoder

typedef std::unique_ptr<Encoder> EncoderUniquePtr;


#endif
