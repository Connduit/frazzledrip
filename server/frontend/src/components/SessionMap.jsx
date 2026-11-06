import React, { useState } from 'react';
import { geoEqualEarth, geoPath } from 'd3-geo';
import { feature } from 'topojson-client';
import atlas from 'world-atlas/countries-110m.json';
import { hostOf, systemOf, validLocation } from '../lib/sessions.js';

const countries = feature(atlas, atlas.objects.countries);
const projection = geoEqualEarth().fitExtent([[12, 12], [628, 302]], { type: 'Sphere' });
const path = geoPath(projection);

function LocationEditor({ client, location, onSave }) {
    const [label, setLabel] = useState(location?.label || '');
    const [lat, setLat] = useState(location?.lat ?? '');
    const [lon, setLon] = useState(location?.lon ?? '');
    const [error, setError] = useState('');
    return <form className="location-editor" onSubmit={e => { e.preventDefault(); const value = { label: label.trim(), lat: Number(lat), lon: Number(lon) };
        if (String(lat).trim() === '' || String(lon).trim() === '' || !validLocation(value)) { setError('Enter latitude −90 to 90 and longitude −180 to 180.'); return; }
        onSave(hostOf(client.IPAddress), value); setError(''); }}>
        <p>Location for <strong>{hostOf(client.IPAddress)}</strong> <span className="muted">· manually assigned</span></p>
        <div className="location-fields"><input aria-label="Location label" placeholder="Location label" maxLength="80" value={label} onChange={e => setLabel(e.target.value)} />
            <input aria-label="Latitude" placeholder="Latitude" type="number" step="any" min="-90" max="90" required value={lat} onChange={e => setLat(e.target.value)} />
            <input aria-label="Longitude" placeholder="Longitude" type="number" step="any" min="-180" max="180" required value={lon} onChange={e => setLon(e.target.value)} />
            <button type="submit">Save location</button>{location && <button type="button" onClick={() => onSave(hostOf(client.IPAddress), null)}>Clear</button>}</div>
        {error && <p role="alert">{error}</p>}
    </form>;
}

export default function SessionMap({ clients, selectedClient, onClientSelect, locations, onSaveLocation }) {
    const located = clients.filter(c => validLocation(locations[hostOf(c.IPAddress)]));
    const selected = clients.find(c => c.IPAddress === selectedClient);
    return <section className="panel"><div className="panel-heading"><div><span className="eyebrow">GEOGRAPHY</span><h2>Session map</h2></div><span className="muted">{located.length} mapped / {clients.length - located.length} unmapped</span></div>
        <svg viewBox="0 0 640 314" className="session-map" role="group" aria-label="World map with manually assigned session locations">
            <path d={path({ type: 'Sphere' })} className="map-ocean" />
            {countries.features.map((country, i) => <path key={i} d={path(country)} className="map-land" />)}
            {located.map(client => { const location = locations[hostOf(client.IPAddress)]; const [x, y] = projection([location.lon, location.lat]); return <g key={client.IPAddress} role="button" tabIndex="0" className="map-marker" aria-label={`Select ${systemOf(client).hostname || client.IPAddress} at ${location.label || 'assigned coordinates'}`}
                onClick={() => onClientSelect(client.IPAddress)} onKeyDown={e => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); onClientSelect(client.IPAddress); } }}>
                <title>{client.IPAddress} · {location.label || `${location.lat}, ${location.lon}`}</title>
                <circle cx={x} cy={y} r="12" className="marker-halo" /><circle cx={x} cy={y} r={selectedClient === client.IPAddress ? 6 : 4} className="marker-core" />
            </g>; })}
        </svg>
        <div className="map-note">Locations are assigned coordinates, not IP geolocation. Private and loopback addresses do not identify a physical location. Map data © Natural Earth, public domain.</div>
        {selected ? <LocationEditor key={`${selected.IPAddress}-${JSON.stringify(locations[hostOf(selected.IPAddress)])}`} client={selected} location={locations[hostOf(selected.IPAddress)]} onSave={onSaveLocation} />
            : <div className="panel-footer">Select a session to assign its location. No addresses are sent to external services.</div>}
    </section>;
}
