import React from 'react';
import { apiService } from '../services/api.js';

function CommandTab({ selectedClient, clients, onClientSelect, setStatus, refreshMessages }) {
    const [command, setCommand] = React.useState('');
    const [sending, setSending] = React.useState(false);

    const sendCommand = async (nativeInfo = false) => {
        if (!selectedClient || sending || (!nativeInfo && !command.trim())) return;
        setSending(true);
        
        try {
            const data = nativeInfo
                ? await apiService.mySendCommand(selectedClient, 'systeminfo')
                : await apiService.sendCommand(selectedClient, command);
            if (data.status?.includes('Failed')) throw new Error(data.status);
            setStatus(nativeInfo ? '✅ System info requested - waiting for response...' : '✅ Command sent - waiting for response...');
            if (!nativeInfo) setCommand('');
            refreshMessages();
        } catch (error) {
            setStatus('Error: ' + error.message);
        } finally {
            setSending(false);
        }
    };

    const quickCommands = [
        { label: 'whoami', command: 'whoami' },
        { label: 'ipconfig', command: 'ipconfig /all' },
        { label: 'systeminfo', command: 'systeminfo' },
        { label: 'dir C:\\', command: 'dir C:\\' },
        { label: 'netstat', command: 'netstat -an' },
        { label: 'tasklist', command: 'tasklist' },
        { label: 'net user', command: 'net user' }
    ];

    return (
        <div>
            <h3>Commands</h3>
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
                aria-label="Shell command"
                onKeyDown={(e) => e.key === 'Enter' && sendCommand()}
            />
            <button onClick={() => sendCommand()} disabled={!selectedClient || !command.trim() || sending}>
                Send Command
            </button>
            
            <div className="native-commands">
                <h4>Host information</h4>
                <button onClick={() => sendCommand(true)} disabled={!selectedClient || sending}>Get system info</button>
                <p className="muted">Request the client's native OS, memory, and processor information. The systeminfo quick command below runs the Windows shell command.</p>
            </div>
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
export default CommandTab;
