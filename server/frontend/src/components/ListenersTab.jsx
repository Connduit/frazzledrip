import React from 'react';

export default function ListenersTab({
    status,
    error,
    checkedAt,
    checking,
    pending,
    onControl,
    onRefresh,
}) {
    const verified = status !== null && !error;
    return (
        <section className="panel workspace-listeners" aria-label="Listeners">
            <div className="workspace-panel-title">
                <strong>Listeners</strong>
                <span>Backend status · refreshes every 2 seconds</span>
                <button
                    type="button"
                    disabled={checking || pending}
                    onClick={onRefresh}
                >
                    {checking ? 'Checking…' : 'Refresh listener status'}
                </button>
            </div>
            <div className="workspace-panel-body">
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
            </div>
        </section>
    );
}
