import React, { useCallback, useEffect, useState, useRef } from 'react';
import { apiService, utils } from '../services/api.js';
import { systemOf, validLocation, hostOf } from '../lib/sessions.js';
import CommandConsole from './CommandConsole.jsx';
import FileExplorer from './FileExplorer.jsx';
import DetailsTools from './DetailsTools.jsx';
import HistoryPanel from './HistoryPanel.jsx';
import EventViewer from './EventViewer.jsx';
import ClientList from './ClientList.jsx';
import SessionGraph from './SessionGraph.jsx';
import SessionContextMenu from './SessionContextMenu.jsx';
import SessionMap from './SessionMap.jsx';
import WorkspaceLayout from './WorkspaceLayout.jsx';
import WorkspaceTabs from './WorkspaceTabs.jsx';
import ListenersTab from './ListenersTab.jsx';
import useListenerStatus from '../lib/useListenerStatus.js';

const locationKey = 'frazzledrip.session-locations.v1';
function readLocations() {
    try {
        const data = JSON.parse(localStorage.getItem(locationKey) || '{}');
        return Object.fromEntries(
            Object.entries(data).filter(([, value]) => validLocation(value)),
        );
    } catch {
        return {};
    }
}

export default function App() {
    const [status, setStatus] = useState('');
    const [clients, setClients] = useState([]);
    const [selectedClient, setSelectedClient] = useState('');
    const selectedAgentID = useRef('');
    const [clientInfo, setClientInfo] = useState(null);
    const [clientMessages, setClientMessages] = useState([]);
    const [messagesClient, setMessagesClient] = useState('');
    const [activeTab, setActiveTab] = useState('commands');
    const [workspaceTab, setWorkspaceTab] = useState('events');
    const [openWorkspaceTabs, setOpenWorkspaceTabs] = useState(['events']);
    const [sessionView, setSessionView] = useState('sessions');
    const [sessionMenu, setSessionMenu] = useState(null);
    const closeSessionMenu = useCallback(() => setSessionMenu(null), []);
    const [showMap, setShowMap] = useState(false);
    const [showTools, setShowTools] = useState(false);
    const [showAudit, setShowAudit] = useState(false);
    const [locations, setLocations] = useState(readLocations);
    const [lastRefresh, setLastRefresh] = useState(null);
    const [fetchError, setFetchError] = useState('');
    const [refreshVersion, setRefreshVersion] = useState(0);
    const [messageVersion, setMessageVersion] = useState(0);
    const [starting, setStarting] = useState(false);
    const {
        status: listenerState,
        error: listenerError,
        checkedAt: listenerCheckedAt,
        checking: listenerChecking,
        refresh: refreshListener,
    } = useListenerStatus();
    const [disconnecting, setDisconnecting] = useState(false);

    useEffect(() => {
        let cancelled = false;
        let timer;
        async function poll() {
            try {
                const data = await apiService.getClients();
                if (!Array.isArray(data))
                    throw new Error('Unexpected session response');
                if (cancelled) return;
                setClients(data);
                setLastRefresh(new Date());
                setFetchError('');
                setSelectedClient(
                    (current) =>
                        data.find((c) => c.ID === selectedAgentID.current)
                            ?.IPAddress ||
                        (data.some((c) => c.IPAddress === current)
                            ? current
                            : ''),
                );
            } catch (error) {
                if (!cancelled) setFetchError(error.message);
            } finally {
                if (!cancelled) timer = setTimeout(poll, 5000);
            }
        }
        poll();
        return () => {
            cancelled = true;
            clearTimeout(timer);
        };
    }, [refreshVersion]);

    useEffect(() => {
        let cancelled = false;
        let timer;
        setClientInfo(null);
        setClientMessages([]);
        setMessagesClient('');
        if (!selectedClient) return;
        async function poll() {
            try {
                const [info, messages] = await Promise.all([
                    apiService.getClientInfo(selectedClient),
                    apiService.getClientMessages(selectedClient),
                ]);
                if (!cancelled) {
                    setClientInfo(info);
                    setClientMessages(messages.messages || []);
                    setMessagesClient(selectedClient);
                }
            } catch (error) {
                if (!cancelled) setStatus('Error: ' + error.message);
            } finally {
                if (!cancelled) timer = setTimeout(poll, 2000);
            }
        }
        poll();
        return () => {
            cancelled = true;
            clearTimeout(timer);
        };
    }, [selectedClient, messageVersion]);

    async function controlListener(action) {
        setStarting(true);
        try {
            const data = await (action === 'start'
                ? apiService.startListener()
                : apiService.stopListener());
            const verified = await refreshListener();
            setStatus(
                verified.running === (action === 'start')
                    ? data.status
                    : 'Warning: backend listener status does not match the requested action. Check the Listeners tab.',
            );
            setRefreshVersion((v) => v + 1);
        } catch (error) {
            setStatus('Error: ' + error.message);
        } finally {
            setStarting(false);
        }
    }
    function openWorkspaceTab(id) {
        setOpenWorkspaceTabs((current) =>
            current.includes(id) ? current : [...current, id],
        );
        setWorkspaceTab(id);
    }
    async function disconnectSession(id) {
        setDisconnecting(true);
        try {
            setStatus((await apiService.disconnectSession(id)).status);
            setRefreshVersion((v) => v + 1);
            setMessageVersion((v) => v + 1);
        } catch (error) {
            setStatus('Error: ' + error.message);
        } finally {
            setDisconnecting(false);
        }
    }
    function saveLocation(host, value) {
        const next = { ...locations };
        if (value) next[host] = value;
        else delete next[host];
        setLocations(next);
        try {
            localStorage.setItem(locationKey, JSON.stringify(next));
            setStatus('Location updated in this browser.');
        } catch {
            setStatus(
                'Warning: location is available for this visit only; browser storage is unavailable.',
            );
        }
    }
    const selected = clients.find((c) => c.IPAddress === selectedClient);
    const active = clients.filter((c) => c.IsActive).length;
    const hosts = new Set(
        clients.map((c) => systemOf(c).hostname).filter(Boolean),
    ).size;
    function selectSession(address) {
        selectedAgentID.current =
            clients.find((client) => client.IPAddress === address)?.ID || '';
        setSelectedClient(address);
    }
    function showSessionMenu(event, client) {
        event.preventDefault();
        event.stopPropagation();
        const bounds = event.currentTarget.getBoundingClientRect();
        const anchor = event.currentTarget.matches('button, [role="button"]')
            ? event.currentTarget
            : event.currentTarget.querySelector('button');
        selectSession(client.IPAddress);
        setSessionMenu({
            id: client.ID,
            anchor,
            x: event.clientX || bounds.left,
            y: event.clientY || bounds.bottom,
        });
    }
    const menuClient =
        sessionMenu && clients.find((client) => client.ID === sessionMenu.id);
    useEffect(() => {
        if (sessionMenu && !menuClient) closeSessionMenu();
    }, [sessionMenu, menuClient, closeSessionMenu]);
    const shared = {
        clients,
        selectedClient,
        onClientSelect: selectSession,
        onSessionContextMenu: showSessionMenu,
    };
    const selectedName = selected
        ? systemOf(selected).hostname || selected.ID
        : 'No session selected';
    const messages = messagesClient === selectedClient ? clientMessages : [];
    return (
        <div className="container workspace-shell">
            <header className="workspace-toolbar">
                <div className="workspace-brand">
                    <span className="brand-mark">F</span>
                    <strong>frazzledrip</strong>
                    <span className="muted">/ workspace</span>
                </div>
                <div className="header-actions">
                    <span
                        className={`badge ${listenerState?.running && !listenerError ? 'live' : ''}`}
                    >
                        <i />
                        {listenerError || listenerState == null
                            ? 'Listener unverified'
                            : listenerState.running
                              ? `Listening · ${listenerState.port}`
                              : 'Listener stopped'}
                    </span>
                    <button
                        disabled={
                            starting ||
                            !!listenerError ||
                            !listenerState ||
                            listenerState.running
                        }
                        onClick={() => controlListener('start')}
                    >
                        Start listener
                    </button>
                    <button
                        disabled={
                            starting ||
                            !!listenerError ||
                            !listenerState?.running
                        }
                        onClick={() => controlListener('stop')}
                    >
                        Stop listener
                    </button>
                    <button onClick={() => setRefreshVersion((v) => v + 1)}>
                        Refresh sessions
                    </button>
                    <button
                        aria-expanded={showAudit}
                        onClick={() => setShowAudit((value) => !value)}
                    >
                        Audit
                    </button>
                    <button
                        aria-expanded={showMap}
                        onClick={() => setShowMap((value) => !value)}
                    >
                        Map
                    </button>
                    <button
                        aria-haspopup="dialog"
                        aria-expanded={showTools}
                        onClick={() => setShowTools((value) => !value)}
                    >
                        Details &amp; tools
                    </button>
                </div>
            </header>
            <nav
                className="workspace-feature-actions"
                aria-label="Workspace features"
            >
                <span>Open view</span>
                {[
                    { id: 'events', label: 'Event Viewer' },
                    { id: 'console', label: 'Console' },
                    { id: 'files', label: 'Files' },
                    { id: 'listeners', label: 'Listeners' },
                    { id: 'processes', label: 'Process Browser' },
                ].map((feature) => (
                    <button
                        type="button"
                        key={feature.id}
                        aria-label={`Open ${feature.label} tab`}
                        aria-pressed={workspaceTab === feature.id}
                        onClick={() => openWorkspaceTab(feature.id)}
                    >
                        {feature.label}
                    </button>
                ))}
            </nav>
            <div className="workspace-context">
                <span>
                    <strong>{clients.length}</strong> sessions{' '}
                    <span className="context-divider">/</span>{' '}
                    <strong className="accent">{active}</strong> active{' '}
                    <span className="context-divider">/</span> {hosts} hosts
                </span>
                <span className="workspace-selection">
                    {selected
                        ? `${selectedName} · ${selected.IsActive ? 'Active' : 'Offline'}`
                        : 'Select a session in the table or graph'}
                    {selected?.IsActive && (
                        <button
                            disabled={disconnecting}
                            onClick={() => disconnectSession(selected.ID)}
                        >
                            {disconnecting
                                ? 'Disconnecting…'
                                : 'Disconnect session'}
                        </button>
                    )}
                </span>
                <span className="muted">
                    {fetchError
                        ? 'Feed unavailable'
                        : lastRefresh
                          ? `Updated ${lastRefresh.toLocaleTimeString()}`
                          : 'Connecting…'}
                </span>
            </div>
            {fetchError && (
                <div className="status error" role="alert">
                    Session refresh failed: {fetchError}. Displayed sessions may
                    be stale; retrying automatically.
                </div>
            )}
            {showAudit && (
                <section className="panel section workspace-tools">
                    <HistoryPanel audit />
                </section>
            )}
            <WorkspaceLayout>
                <WorkspaceTabs
                    label="Sessions window"
                    tabListLabel="Session views"
                    storageKey="frazzledrip.session-tab-order.v1"
                    className="workspace-inventory"
                    activeTab={sessionView}
                    onTabChange={setSessionView}
                    tabs={[
                        {
                            id: 'sessions',
                            label: 'Session table',
                            content: (
                                <div className="workspace-sessions">
                                    <ClientList {...shared} />
                                </div>
                            ),
                        },
                        {
                            id: 'graph',
                            label: 'Session graph',
                            content: (
                                <div className="workspace-graph">
                                    <SessionGraph {...shared} />
                                </div>
                            ),
                        },
                    ]}
                />
                <WorkspaceTabs
                    activeTab={workspaceTab}
                    onTabChange={setWorkspaceTab}
                    tabs={[
                        {
                            id: 'console',
                            label: 'Console',
                            content: (
                                <section
                                    className="panel workspace-console"
                                    aria-label="Console pane"
                                >
                                    <div className="workspace-panel-title">
                                        <strong>Console</strong>
                                        <span>{selectedName}</span>
                                        {selected?.BasicInfoStatus ===
                                            'collecting' && (
                                            <span className="accent">
                                                Collecting initial host info…
                                            </span>
                                        )}
                                    </div>
                                    {selected ? (
                                        <CommandConsole
                                            agent={selected}
                                            messages={messages}
                                            active={workspaceTab === 'console'}
                                            autoFocus={false}
                                            setStatus={setStatus}
                                            refreshMessages={() =>
                                                setMessageVersion((v) => v + 1)
                                            }
                                        />
                                    ) : (
                                        <div className="empty-state">
                                            <strong>No session selected</strong>
                                            <p>
                                                Select an agent to use its
                                                custom-message console.
                                            </p>
                                        </div>
                                    )}
                                </section>
                            ),
                        },
                        {
                            id: 'events',
                            label: 'Event Viewer',
                            content: (
                                <EventViewer
                                    clients={clients.map((client) =>
                                        client.ID === selected?.ID &&
                                        messagesClient === selectedClient
                                            ? { ...client, Messages: messages }
                                            : client,
                                    )}
                                />
                            ),
                        },
                        {
                            id: 'files',
                            label: 'Files',
                            content: (
                                <section
                                    className="panel workspace-files"
                                    aria-label="Files pane"
                                >
                                    <div className="workspace-panel-title">
                                        <strong>Files</strong>
                                        <span>{selectedName}</span>
                                        <span className="muted">Read only</span>
                                    </div>
                                    <div className="workspace-panel-body">
                                        {selected ? (
                                            <FileExplorer
                                                key={selected.ID}
                                                agent={selected}
                                            />
                                        ) : (
                                            <div className="empty-state">
                                                <strong>
                                                    No session selected
                                                </strong>
                                                <p>
                                                    Select an agent to browse
                                                    its directories and file
                                                    metadata.
                                                </p>
                                            </div>
                                        )}
                                    </div>
                                </section>
                            ),
                        },
                        {
                            id: 'listeners',
                            label: 'Listeners',
                            content: (
                                <ListenersTab
                                    status={listenerState}
                                    error={listenerError}
                                    checkedAt={listenerCheckedAt}
                                    checking={listenerChecking}
                                    pending={starting}
                                    onControl={controlListener}
                                    onRefresh={() => {
                                        refreshListener().catch(() => {});
                                    }}
                                />
                            ),
                        },
                        {
                            id: 'processes',
                            label: 'Process Browser',
                            content: (
                                <section
                                    className="panel"
                                    aria-label="Process Browser"
                                >
                                    <div className="workspace-panel-title">
                                        <strong>Process Browser</strong>
                                        <span>{selectedName}</span>
                                    </div>
                                    <div className="empty-state">
                                        <strong>
                                            Process Browser is not implemented
                                            yet.
                                        </strong>
                                        <p>
                                            No process information has been
                                            requested from the client.
                                        </p>
                                    </div>
                                </section>
                            ),
                        },
                    ]
                        .filter((tab) => openWorkspaceTabs.includes(tab.id))
                        .sort(
                            (a, b) =>
                                openWorkspaceTabs.indexOf(a.id) -
                                openWorkspaceTabs.indexOf(b.id),
                        )}
                />
            </WorkspaceLayout>
            {status && (
                <div
                    role="status"
                    className={`workspace-notice ${utils.getStatusClass(status)}`}
                >
                    {status}
                </div>
            )}
            {showMap && (
                <div className="workspace-map">
                    <SessionMap
                        {...shared}
                        locations={locations}
                        onSaveLocation={saveLocation}
                    />
                </div>
            )}
            {showTools && (
                <DetailsTools
                    agent={selected}
                    clients={clients}
                    clientInfo={
                        messagesClient === selectedClient ? clientInfo : null
                    }
                    messages={messages}
                    shared={shared}
                    activeTab={activeTab}
                    setActiveTab={setActiveTab}
                    setStatus={setStatus}
                    refreshMessages={() =>
                        setMessageVersion((value) => value + 1)
                    }
                    onClose={() => setShowTools(false)}
                    onLinked={(saved) => {
                        const savedLocation =
                            locations[hostOf(saved.IPAddress)];
                        if (
                            savedLocation &&
                            !locations[hostOf(selected.IPAddress)]
                        )
                            saveLocation(
                                hostOf(selected.IPAddress),
                                savedLocation,
                            );
                        setStatus('Saved session linked.');
                        setRefreshVersion((value) => value + 1);
                        setMessageVersion((value) => value + 1);
                    }}
                />
            )}
            {sessionMenu && menuClient && (
                <SessionContextMenu
                    {...sessionMenu}
                    client={menuClient}
                    onClose={closeSessionMenu}
                    onOpen={(tab) => {
                        selectSession(menuClient.IPAddress);
                        openWorkspaceTab(tab);
                        closeSessionMenu();
                        requestAnimationFrame(() =>
                            document
                                .querySelector(
                                    `[role="tab"][id$="-tab-${tab}"]`,
                                )
                                ?.focus(),
                        );
                    }}
                />
            )}
            <footer className="workspace-footer">
                <span>FRAZZLEDRIP</span>
                <span>
                    Session facts, task history, and audit events are saved on
                    the server.
                </span>
            </footer>
        </div>
    );
}
