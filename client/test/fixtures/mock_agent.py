"""TEST-ONLY protocol fixture. Returns fixed data; never executes commands or reads files."""
import argparse
import json
import socket
import struct


def read_exact(sock, size):
    data = b""
    while len(data) < size:
        chunk = sock.recv(size - len(data))
        if not chunk:
            raise EOFError
        data += chunk
    return data


def response(kind, payload):
    command = payload.decode("utf-8", errors="replace").strip().lower()
    if kind == 8:
        if command == "test-hang":
            return None
        if command.startswith("dir "):
            return 5, b"Directory of C:\\TEST\n11/06/2026 09:30 PM 4 test.txt\n"
        fixtures = {
            "whoami": "TESTLAB\\test-user\r\n",
            "ipconfig /all": "Host Name . . . : TEST-COMPUTER\nIPv4 Address. . : 192.0.2.10(Preferred)\n",
            "systeminfo": "OS Name: TEST Windows fixture\nOS Version: 10.0.19045 N/A Build 19045\nSystem Type: x64-based PC\n",
        }
        if command == "test-fail":
            return 12, b"TEST failure: no command was executed"
        return 5, fixtures.get(command, "TEST response only: " + command).encode()
    if kind == 4:
        return 5, b"dwMajorVersion: 10\ndwMinorVersion: 0\ndwBuildNumber: 19045\nprocessorArchitecture: 9\n"
    if kind == 13:
        return 14, json.dumps({"path": "C:\\TEST", "parent": "C:\\", "entries": [
            {"name": "test.txt", "path": "C:\\TEST\\test.txt", "is_directory": False, "size": 4}
        ]}).encode()
    return None


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--port", type=int, default=4444)
    args = parser.parse_args()
    with socket.create_connection(("127.0.0.1", args.port)) as sock:
        print("TEST-ONLY mock connected; command execution and filesystem access are disabled", flush=True)
        try:
            while True:
                kind, size, request_id = struct.unpack("<III", read_exact(sock, 12))
                if size > 32 * 1024 * 1024:
                    raise ValueError("oversized test frame")
                result = response(kind, read_exact(sock, size))
                if result:
                    reply_type, data = result
                    sock.sendall(struct.pack("<III", reply_type, len(data), request_id) + data)
        except EOFError:
            pass


if __name__ == "__main__":
    main()
