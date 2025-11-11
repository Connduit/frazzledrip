# TODO
- this is the larger agent/implant that can do everything?
- it's the "monolithic" second stage
- it acts as a module manager enabling true modular loading thereafter? (it should have all of the module logic precompiled inside of it?)
- make stuff multi-threaded
- Folder hierarchy refactor:
  - client
        - src
        - test
        - data (where config files live)
        - docs
        - tools (maybe stuff to help with testing? or just misc stuff)
- add "const" to functions and vars that need it
- change function signature to return by parameter instead of return by value
- switch from child classes to config file? config needs to be evaulated beforing compiling tho
- use templates instead of parent/child classes? it's harder to re

### Class and Ownership Hierarchy
		
		MalwareApp (TOP OWNER)
		│
		├── MessageHandler (OWNS execution + outgoing queue)
		│   ├── CommandExecutor
		│   ├── FileManager  
		│   └── DataHarvester
		│
		└── C2Manager (OWNS multiple C2 connections)
			└── std::vector<Transporter> (each OWNS one server)
				├── ConnectionManager
				├── Serializer
				├── Encoder
				└── Encryptor


### Core
- Beacon/check-in with server
- receive and parse tasks
- execute commands
- send results back to server
- sleep/wait until ready to send "alive" message to server


        ┌─────────────────────────────────────────────────┐
        │                 CORE IMPLANT LOOP               │
        ├─────────────────────────────────────────────────┤
        │                                                 │
        │  1. BEACON & CHECK-IN                           │
        │     → HTTP GET/POST to C2                       │
        │     → "I'm alive" signal                        │
        │     → Wait for tasks                            │
        │                                                 │
        │  2. RECEIVE & PARSE TASKS                       │
        │     ← Get command from C2 response              │
        │     ← Parse command type & parameters           │
        │     ← Prepare for execution                     │
        │                                                 │
        │  3. EXECUTE COMMANDS                            │
        │     → Run shell commands                        │
        │     → Handle file operations                    │
        │     → Capture output/errors                     │
        │                                                 │
        │  4. SEND RESULTS                                │
        │     → Format output for transmission            │
        │     → HTTP POST back to C2                      │
        │     → Confirm task completion                   │
        │                                                 │
        └─────────────────────────────────────────────────┘
                ↓
            SLEEP & REPEAT


### Features (TODO)
- process hollowing, module stomping (reuse existing loaded modules)
- reflective dll, module overloading (load dll into legit module space)

### TODO (files/functions)
- encoding: stealth shouldn't matter too much? it's ok the av knows we're doing encoding. just need to pretend we're a legitimate program sending legitimate messages/data


### References
- https://github.com/HavocFramework/Havoc/tree/main/payloads/Demon
- https://github.com/MythicAgents/Hannibal/tree/main/Payload_Type/hannibal/hannibal/agent_code/Hannibal
- https://github.com/BishopFox/sliver/tree/master/implant
- https://github.com/darkr4y/geacon/tree/master/cmd
