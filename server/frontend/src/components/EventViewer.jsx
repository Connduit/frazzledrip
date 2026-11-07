import React, { useEffect, useState } from 'react';
import { apiService } from '../services/api.js';
import { clientMessageEvent } from '../lib/clientEvents.js';

function levelOf(event) {
    if (
        /failed|error/.test(event.event) ||
        /status=5\d\d/.test(event.detail || '')
    )
        return 'error';
    if (
        /disconnected|interrupted/.test(event.event) ||
        /status=4\d\d/.test(event.detail || '')
    )
        return 'warning';
    return 'info';
}

export default function EventViewer({ clients = [] }) {
    const [events, setEvents] = useState([]);
    const [clientEvents, setClientEvents] = useState([]);
    const [error, setError] = useState('');
    const [query, setQuery] = useState('');
    const [level, setLevel] = useState('all');
    const [paused, setPaused] = useState(false);
    const [refresh, setRefresh] = useState(0);
    const [updated, setUpdated] = useState(null);
    const [cleared, setCleared] = useState(() => new Set());
    useEffect(() => {
        if (paused) return;
        const incoming = clients.flatMap((agent) =>
            (agent.Messages || [])
                .map((message) => clientMessageEvent(agent, message))
                .filter(Boolean),
        );
        if (!incoming.length) return;
        setClientEvents((current) => {
            const unique = new Map(current.map((event) => [event.id, event]));
            incoming.forEach((event) => unique.set(event.id, event));
            return [...unique.values()]
                .sort((a, b) => new Date(b.at) - new Date(a.at))
                .slice(0, 300);
        });
    }, [clients, paused]);
    useEffect(() => {
        if (paused) return;
        let cancelled = false,
            timer;
        async function poll() {
            try {
                const result = await apiService.auditHistory();
                if (!Array.isArray(result))
                    throw new Error('Unexpected event response');
                if (!cancelled) {
                    setEvents(
                        [...result].sort(
                            (a, b) =>
                                new Date(b.at) - new Date(a.at) || b.id - a.id,
                        ),
                    );
                    setUpdated(new Date());
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
        };
    }, [paused, refresh]);
    const allEvents = [
        ...events.map((event) => ({ ...event, id: `audit:${event.id}` })),
        ...clientEvents,
    ];
    const combined = allEvents
        .filter((event) => !cleared.has(event.id))
        .sort(
            (a, b) =>
                new Date(b.at) - new Date(a.at) ||
                String(b.id).localeCompare(String(a.id), undefined, {
                    numeric: true,
                }),
        )
        .slice(0, 300);
    const visible = combined.filter(
        (event) =>
            (level === 'all' || levelOf(event) === level) &&
            [
                event.event,
                event.agent_id,
                event.task_id,
                event.source,
                event.detail,
            ]
                .filter(Boolean)
                .join(' ')
                .toLowerCase()
                .includes(query.toLowerCase()),
    );
    return (
        <section className="panel workspace-events" aria-label="Event Viewer">
            <div className="workspace-panel-title">
                <strong>Event Viewer</strong>
                <span>{paused ? 'Paused' : 'Live · 3s refresh'}</span>
                <button
                    type="button"
                    aria-pressed={paused}
                    onClick={() => setPaused((value) => !value)}
                >
                    {paused ? 'Resume events' : 'Pause events'}
                </button>
                <button
                    type="button"
                    disabled={paused}
                    onClick={() => setRefresh((value) => value + 1)}
                >
                    Refresh events
                </button>
                <button
                    type="button"
                    disabled={!combined.length}
                    title="Clear the displayed events. Saved server history is retained."
                    onClick={() =>
                        setCleared(
                            (current) =>
                                new Set([
                                    ...current,
                                    ...allEvents.map((event) => event.id),
                                ]),
                        )
                    }
                >
                    Clear event history
                </button>
                {cleared.size > 0 && (
                    <button type="button" onClick={() => setCleared(new Set())}>
                        Show cleared events
                    </button>
                )}
            </div>
            <div className="event-filters">
                <input
                    aria-label="Search events"
                    placeholder="Search events, sessions…"
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                />
                <select
                    aria-label="Event severity"
                    value={level}
                    onChange={(event) => setLevel(event.target.value)}
                >
                    <option value="all">All levels</option>
                    <option value="info">Info</option>
                    <option value="warning">Warning</option>
                    <option value="error">Error</option>
                </select>
            </div>
            {error && (
                <p role="alert" className="event-feed-error">
                    Event feed unavailable: {error}. Saved entries remain
                    visible; retrying automatically.
                </p>
            )}
            <ol className="event-feed" aria-label="Server events, newest first">
                {visible.map((event) => (
                    <li
                        key={event.id}
                        className={`event-row event-${levelOf(event)}`}
                    >
                        <time dateTime={event.at} title={event.at}>
                            {new Date(event.at).toLocaleString()}
                        </time>
                        <span className="event-level">{levelOf(event)}</span>
                        <span className="event-description">
                            <strong>{event.event}</strong>
                            {event.agent_id && (
                                <span className="mono">
                                    {' '}
                                    · {event.agent_id}
                                </span>
                            )}
                            {event.detail && <span> · {event.detail}</span>}
                            {event.task_id && (
                                <span className="mono"> · {event.task_id}</span>
                            )}
                            {event.source && (
                                <span className="event-source">
                                    {' '}
                                    · source: {event.source}
                                </span>
                            )}
                        </span>
                    </li>
                ))}
                {!visible.length && (
                    <li className="empty-state">
                        {updated || combined.length
                            ? combined.length
                                ? 'No events match your filters.'
                                : cleared.size
                                  ? 'Event view cleared. New events will appear here.'
                                  : 'No server events recorded yet.'
                            : error
                              ? 'Waiting for the event feed.'
                              : 'Loading events…'}
                    </li>
                )}
            </ol>
            <div className="event-footer">
                {visible.length} of {combined.length} events · newest first ·
                latest 300
                {updated && ` · updated ${updated.toLocaleTimeString()}`}
                {' · Clearing affects this view only.'}
            </div>
        </section>
    );
}
