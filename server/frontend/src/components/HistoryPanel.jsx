import React, { useEffect, useState, useRef } from 'react';
import { apiService } from '../services/api.js';

export default function HistoryPanel({ agentID, audit = false }) {
    const [rows, setRows] = useState([]);
    const [error, setError] = useState('');
    const [query, setQuery] = useState('');
    const [expanded, setExpanded] = useState(null);
    const detailsRequest = useRef(0);
    useEffect(() => {
        let cancelled = false,
            timer;
        setRows([]);
        setExpanded(null);
        detailsRequest.current++;
        async function poll() {
            try {
                const data = audit
                    ? await apiService.auditHistory()
                    : await apiService.taskHistory(agentID);
                if (!cancelled) {
                    setRows(data);
                    setError('');
                }
            } catch (err) {
                if (!cancelled) setError(err.message);
            } finally {
                if (!cancelled) timer = setTimeout(poll, 3000);
            }
        }
        poll();
        return () => {
            cancelled = true;
            clearTimeout(timer);
            detailsRequest.current++;
        };
    }, [agentID, audit]);
    const filtered = rows.filter((row) =>
        JSON.stringify(row).toLowerCase().includes(query.toLowerCase()),
    );
    async function inspect(task) {
        const token = ++detailsRequest.current;
        if (expanded?.id === task.id) {
            setExpanded(null);
            return;
        }
        try {
            const result = await apiService.taskDetails(task.id);
            if (token === detailsRequest.current) setExpanded(result);
        } catch (err) {
            if (token === detailsRequest.current) setError(err.message);
        }
    }
    return (
        <section
            className="history-panel"
            aria-label={audit ? 'Audit history' : 'Task history'}
        >
            <div className="panel-heading">
                <h3>{audit ? 'Audit history' : 'Task history'}</h3>
                <input
                    aria-label={
                        audit ? 'Search audit history' : 'Search task history'
                    }
                    placeholder="Search history…"
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                />
            </div>
            <p className="muted explorer-note">
                {audit
                    ? 'Source is the HTTP peer endpoint or task origin, not an authenticated operator. Latest 300 events.'
                    : 'Latest 500 shell and simulation tasks. Queued means recorded before socket dispatch; sent means written, not acknowledged. Interrupted tasks are never replayed automatically.'}
            </p>
            {error && (
                <p role="alert" className="status error">
                    {error}
                </p>
            )}
            <div className="history-table">
                <table>
                    <thead>
                        <tr>
                            {(audit
                                ? ['Time', 'Event', 'Agent', 'Source', 'Detail']
                                : [
                                      'Created',
                                      'Command',
                                      'Status',
                                      'Origin',
                                      'Request',
                                  ]
                            ).map((label) => (
                                <th key={label}>{label}</th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {filtered.map((row) =>
                            audit ? (
                                <tr key={row.id}>
                                    <td>{new Date(row.at).toLocaleString()}</td>
                                    <td>{row.event}</td>
                                    <td className="mono">
                                        {row.agent_id || '—'}
                                    </td>
                                    <td>{row.source || 'server'}</td>
                                    <td>{row.detail || row.task_id || '—'}</td>
                                </tr>
                            ) : (
                                <tr key={row.id}>
                                    <td>
                                        {new Date(
                                            row.created_at,
                                        ).toLocaleString()}
                                    </td>
                                    <td>
                                        <button
                                            className="file-name"
                                            onClick={() => inspect(row)}
                                        >
                                            {row.command}
                                        </button>
                                        {row.simulation && (
                                            <span className="badge simulation-label">
                                                Simulation
                                            </span>
                                        )}
                                    </td>
                                    <td>
                                        <span
                                            className={`badge ${row.status === 'completed' ? 'live' : ''}`}
                                        >
                                            {row.status}
                                        </span>
                                    </td>
                                    <td>{row.origin}</td>
                                    <td className="mono">{row.request_id}</td>
                                </tr>
                            ),
                        )}
                    </tbody>
                </table>
            </div>
            {!filtered.length && <p className="muted">No matching history.</p>}
            {expanded && (
                <div className="task-details">
                    <strong className="mono">{expanded.id}</strong>
                    <p>
                        {expanded.status} · {expanded.command}
                    </p>
                    {expanded.simulation && (
                        <p className="catalog-banner">
                            Simulation only. No client was contacted.
                        </p>
                    )}
                    {expanded.output_truncated && (
                        <p className="status warning">
                            Archived task output is limited to 1 MiB.
                        </p>
                    )}
                    <pre className="message-content">
                        {expanded.output || '(No output recorded)'}
                    </pre>
                </div>
            )}
        </section>
    );
}
