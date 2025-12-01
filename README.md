# frazzledrip
tcp stub

# TODO (next steps)
- getting flagged by heuristic/ML detection and non signature-based detections
  - will probs have to fix code patterns, structural characteristics, suspicious api sequences, unusual code flow
- recon (mostly just for OS so the stager knows what type of loader/payload to request. however, this kind of implies that the stager is written in a cross platform language and probably not stealthy at all. maybe have a PDF that has embedded scripts that download cross platform stuff... curl on unix powershell cmd?)
- anti av, anti debug, anti edr

### TODO (client)
- rev https stub (use logic from tcp stub)
- check for debugger
  - https://learn.microsoft.com/en-us/windows/win32/debug/debugging-functions
  - https://learn.microsoft.com/en-us/windows/win32/api/debugapi/nf-debugapi-isdebuggerpresent
- use a combination of manual api resolution and direct system calls (direct system calls are stealthier)

### Features (TODO)
- sleep obfuscation
- when server is offline, client should be able to store commands, store results locally, use dead-drop resolvers, and/or retry connections periodically

### TODO (server)
- logic written in golang
- gui written in c++ and qt gui ?

### Workflow
- dropper/installer -> loader -> stub/stager -> beacon/agent/implant (payload) -> specific modules/payloads/tasks -> profit? 
- stub -> payload (establish persistence) -> full install -> profit? 

### stub
- avoid installation / when stub is being transfer to the victim
- avoid detection while idle but installed on victim machine
- avoid detection while stub is executing / alive
- avoid detection after stub is done... clean up or find somewhere to hide?

### Counter Measures
- Ivanti Device & Application Control
- Symantec Endpoint Protection
- RTNotify?
- /////
- Zscaler
- Secret Double Octopus
- BeyondTrust Privilege Management / Avecto Privilege Guard Client
- Forcepoint DLP Endpoint
- Boldon James File Classifier
- Cisco Secure Client

### Counter-Counter Measures
- https://www.cobaltstrike.com/
- dll side-loading
- reflective dll injection
- dll as a payload container
- com / activeX dll abuse
- hide code in different/unusal sections of the PE
- payload/stub encryption (xor)
-   AES-GCM (AES in Galois/Counter Mode)
-   ChaCha20-Poly1305 — faster on devices without AES hardware acceleration
-   -
-   Simple obfuscation for quick hacks: XOR, rolling-XOR, trivial byte transform
-   Packers / crypters: UPX (both legitimate use and abused), commercial crypters and custom packers.
-   Standard symmetric crypto for serious samples: AES (often AES-128 or AES-256), frequently in authenticated modes (GCM is common in newer tooling), or ChaCha20-Poly1305 on toolchains that favor it.
-   Custom multi-stage formats: compressed → encrypted → encoded blobs inside .rsrc or custom sections; small loader/stub in .text unpacks/decrypts in memory.
-   Polymorphism/metamorphism:
- c2/network encryption (https and custom tsl wrappers)
-   TLS 1.3
-   Mutual TLS
-   Use HTTP(S) or WebSockets over TLS if you need easy firewall traversal and visibility via HTTP semantics.
-   Ephemeral key exchange (ECDHE) for Perfect Forward Secrecy (PFS).
-   -
-   HTTPS / TLS (very common): attackers prefer HTTPS because it blends with normal web traffic. Many C2 frameworks use TLS 1.2 or TLS 1.3.
-   Mutual TLS (mTLS) occasionally: more advanced operators use mTLS for client auth (makes detection/control harder).
-   WebSockets over TLS for persistent, bidirectional channels that look like normal web apps.
-   DoH / DoT / DNS tunneling when environments restrict normal web egress. DNS tunneling tools (Iodine, DNScat2 variants) or DoH (DNS over HTTPS) are used for stealthy beacons and tiny data exfil.
-   Custom encrypted TCP/UDP channels: AES/ChaCha20 used over raw sockets when stealth beyond TLS is desired.
-   Abuse of legitimate cloud/web services: GitHub Gists, Google Drive, Slack, Discord, cloud functions, or pastebins as C2 or staging; traffic to these looks normal and is hard to block broadly.
- metamorphic engine

### Windows
- workflow
-   make structure for command parsing on the client to server and retrieve commands sent from the server to the client, as well as parse them and execute the given plugin associated with the command
-   connect to our server on a given port. collect information about the system and encrpyt it. then send it to our server. establish a foothold thread for communication. check if system is a vm machine or sandbox, if it is then don't unencrypt/unpack self. also don't establish persistence, it should actually probably delete itself entirely. at this point, nothing should be written to the disk and everything should be running in memory (ram).
-   spawn zombie process svchost.exe or dllhost.exe. CreateProcessInternal, CREATE_SUSPENDED, SuspendThread
-   %appdata\\Microsoft\\453543.dll. rundll32.exe - execute our dll. rundll32.exe as a runkey.
-   it gets on system. 1. check if it's already running. 2. move to %appdata%. 3. setup persistence.
-   takes the form of a .dll file

### MacOS
- curl, osascript, launchd
- workflow:
-   Delivery → script (AppleScript/Bash/Python)
-   Download → curl or scripting network
-   Execution → script or small Mach-O payload
-   Persistence → launchd agent, cron, login item
-   Stealth → hide files, mimic system apps, automate via osascript

### Goals
- true modular loading: the initial stager downloads modules (DLLs, shared objects, native code blobs, shellcode, or scripts) and the running implant dynamically loads/executes them in memory.

# Resources
- https://attack.mitre.org/
- https://attack.mitre.org/groups/
