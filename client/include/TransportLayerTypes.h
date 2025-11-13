#ifndef TRANSPORT_LAYER_TYPES_H
#define TRANSPORT_LAYER_TYPES_H

typedef enum
{
	//NONE,
	TCP,
	UDP,
	HTTP,
	HTTPS,
	DNS
} TransportLayerType; // TODO: rename to TransportType?

#endif
