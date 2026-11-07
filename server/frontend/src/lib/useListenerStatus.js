import { useCallback, useEffect, useRef, useState } from 'react';
import { apiService } from '../services/api.js';

export default function useListenerStatus() {
    const [status, setStatus] = useState(null);
    const [error, setError] = useState('');
    const [checkedAt, setCheckedAt] = useState(null);
    const [checking, setChecking] = useState(false);
    const mounted = useRef(false);
    const sequence = useRef(0);

    const refresh = useCallback(async () => {
        const token = ++sequence.current;
        if (mounted.current) setChecking(true);
        try {
            const next = await apiService.getListenerStatus();
            if (
                typeof next?.running !== 'boolean' ||
                !Number.isInteger(next.port) ||
                (next.running && (next.port < 1 || next.port > 65535))
            ) {
                throw new Error('Unexpected listener status response');
            }
            if (mounted.current && token === sequence.current) {
                setStatus(next);
                setCheckedAt(new Date());
                setError('');
            }
            return next;
        } catch (err) {
            if (mounted.current && token === sequence.current)
                setError(err.message);
            throw err;
        } finally {
            if (mounted.current && token === sequence.current)
                setChecking(false);
        }
    }, []);

    useEffect(() => {
        mounted.current = true;
        let cancelled = false,
            timer;
        async function poll() {
            try {
                await refresh();
            } catch {
                /* Keep the last known state and display the error. */
            }
            if (!cancelled) timer = setTimeout(poll, 2000);
        }
        poll();
        return () => {
            cancelled = true;
            mounted.current = false;
            sequence.current++;
            clearTimeout(timer);
        };
    }, [refresh]);

    return { status, error, checkedAt, checking, refresh };
}
