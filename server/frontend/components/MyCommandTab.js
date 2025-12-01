function MyCommandTab({ selectedClient, clients, onClientSelect, setStatus, refreshMessages }) {
    const [command, setCommand] = React.useState('');

    const sendCommand = async () => {
        if (!selectedClient || !command) return;
        
        try {
            await apiService.mySendCommand(selectedClient, command);
            setStatus('✅ Command sent - waiting for response...');
            setCommand('');
            refreshMessages();
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    const quickCommands = [
        { label: 'systeminfo', command: 'systeminfo' }//,
        // { label: 'net user', command: 'net user' }
    ];

    return (
        <div>
            <h3>📝 My Execute Commands:</h3>
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
                type="text" 
                value={command}
                onChange={(e) => setCommand(e.target.value)}
                placeholder="Enter command (e.g., whoami, ipconfig, dir)"
                style={{width: '300px'}}
                onKeyPress={(e) => e.key === 'Enter' && sendCommand()}
            />
            <button onClick={sendCommand} disabled={!selectedClient || !command}>
                Send Command
            </button>
            
            <div className="quick-commands">
                <h4>⚡ Quick Commands:</h4>
                {quickCommands.map((cmd, index) => (
                    <button key={index} onClick={() => setCommand(cmd.command)}>
                        {cmd.label}
                    </button>
                ))}
            </div>
        </div>
    );
}