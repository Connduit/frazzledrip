#ifndef UTILS_H
#define UTILS_H

#include "MessageTypes.h"

#include <string>
#include <iostream>

inline RawByteBuffer string2byte(std::string& msg)
{
	return RawByteBuffer(msg.begin(), msg.end());
}

inline std::string byte2string(RawByteBuffer& msg)
{
	return std::string(msg.begin(), msg.end());
}

inline RawByteBuffer string2byte(const std::string& msg)
{
	return RawByteBuffer(msg.begin(), msg.end());
}

inline std::string byte2string(const RawByteBuffer& msg)
{
	return std::string(msg.begin(), msg.end());
}

inline void printRawByteBuffer(RawByteBuffer& msg)
{
	for (int c : msg)
	{
		std::cout << c << " ";
	}
	std::cout << std::endl;
}

inline void printInternalMessage(InternalMessage& msg)
{
	std::cout << "messageId_ = " << msg.header_.messageId_ << std::endl;
	std::cout << "messageType_ = " << msg.header_.messageType_ << std::endl;
	std::cout << "dataSize_ = " << msg.header_.dataSize_ << std::endl;
	std::cout << "data_ = " << byte2string(msg.data_) << std::endl;
}

#endif
