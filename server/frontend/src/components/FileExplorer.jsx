import React, { useEffect, useRef, useState } from 'react';
import { apiService } from '../services/api.js';

import { formatFileSize, formatFileDate, fileType, sortFiles } from '../lib/fileExplorer.js';

export default function FileExplorer({ agent }) {
    const [listing, setListing] = useState(null);
    const [path, setPath] = useState('');
    const [source, setSource] = useState('dir');
    const [query, setQuery] = useState('');
    const [sort, setSort] = useState({ field: 'name', direction: 'asc' });
    const [selected, setSelected] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [contextMenu, setContextMenu] = useState(null);
    const [notice, setNotice] = useState('');
    const menuRef = useRef(null);
    const menuOrigin = useRef(null);
    const downloadRequest = useRef(null);
    useEffect(() => {
        setContextMenu(null); setNotice(''); setSelected(null);
        return () => downloadRequest.current?.abort();
    }, [agent.ID]);
    useEffect(() => {
        if (!contextMenu) return;
        menuRef.current?.querySelector('button')?.focus();
        const dismiss = event => {
            if (!menuRef.current?.contains(event.target)) setContextMenu(null);
        };
        const escape = event => {
            if (event.key === 'Escape') {
                event.preventDefault(); setContextMenu(null); menuOrigin.current?.focus();
            }
        };
        const close = () => setContextMenu(null);
        document.addEventListener('pointerdown', dismiss);
        document.addEventListener('keydown', escape);
        window.addEventListener('resize', close);
        window.addEventListener('scroll', close, true);
        return () => {
            document.removeEventListener('pointerdown', dismiss);
            document.removeEventListener('keydown', escape);
            window.removeEventListener('resize', close);
            window.removeEventListener('scroll', close, true);
        };
    }, [contextMenu]);
    function openContextMenu(event, entry) {
        event.preventDefault();
        menuOrigin.current = event.currentTarget.querySelector('button');
        const bounds = event.currentTarget.getBoundingClientRect();
        setSelected(entry); setNotice('');
        setContextMenu({ entry, x: Math.max(8, Math.min(event.clientX || bounds.left, window.innerWidth - 256)), y: Math.max(8, Math.min(event.clientY || bounds.bottom, window.innerHeight - 132)) });
    }
    async function requestDownloadPlaceholder(entry) {
        if (entry.is_directory) {
            setContextMenu(null); menuOrigin.current?.focus();
            setNotice(`“${entry.name}” is a folder, not a file. Folder downloads are not implemented yet and need to be added in the future. No download request was sent.`);
            return;
        }
        downloadRequest.current?.abort();
        const controller = new AbortController(); downloadRequest.current = controller;
        setContextMenu(null); menuOrigin.current?.focus();
        setNotice('Sending download request…');
        try {
            const result = await apiService.downloadFilePlaceholder(agent.ID, entry.path, controller.signal);
            if (!controller.signal.aborted) setNotice(result.status || 'Request acknowledged. File receipt has not been confirmed.');
        } catch (err) {
            if (!controller.signal.aborted) setNotice(err.message);
        }
    }
    const request = useRef(null);
    useEffect(() => () => request.current?.abort(), []);
    useEffect(() => {
        const controller = new AbortController(); request.current = controller;
        apiService.browseFiles(agent.IPAddress, '', controller.signal, 'dir', false, true).then(result => {
            if (!controller.signal.aborted) { setListing(result); setPath(current => current || result.path); }
        }).catch(() => {}); // A cache miss leaves the explicit Open action available.
        return () => controller.abort();
    }, [agent.ID]);
    async function browse(destination, refresh = false) {
        setContextMenu(null); setNotice('');
        if (refresh && !agent.IsActive) { setError('Reconnect to refresh this directory. Saved listings remain available.'); return; }
        request.current?.abort();
        const controller = new AbortController(); request.current = controller;
        setLoading(true); setError('');
        try {
            const result = await apiService.browseFiles(agent.IPAddress, destination, controller.signal, source, refresh);
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
            <select aria-label="Listing method" disabled={loading} value={source} onChange={event => { setSource(event.target.value); setError(''); }}><option value="dir">Windows dir</option><option value="native">Native listing</option></select>
            <input aria-label="Client directory path" placeholder="Empty path opens the client's current directory, or enter C:\\" value={path} onChange={event => setPath(event.target.value)} maxLength="2048" />
            <button type="submit" disabled={loading}>{loading ? 'Loading…' : 'Open directory'}</button>
            <button type="button" disabled={loading || !listing || !agent.IsActive} onClick={() => browse(listing.path, true)}>Refresh directory</button>
        </form>
        <p className="muted explorer-note">Windows dir works through the existing command handler. Native listing requires the updated client. Dates from dir use the client’s local format.</p>
        {error && <div role="alert" className="status error">{error}{listing && <span> The previous directory listing is still displayed.</span>}</div>}
        {notice && <p role="status" className="status warning">{notice}</p>}
        {listing && <><div className="explorer-heading"><span className="mono explorer-current">{listing.path}</span><input aria-label="Filter files" placeholder="Filter names…" value={query} onChange={event => setQuery(event.target.value)} /></div>
            {listing.captured_at && <p className="explorer-note muted">{listing.cached ? 'Cached snapshot' : 'Snapshot captured'} · {new Date(listing.captured_at).toLocaleString()} · Refresh queries the client again.</p>}
            {listing.cache_warning && <p role="alert" className="status warning">{listing.cache_warning}</p>}
            {listing.truncated && <p className="status warning">This listing is truncated. At most 2,000 entries are shown; additional entries remain.</p>}
            <div className="explorer-table"><table><thead><tr>{columns.map(column => <th key={column.field} aria-sort={sort.field === column.field ? (sort.direction === 'asc' ? 'ascending' : 'descending') : 'none'}><button className="file-sort" onClick={() => sortBy(column.field)} aria-label={`Sort by ${column.label}`}>{column.label}<span aria-hidden="true">{sort.field === column.field ? (sort.direction === 'asc' ? ' ↑' : ' ↓') : ''}</span></button></th>)}</tr></thead><tbody>
                {rows.map(entry => <tr key={entry.path} className={selected?.path === entry.path ? 'selected' : ''} onContextMenu={event => openContextMenu(event, entry)} onKeyDown={event => { if (event.key === 'ContextMenu' || (event.shiftKey && event.key === 'F10')) openContextMenu(event, entry); }}>
                    <td><button className="file-name" disabled={loading} onClick={() => entry.is_directory ? browse(entry.path) : setSelected(entry)}><span aria-hidden="true">{entry.is_directory ? '📁' : '📄'}</span> {entry.name}</button></td>
                    <td className="mono">{entry.is_directory ? '—' : formatFileSize(entry.size)}</td><td>{formatFileDate(entry.modified_ms, entry.modified_text)}</td><td>{fileType(entry)}</td>
                </tr>)}
            </tbody></table></div>
            {!rows.length && <div className="empty-state">{query ? 'No matching names.' : 'This directory is empty.'}</div>}
            <p className="muted explorer-note">{rows.length} of {listing.entries.length} entries · folders first · listing changes only when refreshed</p></>}
        {!listing && !loading && <div className="empty-state">Open a directory to browse the selected client's filesystem.</div>}
        {selected && <div className="system-info"><strong>{selected.name}</strong><span className="mono">{selected.path}</span><span>Size: {formatFileSize(selected.size)}</span><span>Date modified: {formatFileDate(selected.modified_ms, selected.modified_text)}</span><span>Type: {fileType(selected)}</span></div>}
        {contextMenu && <div ref={menuRef} className="file-context-menu" aria-label={`Actions for ${contextMenu.entry.name}`} style={{ left: contextMenu.x, top: contextMenu.y }}>
            <strong>{contextMenu.entry.name}</strong>
            <button type="button" onClick={() => requestDownloadPlaceholder(contextMenu.entry)}>Download to server</button>
            <button type="button" onClick={() => { setContextMenu(null); menuOrigin.current?.focus(); }}>Close</button>
        </div>}
    </section>;
}
