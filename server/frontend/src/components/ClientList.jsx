import React, { useState } from 'react';
import { utils } from '../services/api.js';
import { hostOf, systemOf } from '../lib/sessions.js';

const columns = ['Agent ID', 'Type', 'External', 'Listener', 'Internal', 'Domain', 'Computer', 'User', 'OS', 'Process', 'PID', 'Tags', 'Last', 'Sleep'];

export default function ClientList({ clients, selectedClient, onClientSelect }) {
    const [query, setQuery] = useState('');
    const [sort, setSort] = useState('host');
    const rows = clients.filter(client => {
        const info = systemOf(client);
        return [client.ID, client.IPAddress, client.Type, client.Listener, client.SystemInfo?.domain,
            client.SystemInfo?.process_name, client.SystemInfo?.process_id, ...(client.Tags || []),
            ...(client.SystemInfo?.ip_addresses || []), ...Object.values(info)].join(' ').toLowerCase().includes(query.toLowerCase());
    }).sort((a, b) => sort === 'seen' ? new Date(b.LastSeen) - new Date(a.LastSeen)
        : (systemOf(a).hostname || a.IPAddress).localeCompare(systemOf(b).hostname || b.IPAddress));
    return <section className="panel sessions-panel">
        <div className="panel-heading"><div><span className="eyebrow">INVENTORY</span><h2>Sessions <span className="count">{clients.length}</span></h2></div>
            <div className="table-tools"><input aria-label="Search sessions" placeholder="Search sessions…" value={query} onChange={e => setQuery(e.target.value)} />
                <select aria-label="Sort sessions" value={sort} onChange={e => setSort(e.target.value)}><option value="host">Computer name</option><option value="seen">Last seen</option></select></div></div>
        <div className="table-scroll"><table className="sessions-table"><thead><tr>{columns.map(column => <th key={column} scope="col">{column}</th>)}</tr></thead>
            <tbody>{rows.map(client => {
                const info = systemOf(client);
                const details = client.SystemInfo || {};
                return <tr key={client.ID} className={selectedClient === client.IPAddress ? 'selected' : ''}>
                    <td><button className="agent-link mono" aria-label={`Inspect session ${client.ID}`} aria-pressed={selectedClient === client.IPAddress} onClick={() => onClientSelect(client.IPAddress)}>
                        <span className={`session-dot ${client.IsActive ? 'active' : ''}`} role="img" aria-label={client.IsActive ? 'Active' : 'Offline'} title={client.IsActive ? 'Active' : 'Offline'} />{client.ID}</button></td>
                    <td>{client.Type || '—'}</td>
                    <td className="mono" title={`Observed remote endpoint: ${client.IPAddress}`}>{hostOf(client.IPAddress) || '—'}</td>
                    <td className="mono">{client.Listener || '—'}</td>
                    <td className="mono">{details.ip_addresses?.join(', ') || '—'}</td>
                    <td>{details.domain || '—'}</td><td>{info.hostname || '—'}</td><td>{info.user || '—'}</td>
                    <td>{[info.os, info.arch].filter(Boolean).join(' / ') || '—'}</td>
                    <td>{details.process_name || '—'}</td><td className="mono">{details.process_id ?? '—'}</td>
                    <td>{client.Tags?.length ? client.Tags.map(tag => <span className="session-tag" key={tag}>{tag}</span>) : '—'}</td>
                    <td title="Last message received">{utils.formatLastSeen(client.LastSeen)}</td>
                    <td>{client.SleepSeconds == null ? '—' : `${client.SleepSeconds}s`}</td>
                </tr>;
            })}</tbody></table></div>
        {rows.length === 0 && <div className="empty-state"><strong>{clients.length ? 'No matching sessions' : 'No connected sessions'}</strong><p>{clients.length ? 'Try another computer, user, or address.' : 'Start the listener to accept connections. Sessions appear here automatically.'}</p></div>}
        <footer className="panel-footer">Showing {rows.length} of {clients.length} sessions · refreshes every 5 seconds · External is the observed peer IP · — indicates unreported data</footer>
    </section>;
}
