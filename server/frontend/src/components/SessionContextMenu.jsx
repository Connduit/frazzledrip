import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { systemOf } from '../lib/sessions.js';

export default function SessionContextMenu({
    client,
    anchor,
    x,
    y,
    onClose,
    onOpen,
}) {
    const menu = useRef(null);
    const [position, setPosition] = useState({ left: x, top: y });

    useLayoutEffect(() => {
        const bounds = menu.current.getBoundingClientRect();
        setPosition({
            left: Math.max(8, Math.min(x, innerWidth - bounds.width - 8)),
            top: Math.max(8, Math.min(y, innerHeight - bounds.height - 8)),
        });
        menu.current.querySelector('[role="menuitem"]')?.focus();
    }, [x, y, anchor]);

    useEffect(() => {
        const dismiss = (event) => {
            if (!menu.current?.contains(event.target)) onClose();
        };
        const close = (event) => {
            if (event.type === 'scroll' && menu.current?.contains(event.target))
                return;
            onClose();
        };
        document.addEventListener('pointerdown', dismiss);
        window.addEventListener('resize', close);
        window.addEventListener('scroll', close, true);
        return () => {
            document.removeEventListener('pointerdown', dismiss);
            window.removeEventListener('resize', close);
            window.removeEventListener('scroll', close, true);
        };
    }, [onClose]);

    function keyDown(event) {
        if (event.key === 'Escape' || event.key === 'Tab') {
            event.preventDefault();
            onClose();
            anchor?.focus();
            return;
        }
        const items = [...menu.current.querySelectorAll('[role="menuitem"]')];
        const current = items.indexOf(document.activeElement);
        let next;
        if (event.key === 'ArrowDown') next = (current + 1) % items.length;
        else if (event.key === 'ArrowUp')
            next = (current + items.length - 1) % items.length;
        else if (event.key === 'Home') next = 0;
        else if (event.key === 'End') next = items.length - 1;
        else return;
        event.preventDefault();
        items[next].focus();
    }

    return createPortal(
        <div
            ref={menu}
            role="menu"
            aria-label={`Actions for session ${client.ID}`}
            className="session-context-menu"
            style={position}
            onKeyDown={keyDown}
            onContextMenu={(event) => event.preventDefault()}
        >
            <div className="session-menu-title">
                {systemOf(client).hostname || client.ID}
                <small>
                    {client.ID} · {client.IsActive ? 'Active' : 'Offline'}
                </small>
            </div>
            <button
                type="button"
                role="menuitem"
                onClick={() => onOpen('console')}
            >
                Open console
            </button>
            <div
                role="group"
                aria-label="Browsers"
                className="session-menu-browsers"
            >
                <span>Browsers</span>
                <button
                    type="button"
                    role="menuitem"
                    onClick={() => onOpen('files')}
                >
                    File Browser
                </button>
                <button
                    type="button"
                    role="menuitem"
                    onClick={() => onOpen('processes')}
                >
                    Process Browser
                </button>
            </div>
        </div>,
        document.body,
    );
}
