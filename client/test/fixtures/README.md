# Mock agent (test only)

`mock_agent.py` is a loopback-only protocol fixture. It returns fixed TEST data
and echoes request IDs. It never executes commands, accesses the filesystem,
establishes persistence, or modifies the Windows client.

1. Start the Go server using a separate test data directory:
   `go run ./backend -cache-dir /tmp/frazzledrip-test` from `server`.
2. Start the listener from the dashboard.
3. Run `python3 client/test/fixtures/mock_agent.py` from the repository root.
4. Inspect TEST-COMPUTER. Initial information queries use synthetic responses.
5. Send `test-command` for a completed test task or `test-fail` for a failed one.
6. Stop the mock and restart the server to inspect persisted offline records,
   task outcomes, and audit events.

The fixture does not reconnect automatically. Restart it explicitly for reconnect
tests. Use `--port` to select another loopback TCP port. The Windows C++ client
remains unchanged in this milestone.
