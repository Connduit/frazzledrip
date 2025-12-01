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




/*
TransportLayer::TransportLayer(
    ApiManager* apiManager)
    :
    apiManager_(apiManager)
{
}*/


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
