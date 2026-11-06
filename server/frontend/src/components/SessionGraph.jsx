import React from 'react';
import { systemOf } from '../lib/sessions.js';

export default function SessionGraph({ clients, selectedClient, onClientSelect }) {
    const height = Math.max(240, clients.length * 66 + 40);
    return <section className="panel"><div className="panel-heading"><div><span className="eyebrow">TOPOLOGY</span><h2>Connection graph</h2></div><span className="badge live"><i />Direct TCP</span></div>
        <div className="graph-scroll"><svg viewBox={`0 0 620 ${height}`} className="session-graph" role="group" aria-label="Sessions connected directly to the controller">
            {clients.map((client, i) => { const y = 40 + i * 66; return <g key={client.IPAddress}>
                <path className={`graph-edge ${client.IsActive ? '' : 'disconnected'}`} d={`M 164 ${height / 2} C 290 ${height / 2}, 280 ${y}, 352 ${y}`} />
                <g className={`graph-node ${selectedClient === client.IPAddress ? 'chosen' : ''}`} role="button" tabIndex="0" aria-label={`Select ${systemOf(client).hostname || client.IPAddress}`} aria-pressed={selectedClient === client.IPAddress}
                    onClick={() => onClientSelect(client.IPAddress)} onKeyDown={e => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); onClientSelect(client.IPAddress); } }}>
                    <rect x="352" y={y - 23} width="245" height="48" rx="7" /><circle cx="368" cy={y} r="4" className={client.IsActive ? 'active-dot' : 'inactive-dot'} />
                    <text x="382" y={y - 3}>{(systemOf(client).hostname || 'Unknown host').slice(0, 25)}</text><text x="382" y={y + 14} className="graph-caption">{client.IPAddress}</text>
                </g></g>; })}
            <rect x="20" y={height / 2 - 32} width="144" height="64" rx="8" className="controller-node" /><text x="40" y={height / 2 - 4}>CONTROLLER</text><text x="40" y={height / 2 + 16} className="graph-caption">TCP listener · 4444</text>
            {!clients.length && <text x="240" y={height / 2} className="graph-caption">Waiting for connections</text>}
        </svg></div><footer className="panel-footer">Observed controller connections · select a device to inspect it</footer>
    </section>;
}
