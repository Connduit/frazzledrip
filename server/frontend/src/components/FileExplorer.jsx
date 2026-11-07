import React, { useEffect, useRef, useState } from 'react';
import { apiService } from '../services/api.js';

import { formatFileSize, formatFileDate, fileType, sortFiles } from '../lib/fileExplorer.js';

export default function FileExplorer({ agent }) {
    const [listing, setListing] = useState(null);
    const [path, setPath] = useState('');
    const [query, setQuery] = useState('');
    const [sort, setSort] = useState({ field: 'name', direction: 'asc' });
    const [selected, setSelected] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const request = useRef(null);
    useEffect(() => () => request.current?.abort(), []);
    async function browse(destination) {
        request.current?.abort();
        const controller = new AbortController(); request.current = controller;
        setLoading(true); setError('');
        try {
            const result = await apiService.browseFiles(agent.IPAddress, destination, controller.signal);
            if (controller.signal.aborted) return;
            setListing(result); setPath(result.path); setSelected(null); setQuery('');
        } catch (err) { if (!controller.signal.aborted) setError(err.message); }
        finally { if (!controller.signal.aborted) setLoading(false); }
    }
    const rows = sortFiles((listing?.entries || []).filter(entry => entry.name.toLowerCase().includes(query.toLowerCase())), sort.field, sort.direction);
    const columns = [{ field: 'name', label: 'Name' }, { field: 'size', label: 'Size' }, { field: 'modified_ms', label: 'Date modified' }, { field: 'type', label: 'Type' }];
    function sortBy(field) { setSort(current => ({ field, direction: current.field === field && current.direction === 'asc' ? 'desc' : 'asc' })); }
    return <section className="file-explorer" aria-label="Remote file explorer">
        <div className="explorer-heading"><div><h3>File Explorer</h3><p className="muted">Browse directory listings on {agent.SystemInfo?.hostname || agent.ID}. Select a file to inspect its metadata.</p></div><span className="badge">Read only</span></div>
        <form className="explorer-path" onSubmit={event => { event.preventDefault(); browse(path); }}>
            <button type="button" disabled={loading || !listing?.parent || listing.parent === listing.path} onClick={() => browse(listing.parent)}>Up</button>
            <input aria-label="Client directory path" placeholder="Empty path opens the client's current directory, or enter C:\\" value={path} onChange={event => setPath(event.target.value)} maxLength="2048" />
            <button type="submit" disabled={loading || !agent.IsActive}>{loading ? 'Loading…' : 'Open directory'}</button>
            <button type="button" disabled={loading || !listing} onClick={() => browse(listing.path)}>Refresh directory</button>
        </form>
        <p className="muted explorer-note">Requires the updated Windows client. Initial host queries must finish before browsing. Browse shows names and metadata; file contents and transfers stay separate.</p>
        {error && <div role="alert" className="status error">{error}{listing && <span> The previous directory listing is still displayed.</span>}</div>}
        {listing && <><div className="explorer-heading"><span className="mono explorer-current">{listing.path}</span><input aria-label="Filter files" placeholder="Filter names…" value={query} onChange={event => setQuery(event.target.value)} /></div>
            {listing.truncated && <p className="status warning">This listing is truncated. At most 2,000 entries are shown; additional entries remain.</p>}
            <div className="explorer-table"><table><thead><tr>{columns.map(column => <th key={column.field} aria-sort={sort.field === column.field ? (sort.direction === 'asc' ? 'ascending' : 'descending') : 'none'}><button className="file-sort" onClick={() => sortBy(column.field)} aria-label={`Sort by ${column.label}`}>{column.label}<span aria-hidden="true">{sort.field === column.field ? (sort.direction === 'asc' ? ' ↑' : ' ↓') : ''}</span></button></th>)}</tr></thead><tbody>
                {rows.map(entry => <tr key={entry.path} className={selected?.path === entry.path ? 'selected' : ''}>
                    <td><button className="file-name" disabled={loading} onClick={() => entry.is_directory ? browse(entry.path) : setSelected(entry)}><span aria-hidden="true">{entry.is_directory ? '📁' : '📄'}</span> {entry.name}</button></td>
                    <td className="mono">{entry.is_directory ? '—' : formatFileSize(entry.size)}</td><td>{formatFileDate(entry.modified_ms)}</td><td>{fileType(entry)}</td>
                </tr>)}
            </tbody></table></div>
            {!rows.length && <div className="empty-state">{query ? 'No matching names.' : 'This directory is empty.'}</div>}
            <p className="muted explorer-note">{rows.length} of {listing.entries.length} entries · folders first · listing changes only when refreshed</p></>}
        {!listing && !loading && <div className="empty-state">Open a directory to browse the selected client's filesystem.</div>}
        {selected && <div className="system-info"><strong>{selected.name}</strong><span className="mono">{selected.path}</span><span>Size: {formatFileSize(selected.size)}</span><span>Date modified: {formatFileDate(selected.modified_ms)}</span><span>Type: {fileType(selected)}</span></div>}
    </section>;
}
