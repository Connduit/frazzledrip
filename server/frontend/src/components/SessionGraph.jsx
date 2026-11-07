import React, { useState, useRef, useEffect } from 'react';
import { systemOf } from '../lib/sessions.js';

export default function SessionGraph({
    clients,
    selectedClient,
    onClientSelect,
    onSessionContextMenu,
}) {
    const [camera, setCamera] = useState({ zoom: 1, x: 0, y: 0 });
    const svg = useRef(null);
    const drag = useRef(null);
    const suppressClick = useRef(false);
    const [dragging, setDragging] = useState(false);
    const height = Math.max(240, clients.length * 66 + 40);
    const fitScale = Math.min(1, 360 / height);
    const zoom = camera.zoom;
    function point(clientX, clientY) {
        const p = svg.current.createSVGPoint();
        p.x = clientX;
        p.y = clientY;
        const matrix = svg.current.getScreenCTM();
        return matrix
            ? p.matrixTransform(matrix.inverse())
            : { x: 310, y: 200 };
    }
    function adjustZoom(factor, anchor = { x: 310, y: 200 }) {
        setCamera((current) => {
            const next = Math.max(
                0.5,
                Math.min(3, Math.round(current.zoom * factor * 100) / 100),
            );
            const ratio = next / current.zoom;
            return {
                zoom: next,
                x: anchor.x - 310 - (anchor.x - 310 - current.x) * ratio,
                y: anchor.y - 200 - (anchor.y - 200 - current.y) * ratio,
            };
        });
    }
    useEffect(() => {
        const element = svg.current;
        const wheel = (event) => {
            event.preventDefault();
            if (!drag.current)
                adjustZoom(
                    Math.exp(
                        -Math.max(-100, Math.min(100, event.deltaY)) * 0.003,
                    ),
                    point(event.clientX, event.clientY),
                );
        };
        element.addEventListener('wheel', wheel, { passive: false });
        return () => element.removeEventListener('wheel', wheel);
    }, []);
    function pointerDown(event) {
        if (event.button !== 0) return;
        drag.current = {
            id: event.pointerId,
            start: point(event.clientX, event.clientY),
            clientX: event.clientX,
            clientY: event.clientY,
            x: camera.x,
            y: camera.y,
            moved: false,
        };
    }
    function pointerMove(event) {
        const current = drag.current;
        if (!current || current.id !== event.pointerId) return;
        if (
            !current.moved &&
            Math.hypot(
                event.clientX - current.clientX,
                event.clientY - current.clientY,
            ) < 4
        )
            return;
        current.moved = true;
        setDragging(true);
        svg.current.setPointerCapture(event.pointerId);
        const p = point(event.clientX, event.clientY);
        setCamera((previous) => ({
            ...previous,
            x: current.x + p.x - current.start.x,
            y: current.y + p.y - current.start.y,
        }));
    }
    function pointerUp(event) {
        if (!drag.current || drag.current.id !== event.pointerId) return;
        suppressClick.current = drag.current.moved;
        drag.current = null;
        setDragging(false);
        if (svg.current.hasPointerCapture(event.pointerId))
            svg.current.releasePointerCapture(event.pointerId);
        setTimeout(() => {
            suppressClick.current = false;
        }, 0);
    }
    function fit() {
        setCamera({ zoom: 1, x: 0, y: 0 });
    }
    return (
        <section className="panel">
            <div className="panel-heading">
                <div>
                    <span className="eyebrow">TOPOLOGY</span>
                    <h2>Connection graph</h2>
                </div>
                <div
                    className="graph-controls"
                    role="group"
                    aria-label="Graph zoom controls"
                >
                    <button
                        type="button"
                        aria-label="Zoom out graph"
                        disabled={zoom <= 0.5}
                        onClick={() => adjustZoom(1 / 1.25)}
                    >
                        −
                    </button>
                    <output aria-live="polite" aria-label="Graph zoom">
                        {Math.round(zoom * 100)}%
                    </output>
                    <button
                        type="button"
                        aria-label="Zoom in graph"
                        disabled={zoom >= 3}
                        onClick={() => adjustZoom(1.25)}
                    >
                        +
                    </button>
                    <button type="button" onClick={fit}>
                        Fit graph
                    </button>
                </div>
            </div>
            <div className="graph-viewport">
                <svg
                    ref={svg}
                    viewBox="0 0 620 400"
                    className={`session-graph ${dragging ? 'dragging' : ''}`}
                    role="group"
                    aria-label="Sessions connected directly to the controller"
                    onPointerDown={pointerDown}
                    onPointerMove={pointerMove}
                    onPointerUp={pointerUp}
                    onPointerCancel={pointerUp}
                    onClickCapture={(event) => {
                        if (suppressClick.current) {
                            event.preventDefault();
                            event.stopPropagation();
                        }
                    }}
                >
                    <g
                        className="graph-scene"
                        transform={`translate(${310 + camera.x} ${200 + camera.y}) scale(${fitScale * zoom}) translate(-310 ${-height / 2})`}
                    >
                        {clients.map((client, i) => {
                            const y = 40 + i * 66;
                            return (
                                <g key={client.IPAddress}>
                                    <path
                                        className={`graph-edge ${client.IsActive ? '' : 'disconnected'}`}
                                        d={`M 164 ${height / 2} C 290 ${height / 2}, 280 ${y}, 352 ${y}`}
                                    />
                                    <g
                                        className={`graph-node ${selectedClient === client.IPAddress ? 'chosen' : ''}`}
                                        role="button"
                                        tabIndex="0"
                                        aria-label={`Select ${systemOf(client).hostname || client.IPAddress}`}
                                        aria-pressed={
                                            selectedClient === client.IPAddress
                                        }
                                        onClick={() =>
                                            onClientSelect(client.IPAddress)
                                        }
                                        onContextMenu={(event) =>
                                            onSessionContextMenu?.(
                                                event,
                                                client,
                                            )
                                        }
                                        onKeyDown={(e) => {
                                            if (
                                                e.key === 'ContextMenu' ||
                                                (e.shiftKey && e.key === 'F10')
                                            ) {
                                                onSessionContextMenu?.(
                                                    e,
                                                    client,
                                                );
                                                return;
                                            }
                                            if (
                                                e.key === 'Enter' ||
                                                e.key === ' '
                                            ) {
                                                e.preventDefault();
                                                onClientSelect(
                                                    client.IPAddress,
                                                );
                                            }
                                        }}
                                    >
                                        <rect
                                            x="352"
                                            y={y - 23}
                                            width="245"
                                            height="48"
                                            rx="7"
                                        />
                                        <circle
                                            cx="368"
                                            cy={y}
                                            r="4"
                                            className={
                                                client.IsActive
                                                    ? 'active-dot'
                                                    : 'inactive-dot'
                                            }
                                        />
                                        <text x="382" y={y - 3}>
                                            {(
                                                systemOf(client).hostname ||
                                                'Unknown host'
                                            ).slice(0, 25)}
                                        </text>
                                        <text
                                            x="382"
                                            y={y + 14}
                                            className="graph-caption"
                                        >
                                            {client.IPAddress}
                                        </text>
                                    </g>
                                </g>
                            );
                        })}
                        <rect
                            x="20"
                            y={height / 2 - 32}
                            width="144"
                            height="64"
                            rx="8"
                            className="controller-node"
                        />
                        <text x="40" y={height / 2 - 4}>
                            CONTROLLER
                        </text>
                        <text
                            x="40"
                            y={height / 2 + 16}
                            className="graph-caption"
                        >
                            TCP listener · 4444
                        </text>
                        {!clients.length && (
                            <text
                                x="240"
                                y={height / 2}
                                className="graph-caption"
                            >
                                Waiting for connections
                            </text>
                        )}
                    </g>
                </svg>
            </div>
            <footer className="panel-footer">
                Direct TCP connections · wheel to zoom · left-drag to pan · Fit
                graph to reset · select a device to inspect it
            </footer>
        </section>
    );
}
