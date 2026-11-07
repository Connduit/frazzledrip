import React, { useId } from 'react';

export default function WorkspaceTabs({ tabs, activeTab, onTabChange }) {
    const id = useId();

    function navigate(event, index) {
        let next;
        if (event.key === 'ArrowRight') next = (index + 1) % tabs.length;
        else if (event.key === 'ArrowLeft')
            next = (index + tabs.length - 1) % tabs.length;
        else if (event.key === 'Home') next = 0;
        else if (event.key === 'End') next = tabs.length - 1;
        else return;
        event.preventDefault();
        onTabChange(tabs[next].id);
        event.currentTarget.parentElement
            .querySelectorAll('[role="tab"]')
            [next].focus();
    }

    return (
        <section
            className="panel workspace-tabs"
            aria-label="Session workspace"
        >
            <div
                className="workspace-tab-list"
                role="tablist"
                aria-label="Workspace views"
            >
                {tabs.map((tab, index) => (
                    <button
                        key={tab.id}
                        type="button"
                        role="tab"
                        id={`${id}-tab-${tab.id}`}
                        aria-controls={`${id}-panel-${tab.id}`}
                        aria-selected={activeTab === tab.id}
                        tabIndex={activeTab === tab.id ? 0 : -1}
                        onClick={() => onTabChange(tab.id)}
                        onKeyDown={(event) => navigate(event, index)}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>
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
