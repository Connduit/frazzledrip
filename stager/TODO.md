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
