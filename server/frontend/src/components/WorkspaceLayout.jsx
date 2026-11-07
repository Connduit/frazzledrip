import React, { useEffect, useRef, useState } from 'react';

const storageKey = 'frazzledrip.workspace-sizes.v3';
const defaultSize = 0.475;
const clamp = (value, min, max) => Math.max(min, Math.min(max, value));

function readSize() {
    try {
        const saved = JSON.parse(localStorage.getItem(storageKey));
        if (Number.isFinite(saved) && saved >= 0.2 && saved <= 0.78)
            return saved;
    } catch {
        /* Use the default size when storage is unavailable. */
    }
    return defaultSize;
}

export default function WorkspaceLayout({ children, resetVersion = 0 }) {
    const [size, setSize] = useState(readSize);
    const [dragging, setDragging] = useState(false);
    const grid = useRef(null);
    const offset = useRef(null);
    useEffect(() => {
        if (resetVersion) setSize(defaultSize);
    }, [resetVersion]);

    useEffect(() => {
        try {
            localStorage.setItem(storageKey, JSON.stringify(size));
        } catch {
            /* Resizing still works without persistence. */
        }
    }, [size]);

    function update(value) {
        const height = Math.max(
            1,
            grid.current.getBoundingClientRect().height - 8,
        );
        setSize(
            clamp(
                value,
                Math.max(0.2, 180 / height),
                Math.min(0.78, 1 - 220 / height),
            ),
        );
    }

    function start(event) {
        if (event.button !== 0) return;
        event.preventDefault();
        event.currentTarget.setPointerCapture(event.pointerId);
        const bounds = grid.current.getBoundingClientRect();
        offset.current =
            event.clientY - (bounds.top + size * (bounds.height - 8) + 4);
        setDragging(true);
    }

    function move(event) {
        if (offset.current === null) return;
        const bounds = grid.current.getBoundingClientRect();
        update(
            (event.clientY - offset.current - bounds.top - 4) /
                (bounds.height - 8),
        );
    }

    function stop() {
        offset.current = null;
        setDragging(false);
    }

    return (
        <>
            <div className="workspace-layout-toolbar">
                <span>Drag the divider to resize windows</span>
                <button type="button" onClick={() => setSize(defaultSize)}>
                    Reset panel sizes
                </button>
            </div>
            <main
                ref={grid}
                className={`workspace-grid resizable-workspace ${dragging ? 'workspace-resizing' : ''}`}
                style={{ '--workspace-rows': `${size}fr ${1 - size}fr` }}
            >
                {children}
                <div
                    className="workspace-divider divider-row"
                    role="separator"
                    tabIndex={0}
                    aria-label="Resize upper and lower panels"
                    aria-orientation="horizontal"
                    aria-valuemin={0}
                    aria-valuemax={100}
                    aria-valuenow={Math.round(size * 100)}
                    title="Drag or use Up/Down to resize windows. Home resets the divider."
                    style={{
                        left: 0,
                        right: 0,
                        top: `calc(${size * 100}% + ${4 - size * 8}px)`,
                    }}
                    onPointerDown={start}
                    onPointerMove={move}
                    onPointerUp={stop}
                    onPointerCancel={stop}
                    onLostPointerCapture={stop}
                    onKeyDown={(event) => {
                        if (
                            ['ArrowUp', 'ArrowDown', 'Home'].includes(event.key)
                        ) {
                            event.preventDefault();
                            update(
                                event.key === 'Home'
                                    ? defaultSize
                                    : size +
                                          (event.key === 'ArrowDown'
                                              ? 0.02
                                              : -0.02),
                            );
                        }
                    }}
                />
            </main>
        </>
    );
}
