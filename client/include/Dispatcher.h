#ifndef DISPATCHER_H
#define DISPATCHER_H
// message dispatcher... where messages should go

#include "MessageTypes.h"

#include <functional>



// TODO: rename to inboundRouter
// and then make another class called outboundRouter

class Dispatcher
{
public:
    using Handler = std::function<void(const InternalMessage&)>;
    void registerHandler(MessageType type, Handler handler);

	// TODO: only need one of these. find out which one is needed and remove the others
    void dispatch(const InternalMessage& msg);
    void dispatch(InternalMessage& msg);
private:
    std::unordered_map<MessageType, Handler> handlers_;
};

#endif
