import React, { useEffect, useRef, useState } from 'react';
import { apiService } from '../services/api.js';
import {
    isConsoleMessage,
    parseConsoleCommand,
} from '../lib/consoleCommands.js';

const emptySession = () => ({
    entries: [],
    history: [],
    seen: new Set(),
    draft: '',
});

export default function CommandConsole({
    agent,
    messages,
    active,
    setStatus,
    refreshMessages,
}) {
    const [sessions, setSessions] = useState({});
    const [sending, setSending] = useState(false);
    const input = useRef(null);
    const output = useRef(null);
    const followOutput = useRef(true);
    const navigation = useRef({ index: -1, draft: '' });
    const sequence = useRef(0);
    const session = sessions[agent.ID] || emptySession();
    function update(id, change) {
        setSessions((previous) => ({
            ...previous,
            [id]: change(previous[id] || emptySession()),
        }));
    }
    useEffect(() => {
        update(agent.ID, (current) => {
            const seen = new Set(current.seen);
            const entries = [];
            for (const message of messages) {
                if (!isConsoleMessage(message)) continue;
                const key = `${message.timestamp}:${message.message_id}:${message.type}`;
                if (seen.has(key)) continue;
                seen.add(key);
                entries.push({
                    id: key,
                    type:
                        message.type === 'ERROR_REPORT' ? 'error' : 'response',
                    text: message.content || '(Empty response)',
                });
            }
            return {
                ...current,
                entries: [...current.entries, ...entries].slice(-300),
                seen: new Set([...seen].slice(-1000)),
            };
        });
    }, [agent.ID, messages]);
    useEffect(() => {
        navigation.current = { index: -1, draft: '' };
        followOutput.current = true;
        if (active) input.current?.focus({ preventScroll: true });
    }, [agent.ID, active]);
    useEffect(() => {
        if (active && !sending) input.current?.focus({ preventScroll: true });
    }, [sending, active, agent.ID]);
    useEffect(() => {
        if (active && followOutput.current && output.current)
            output.current.scrollTop = output.current.scrollHeight;
    }, [session.entries[session.entries.length - 1]?.id, active, agent.ID]);
    function clear() {
        navigation.current = { index: -1, draft: '' };
        update(agent.ID, (current) => ({
            ...current,
            entries: [],
            history: [],
        }));
    }
    async function submit(event) {
        event.preventDefault();
        const command = session.draft.trim();
        if (!command || sending) return;
        const id = agent.ID;
        const entryID = `local-${++sequence.current}`;
        navigation.current = { index: -1, draft: '' };
        followOutput.current = true;
        update(id, (current) => ({
            ...current,
            draft: '',
            history: [...current.history, command].slice(-100),
            entries: [
                ...current.entries,
                { id: entryID, type: 'command', text: command },
            ].slice(-300),
        }));
        const parsed = parseConsoleCommand(command);
        if (parsed.kind === 'clear') {
            clear();
            return;
        }
        if (parsed.kind !== 'request' || !agent.IsActive) {
            const text =
                parsed.kind === 'help'
                    ? 'Custom commands:\nSYS_INFO — native OS, memory, and processor information\nhelp — show this list\nclear — clear console history\nShell commands are not accepted.'
                    : parsed.kind === 'request'
                      ? 'This session is offline. Reconnect to send SYS_INFO.'
                      : 'Unknown custom command. Type help for available commands. Shell commands are not accepted.';
            update(id, (current) => ({
                ...current,
                entries: [
                    ...current.entries,
                    {
                        id: `${entryID}-local`,
                        type: parsed.kind === 'help' ? 'response' : 'error',
                        text,
                    },
                ].slice(-300),
            }));
            return;
        }
        setSending(true);
        try {
            const result = await apiService.mySendCommand(
                agent.IPAddress,
                'systeminfo',
            );
            if (!result.status?.startsWith('✅'))
                throw new Error(
                    result.error || result.status || 'Command was not accepted',
                );
            setStatus('✅ SYS_INFO sent - waiting for response...');
            refreshMessages();
        } catch (error) {
            update(id, (current) => ({
                ...current,
                entries: [
                    ...current.entries,
                    {
                        id: `${entryID}-error`,
                        type: 'error',
                        text: error.message,
                    },
                ].slice(-300),
            }));
            setStatus('Error: ' + error.message);
        } finally {
            setSending(false);
        }
    }
    function keyDown(event) {
        if (event.ctrlKey && event.key.toLowerCase() === 'l') {
            event.preventDefault();
            clear();
            return;
        }
        if (
            !['ArrowUp', 'ArrowDown'].includes(event.key) ||
            !session.history.length
        )
            return;
        event.preventDefault();
        const nav = navigation.current;
        if (nav.index === -1) nav.draft = session.draft;
        nav.index =
            event.key === 'ArrowUp'
                ? Math.min(session.history.length - 1, nav.index + 1)
                : Math.max(-1, nav.index - 1);
        const draft =
            nav.index === -1
                ? nav.draft
                : session.history[session.history.length - 1 - nav.index];
        update(agent.ID, (current) => ({ ...current, draft }));
    }
    return (
        <section className="terminal-panel" aria-label="Custom command console">
            <div className="terminal-toolbar">
                <span className="mono">
                    {agent.SystemInfo?.hostname || agent.ID} · {agent.IPAddress}
                </span>
                <button
                    type="button"
                    aria-label="Clear console history"
                    title="Clear this session's custom-command output and recalled commands. Saved server history is retained."
                    onClick={() => {
                        clear();
                        input.current?.focus({ preventScroll: true });
                    }}
                >
                    Clear history
                </button>
            </div>
            <p className="terminal-hint">
                Custom message commands only: SYS_INFO · help · clear. Shell
                commands are not accepted. ↑ / ↓ recalls commands; Ctrl+L clears
                this view.
            </p>
            <div
                className="terminal-output"
                ref={output}
                role="log"
                aria-label="Console output"
                tabIndex="0"
                onScroll={(event) => {
                    const el = event.currentTarget;
                    followOutput.current =
                        el.scrollHeight - el.scrollTop - el.clientHeight < 40;
                }}
            >
                {!session.entries.length && (
                    <div className="muted">
                        Type help for available custom commands. SYS_INFO
                        requests native host information.
                    </div>
                )}
                {session.entries.map((entry) => (
                    <pre
                        key={entry.id}
                        className={`terminal-entry terminal-${entry.type}`}
                    >
                        {entry.type === 'command'
                            ? '> '
                            : entry.type === 'error'
                              ? 'ERROR: '
                              : ''}
                        {entry.text}
                    </pre>
                ))}
            </div>
            <form className="terminal-prompt" onSubmit={submit}>
                <span aria-hidden="true">&gt;</span>
                <input
                    ref={input}
                    aria-label="Console command"
                    value={session.draft}
                    maxLength={128}
                    placeholder="SYS_INFO, help, clear"
                    autoComplete="off"
                    spellCheck="false"
                    disabled={sending}
                    onKeyDown={keyDown}
                    onChange={(event) => {
                        navigation.current.index = -1;
                        update(agent.ID, (current) => ({
                            ...current,
                            draft: event.target.value,
                        }));
                    }}
                />
                <button
                    type="submit"
                    disabled={sending || !session.draft.trim()}
                >
                    {sending ? 'Sending…' : 'Run'}
                </button>
            </form>
        </section>
    );
}
