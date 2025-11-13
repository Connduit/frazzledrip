/* Serializer.cpp */

#include "Serializer.h"
#include <iostream>

std::vector<uint8_t> BinarySerializer::serialize(const InternalMessage& msg)
{
	std::vector<uint8_t> serialized;
	//serialized.clear();

	const uint8_t* header_bytes = reinterpret_cast<const uint8_t*>(&msg.header);
	serialized.insert(serialized.end(), header_bytes, header_bytes + sizeof(MessageHeader));
	serialized.insert(serialized.end(), msg.data.begin(), msg.data.end());
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
	memcpy(&msg.header, data.data(), sizeof(MessageHeader));

	// Deserialize payload (if any)
	size_t payload_size = data.size() - sizeof(MessageHeader);
	if (payload_size > 0)
	{
		msg.data.assign(data.begin() + sizeof(MessageHeader), data.end());
	}

	// Verify data size matches header
	if (msg.header.dataSize != msg.data.size())
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


