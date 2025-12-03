# TODO
- maintain tcp connection even after sending commands via server
- get lower score on virus total
  - try adding a sleep or benign noise in between kernel32.dll functions?
    - https://github.com/vxunderground/VX-API/blob/main/VX-API/SleepObfuscationViaVirtualProtect.cpp
- https://github.com/MythicAgents/Hannibal
- integrity checks done here to prevent adversary-in-the-middle tampering

- stager.exe should be "truly modular"... receveies dlls or shell code and executes

# NEW NOTES
- stager = downloader + loader + bootstrapper
- downloader = downloads the core payload/implant
    - the core implant can either be downloaded by connecting back to the server through the internet or the core implant can already be "packaged" inside the stager
- loader = loads the core implant into memory or downloads/writes it onto the disk
- bootstrapper = something that executes/starts/runs the downloader or loader

# Requirements
- able to load normal dlls, position independent code, and my custom dlls/modules
- able to unload normal dlls and my custom dlls/modules
- able to walk peb to resolve apis
- able to do function export forwarding
- able to parse and resolve imports, find exported function, patch the payload so it can use exported functions (even when no mz header, pe, header or import table is present and has hashed imports in the import table instead of names)
- able to do relocation table resolution
- is able to call DllMain DLL_PROCESS_ATTACH / DLL_THREAD_ATTACH / DLL PROCESS_DETACH / DLL_THREAD_DETACH
