import { defineConfig } from 'vite';

const apiPaths = [
    '/browse-files',
    '/stop-listener',
    '/listener-status',
    '/disconnect-session',
    '/link-session',
    '/start-listener',
    '/get-clients',
    '/get-client-messages',
    '/send-command',
    '/my-send-command',
    '/send-binary',
    '/send-shellcode',
    '/get-client-info',
];

const backendURL = process.env.BACKEND_URL || 'http://127.0.0.1:8081';

export default defineConfig({
    server: {
        proxy: Object.fromEntries(apiPaths.map(path => [path, backendURL])),
    },
});
