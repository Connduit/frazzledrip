import test from 'node:test';
import assert from 'node:assert/strict';
import { activeTabAfterClose, moveTab, orderedTabs } from './workspaceTabs.js';

test('saved order ignores stale and duplicate IDs and includes new tabs', () => {
    const tabs = [{ id: 'console' }, { id: 'events' }, { id: 'files' }];
    assert.deepEqual(
        orderedTabs(tabs, ['files', 'removed', 'files']).map((tab) => tab.id),
        ['files', 'console', 'events'],
    );
});

test('tabs move before and after targets without mutating their original order', () => {
    const ids = ['console', 'events', 'files'];
    assert.deepEqual(moveTab(ids, 'files', 'console'), [
        'files',
        'console',
        'events',
    ]);
    assert.deepEqual(moveTab(ids, 'console', 'files', 'after'), [
        'events',
        'files',
        'console',
    ]);
    assert.deepEqual(moveTab(ids, 'console', 'events', 'before'), ids);
    assert.deepEqual(moveTab(ids, 'events', 'events'), ids);
    assert.deepEqual(moveTab(ids, 'missing', 'console'), ids);
    assert.deepEqual(ids, ['console', 'events', 'files']);
});

test('closing selects the next visible neighbor only when the active tab closes', () => {
    const ids = ['events', 'files', 'console'];
    assert.equal(activeTabAfterClose(ids, 'events', 'files'), 'events');
    assert.equal(activeTabAfterClose(ids, 'files', 'files'), 'console');
    assert.equal(activeTabAfterClose(ids, 'console', 'console'), 'files');
    assert.equal(activeTabAfterClose(['events'], 'events', 'events'), '');
    assert.equal(activeTabAfterClose(ids, 'events', 'missing'), 'events');
});
