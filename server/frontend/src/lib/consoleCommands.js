import { clientMessageCategory } from './clientEvents.js';

// Explicit allowlist: never fall back to EXECUTE_COMMAND or forward input text.
export function parseConsoleCommand(input) {
    switch (input.trim().toUpperCase()) {
        case 'HELP':
            return { kind: 'help' };
        case 'CLEAR':
            return { kind: 'clear' };
        case 'SYS_INFO':
            return { kind: 'request', messageType: 'SYS_INFO' };
        default:
            return { kind: 'unknown' };
    }
}

export function isConsoleMessage(message) {
    if (clientMessageCategory(message)) return false;
    if (message.type === 'SYS_INFO') return true;
    // The current Windows SYS_INFO handler replies as COMMAND_RESULT with ID 99.
    return (
        message.message_id === 99 &&
        ['COMMAND_RESULT', 'ERROR_REPORT'].includes(message.type)
    );
}
