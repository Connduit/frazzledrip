export function hostOf(address = '') {
    if (address.startsWith('[')) return address.slice(1, address.indexOf(']'));
    return address.split(':').length === 2 ? address.split(':')[0] : address;
}

export function systemOf(client) {
    const info = client.SystemInfo || {};
    return {
        hostname: info.hostname || info.Hostname || '',
        user: info.user || info.User || '',
        os: info.os || info.OS || '',
        arch: info.arch || info.Arch || '',
    };
}

export function validLocation(location) {
    return location && typeof location.lat === 'number' && typeof location.lon === 'number'
        && Number.isFinite(location.lat) && Number.isFinite(location.lon)
        && Math.abs(location.lat) <= 90 && Math.abs(location.lon) <= 180;
}

export function formatMemory(bytes) {
    if (typeof bytes !== 'number' || !Number.isFinite(bytes) || bytes < 0) return '—';
    return `${(bytes / (1024 ** 3)).toFixed(1)} GiB`;
}
