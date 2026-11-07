import test from 'node:test';
import assert from 'node:assert/strict';
import { isConsoleMessage, parseConsoleCommand } from './consoleCommands.js';

test('console accepts only its explicit custom command and local controls', () => {
    assert.deepEqual(parseConsoleCommand('  sys_info '), {
        kind: 'request',
        messageType: 'SYS_INFO',
    });
    assert.equal(parseConsoleCommand('help').kind, 'help');
    assert.equal(parseConsoleCommand('clear').kind, 'clear');
    for (const input of [
        'whoami',
        'ipconfig',
        'systeminfo',
        'SYS_INFO && whoami',
        'SYS_INFO\nwhoami',
        'EXECUTE_COMMAND',
        'EXECUTE_SHELL_CODE',
        'COMMAND_RESULT',
        'UPLOAD_FILE',
    ]) {
        assert.equal(parseConsoleCommand(input).kind, 'unknown', input);
    }
});

test('console shows native responses, excluding shell output and transfer status', () => {
    assert.equal(isConsoleMessage({ type: 'SYS_INFO' }), true);
    assert.equal(
        isConsoleMessage({ type: 'COMMAND_RESULT', message_id: 99 }),
        true,
    );
    assert.equal(
        isConsoleMessage({
            type: 'ERROR_REPORT',
            message_id: 99,
            content: 'Native request failed',
        }),
        true,
    );
    assert.equal(
        isConsoleMessage({ type: 'COMMAND_RESULT', message_id: 105 }),
        false,
    );
    assert.equal(
        isConsoleMessage({
            type: 'ERROR_REPORT',
            message_id: 105,
            content: 'Command exited with code 1',
        }),
        false,
    );
    assert.equal(
        isConsoleMessage({
            type: 'ERROR_REPORT',
            message_id: 99,
            content: 'UPLOAD_FILE is not implemented',
        }),
        false,
    );
});
