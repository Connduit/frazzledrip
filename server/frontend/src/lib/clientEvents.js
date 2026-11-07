// Legacy clients report transfer acknowledgments using command/error message types.
// Keep these narrow matches until the protocol provides explicit status types.
export function clientMessageCategory(message) {
    if (message.type === 'FILE_UPLOAD') return 'file.transfer.received';
    const text = (message.content || '').trim();
    if (
        message.type === 'COMMAND_RESULT' &&
        text === 'File downloaded successfully.'
    )
        return 'file.transfer.completed';
    if (message.type !== 'ERROR_REPORT') return null;
    if (
        /^(?:ERROR:\s*)?(?:UPLOAD_FILE|DOWNLOAD_FILE)\b/.test(text) ||
        text === 'Failed to save downloaded file.'
    )
        return 'file.transfer.error';
    if (text.startsWith('Directory listing failed:')) return 'directory.error';
    return null;
}

export function clientMessageEvent(agent, message) {
    const category = clientMessageCategory(message);
    if (!category) return null;
    return {
        id: `client:${JSON.stringify([agent.ID, message.timestamp, message.message_id, message.type])}`,
        at: message.timestamp,
        event: category,
        agent_id: agent.ID,
        source: 'client report',
        detail:
            message.type === 'FILE_UPLOAD'
                ? `File upload received (${message.data_size || 0} bytes).`
                : message.content,
    };
}
