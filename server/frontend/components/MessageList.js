function MessageList({ messages }) {
    if (messages.length === 0) {
        return (
            <div className="message-list">
                <p>No messages yet. Send a command to see results.</p>
            </div>
        );
    }

    return (
        <div className="message-list">
            {messages.map((msg, index) => (
                <div key={index} className="message-item">
                    <div className="message-header">
                        <strong>{msg.type}</strong> 
                        <span className="message-time">
                            {utils.formatMessageTime(msg.timestamp)}
                        </span>
                    </div>
                    <div className="message-content">
                        {msg.content || `Binary data: ${msg.data_size} bytes`}
                    </div>
                </div>
            ))}
        </div>
    );
}