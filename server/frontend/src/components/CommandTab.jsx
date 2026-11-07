import React, { useState } from 'react';
import { apiService } from '../services/api.js';

export default function CommandTab({
    selectedClient,
    clients,
    onClientSelect,
    setStatus,
    refreshMessages,
}) {
    const [sending, setSending] = useState(false);
    const selected = clients.find(
        (client) => client.IPAddress === selectedClient,
    );

    async function sendSystemInfo() {
        if (!selected?.IsActive || sending) return;
        setSending(true);
        try {
            const result = await apiService.mySendCommand(
                selectedClient,
                'systeminfo',
            );
            if (!result.status?.startsWith('✅')) {
                throw new Error(
                    result.error ||
                        result.status ||
                        'SYS_INFO was not accepted',
                );
            }
            setStatus('✅ SYS_INFO sent - waiting for response...');
            refreshMessages();
        } catch (error) {
            setStatus('Error: ' + error.message);
        } finally {
            setSending(false);
        }
    }

    return (
        <div>
            <h3>Custom commands</h3>
            <select
                aria-label="Custom command session"
                value={selectedClient}
                onChange={(event) => onClientSelect(event.target.value)}
            >
                <option value="">Select Client</option>
                {clients.map((client) => (
                    <option key={client.ID} value={client.IPAddress}>
                        {client.ID} ({client.IPAddress})
                    </option>
                ))}
            </select>
            <div className="native-commands">
                <h4>SYS_INFO</h4>
                <p className="muted">
                    Request native OS, memory, and processor information using
                    the SYS_INFO message type. This interface does not execute
                    shell commands.
                </p>
                <button
                    type="button"
                    onClick={sendSystemInfo}
                    disabled={!selected?.IsActive || sending}
                >
                    {sending ? 'Sending…' : 'Send SYS_INFO'}
                </button>
            </div>
        </div>
    );
}
