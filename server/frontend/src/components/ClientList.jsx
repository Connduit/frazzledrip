import React, { useState } from 'react';
import { utils } from '../services/api.js';
import { hostOf, systemOf, formatMemory } from '../lib/sessions.js';

export default function ClientList({ clients, selectedClient, onClientSelect, locations }) {
    const [query, setQuery] = useState('');
    const [sort, setSort] = useState('host');
    const rows = clients.filter(client => {
        const info = systemOf(client);
        return [client.ID, client.IPAddress, ...Object.values(info)].join(' ').toLowerCase().includes(query.toLowerCase());
    }).sort((a, b) => sort === 'seen' ? new Date(b.LastSeen) - new Date(a.LastSeen)
        : (systemOf(a).hostname || a.IPAddress).localeCompare(systemOf(b).hostname || b.IPAddress));
    return <section className="panel sessions-panel">
        <div className="panel-heading"><div><span className="eyebrow">INVENTORY</span><h2>Sessions <span className="count">{clients.length}</span></h2></div>
            <div className="table-tools"><input aria-label="Search sessions" placeholder="Search host, user, address…" value={query} onChange={e => setQuery(e.target.value)} />
                <select aria-label="Sort sessions" value={sort} onChange={e => setSort(e.target.value)}><option value="host">Host name</option><option value="seen">Last seen</option></select></div></div>
        <div className="table-scroll"><table><thead><tr><th>Status</th><th>Device / session</th><th>Remote endpoint</th><th>User</th><th>Platform</th><th>Processors</th><th>Memory</th><th>Location</th><th>Last seen</th><th><span className="sr-only">Action</span></th></tr></thead>
            <tbody>{rows.map(client => { const info = systemOf(client); return <tr key={client.IPAddress} className={selectedClient === client.IPAddress ? 'selected' : ''}>
                <td><span className={`badge ${client.IsActive ? 'live' : ''}`}><i />{client.IsActive ? 'Active' : 'Inactive'}</span></td>
                <td><strong>{info.hostname || 'Unknown host'}</strong><small className="mono">{client.ID}</small></td>
                <td className="mono">{client.IPAddress}</td><td>{info.user || '—'}</td><td>{[info.os, info.arch].filter(Boolean).join(' / ') || '—'}</td>
                <td>{client.SystemInfo?.number_of_processors ?? '—'}</td>
                <td>{formatMemory(client.SystemInfo?.total_physical_memory)}<small>{formatMemory(client.SystemInfo?.available_physical_memory)} available</small></td>
                <td>{locations[hostOf(client.IPAddress)]?.label || 'Unmapped'}</td><td>{utils.formatLastSeen(client.LastSeen)}</td>
                <td><button aria-pressed={selectedClient === client.IPAddress} onClick={() => onClientSelect(client.IPAddress)}>Inspect</button></td>
            </tr>; })}</tbody></table></div>
        {rows.length === 0 && <div className="empty-state"><strong>{clients.length ? 'No matching sessions' : 'No connected sessions'}</strong><p>{clients.length ? 'Try another host, user, or address.' : 'Start the listener to accept connections. Sessions appear here automatically.'}</p></div>}
        <footer className="panel-footer">Showing {rows.length} of {clients.length} sessions · refreshes every 5 seconds</footer>
    </section>;
}
