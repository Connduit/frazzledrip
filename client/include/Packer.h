#ifndef PACKER_H
#define PACKER_H

#include "MessageTypes.h"

class Packer
{
public:
	void pack();
	void unpack();
private:
	// TODO: figure out if this is needed and if it is the correct type
	RawByteBuffer buffer_;
};

#endif
