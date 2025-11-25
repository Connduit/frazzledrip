/* Serializer.cpp */

#include "Serializer.h"
#include <iostream>

std::vector<uint8_t> BinarySerializer::serialize(const InternalMessage& msg)
{
	std::vector<uint8_t> serialized;
	//serialized.clear();

	const uint8_t* header_bytes = reinterpret_cast<const uint8_t*>(&msg.header_);
	serialized.insert(serialized.end(), header_bytes, header_bytes + sizeof(MessageHeader));
	serialized.insert(serialized.end(), msg.data_.begin(), msg.data_.end());
	return serialized;


}

InternalMessage BinarySerializer::deserialize(const RawByteBuffer& data)
{
	if (data.size() < sizeof(MessageHeader))
	{
		//throw std::runtime_error("Data too small for header");
	}

	InternalMessage msg;

	// Deserialize header
	memcpy(&msg.header_, data.data(), sizeof(MessageHeader));

	// Deserialize payload (if any)
	size_t payload_size = data.size() - sizeof(MessageHeader);
	if (payload_size > 0)
	{
		msg.data_.assign(data.begin() + sizeof(MessageHeader), data.end());
	}
	else
	{
		std::cout << "RawByteBuffer failed validation" << std::endl;
	}

	// Verify data size matches header
	if (msg.header_.dataSize_ != msg.data_.size())
	{
		//throw std::runtime_error("Data size mismatch");
		std::cout << "data size doesn't match header" << std::endl;
	}

	return msg;
}



/*


int serializeReconMessage(const ReconMessage* msg, BYTE* buffer) {
	BYTE* ptr = buffer; // uint8_t*

	// Copy each field from msg struct to buffer
	memcpy(ptr, &msg->dwMajorVersion, sizeof(ULONG)); ptr += sizeof(ULONG);
	memcpy(ptr, &msg->dwMinorVersion, sizeof(ULONG)); ptr += sizeof(ULONG);
	memcpy(ptr, &msg->dwBuildNumber, sizeof(ULONG)); ptr += sizeof(ULONG);
	memcpy(ptr, &msg->wProductType, sizeof(UCHAR)); ptr += sizeof(UCHAR);
	// ... other fields
	return ptr - buffer; // Returns total bytes copied
}

// int deserialize() {}
*/


