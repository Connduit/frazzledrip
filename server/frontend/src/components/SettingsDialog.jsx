import React, { useEffect, useRef, useState } from 'react';

export default function SettingsDialog({
    onClose,
    onResetSizes,
    onResetOrder,
}) {
    const dialog = useRef(null);
    const [notice, setNotice] = useState('');
    useEffect(() => {
        const previous = document.activeElement;
        const element = dialog.current;
        element.showModal();
        return () => {
            element.close();
            previous?.focus?.({ preventScroll: true });
        };
    }, []);

    return (
        <dialog
            ref={dialog}
            className="settings-dialog"
            aria-labelledby="settings-title"
            onCancel={(event) => {
                event.preventDefault();
                onClose();
            }}
        >
            <div className="workspace-panel-title">
                <h2 id="settings-title">Settings</h2>
                <button
                    type="button"
                    aria-label="Close settings"
                    onClick={onClose}
                >
                    Close ×
                </button>
            </div>
            <div className="settings-body">
                <p className="muted">Workspace preferences for this browser.</p>
                <div className="settings-row">
                    <div>
                        <strong>Panel sizes</strong>
                        <p className="muted">
                            Restore the default space split between the Sessions
                            and feature windows.
                        </p>
                    </div>
                    <button
                        type="button"
                        onClick={() => {
                            onResetSizes();
                            setNotice('Panel sizes reset.');
                        }}
                    >
                        Reset sizes
                    </button>
                </div>
                <div className="settings-row">
                    <div>
                        <strong>Tab order</strong>
                        <p className="muted">
                            Reset the tab order in both windows. Open tabs and
                            the selected views are preserved.
                        </p>
                    </div>
                    <button
                        type="button"
                        onClick={() => {
                            onResetOrder();
                            setNotice('Tab order reset.');
                        }}
                    >
                        Reset tab order
                    </button>
                </div>
                {notice && (
                    <p role="status" className="status">
                        {notice}
                    </p>
                )}
            </div>
        </dialog>
    );
}
