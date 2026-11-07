import test from 'node:test';
import assert from 'node:assert/strict';
import {
    parseListenerDrafts,
    validateListenerDraft,
} from './listenerDrafts.js';

const valid = {
    name: 'Local test',
    protocol: 'TCP',
    address: '127.0.0.1',
    port: '5555',
};

test('drafts normalize names, ports, and IP addresses', () => {
    assert.deepEqual(
        validateListenerDraft({ ...valid, name: ' Local test ' }),
        { ...valid, port: 5555 },
    );
    assert.equal(
        validateListenerDraft({ ...valid, address: '2001:0DB8:0:0::1' })
            .address,
        '2001:db8::1',
    );
    assert.equal(
        validateListenerDraft({ ...valid, address: 'LOCALHOST' }).address,
        'localhost',
    );
});

test('drafts reject invalid fields, duplicate names, and duplicate endpoints', () => {
    for (const changes of [
        { name: '' },
        { name: 'a'.repeat(65) },
        { name: 'a\nb' },
        { protocol: 'UDP' },
        { address: '256.1.1.1' },
        { address: 'example.com' },
        { address: '127.1' },
        { address: 'not:ipv6' },
        { port: '0' },
        { port: '65536' },
        { port: '1.5' },
        { port: '1e3' },
    ]) {
        assert.throws(() => validateListenerDraft({ ...valid, ...changes }));
    }
    const saved = validateListenerDraft(valid);
    assert.throws(
        () =>
            validateListenerDraft(
                { ...valid, name: 'LOCAL TEST', port: 5556 },
                [saved],
            ),
        /name already exists/,
    );
    assert.throws(
        () => validateListenerDraft({ ...valid, name: 'Other' }, [saved]),
        /bind address, and port already exists/,
    );
});

test('stored drafts tolerate corrupted data and discard invalid or duplicate entries', () => {
    assert.deepEqual(parseListenerDrafts('broken'), []);
    assert.deepEqual(parseListenerDrafts('{}'), []);
    const item = { ...valid, id: 'test-id' };
    assert.deepEqual(
        parseListenerDrafts(
            JSON.stringify([
                null,
                item,
                item,
                { ...valid, id: 'another', port: 0 },
            ]),
        ),
        [{ ...valid, port: 5555, id: 'test-id' }],
    );
});
