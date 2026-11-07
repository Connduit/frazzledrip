import React, { useEffect, useRef, useState } from 'react';
import ClientInfo from './ClientInfo.jsx';
import HistoryPanel from './HistoryPanel.jsx';
import CommandResults from './CommandResults.jsx';
import MessageList from './MessageList.jsx';
import SessionLink from './SessionLink.jsx';
import TabContainer from './TabContainer.jsx';
import ToolCatalog from './ToolCatalog.jsx';
import { systemOf } from '../lib/sessions.js';

export default function DetailsTools({
    agent,
    clients,
    clientInfo,
    messages,
    shared,
    activeTab,
    setActiveTab,
    setStatus,
    refreshMessages,
    onLinked,
    onClose,
}) {
    const dialog = useRef(null);
    const [section, setSection] = useState('overview');
    useEffect(() => {
        const previous = document.activeElement;
        dialog.current.showModal();
        return () => {
            dialog.current?.close();
            previous?.focus?.({ preventScroll: true });
        };
    }, []);
    const sections = [
        ['overview', 'Overview'],
        ['tasks', 'Task history'],
        ['catalog', 'Attacks & scripts'],
        ['controls', 'Commands & transfers'],
        ['messages', 'Messages'],
    ];
    return (
        <dialog
            className="tools-dialog"
            ref={dialog}
            aria-labelledby="tools-title"
            onCancel={(event) => {
                event.preventDefault();
                onClose();
            }}
        >
            <header className="tools-dialog-header">
                <div>
                    <span className="eyebrow">SESSION WORKSPACE</span>
                    <h2 id="tools-title">Details &amp; tools</h2>
                </div>
                <select
                    aria-label="Tools session"
                    value={agent?.IPAddress || ''}
                    onChange={(event) =>
                        shared.onClientSelect(event.target.value)
                    }
                >
                    <option value="">Select session…</option>
                    {clients.map((client) => (
                        <option key={client.ID} value={client.IPAddress}>
                            {systemOf(client).hostname || client.ID} ·{' '}
                            {client.IsActive ? 'Active' : 'Offline'} ·{' '}
                            {client.ID}
                        </option>
                    ))}
                </select>
                <button onClick={onClose} aria-label="Close details and tools">
                    Close ✕
                </button>
            </header>
            <div className="tools-dialog-layout">
                <nav
                    className="tools-navigation"
                    aria-label="Details and tools sections"
                >
                    {sections.map(([id, name]) => (
                        <button
                            key={id}
                            aria-current={section === id ? 'page' : undefined}
                            onClick={() => setSection(id)}
                        >
                            {name}
                        </button>
                    ))}
                </nav>
                <div
                    className="tools-dialog-content"
                    key={`${agent?.ID || 'none'}-${section}`}
                >
                    {section === 'catalog' ? (
                        <ToolCatalog agent={agent} />
                    ) : !agent ? (
                        <div className="empty-state">
                            Select a session to see its details and tools.
                        </div>
                    ) : (
                        <>
                            {section === 'overview' && (
                                <>
                                    <h3>
                                        {systemOf(agent).hostname || agent.ID}
                                    </h3>
                                    <p className="mono muted">
                                        {agent.ID} · {agent.IPAddress} ·{' '}
                                        {agent.IsActive ? 'Active' : 'Offline'}
                                    </p>
                                    <ClientInfo
                                        clientInfo={
                                            clientInfo || agent.SystemInfo
                                        }
                                    />
                                    {agent.IsActive && (
                                        <SessionLink
                                            live={agent}
                                            clients={clients}
                                            onLinked={onLinked}
                                        />
                                    )}
                                </>
                            )}
                            {section === 'tasks' && (
                                <>
                                    <HistoryPanel agentID={agent.ID} />
                                    <CommandResults
                                        results={agent.CommandResults || []}
                                    />
                                </>
                            )}
                            {section === 'controls' &&
                                (agent.IsActive ? (
                                    <TabContainer
                                        activeTab={activeTab}
                                        setActiveTab={setActiveTab}
                                        {...shared}
                                        setStatus={setStatus}
                                        refreshMessages={refreshMessages}
                                    />
                                ) : (
                                    <p className="status warning">
                                        This session is offline. Live controls
                                        require a connection.
                                    </p>
                                ))}
                            {section === 'messages' && (
                                <>
                                    <div className="panel-heading">
                                        <h3>Session messages</h3>
                                        <button onClick={refreshMessages}>
                                            Refresh messages
                                        </button>
                                    </div>
                                    <MessageList messages={messages} />
                                </>
                            )}
                        </>
                    )}
                </div>
            </div>
        </dialog>
    );
}
