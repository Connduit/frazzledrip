import React, { useEffect, useState, useRef } from 'react';
import { apiService, utils } from '../services/api.js';
import { systemOf, validLocation, hostOf } from '../lib/sessions.js';
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
    const [activeTab, setActiveTab] = useState('commands');
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
        setClientInfo(null); setClientMessages([]);
        if (!selectedClient) return;
        async function poll() {
            try {
                const [info, messages] = await Promise.all([apiService.getClientInfo(selectedClient), apiService.getClientMessages(selectedClient)]);
                if (!cancelled) { setClientInfo(info); setClientMessages(messages.messages || []); }
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
    return <div className="container">
        <header className="app-header"><div><span className="eyebrow">FRAZZLEDRIP / CONTROLLER</span><h1>Session workspace</h1><p className="muted">Session inventory, topology, and geography.</p></div>
            <div className="header-actions"><button onClick={() => setRefreshVersion(v => v + 1)}>Refresh sessions</button><span className={`badge ${listenerState?.running ? 'live' : ''}`}><i />{listenerState == null ? 'Listener unknown' : listenerState.running ? `Listening · ${listenerState.port}` : 'Listener stopped'}</span>
                <button className="primary" disabled={starting || listenerState?.running} onClick={() => controlListener('start')}>Start listener</button>
                <button disabled={starting || !listenerState?.running} onClick={() => controlListener('stop')}>Stop listener</button></div></header>
        <div className="summary-grid"><div className="metric"><span>Known sessions</span><strong>{clients.length.toString().padStart(2, '0')}</strong></div>
            <div className="metric"><span>Active sessions</span><strong className="accent">{active.toString().padStart(2, '0')}</strong></div>
            <div className="metric"><span>Reported hosts</span><strong>{hosts.toString().padStart(2, '0')}</strong></div>
            <div className="metric"><span>Session feed</span><strong className="feed-state">{fetchError ? 'Unavailable' : lastRefresh ? 'Connected' : 'Connecting…'}</strong><small>{lastRefresh ? `Last update ${lastRefresh.toLocaleTimeString()}` : 'Waiting for first update'}</small></div></div>
        {fetchError && <div className="status error" role="alert">Session refresh failed: {fetchError}. Displayed sessions may be stale; retrying automatically.</div>}
        <ClientList {...shared} locations={locations} />
        <div className="visual-grid"><SessionGraph {...shared} /><SessionMap {...shared} locations={locations} onSaveLocation={saveLocation} /></div>
        {selected ? <><section className="panel section"><span className="eyebrow">SESSION INSPECTOR</span><h2>{systemOf(selected).hostname || selectedClient}</h2><p className="mono muted">{selected.ID} · {selectedClient}</p>
            {selected.IsActive && <SessionLink key={selected.ID} live={selected} clients={clients} onLinked={saved => {
                const savedLocation = locations[hostOf(saved.IPAddress)];
                if (savedLocation && !locations[hostOf(selected.IPAddress)]) saveLocation(hostOf(selected.IPAddress), savedLocation);
                setStatus('Saved session linked.');
                setRefreshVersion(v => v + 1); setMessageVersion(v => v + 1);
            }} />}
            {selected.IsActive && <div className="inspector-actions"><button disabled={disconnecting} onClick={() => disconnectSession(selected.ID)}>{disconnecting ? 'Disconnecting…' : 'Disconnect session'}</button></div>}
            <ClientInfo clientInfo={clientInfo} /><CommandResults results={selected.CommandResults || []} /><div className="panel-heading"><h3>Messages</h3><button onClick={() => setMessageVersion(v => v + 1)}>Refresh messages</button></div><MessageList messages={clientMessages} /></section>
            {selected.IsActive ? <TabContainer activeTab={activeTab} setActiveTab={setActiveTab} {...shared} setStatus={setStatus} refreshMessages={() => setMessageVersion(v => v + 1)} /> : <div className="status warning">This session is offline. Collected information is retained; commands require an active connection.</div>}</>
            : <section className="panel empty-state"><strong>Select a session to inspect</strong><p>Use the session table, connection graph, or map to open its details.</p></section>}
        {status && <div role="status" className={`status ${utils.getStatusClass(status)}`}>{status}</div>}
        <footer className="app-footer">FRAZZLEDRIP <span>Session inventory reflects the current server registry. Offline sessions retain their collected information while the server is running.</span></footer>
    </div>;
}
