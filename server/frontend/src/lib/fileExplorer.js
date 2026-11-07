export function formatFileSize(size) {
    if (typeof size !== 'number' || !Number.isFinite(size) || size < 0)
        return '—';
    const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB'];
    let value = size,
        unit = 0;
    while (value >= 1024 && unit < units.length - 1) {
        value /= 1024;
        unit++;
    }
    return `${unit ? value.toFixed(1) : value} ${units[unit]}`;
}

export function fileType(entry) {
    if (entry.is_directory) return entry.is_link ? 'Folder link' : 'Folder';
    const dot = entry.name.lastIndexOf('.');
    const extension =
        dot > 0 && dot < entry.name.length - 1
            ? entry.name.slice(dot + 1).toUpperCase()
            : '';
    return `${extension ? `${extension} file` : 'File'}${entry.is_link ? ' · Link' : ''}`;
}

export function formatFileDate(milliseconds, reportedText) {
    if (reportedText) return reportedText;
    if (typeof milliseconds !== 'number' || !Number.isFinite(milliseconds))
        return '—';
    const date = new Date(milliseconds);
    return Number.isNaN(date.getTime()) ? '—' : date.toLocaleString();
}

export function sortFiles(entries, field, direction) {
    const value = (entry) =>
        field === 'type'
            ? fileType(entry)
            : field === 'size'
              ? entry.is_directory
                  ? null
                  : entry.size
              : field === 'modified_ms'
                ? entry.modified_sort ||
                  (Number.isFinite(entry.modified_ms) &&
                  !Number.isNaN(new Date(entry.modified_ms).getTime())
                      ? new Date(entry.modified_ms).toISOString()
                      : null)
                : entry.name;
    return entries.slice().sort((a, b) => {
        const foldersFirst =
            Number(Boolean(b.is_directory)) - Number(Boolean(a.is_directory));
        if (foldersFirst) return foldersFirst;
        const av = value(a),
            bv = value(b);
        if (av == null && bv != null) return 1;
        if (bv == null && av != null) return -1;
        const comparison =
            av == null
                ? 0
                : typeof av === 'number'
                  ? av - bv
                  : av.localeCompare(bv, undefined, {
                        numeric: true,
                        sensitivity: 'base',
                    });
        return (
            comparison * (direction === 'asc' ? 1 : -1) ||
            a.name.localeCompare(b.name)
        );
    });
}
