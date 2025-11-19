// message dispatcher... where messages should go
class Dispatcher
{
public:
    using Handler = std::function<void(const InternalMessage&)>;
    void registerHandler(MessageType type, Handler handler);
    void dispatch(const InternalMessage& msg);
private:
    std::unordered_map<MessageType, Handler> handlers_;
};
