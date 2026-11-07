import React, { useState } from 'react';
import { apiService } from '../services/api.js';
import { systemOf } from '../lib/sessions.js';

export default function SessionLink({ live, clients, onLinked }) {
    const [savedID, setSavedID] = useState('');
    const [busy, setBusy] = useState(false);
    const [error, setError] = useState('');
    const offline = clients.filter(
        (client) => !client.IsActive && client.ID !== live.ID,
    );
    if (!offline.length) return null;
    const saved = offline.find((client) => client.ID === savedID);
    async function link(event) {
        event.preventDefault();
        if (!saved || busy) return;
        setBusy(true);
        setError('');
        try {
            await apiService.linkSession(saved.ID, live.ID);
            onLinked(saved);
            setSavedID('');
        } catch (err) {
            setError(err.message);
        } finally {
            setBusy(false);
        }
    }
    return (
        <form className="session-link" onSubmit={link}>
            <h3>Link a saved session</h3>
            <p>
                Choose the offline record belonging to this device. Its saved
                facts and history will move into this active connection, and its
                offline row will be removed. Existing live facts take priority.
            </p>
            <label htmlFor="saved-session">Offline session</label>
            <div className="link-controls">
                <select
                    id="saved-session"
                    value={savedID}
                    disabled={busy}
                    onChange={(event) => setSavedID(event.target.value)}
                >
                    <option value="">Choose a saved session…</option>
                    {offline.map((client) => (
                        <option key={client.ID} value={client.ID}>
                            {systemOf(client).hostname || 'Unknown host'} ·{' '}
                            {client.IPAddress} · {client.ID}
                        </option>
                    ))}
                </select>
                <button type="submit" disabled={!saved || busy}>
                    {busy ? 'Linking…' : 'Link saved session'}
                </button>
            </div>
            {saved && (
                <p className="mono">
                    {saved.ID} → {live.ID}
                </p>
            )}
            {error && (
                <p className="error status" role="alert">
                    {error}
                </p>
            )}
        </form>
    );
}
