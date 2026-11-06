import React from 'react';

function ClientInfo({ clientInfo }) {
    if (!clientInfo) {
        return <div>No client information available.</div>;
    }

    return (
        <div className="system-info">
            <strong>System Information:</strong>
            <div>Hostname: {clientInfo.hostname || 'Unknown'}</div>
            <div>User: {clientInfo.user || 'Unknown'}</div>
            <div>OS: {clientInfo.os || 'Unknown'} {clientInfo.arch || ''}</div>
        </div>
    );
}
export default ClientInfo;
