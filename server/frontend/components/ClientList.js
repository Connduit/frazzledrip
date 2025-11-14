function ClientList({ clients, selectedClient, onClientSelect }) {
    if (clients.length === 0) {
        return (
            <div className="section">
                <h3>🔗 Connected Clients:</h3>
                <p>No clients connected. Start listener and wait for beacons...</p>
            </div>
        );
    }

    return (
        <div className="section">
            <h3>🔗 Connected Clients:</h3>
            <div className="client-list">
                {clients.map(client => (
                    <div key={client.IPAddress} className="client-item">
                        <div>
                            <span className="client-status online"></span>
                            <span className="client-id">{client.ID}</span>
                            <span className="client-ip"> ({client.IPAddress})</span>
                        </div>
                        <div className="client-info">
                            Last seen: {utils.formatLastSeen(client.LastSeen)}
                            {client.SystemInfo && client.SystemInfo.Hostname && (
                                <div>
                                    System: {client.SystemInfo.Hostname} | 
                                    User: {client.SystemInfo.User} | 
                                    OS: {client.SystemInfo.OS} {client.SystemInfo.Arch}
                                </div>
                            )}
                        </div>
                        <button 
                            onClick={() => onClientSelect(client.IPAddress)}
                            style={{fontSize: '0.8em', padding: '4px 8px'}}
                        >
                            Select
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
}