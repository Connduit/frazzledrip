function normalizeAddress(value) {
    const address = value.trim().toLowerCase();
    if (address === 'localhost') return address;
    const parts = address.split('.');
    if (
        parts.length === 4 &&
        parts.every(
            (part) => /^(0|[1-9]\d{0,2})$/.test(part) && Number(part) <= 255,
        )
    )
        return address;
    if (address.includes(':')) {
        try {
            const url = new URL(`http://[${address}]/`);
            return url.hostname.slice(1, -1);
        } catch {
            /* Report invalid IPv6 below. */
        }
    }
    throw new Error(
        'Bind address must be an IPv4 address, IPv6 address, or localhost.',
    );
}

export function validateListenerDraft(input, drafts = []) {
    const name = String(input.name || '').trim();
    if (!name || name.length > 64 || /[\x00-\x1f\x7f]/.test(name))
        throw new Error(
            'Enter a name of 1–64 characters without control characters.',
        );
    if (input.protocol !== 'TCP')
        throw new Error('Only TCP drafts are supported.');
    const address = normalizeAddress(String(input.address || ''));
    const text = String(input.port);
    if (!/^\d+$/.test(text) || Number(text) < 1 || Number(text) > 65535)
        throw new Error('Port must be a whole number between 1 and 65535.');
    const port = Number(text);
    if (drafts.some((draft) => draft.name.toLowerCase() === name.toLowerCase()))
        throw new Error('A draft with this name already exists.');
    if (
        drafts.some(
            (draft) =>
                draft.protocol === input.protocol &&
                draft.address === address &&
                draft.port === port,
        )
    )
        throw new Error(
            'A draft with this protocol, bind address, and port already exists.',
        );
    return { name, protocol: input.protocol, address, port };
}

export function parseListenerDrafts(raw) {
    try {
        const saved = JSON.parse(raw);
        if (!Array.isArray(saved)) return [];
        const drafts = [];
        for (const item of saved.slice(0, 50)) {
            if (
                !item ||
                typeof item.id !== 'string' ||
                !item.id ||
                item.id.length > 100 ||
                drafts.some((draft) => draft.id === item.id)
            )
                continue;
            try {
                drafts.push({
                    ...validateListenerDraft(item, drafts),
                    id: item.id,
                });
            } catch {
                /* Ignore invalid or duplicate stored drafts. */
            }
        }
        return drafts;
    } catch {
        return [];
    }
}
