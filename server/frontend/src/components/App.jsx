import React, { useEffect, useState, useRef } from 'react';
import { apiService, utils } from '../services/api.js';
import { systemOf, validLocation, hostOf } from '../lib/sessions.js';
import TerminalTab from './TerminalTab.jsx';
import FileExplorer from './FileExplorer.jsx';
import SessionLink from './SessionLink.jsx';
import ClientList from './ClientList.jsx';
import SessionGraph from './SessionGraph.jsx';
import SessionMap from './SessionMap.jsx';
import CommandResults from './CommandResults.jsx';
import ClientInfo from './ClientInfo.jsx';
import MessageList from './MessageList.jsx';
import TabContainer from './TabContainer.jsx';

const locationKey = 'frazzledrip.session-locations.v1';
function readLocations() {
    try { const data = JSON.parse(localStorage.getItem(locationKey) || '{}');
        return Object.fromEntries(Object.entries(data).filter(([, value]) => validLocation(value)));
    } catch { return {}; }
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
    const [showMap, setShowMap] = useState(false);
    const [showTools, setShowTools] = useState(false);
    const [locations, setLocations] = useState(readLocations);
    const [lastRefresh, setLastRefresh] = useState(null);
    const [fetchError, setFetchError] = useState('');
    const [refreshVersion, setRefreshVersion] = useState(0);
    const [messageVersion, setMessageVersion] = useState(0);
    const [starting, setStarting] = useState(false);
    const [listenerState, setListenerState] = useState(null);
    const [disconnecting, setDisconnecting] = useState(false);

    useEffect(() => {
        let cancelled = false;
        let timer;
        async function poll() {
            try {
                const [data, listener] = await Promise.all([apiService.getClients(), apiService.getListenerStatus()]);
                if (!Array.isArray(data)) throw new Error('Unexpected session response');
                if (cancelled) return;
                setListenerState(listener); setClients(data); setLastRefresh(new Date()); setFetchError('');
                setSelectedClient(current => data.find(c => c.ID === selectedAgentID.current)?.IPAddress || (data.some(c => c.IPAddress === current) ? current : ''));
            } catch (error) { if (!cancelled) setFetchError(error.message); }
            finally { if (!cancelled) timer = setTimeout(poll, 5000); }
        }
        poll();
        return () => { cancelled = true; clearTimeout(timer); };
    }, [refreshVersion]);

    useEffect(() => {
        let cancelled = false;
        let timer;
        setClientInfo(null); setClientMessages([]); setMessagesClient('');
        if (!selectedClient) return;
        async function poll() {
            try {
                const [info, messages] = await Promise.all([apiService.getClientInfo(selectedClient), apiService.getClientMessages(selectedClient)]);
                if (!cancelled) { setClientInfo(info); setClientMessages(messages.messages || []); setMessagesClient(selectedClient); }
            } catch (error) { if (!cancelled) setStatus('Error: ' + error.message); }
            finally { if (!cancelled) timer = setTimeout(poll, 2000); }
        }
        poll();
        return () => { cancelled = true; clearTimeout(timer); };
    }, [selectedClient, messageVersion]);

    async function controlListener(action) {
        setStarting(true);
        try {
            const data = await (action === 'start' ? apiService.startListener() : apiService.stopListener());
            setStatus(data.status); setListenerState(await apiService.getListenerStatus());
            setRefreshVersion(v => v + 1);
        } catch (error) { setStatus('Error: ' + error.message); }
        finally { setStarting(false); }
    }
    async function disconnectSession(id) {
        setDisconnecting(true);
        try {
            setStatus((await apiService.disconnectSession(id)).status);
            setRefreshVersion(v => v + 1); setMessageVersion(v => v + 1);
        } catch (error) { setStatus('Error: ' + error.message); }
        finally { setDisconnecting(false); }
    }
    function saveLocation(host, value) {
        const next = { ...locations };
        if (value) next[host] = value; else delete next[host];
        setLocations(next);
        try { localStorage.setItem(locationKey, JSON.stringify(next)); setStatus('Location updated in this browser.'); }
        catch { setStatus('Warning: location is available for this visit only; browser storage is unavailable.'); }
    }
    const selected = clients.find(c => c.IPAddress === selectedClient);
    const active = clients.filter(c => c.IsActive).length;
    const hosts = new Set(clients.map(c => systemOf(c).hostname).filter(Boolean)).size;
    const shared = { clients, selectedClient, onClientSelect: address => { selectedAgentID.current = clients.find(c => c.IPAddress === address)?.ID || ''; setSelectedClient(address); } };
    const selectedName = selected ? systemOf(selected).hostname || selected.ID : 'No session selected';
    const messages = messagesClient === selectedClient ? clientMessages : [];
    return <div className="container workspace-shell">
        <header className="workspace-toolbar">
            <div className="workspace-brand"><span className="brand-mark">F</span><strong>frazzledrip</strong><span className="muted">/ workspace</span></div>
            <div className="header-actions"><span className={`badge ${listenerState?.running ? 'live' : ''}`}><i />{listenerState == null ? 'Listener unknown' : listenerState.running ? `Listening · ${listenerState.port}` : 'Listener stopped'}</span>
                <button disabled={starting || listenerState?.running} onClick={() => controlListener('start')}>Start listener</button>
                <button disabled={starting || !listenerState?.running} onClick={() => controlListener('stop')}>Stop listener</button>
                <button onClick={() => setRefreshVersion(v => v + 1)}>Refresh sessions</button>
                <button aria-expanded={showMap} onClick={() => setShowMap(value => !value)}>Map</button>
                <button aria-expanded={showTools} onClick={() => setShowTools(value => !value)}>Details &amp; tools</button>
            </div>
        </header>
        <div className="workspace-context"><span><strong>{clients.length}</strong> sessions <span className="context-divider">/</span> <strong className="accent">{active}</strong> active <span className="context-divider">/</span> {hosts} hosts</span>
            <span className="workspace-selection">{selected ? `${selectedName} · ${selected.IsActive ? 'Active' : 'Offline'}` : 'Select a session in the table or graph'}
                {selected?.IsActive && <button disabled={disconnecting} onClick={() => disconnectSession(selected.ID)}>{disconnecting ? 'Disconnecting…' : 'Disconnect session'}</button>}</span>
            <span className="muted">{fetchError ? 'Feed unavailable' : lastRefresh ? `Updated ${lastRefresh.toLocaleTimeString()}` : 'Connecting…'}</span></div>
        {fetchError && <div className="status error" role="alert">Session refresh failed: {fetchError}. Displayed sessions may be stale; retrying automatically.</div>}
        <main className="workspace-grid">
            <div className="workspace-graph"><SessionGraph {...shared} /></div>
            <div className="workspace-sessions"><ClientList {...shared} /></div>
            <section className="panel workspace-console" aria-label="Console pane"><div className="workspace-panel-title"><strong>Console</strong><span>{selectedName}</span>{selected?.BasicInfoStatus === 'collecting' && <span className="accent">Collecting initial host info…</span>}</div>
                {selected ? <TerminalTab agent={selected} messages={messages} active={true} setStatus={setStatus} refreshMessages={() => setMessageVersion(v => v + 1)} /> : <div className="empty-state"><strong>No session selected</strong><p>Select an agent to open its command prompt and output.</p></div>}
            </section>
            <section className="panel workspace-files" aria-label="Files pane"><div className="workspace-panel-title"><strong>Files</strong><span>{selectedName}</span><span className="muted">Read only</span></div>
                <div className="workspace-panel-body">{selected ? <FileExplorer key={selected.ID} agent={selected} /> : <div className="empty-state"><strong>No session selected</strong><p>Select an agent to browse its directories and file metadata.</p></div>}</div>
            </section>
        </main>
        {status && <div role="status" className={`workspace-notice ${utils.getStatusClass(status)}`}>{status}</div>}
        {showMap && <div className="workspace-map"><SessionMap {...shared} locations={locations} onSaveLocation={saveLocation} /></div>}
        {showTools && <section className="panel section workspace-tools" aria-label="Session details and tools"><h2>{selectedName}</h2>
            {selected ? <>
                {selected.IsActive && <SessionLink key={selected.ID} live={selected} clients={clients} onLinked={saved => {
                    const savedLocation = locations[hostOf(saved.IPAddress)];
                    if (savedLocation && !locations[hostOf(selected.IPAddress)]) saveLocation(hostOf(selected.IPAddress), savedLocation);
                    setStatus('Saved session linked.'); setRefreshVersion(v => v + 1); setMessageVersion(v => v + 1);
                }} />}
                <ClientInfo clientInfo={clientInfo} /><CommandResults results={selected.CommandResults || []} /><MessageList messages={messages} />
                {selected.IsActive && <TabContainer activeTab={activeTab} setActiveTab={setActiveTab} {...shared} setStatus={setStatus} refreshMessages={() => setMessageVersion(v => v + 1)} />}
            </> : <p className="muted">Select a session to inspect its details and tools.</p>}
        </section>}
        <footer className="workspace-footer"><span>FRAZZLEDRIP</span><span>Offline sessions retain collected information while the server is running.</span></footer>
    </div>;
}
