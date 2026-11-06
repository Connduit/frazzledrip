import { defineConfig } from 'vite';

const apiPaths = [
    '/start-listener',
    '/get-clients',
    '/get-client-messages',
    '/send-command',
    '/my-send-command',
    '/send-binary',
    '/send-shellcode',
    '/get-client-info',
];

export default defineConfig({
    server: {
        proxy: Object.fromEntries(apiPaths.map(path => [path, 'http://127.0.0.1:8080'])),
    },
});
