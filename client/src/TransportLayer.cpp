#include "TransportLayer.h"
#include "MessageTypes.h"
#include "MessageHandler.h"
#include "ComponentFactory.h"
//
#include <winsock2.h>
#include <ws2tcpip.h>
#include <iomanip>  // for std::hex, std::setw, std::setfill
#include <iostream>
#pragma comment(lib, "ws2_32.lib")


TransportLayer::TransportLayer(Encryptor* encryptor) :
    encryptor_(encryptor)
{
}

bool TransportLayer::sendMessage(const InternalMessage& msg)
{
    //auto serialized = serializer_->serialize(msg);
    //auto serialized = messageHandler_->serializer_->serialize(msg);
    //return send(serialized);
	// TODO: 
	return false;
}

//InternalMessage TransportLayer::receiveMessage()
RawByteBuffer TransportLayer::receiveMessage()
{
    auto data = receive();
    if (data.empty()) 
	{
		//return InternalMessage();
        return {};
	}

    //auto decrypted = encryptor_.decrypt(data);
    //auto decoded = encoder_.decode(decrypted);
    //return serializer_.deserialize(decoded);
    return data;
    //return serializer_->deserialize(data);
}

void TransportLayer::beacon()
{
    //auto heartbeat = createHeartbeat();
    //sendMessage(heartbeat);


	// TODO: 
    auto incoming = receiveMessage();
    //auto incoming = encoder_->decode(incoming);
    // incoming = messageHandler_->encoder_->decode(incoming);
    // InternalMessage internalMessage = messageHandler_->serializer_->deserialize(incoming);
    // if (internalMessage.header.messageType != DEFAULT)
    // {
    //     messageHandler_->processMessage(internalMessage); // change to not use ptr?
    // }
}

InternalMessage TransportLayer::createHeartbeat()
{
    InternalMessage msg;
    msg.header.messageType = MessageType::HEARTBEAT;
    msg.header.messageId = generateId();
    msg.header.dataSize = 0;
    return msg;
}

uint32_t TransportLayer::generateId()
{
    return rand();
}





/*
int handleHTTPS(uint8_t* data)
{
    // HTTPS Text Message

    // PARSE HTTP HEADERS

    // EXTRACT JSON BODY

    // EXTRACT BASE64 FIELD

    // BASE64 DECODE -> Raw Bytes

    // DECRYPT -> Decrypted Bytes

    // DESERIALIZE -> InternalMessage

    // VALIDATE

    // EXECUTE (InternalMessage)

}*/


///////////////////////



/*
bool TransportLayer::sendMessage() // overload this function?
{
    // TODO: where does parsing go in this? (if it's even needed)

    // InternalMessage
    // Validate
    // Serialize
    // Encode
    // Compress
    // Encrypt
    // Frame (if needed)
    // Raw Bytes
    // Send
    return false;
}

// TODO: setup a listener function that calls this function?
// listener function simply calls recv and handles socket stuff?
// recvMessage would just serve as a wrapper for that maybe?
bool TransportLayer::recvMessage() // overload this function?
{
    // TODO: where does parsing go in this? (if it's even needed)

    // Recv
    // Raw Bytes
    // Decrypt
    // Decompress
    // Decode
    // Deserialize
    // Obtain InternalMessage object
    // Validate
    // handle InternalMessage (based on its header)
    return false;
}
*/
