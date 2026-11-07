export function orderedTabs(tabs, preferred) {
    const byID = new Map(tabs.map((tab) => [tab.id, tab]));
    const ids = [...new Set([...preferred, ...byID.keys()])];
    return ids.filter((id) => byID.has(id)).map((id) => byID.get(id));
}

export function moveTab(ids, source, target, position = 'before') {
    if (source === target || !ids.includes(source) || !ids.includes(target))
        return ids;
    const remaining = ids.filter((id) => id !== source);
    remaining.splice(
        remaining.indexOf(target) + (position === 'after' ? 1 : 0),
        0,
        source,
    );
    return remaining;
}

export function activeTabAfterClose(ids, active, closed) {
    const index = ids.indexOf(closed);
    if (index < 0 || active !== closed) return active;
    return ids[index + 1] || ids[index - 1] || '';
}
