/* Encoder.cpp */
/* encode json into base64 */
/* encoding is needed for http, smtp, json, xml, html, urls,  (but not tcp) */
/* tcp/udp, smb, dns, ftp, and other binary traffic protocols don't need encoding */


/* define differnt encoding formats in here... start with TLV */


// TODO:
// char table eventually obfuscated by maybe building the table dynamically at runtime,
// splitting the table into separate parts and then reconstructing them at runtime,
// encrypting the table, or mathematically generating the characters for the table?
// https://github.com/Adaptix-Framework/AdaptixC2/blob/main/Extenders/agent_beacon/src_beacon/beacon/Encoders.cpp

#include <Encoder.h>
// base64 or hex encoding
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

//#include <iostream>

//const char b64_table[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"; // default_table


// if i want to rewrite this with the input length unknown, will have to use strlen... 
// hopefully by then ill have my own strlen func written
//char* B64Encoder::encode(const unsigned char* data, size_t input_len)
RawByteBuffer Base64Encoder::encode(RawByteBuffer& msg)
{


	RawByteBuffer encoded;
	size_t i = 0;

	while (i < msg.size())
	{
		uint32_t octet_a = i < msg.size() ? msg[i++] : 0;
		uint32_t octet_b = i < msg.size() ? msg[i++] : 0;
		uint32_t octet_c = i < msg.size() ? msg[i++] : 0;

		uint32_t triple = (octet_a << 0x10) + (octet_b << 0x08) + octet_c;

		encoded.push_back(base64_chars[(triple >> 3 * 6) & 0x3F]);
		encoded.push_back(base64_chars[(triple >> 2 * 6) & 0x3F]);
		encoded.push_back(base64_chars[(triple >> 1 * 6) & 0x3F]);
		encoded.push_back(base64_chars[(triple >> 0 * 6) & 0x3F]);
	}

	// Add padding
	size_t pad_count = msg.size() % 3;
	if (pad_count > 0)
	{
		for (size_t j = 0; j < 3 - pad_count; j++)
		{
			encoded[encoded.size() - 1 - j] = '=';
		}
	}
	return encoded;
}



/* decode.h */
// break decoding across different functions to make it harder to detect?
// maybe i don't have to worry about b64 being easy to detect cuz
// it's ok for av to know we're sending a message? just have to make 
// the message/sending look like a normal program

// https://github.com/Adaptix-Framework/AdaptixC2/blob/main/Extenders/agent_beacon/src_beacon/beacon/Encoders.cpp

//#include <stdint.h>

// Simple LUT as a string (very obvious in binary)
//const char b64_table[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"; // default_table

// TODO: pretty code water way of encoding
RawByteBuffer Base64Encoder::decode(RawByteBuffer& msg)
{
	RawByteBuffer decoded;
	size_t i = 0;

	// Build reverse lookup table
	uint8_t base64_lookup[256] = { 0 };
	for (uint8_t j = 0; j < 64; j++)
	{
		base64_lookup[base64_chars[j]] = j;
	}

	while (i < msg.size() && msg[i] != '=')
	{
		uint8_t char1 = msg[i++];
		uint8_t char2 = i < msg.size() ? msg[i++] : 'A';
		uint8_t char3 = i < msg.size() ? msg[i++] : 'A';
		uint8_t char4 = i < msg.size() ? msg[i++] : 'A';

		if (char1 == '=' || char2 == '=') break;

		uint32_t triple = (base64_lookup[char1] << 18) +
			(base64_lookup[char2] << 12) +
			(base64_lookup[char3] << 6) +
			base64_lookup[char4];

		decoded.push_back((triple >> 16) & 0xFF);
		if (char3 != '=') decoded.push_back((triple >> 8) & 0xFF);
		if (char4 != '=') decoded.push_back(triple & 0xFF);
	}

	return decoded;
}

