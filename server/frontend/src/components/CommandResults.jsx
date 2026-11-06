import React from 'react';

export default function CommandResults({ results }) {
    return <div className="command-results"><h3>Collected command output</h3>
        <p className="muted">Latest response per command. Host facts above update from recognized whoami, ipconfig, and systeminfo output.</p>
        {!results.length && <p className="muted">No correlated command responses yet.</p>}
        {results.map(result => <details key={result.command}>
            <summary><span className="mono">{result.command}</span><time dateTime={result.received_at}>{new Date(result.received_at).toLocaleString()}</time></summary>
            <pre>{result.output || '(Empty response)'}</pre>
        </details>)}
    </div>;
}
