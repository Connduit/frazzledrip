function FileTransferTab({ selectedClient, clients, onClientSelect, setStatus }) {
    const [selectedFile, setSelectedFile] = React.useState(null);

    const sendBinary = async () => {
        if (!selectedClient || !selectedFile) return;
        
        try {
            const data = await apiService.sendFile(selectedClient, selectedFile);
            setStatus(data.status);
            setSelectedFile(null);
            document.getElementById('binary-input').value = '';
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    const handleFileSelect = (event) => {
        setSelectedFile(event.target.files[0]);
    };

    return (
        <div>
            <h3>📦 File Transfer:</h3>
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
            <input 
                id="binary-input"
                type="file" 
                onChange={handleFileSelect}
            />
            <button onClick={sendBinary} disabled={!selectedClient || !selectedFile}>
                Send File to Client
            </button>
            {selectedFile && (
                <div style={{marginTop: '10px'}}>
                    <strong>Selected:</strong> {selectedFile.name} 
                    ({(selectedFile.size / 1024).toFixed(2)} KB)
                </div>
            )}
            <p style={{fontSize: '0.9em', color: '#666', marginTop: '10px'}}>
                📁 Files will be downloaded to the client's current directory
            </p>
        </div>
    );
}