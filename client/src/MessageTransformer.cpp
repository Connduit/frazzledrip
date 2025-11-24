#include "MessageTransformer.h"

MessageTransformer::MessageTransformer(
		Serializer* serializer, 
		Encoder* encoder, 
		Encryptor* encryptor)
		:
		serializer_(serializer),
		encoder_(encoder),
		encryptor_(encryptor)
{}

// inbound transform
InternalMessage MessageTransformer::transform(const RawByteBuffer& msg)
{
	//InternalMessage inMsg;
	auto decrypted = encryptor_->decrypt(msg);
	auto decoded = encoder_->decode(decrypted);
	auto deserialized = serializer_->deserialize(decoded);

	return deserialized;
}

// outbound transform
RawByteBuffer MessageTransformer::transform(const InternalMessage& msg)
{
	auto serialized = serializer_->serialize(msg);
	auto encoded = encoder_->encode(serialized);
	auto encrypted = encryptor_->encrypt(encoded);

	return encrypted;
}
