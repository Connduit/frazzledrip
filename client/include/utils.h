#ifndef UTILS_H
#define UTILS_H

#include "MessageTypes.h"

#include <string>

RawByteBuffer string2byte(std::string& msg)
{
	return RawByteBuffer(msg.begin(), msg.end());
}

std::string byte2string(RawByteBuffer& msg)
{
	return std::string(msg.begin(), msg.end());
}

RawByteBuffer string2byte(const std::string& msg)
{
	return RawByteBuffer(msg.begin(), msg.end());
}

std::string byte2string(const RawByteBuffer& msg)
{
	return std::string(msg.begin(), msg.end());
}

#endif
