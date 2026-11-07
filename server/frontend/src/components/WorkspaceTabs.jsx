import React, { useEffect, useId, useRef, useState } from 'react';
import {
    activeTabAfterClose,
    moveTab,
    orderedTabs,
} from '../lib/workspaceTabs.js';

const defaultStorageKey = 'frazzledrip.workspace-tab-order.v1';

function readOrder(storageKey) {
    try {
        const saved = JSON.parse(localStorage.getItem(storageKey));
        if (Array.isArray(saved) && saved.every((id) => typeof id === 'string'))
            return saved;
    } catch {
        /* Use the default order when storage is unavailable. */
    }
    return [];
}

export default function WorkspaceTabs({
    tabs,
    activeTab,
    onTabChange,
    onTabClose,
    emptyContent = 'No tabs open. Use the toolbar to open a view.',
    label = 'Session workspace',
    tabListLabel = 'Workspace views',
    storageKey = defaultStorageKey,
    className = '',
    resetVersion = 0,
}) {
    const id = useId();
    const [order, setOrder] = useState(() => readOrder(storageKey));
    const [dragging, setDragging] = useState(null);
    const [dropTarget, setDropTarget] = useState(null);
    const [announcement, setAnnouncement] = useState('');
    const container = useRef(null);
    const draggedTab = useRef(null);
    const displayed = orderedTabs(tabs, order);
    useEffect(() => {
        if (resetVersion) setOrder([]);
    }, [resetVersion]);

    useEffect(() => {
        try {
            localStorage.setItem(storageKey, JSON.stringify(order));
        } catch {
            /* Reordering still works without persistence. */
        }
    }, [order, storageKey]);

    function reorder(source, target, position) {
        const next = moveTab(
            displayed.map((tab) => tab.id),
            source,
            target,
            position,
        );
        setOrder(next);
        const label = tabs.find((tab) => tab.id === source)?.label;
        setAnnouncement(
            `${label} moved to position ${next.indexOf(source) + 1} of ${next.length}.`,
        );
    }

    function finishDrag() {
        draggedTab.current = null;
        setDragging(null);
        setDropTarget(null);
    }

    function closeTab(tab) {
        const next = activeTabAfterClose(
            displayed.map((item) => item.id),
            activeTab,
            tab.id,
        );
        finishDrag();
        onTabClose(tab.id, next);
        setAnnouncement(`${tab.label} closed.`);
        requestAnimationFrame(() => {
            const target = next && document.getElementById(`${id}-tab-${next}`);
            (target || container.current)?.focus();
        });
    }

    function navigate(event, index) {
        if (event.key === 'Delete' && onTabClose) {
            event.preventDefault();
            closeTab(displayed[index]);
            return;
        }
        if (
            event.ctrlKey &&
            event.shiftKey &&
            ['ArrowLeft', 'ArrowRight'].includes(event.key)
        ) {
            event.preventDefault();
            const next = index + (event.key === 'ArrowRight' ? 1 : -1);
            if (next < 0 || next >= displayed.length) return;
            const source = displayed[index].id;
            reorder(
                source,
                displayed[next].id,
                next > index ? 'after' : 'before',
            );
            requestAnimationFrame(() =>
                document.getElementById(`${id}-tab-${source}`)?.focus(),
            );
            return;
        }
        let next;
        if (event.key === 'ArrowRight') next = (index + 1) % displayed.length;
        else if (event.key === 'ArrowLeft')
            next = (index + displayed.length - 1) % displayed.length;
        else if (event.key === 'Home') next = 0;
        else if (event.key === 'End') next = displayed.length - 1;
        else return;
        event.preventDefault();
        onTabChange(displayed[next].id);
        event.currentTarget
            .closest('[role="tablist"]')
            .querySelectorAll('[role="tab"]')
            [next].focus();
    }

    return (
        <section
            ref={container}
            tabIndex={-1}
            className={`panel workspace-tabs ${className}`}
            aria-label={label}
        >
            <div
                className="workspace-tab-list"
                role="tablist"
                aria-label={tabListLabel}
            >
                {displayed.map((tab, index) => (
                    <div
                        key={tab.id}
                        role="presentation"
                        className={`workspace-tab-item ${onTabClose ? 'closable' : ''}`}
                        onDragOver={(event) => {
                            if (
                                !draggedTab.current ||
                                draggedTab.current === tab.id
                            )
                                return;
                            event.preventDefault();
                            event.dataTransfer.dropEffect = 'move';
                            const bounds =
                                event.currentTarget.getBoundingClientRect();
                            setDropTarget({
                                id: tab.id,
                                position:
                                    event.clientX <
                                    bounds.left + bounds.width / 2
                                        ? 'before'
                                        : 'after',
                            });
                        }}
                        onDragLeave={() => setDropTarget(null)}
                        onDrop={(event) => {
                            if (!draggedTab.current) return;
                            event.preventDefault();
                            const bounds =
                                event.currentTarget.getBoundingClientRect();
                            reorder(
                                draggedTab.current,
                                tab.id,
                                event.clientX < bounds.left + bounds.width / 2
                                    ? 'before'
                                    : 'after',
                            );
                            finishDrag();
                        }}
                    >
                        <button
                            type="button"
                            role="tab"
                            draggable
                            title="Drag to reorder. Ctrl+Shift+Left/Right moves this tab. Delete closes it."
                            className={
                                dragging === tab.id
                                    ? 'workspace-tab-dragging'
                                    : ''
                            }
                            data-drop-position={
                                dropTarget?.id === tab.id
                                    ? dropTarget.position
                                    : undefined
                            }
                            id={`${id}-tab-${tab.id}`}
                            aria-controls={`${id}-panel-${tab.id}`}
                            aria-selected={activeTab === tab.id}
                            tabIndex={activeTab === tab.id ? 0 : -1}
                            onClick={() => onTabChange(tab.id)}
                            onKeyDown={(event) => navigate(event, index)}
                            onDragStart={(event) => {
                                draggedTab.current = tab.id;
                                setDragging(tab.id);
                                event.dataTransfer.effectAllowed = 'move';
                                event.dataTransfer.setData(
                                    'text/plain',
                                    tab.id,
                                );
                            }}
                            onDragEnd={finishDrag}
                        >
                            {tab.label}
                        </button>
                        {onTabClose && (
                            <button
                                type="button"
                                className="workspace-tab-close"
                                aria-label={`Close ${tab.label} tab`}
                                title={`Close ${tab.label} (Delete)`}
                                tabIndex={-1}
                                onClick={(event) => {
                                    event.stopPropagation();
                                    closeTab(tab);
                                }}
                            >
                                ×
                            </button>
                        )}
                    </div>
                ))}
            </div>
            <span className="sr-only" role="status">
                {announcement}
            </span>
            {!tabs.length && (
                <div className="empty-state workspace-tabs-empty">
                    {emptyContent}
                </div>
            )}
            {tabs.map((tab) => (
                <div
                    key={tab.id}
                    className="workspace-tab-panel"
                    role="tabpanel"
                    id={`${id}-panel-${tab.id}`}
                    aria-labelledby={`${id}-tab-${tab.id}`}
                    hidden={activeTab !== tab.id}
                >
                    {tab.content}
                </div>
            ))}
        </section>
    );
}
