# ExtraSysInfoTODO

    TODO: get locatoin of pc/person

    TODO: also maybe check if these processes are even installed
    processes like ollydbg, wireshark, procmon, processhacker, 
    idaq, x32dbg, x64dbg", fiddler, vboxservice

    isRunningInVM() {
    // Check common VM process names
    const char* vmProcesses[] = {
        "vboxservice.exe", "vboxtray.exe", 
        "vmwaretray.exe", "vmwareuser.exe",
        "xenservice.exe", "qemu-ga.exe"
    };
    
    // Check for VM-specific hardware
    HKEY hKey;
    if (RegOpenKeyExA(HKEY_LOCAL_MACHINE, 
        "HARDWARE\\DESCRIPTION\\System", 0, KEY_READ, &hKey) == ERROR_SUCCESS) {
        // Check system BIOS/vendor info
        RegCloseKey(hKey);
    }
    
    return false; // Placeholder


    DWORD getWindowsBuildNumber();

    // GetSystemInfo() // get cpu info
    typedef struct _RTL_OSVERSIONINFOW {
        DWORD dwOSVersionInfoSize;
        DWORD dwMajorVersion;
        DWORD dwMinorVersion;
        DWORD dwBuildNumber;
        DWORD dwPlatformId;
        WCHAR szCSDVersion[128];
    } RTL_OSVERSIONINFOW;


    StealthSystemInfo getStealthSystemInfo() {
        StealthSystemInfo info = {0};
        SYSTEM_INFO sysInfo;
    
    // Single API call - very stealthy
    GetNativeSystemInfo(&sysInfo);
    
    info.processorCount = sysInfo.dwNumberOfProcessors;
    info.architecture = sysInfo.wProcessorArchitecture;
    info.pageSize = sysInfo.dwPageSize;
    info.minAppAddr = (DWORD)sysInfo.lpMinimumApplicationAddress;
    info.maxAppAddr = (DWORD)sysInfo.lpMaximumApplicationAddress;
    
    return info;



    // TODO: getTickCount() // (amount of time the system has been running?)

    // TODO: getGPU()
