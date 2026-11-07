import test from 'node:test';
import assert from 'node:assert/strict';
import { clientMessageCategory, clientMessageEvent } from './clientEvents.js';

test('transfer reports and directory errors belong in the event viewer', () => {
    for (const content of [
        'UPLOAD_FILE is not implemented; no file transfer was performed.',
        'DOWNLOAD_FILE is not implemented; no file transfer was performed.',
        'Failed to save downloaded file.',
    ]) {
        assert.equal(
            clientMessageCategory({ type: 'ERROR_REPORT', content }),
            'file.transfer.error',
        );
    }
    assert.equal(
        clientMessageCategory({
            type: 'COMMAND_RESULT',
            content: 'File downloaded successfully.',
        }),
        'file.transfer.completed',
    );
    assert.equal(
        clientMessageCategory({
            type: 'ERROR_REPORT',
            content: 'Directory listing failed: Access denied',
        }),
        'directory.error',
    );
});

test('ordinary command output and errors stay in the terminal', () => {
    for (const message of [
        { type: 'ERROR_REPORT', content: 'Command exited with code 1' },
        {
            type: 'COMMAND_RESULT',
            content: 'UPLOAD_FILE is mentioned in command output',
        },
        {
            type: 'COMMAND_RESULT',
            content: 'ordinary output\nFile downloaded successfully.',
        },
    ])
        assert.equal(clientMessageCategory(message), null);
});

test('binary uploads produce metadata only and event IDs distinguish sessions', () => {
    const message = {
        type: 'FILE_UPLOAD',
        content: 'PRIVATE_BYTES',
        data_size: 13,
        timestamp: '2026-11-07T10:00:00Z',
        message_id: 5,
    };
    const event = clientMessageEvent({ ID: 'test-agent' }, message);
    assert.equal(event.detail, 'File upload received (13 bytes).');
    assert.equal(event.agent_id, 'test-agent');
    assert.equal(JSON.stringify(event).includes('PRIVATE_BYTES'), false);
    assert.equal(
        event.id,
        clientMessageEvent({ ID: 'test-agent' }, message).id,
    );
    assert.notEqual(
        event.id,
        clientMessageEvent({ ID: 'another-agent' }, message).id,
    );
});
