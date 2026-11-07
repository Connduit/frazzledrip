import React, { useEffect, useRef, useState } from 'react';

const storageKey = 'frazzledrip.workspace-sizes.v1';
const defaults = { column: 0.34, firstRow: 0.25, upperRows: 0.475 };
const clamp = (value, min, max) => Math.max(min, Math.min(max, value));

function readSizes() {
    try {
        const saved = JSON.parse(localStorage.getItem(storageKey));
        if (
            saved &&
            [saved.column, saved.firstRow, saved.upperRows].every(
                Number.isFinite,
            ) &&
            saved.column >= 0.18 &&
            saved.column <= 0.78 &&
            saved.firstRow >= 0.12 &&
            saved.upperRows <= 0.78 &&
            saved.upperRows - saved.firstRow >= 0.12
        )
            return saved;
    } catch {
        /* Use defaults when browser storage is unavailable. */
    }
    return defaults;
}

export default function WorkspaceLayout({ children }) {
    const [sizes, setSizes] = useState(readSizes);
    const [dragging, setDragging] = useState(false);
    const grid = useRef(null);
    const drag = useRef(null);

    useEffect(() => {
        try {
            localStorage.setItem(storageKey, JSON.stringify(sizes));
        } catch {
            /* Resizing still works without persistence. */
        }
    }, [sizes]);

    function update(kind, value) {
        const bounds = grid.current.getBoundingClientRect();
        const width = Math.max(1, bounds.width - 8);
        const height = Math.max(1, bounds.height - 16);
        const minTop = Math.max(0.12, 100 / height);
        const minMiddle = Math.max(0.12, 100 / height);
        const maxUpper = Math.min(0.78, 1 - 180 / height);
        setSizes((current) => {
            if (kind === 'column')
                return {
                    ...current,
                    column: clamp(
                        value,
                        Math.max(0.18, 220 / width),
                        Math.min(0.78, 1 - 280 / width),
                    ),
                };
            if (kind === 'firstRow')
                return {
                    ...current,
                    firstRow: clamp(
                        value,
                        minTop,
                        current.upperRows - minMiddle,
                    ),
                };
            return {
                ...current,
                upperRows: clamp(value, current.firstRow + minMiddle, maxUpper),
            };
        });
    }

    function start(event, kind) {
        if (event.button !== 0) return;
        event.preventDefault();
        event.currentTarget.setPointerCapture(event.pointerId);
        const bounds = grid.current.getBoundingClientRect();
        const boundary =
            kind === 'column'
                ? bounds.left + sizes.column * (bounds.width - 8) + 4
                : bounds.top +
                  sizes[kind] * (bounds.height - 16) +
                  (kind === 'firstRow' ? 4 : 12);
        drag.current = {
            kind,
            offset:
                (kind === 'column' ? event.clientX : event.clientY) - boundary,
        };
        setDragging(true);
    }

    function move(event) {
        if (!drag.current) return;
        const { kind, offset } = drag.current;
        const bounds = grid.current.getBoundingClientRect();
        const value =
            kind === 'column'
                ? (event.clientX - offset - bounds.left - 4) /
                  (bounds.width - 8)
                : (event.clientY -
                      offset -
                      bounds.top -
                      (kind === 'firstRow' ? 4 : 12)) /
                  (bounds.height - 16);
        update(kind, value);
    }

    function stop() {
        drag.current = null;
        setDragging(false);
    }

    function divider(kind, label, style) {
        const vertical = kind === 'column';
        return (
            <div
                className={`workspace-divider ${vertical ? 'divider-column' : 'divider-row'}`}
                role="separator"
                tabIndex={0}
                aria-label={label}
                aria-orientation={vertical ? 'vertical' : 'horizontal'}
                aria-valuemin={0}
                aria-valuemax={100}
                aria-valuenow={Math.round(sizes[kind] * 100)}
                title={`${label}: drag or use arrow keys. Home resets this divider.`}
                style={style}
                onPointerDown={(event) => start(event, kind)}
                onPointerMove={move}
                onPointerUp={stop}
                onPointerCancel={stop}
                onLostPointerCapture={stop}
                onKeyDown={(event) => {
                    const decrease = vertical ? 'ArrowLeft' : 'ArrowUp';
                    const increase = vertical ? 'ArrowRight' : 'ArrowDown';
                    if (
                        event.key === decrease ||
                        event.key === increase ||
                        event.key === 'Home'
                    ) {
                        event.preventDefault();
                        update(
                            kind,
                            event.key === 'Home'
                                ? defaults[kind]
                                : sizes[kind] +
                                      (event.key === increase ? 0.02 : -0.02),
                        );
                    }
                }}
            />
        );
    }

    const columnPosition = `calc(${sizes.column * 100}% + ${4 - sizes.column * 8}px)`;
    return (
        <>
            <div className="workspace-layout-toolbar">
                <span>Drag dividers to resize panels</span>
                <button type="button" onClick={() => setSizes(defaults)}>
                    Reset panel sizes
                </button>
            </div>
            <main
                ref={grid}
                className={`workspace-grid resizable-workspace ${dragging ? 'workspace-resizing' : ''}`}
                style={{
                    '--workspace-columns': `${sizes.column}fr ${1 - sizes.column}fr`,
                    '--workspace-rows': `${sizes.firstRow}fr ${sizes.upperRows - sizes.firstRow}fr ${1 - sizes.upperRows}fr`,
                }}
            >
                {children}
                {divider('column', 'Resize left and right panels', {
                    left: columnPosition,
                })}
                {divider('firstRow', 'Resize events and sessions', {
                    left: columnPosition,
                    right: 0,
                    top: `calc(${sizes.firstRow * 100}% + ${4 - sizes.firstRow * 16}px)`,
                })}
                {divider('upperRows', 'Resize upper and lower panels', {
                    left: 0,
                    right: 0,
                    top: `calc(${sizes.upperRows * 100}% + ${12 - sizes.upperRows * 16}px)`,
                })}
            </main>
        </>
    );
}
