function ConfigTab({ selectedClient, clients, onClientSelect, setStatus }) {
    const [configUpdate, setConfigUpdate] = React.useState('');

    const updateConfig = async () => {
        if (!selectedClient || !configUpdate.trim()) return;
        
        try {
            const data = await apiService.updateConfig(selectedClient, configUpdate);
            setStatus(data.status);
            setConfigUpdate('');
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    return (
        <div>
            <h3>⚙️ Client Configuration:</h3>
            <select 
                value={selectedClient} 
                onChange={(e) => onClientSelect(e.target.value)}
            >
                <option value="">Select Client</option>
                {clients.map(client => (
                    <option key={client.IPAddress} value={client.IPAddress}>
                        {client.ID} ({client.IPAddress})
                    </option>
                ))}
            </select>
            <textarea 
                value={configUpdate}
                onChange={(e) => setConfigUpdate(e.target.value)}
                placeholder="Enter configuration (JSON format recommended)"
                style={{width: '100%', height: '100px', padding: '10px', margin: '10px 0'}}
            />
            <button onClick={updateConfig} disabled={!selectedClient || !configUpdate.trim()}>
                Update Configuration
            </button>
            <p style={{fontSize: '0.9em', color: '#666'}}>
                📋 Send configuration updates to the client (beacon intervals, C2 settings, etc.)
            </p>
        </div>
    );
}