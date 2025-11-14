function App() {
    const [status, setStatus] = React.useState('');
    const [clients, setClients] = React.useState([]);
    const [selectedClient, setSelectedClient] = React.useState('');
    const [clientInfo, setClientInfo] = React.useState(null);
    const [clientMessages, setClientMessages] = React.useState([]);
    const [activeTab, setActiveTab] = React.useState('commands');

    const startListener = async () => {
        try {
            const data = await apiService.startListener();
            setStatus(data.status);
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    const refreshClients = async () => {
        try {
            const data = await apiService.getClients();
            setClients(data);
            if (data.length > 0 && !selectedClient) {
                handleClientSelect(data[0].IPAddress);
            }
        } catch (error) {
            console.error('Error fetching clients:', error);
        }
    };

    const getClientInfo = async (clientIP) => {
        try {
            const data = await apiService.getClientInfo(clientIP);
            setClientInfo(data);
        } catch (error) {
            console.error('Error fetching client info:', error);
            setClientInfo(null);
        }
    };

    const fetchClientMessages = async (clientIP) => {
        try {
            const data = await apiService.getClientMessages(clientIP);
            setClientMessages(data.messages || []);
        } catch (error) {
            console.error('Error fetching messages:', error);
            setClientMessages([]);
        }
    };

    const handleClientSelect = (clientIP) => {
        setSelectedClient(clientIP);
        getClientInfo(clientIP);
        fetchClientMessages(clientIP);
    };

    const refreshMessages = async () => {
        if (selectedClient) {
            await fetchClientMessages(selectedClient);
        }
    };

    // Auto-refresh clients every 5 seconds
    React.useEffect(() => {
        refreshClients();
        const clientInterval = setInterval(refreshClients, 5000);
        return () => clearInterval(clientInterval);
    }, []);

    // Auto-refresh messages for selected client every 2 seconds
    React.useEffect(() => {
        if (selectedClient) {
            const messageInterval = setInterval(refreshMessages, 2000);
            return () => clearInterval(messageInterval);
        }
    }, [selectedClient]);

    return (
        <div className="container">
            <h1>🖥️ Advanced C2 Controller</h1>
            
            <div className="section">
                <button onClick={startListener}>Start C2 Listener on Port 4444</button>
                <button onClick={refreshClients}>Refresh Clients</button>
            </div>
            
            <ClientList 
                clients={clients}
                selectedClient={selectedClient}
                onClientSelect={handleClientSelect}
            />
            
            {selectedClient && (
                <div className="section">
                    <h3>🎯 Selected Client: {selectedClient}</h3>
                    <ClientInfo clientInfo={clientInfo} />
                    
                    <div style={{marginTop: '15px'}}>
                        <h4>📨 Client Messages:</h4>
                        <button onClick={refreshMessages}>Refresh Messages</button>
                        <MessageList messages={clientMessages} />
                    </div>
                </div>
            )}
            
            <TabContainer 
                activeTab={activeTab}
                setActiveTab={setActiveTab}
                selectedClient={selectedClient}
                clients={clients}
                onClientSelect={handleClientSelect}
                setStatus={setStatus}
                refreshMessages={refreshMessages}
            />
            
            {status && (
                <div className={`status ${utils.getStatusClass(status)}`}>
                    <strong>Status:</strong> {status}
                </div>
            )}
        </div>
    );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);