import React, { useEffect, useState } from 'react';
import { apiService } from '../services/api.js';

export default function ToolCatalog({ agent }) {
    const [tools, setTools] = useState([]);
    const [query, setQuery] = useState('');
    const [category, setCategory] = useState('All');
    const [kind, setKind] = useState('all');
    const [note, setNote] = useState('');
    const [busy, setBusy] = useState('');
    const [error, setError] = useState('');
    const [result, setResult] = useState(null);
    useEffect(() => {
        let cancelled = false;
        apiService
            .toolCatalog()
            .then((data) => {
                if (!cancelled) setTools(data);
            })
            .catch((err) => {
                if (!cancelled) setError(err.message);
            });
        return () => {
            cancelled = true;
        };
    }, []);
    const categories = ['All', ...new Set(tools.map((tool) => tool.category))];
    const visible = tools.filter(
        (tool) =>
            (category === 'All' || tool.category === category) &&
            (kind === 'all' || tool.kind === kind) &&
            `${tool.name} ${tool.description} ${tool.category}`
                .toLowerCase()
                .includes(query.toLowerCase()),
    );
    async function simulate(tool) {
        if (!agent || busy) return;
        setBusy(tool.id);
        setError('');
        setResult(null);
        try {
            setResult(await apiService.simulateTool(agent.ID, tool.id, note));
        } catch (err) {
            setError(err.message);
        } finally {
            setBusy('');
        }
    }
    return (
        <section aria-label="Attacks and scripts catalog">
            <div className="catalog-banner">
                <strong>Development placeholders · simulation only</strong>
                <p>
                    These entries test the server's task history and audit
                    workflow. They do not contact a client, run a script, or
                    change a host. Simulations are labeled in Task History.
                </p>
            </div>
            <div className="catalog-filters">
                <input
                    aria-label="Search tools"
                    placeholder="Search attacks or scripts…"
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                />
                <select
                    aria-label="Tool category"
                    value={category}
                    onChange={(event) => setCategory(event.target.value)}
                >
                    {categories.map((value) => (
                        <option key={value}>{value}</option>
                    ))}
                </select>
                <select
                    aria-label="Tool kind"
                    value={kind}
                    onChange={(event) => setKind(event.target.value)}
                >
                    <option value="all">Attacks and scripts</option>
                    <option value="attack">Attack placeholders</option>
                    <option value="script">Script placeholders</option>
                </select>
            </div>
            <label className="catalog-note">
                Simulation note{' '}
                <input
                    aria-label="Simulation note"
                    maxLength="500"
                    placeholder="Optional test or engagement note"
                    value={note}
                    onChange={(event) => setNote(event.target.value)}
                />
            </label>
            {error && (
                <div role="alert" className="status error">
                    {error}
                </div>
            )}
            {result && (
                <div role="status" className="simulation-result">
                    <strong>Simulation recorded</strong>
                    <p className="mono">{result.id}</p>
                    <pre>{result.output}</pre>
                </div>
            )}
            <div className="catalog-grid">
                {visible.map((tool) => (
                    <article key={tool.id} className="catalog-card">
                        <div className="catalog-card-heading">
                            <span className="eyebrow">{tool.category}</span>
                            <span className="badge">
                                {tool.kind === 'script'
                                    ? 'Script template'
                                    : 'Attack placeholder'}
                            </span>
                        </div>
                        <h3>{tool.name}</h3>
                        <p>{tool.description}</p>
                        <button
                            disabled={!agent || Boolean(busy)}
                            onClick={() => simulate(tool)}
                        >
                            {busy === tool.id ? 'Recording…' : 'Simulate task'}
                        </button>
                    </article>
                ))}
            </div>
            {!visible.length && (
                <p className="empty-state">No tools match these filters.</p>
            )}
            {!agent && (
                <p className="status warning">
                    Choose a session in the panel header to record a simulation.
                </p>
            )}
        </section>
    );
}
