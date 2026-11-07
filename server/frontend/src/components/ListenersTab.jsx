import React, { useId, useRef, useState } from 'react';
import {
    parseListenerDrafts,
    validateListenerDraft,
} from '../lib/listenerDrafts.js';

const draftKey = 'frazzledrip.listener-drafts.v1';
const emptyForm = () => ({
    name: '',
    protocol: 'TCP',
    address: '127.0.0.1',
    port: '',
});

function readDrafts() {
    try {
        return parseListenerDrafts(localStorage.getItem(draftKey));
    } catch {
        return [];
    }
}

export default function ListenersTab({
    status,
    error,
    checkedAt,
    checking,
    pending,
    onControl,
    onRefresh,
}) {
    const formID = useId();
    const createButton = useRef(null);
    const [drafts, setDrafts] = useState(readDrafts);
    const [formOpen, setFormOpen] = useState(false);
    const [form, setForm] = useState(emptyForm);
    const [formError, setFormError] = useState('');
    const [notice, setNotice] = useState('');
    const [storageError, setStorageError] = useState('');

    function saveDrafts(next) {
        setDrafts(next);
        try {
            localStorage.setItem(draftKey, JSON.stringify(next));
            setStorageError('');
        } catch {
            setStorageError(
                'Browser storage is unavailable. Draft changes are kept for this visit only.',
            );
        }
    }

    function closeForm() {
        setFormOpen(false);
        setFormError('');
        setForm(emptyForm());
        createButton.current?.focus();
    }

    function submitDraft(event) {
        event.preventDefault();
        try {
            if (drafts.length >= 50)
                throw new Error(
                    'At most 50 drafts can be saved. Remove a draft before adding another.',
                );
            const draft = {
                ...validateListenerDraft(form, drafts),
                id: crypto.randomUUID(),
            };
            saveDrafts([...drafts, draft]);
            setNotice(`Draft “${draft.name}” saved. No listener was started.`);
            closeForm();
        } catch (err) {
            setFormError(err.message);
        }
    }

    function field(name, value) {
        setForm((current) => ({ ...current, [name]: value }));
        setFormError('');
    }
    const verified = status !== null && !error;
    return (
        <section className="panel workspace-listeners" aria-label="Listeners">
            <div className="workspace-panel-title">
                <strong>Listeners</strong>
                <span>Backend status · refreshes every 2 seconds</span>
                <button
                    type="button"
                    ref={createButton}
                    aria-expanded={formOpen}
                    aria-controls={formID}
                    onClick={() => {
                        setFormOpen(true);
                        setNotice('');
                    }}
                >
                    Create listener
                </button>
                <button
                    type="button"
                    disabled={checking || pending}
                    onClick={onRefresh}
                >
                    {checking ? 'Checking…' : 'Refresh listener status'}
                </button>
            </div>
            <div className="workspace-panel-body">
                {formOpen && (
                    <form
                        id={formID}
                        className="listener-draft-form"
                        aria-label="Create listener draft"
                        onSubmit={submitDraft}
                        noValidate
                    >
                        <h3>Create listener draft</h3>
                        <p className="muted">
                            Save a test configuration in this browser. Drafts do
                            not open sockets or create running listeners.
                        </p>
                        <div className="listener-form-fields">
                            <label>
                                Name
                                <input
                                    autoFocus
                                    value={form.name}
                                    maxLength={64}
                                    onChange={(event) =>
                                        field('name', event.target.value)
                                    }
                                    required
                                />
                            </label>
                            <label>
                                Protocol
                                <select
                                    value={form.protocol}
                                    onChange={(event) =>
                                        field('protocol', event.target.value)
                                    }
                                >
                                    <option value="TCP">TCP</option>
                                </select>
                            </label>
                            <label>
                                Bind address
                                <input
                                    value={form.address}
                                    onChange={(event) =>
                                        field('address', event.target.value)
                                    }
                                    maxLength={128}
                                    placeholder="127.0.0.1"
                                    required
                                />
                            </label>
                            <label>
                                Port
                                <input
                                    type="text"
                                    inputMode="numeric"
                                    value={form.port}
                                    maxLength={5}
                                    onChange={(event) =>
                                        field('port', event.target.value)
                                    }
                                    placeholder="1–65535"
                                    required
                                />
                            </label>
                        </div>
                        {formError && (
                            <p role="alert" className="status error">
                                {formError}
                            </p>
                        )}
                        <div className="listener-actions">
                            <button type="submit">Save draft</button>
                            <button type="button" onClick={closeForm}>
                                Cancel
                            </button>
                        </div>
                    </form>
                )}
                {notice && (
                    <p role="status" className="status">
                        {notice}
                    </p>
                )}
                {storageError && (
                    <p role="alert" className="status warning">
                        {storageError}
                    </p>
                )}
                <p className="muted">
                    The server currently supports one TCP listener. Status below
                    comes from the backend listener registry, rather than
                    assuming a button click succeeded.
                </p>
                {error && (
                    <p role="alert" className="status error">
                        Could not verify listener status: {error}. Last known
                        details may be stale; retrying automatically.
                    </p>
                )}
                <div className="listeners-table">
                    <table>
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Type</th>
                                <th>Status</th>
                                <th>Port</th>
                                <th>Last verified</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>Default TCP listener</td>
                                <td>TCP</td>
                                <td>
                                    <span
                                        className={`badge ${verified && status.running ? 'live' : ''}`}
                                    >
                                        {!verified
                                            ? 'Unverified'
                                            : status.running
                                              ? 'Running'
                                              : 'Stopped'}
                                    </span>
                                    {error && status && (
                                        <small>
                                            Last known:{' '}
                                            {status.running
                                                ? 'Running'
                                                : 'Stopped'}
                                        </small>
                                    )}
                                </td>
                                <td>{status?.running ? status.port : '—'}</td>
                                <td>
                                    {checkedAt ? (
                                        <time
                                            dateTime={checkedAt.toISOString()}
                                        >
                                            {checkedAt.toLocaleString()}
                                        </time>
                                    ) : (
                                        'Not checked yet'
                                    )}
                                </td>
                                <td>
                                    <div className="listener-actions">
                                        <button
                                            type="button"
                                            disabled={
                                                pending ||
                                                !verified ||
                                                status.running
                                            }
                                            onClick={() => onControl('start')}
                                        >
                                            Start
                                        </button>
                                        <button
                                            type="button"
                                            disabled={
                                                pending ||
                                                !verified ||
                                                !status.running
                                            }
                                            onClick={() => onControl('stop')}
                                        >
                                            Stop
                                        </button>
                                        {pending && (
                                            <span role="status">
                                                Applying and verifying…
                                            </span>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <p className="muted explorer-note">
                    Start uses the existing configured TCP port (4444). The Port
                    column shows the backend-reported bound port while running.
                    Stopping the listener stops new connections; existing
                    sessions remain connected.
                </p>
                <section
                    className="listener-drafts"
                    aria-label="Listener drafts"
                >
                    <h3>
                        Listener drafts{' '}
                        <span className="badge">{drafts.length}</span>
                    </h3>
                    <p className="muted">
                        Browser-local test configurations. These entries are
                        separate from the running listener above.
                    </p>
                    {!drafts.length && (
                        <p className="muted">
                            No drafts yet. Choose Create listener to add one.
                        </p>
                    )}
                    <ul className="listener-draft-list">
                        {drafts.map((draft) => (
                            <li key={draft.id}>
                                <div>
                                    <strong>{draft.name}</strong>
                                    <span className="badge">Draft only</span>
                                    <p className="mono">
                                        {draft.protocol} ·{' '}
                                        {draft.address.includes(':')
                                            ? `[${draft.address}]`
                                            : draft.address}
                                        :{draft.port}
                                    </p>
                                </div>
                                <button
                                    type="button"
                                    aria-label={`Remove draft ${draft.name}`}
                                    onClick={() => {
                                        saveDrafts(
                                            drafts.filter(
                                                (item) => item.id !== draft.id,
                                            ),
                                        );
                                        setNotice(
                                            `Draft “${draft.name}” removed.`,
                                        );
                                    }}
                                >
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                </section>
            </div>
        </section>
    );
}
