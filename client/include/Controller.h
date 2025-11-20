#ifndef CONTROLLER_H
#define CONTROLLER_H
// message controller... should honestly be renamed ClientController.h and copy and paste the ClientSubsystem in here

#include "MessageTypes.h"

#include <functional>

class Controller
{
public:
    using Handler = std::function<void(const InternalMessage&)>;
    void registerHandler(MessageType type, Handler handler);
    void dispatch(const InternalMessage& msg);


	void handleDefault(const InternalMessage& msg);
private:
    std::unordered_map<MessageType, Handler> handlers_;
};

#endif
