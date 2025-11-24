#ifndef MESSAGE_TRANSFORMER_H
#define MESSAGE_TRANSFORMER_H

#include "MessageTypes.h"
#include "Serializer.h"
#include "Encoder.h"
#include "Encryptor.h"

// TODO: rename to message transformer?
class MessageTransformer
{
public:
	MessageTransformer(Serializer* serializer, Encoder* encoder, Encryptor* encryptor);

	// inbound transform
	InternalMessage transform(RawByteBuffer& msg);

	// outbound transform
	RawByteBuffer transform(InternalMessage& msg);
private:
	Serializer* serializer_;
	Encoder* encoder_;
	Encryptor* encryptor_;
};

#endif
