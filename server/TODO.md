# Server

### TODO
- backend: written in go or python?
- gui/frontend (website): react
- gui/frontend (exe): written in c++ qt gui? (do this after react implementation?)
- TODO: write in c# so i can become more familiar with .net which will make it easier to learn powershell scripting

### Diagram

          ┌───────────────────────────────┐
          │      Frontend (Mythic UI)     │
          │  - Web interface for operator │
          └──────────────┬────────────────┘
                         │  REST/WebSocket
                         ▼
          ┌───────────────────────────────┐
          │   Mythic Core (Backend)       │
          │  - Handles tasking, DB, auth  │
          │  - Loads C2 profiles + agents │
          └──────────────┬────────────────┘
                         │
          ┌───────────────────────────────┐
          │     C2 Profile (Python)       │
          │  - Defines network layer      │
          │  - Controls how agents talk   │
          │    to Mythic (format, routes) │
          └──────────────┬────────────────┘
                         │
                  Network / Internet
                         │
          ┌───────────────────────────────┐
          │      Agent (Payload)          │
          │  - Written in C/C++/Go/etc.   │
          │  - Communicates via C2 Profile│
          └───────────────────────────────┘

##### Explanation (in data flow order)
- server ui (frontend)
- server core (backend)
  - handles commands from user input from the frontend
- c2 profile (backend)
  - does logic / preps commands and data coming in or going out of the server?
  - integrity checks done here to prevent adversary-in-the-middle tampering?
- translator (backend)
  - the translator is literally just a parser
  - converts data sent by victim machine into data that the server core actually knows how to interpret
  - converts commands sent by the server core into commands that the victim will actually know what to do with
- agent/implant (runs on victim machine)


### Diagram 2

                       +---------------------+
                       |  Frontend (Mythic)  |
                       |  - Operator UI      |
                       +----------+----------+
                                  |
                                  | issue task (UI)
                                  v
                       +---------------------+
                       |   Mythic Core       |
                       |  - Core API / DB    |
                       |  - Task queue       |
                       +---+-------------+---+
                           |             ^
           dequeue task    |             | store/display
                           v             |
               +---------------------------+
               | Agent Worker / Translator |
               | - create_tasking()        |  <-- parses & formats
               | - process_response()      |      human <> agent bytes
               +-----+---------------------+
                               |
             serialized bytes  | (agent-specific payload)
                               v
            +-------------------------+
            |   C2 Profile (Python)   |
            | - framing / headers     |
            | - encrypt/decrypt       |
            | - beaconing / jitter    |
            | - endpoint config (CDN) |
            +-----+-------------------+
                  | (transported bytes)
                  v
         +-----------------------------+
         |   Network / Proxy / CDN     |
         | - optional TLS termination  |
         | - proxying / egress infra   |
         +-----+-----------------------+
                  |
                  v
         +-----------------------------+
         |   Hannibal Agent (C)        |
         | - receive, decrypt, parse   |
         | - execute tasks             |
         | - optional staging/modules  |
         | - send back raw results     |
         +-----+-----------------------+
                  |
           raw result bytes
                  |
                  v
         +-----------------------------+
         |   Network / Proxy / CDN     |
         +-----+-----------------------+
                  |
                  v
            +--------------------+
            |  C2 Profile (cont) |
            | - unwrap/frame     |
            | - decrypt if needed|
            +-----+--------------+
                  |
                  v
        +------------------------------+
        |  Agent Worker / Translator   |
        | - process_response()         |
        | - store result, update Core  |
        +-----+------------------------+
              |
              v
         +---------------------+
         |   Mythic Core / UI  |
         | - show result to op |
         +---------------------+

