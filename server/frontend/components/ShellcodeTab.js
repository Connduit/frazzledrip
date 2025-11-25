function ShellcodeTab({ selectedClient, clients, onClientSelect, setStatus }) {
    const [selectedShellcodeFile, setSelectedShellcodeFile] = React.useState(null);

    const sendShellcode = async () => {
        if (!selectedClient || !selectedShellcodeFile) {
            setStatus('❌ Please select a client and shellcode file');
            return;
        }

        if (!window.confirm('⚠️ DANGER: Executing shellcode can crash the target system or cause unexpected behavior. Are you sure you want to continue?')) {
            setStatus('Shellcode execution cancelled');
            return;
        }
        
        try {
            const data = await apiService.sendShellcode(selectedClient, selectedShellcodeFile);
            setStatus(data.status);
            setSelectedShellcodeFile(null);
            document.getElementById('shellcode-input').value = '';
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    const handleShellcodeFileSelect = (event) => {
        setSelectedShellcodeFile(event.target.files[0]);
    };

    const loadExampleShellcode = (exampleName) => {
        setStatus(`⚠️ ${exampleName} - Remember to upload actual shellcode file`);
    };

    const examples = [
        { label: 'MessageBox', name: 'MessageBox' },
        { label: 'Calculator', name: 'Calc' },
        { label: 'Reverse Shell', name: 'Reverse Shell' },
        { label: 'Beacon', name: 'Beacon' }
    ];

    return (
        <div className="danger-zone">
            <h3>⚡ Shellcode Execution</h3>
            <div className="shellcode-warning">
                <strong>⚠️ WARNING: DANGEROUS OPERATION</strong><br/>
                • Shellcode executes directly in memory<br/>
                • Can crash the target system<br/>
                • May trigger antivirus detection<br/>
                • Use only with proper authorization
            </div>
            
            <select 
                value={selectedClient} 
                onChange={(e) => onClientSelect(e.target.value)}
            >
                <option value="">Select Client</option>
                {clients.map(client => (
                    <option key={client.IPAddress} value={client.IPAddress}>
                        {client.ID} ({client.IPAddress}) - {client.SystemInfo?.OS}
                    </option>
                ))}
            </select>
            <input 
                id="shellcode-input"
                type="file" 
                onChange={handleShellcodeFileSelect}
                accept=".bin,.raw,.sc,.txt"
            />
            <button 
                onClick={sendShellcode} 
                disabled={!selectedClient || !selectedShellcodeFile}
                style={{background: '#dc3545'}}
            >
                🚀 Execute Shellcode
            </button>
            
            {selectedShellcodeFile && (
                <div style={{marginTop: '10px'}}>
                    <strong>Selected Shellcode:</strong> {selectedShellcodeFile.name} 
                    ({(selectedShellcodeFile.size / 1024).toFixed(2)} KB)
                </div>
            )}
            
            <div style={{marginTop: '15px'}}>
                <h4>💡 Example Shellcode Types:</h4>
                {examples.map((example, index) => (
                    <button key={index} onClick={() => loadExampleShellcode(example.name)}>
                        {example.label}
                    </button>
                ))}
            </div>
            
            <p style={{fontSize: '0.8em', color: '#666', marginTop: '10px'}}>
                📋 Upload raw binary shellcode files (.bin, .raw) or text format
            </p>
        </div>
    );
}