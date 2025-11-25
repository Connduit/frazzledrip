function ClientInfo({ clientInfo }) {
    if (!clientInfo) {
        return <div>No client information available.</div>;
    }

    return (
        <div style={{background: '#f8f9fa', padding: '10px', borderRadius: '4px'}}>
            <strong>System Information:</strong>
            <div>Hostname: {clientInfo.Hostname || 'Unknown'}</div>
            <div>User: {clientInfo.User || 'Unknown'}</div>
            <div>OS: {clientInfo.OS || 'Unknown'} {clientInfo.Arch || ''}</div>
        </div>
    );
}