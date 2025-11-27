#include "ApiManager.h"
#include <iostream>


HMODULE ApiManager::GetModuleHandleManual(LPCWSTR lpModuleName)
{
    PPEB PebAddress = getPeb();
    CHAR ModuleName[MAX_PATH] = {0}; // TODO: is this = {0} needed?


    //PVOID pModule = nullptr;

    PLIST_ENTRY pListHead = &PebAddress->Ldr->InMemoryOrderModuleList;
    PLIST_ENTRY pList = PebAddress->Ldr->InMemoryOrderModuleList.Flink;
    PLDR_DATA_TABLE_ENTRY pDataTableEntry; // TODO: does this need to be initialized to nullptr?

    while (pList != pListHead)
    {
        pDataTableEntry = CONTAINING_RECORD(pList, LDR_DATA_TABLE_ENTRY, InMemoryOrderLinks);

        //if (pDataTableEntry->BaseDllName.Buffer == lpModuleName)
        if (_wcsicmp(pDataTableEntry->BaseDllName.Buffer, lpModuleName) == 0)
        {
            std::cout << "hModule was found for: " << pDataTableEntry->BaseDllName.Buffer << std::endl;
            return (HMODULE)pDataTableEntry->DllBase;
        }
        pList = pList->Flink;
    }

    return nullptr;
    //return pModule;
}

HMODULE ApiManager::GetModuleHandleManual(DWORD moduleHash)
{
    std::cout << "GetModuleHandleManual(DWORD moduleHash) is not implemented" << std::endl;
    return nullptr;
}

// TODO: rename pModule to hModule?
FARPROC ApiManager::GetProcAddressManual(PVOID pModule, LPCSTR lpProcName)
{
	PIMAGE_NT_HEADERS pNTHeader = (PIMAGE_NT_HEADERS)((ULONG_PTR)pModule + ((PIMAGE_DOS_HEADER)pModule)->e_lfanew);
	DWORD dwExportDirRVA = pNTHeader->OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress;
	PIMAGE_EXPORT_DIRECTORY pExportDir = (PIMAGE_EXPORT_DIRECTORY)((ULONG_PTR)pModule + dwExportDirRVA);

	DWORD* arrayOfFunctionRVAs = (DWORD*)((ULONG_PTR)pModule + pExportDir->AddressOfFunctions);
	DWORD* arrayOfNamesRVAs = (DWORD*)((ULONG_PTR)pModule + pExportDir->AddressOfNames);
	WORD* arrayOfNameOrdinals = (WORD*)((ULONG_PTR)pModule + pExportDir->AddressOfNameOrdinals);

	DWORD numNames = pExportDir->NumberOfNames;
	//for (DWORD i = 0; i < pExportDir->NumberOfNames; ++i)
	for (DWORD i = 0; i < numNames; ++i)
	{
		char* prodName = (char*)((ULONG_PTR)pModule + arrayOfNamesRVAs[i]);
		//WORD ordinalIndex = arrayOfNameOrdinals[i];
		//FARPROC functionAddress = (FARPROC)((ULONG_PTR)pModule + arrayOfFunctionRVAs[ordinalIndex]);

		// TODO: idk if case-insensitive strcmp is needed on here
		// printf("prodName = %s\n", prodName);
		if (_stricmp(lpProcName, prodName) == 0)
		{
            std::cout << "pModule was found for: " << prodName << std::endl;
			return (FARPROC)((ULONG_PTR)pModule + arrayOfFunctionRVAs[arrayOfNameOrdinals[i]]);
			//return (FARPROC)((ULONG_PTR)pModule + arrayOfFunctionRVAs[ordinalIndex]);
			//return functionAddress;
		}
	}
    return nullptr;
}

FARPROC ApiManager::GetProcAddressManual(HMODULE hModule, DWORD procHash)
{
    std::cout << "GetProcAddressManual(HMODULE hModule, DWORD procHash) is not implemented" << std::endl;
    return nullptr;
}

void ApiManager::loadAPIs()
{
    // TODO: hardcoded hModule names 
    hModules_[L"kernel32.dll"] = GetModuleHandleManual(L"kernel32.dll");
    hModules_[L"ntdll.dll"] = GetModuleHandleManual(L"ntdll.dll");
    hModules_[L"ws2_32.dll"] = GetModuleHandleManual(L"ws2_32.dll");

    // TransportLayer functions
    fProcedures_["socket"] = GetProcAddressManual(hModules_[L"ws2_32.dll"], "socket");
    fProcedures_["getaddrinfo"] = GetProcAddressManual(hModules_[L"ws2_32.dll"], "getaddrinfo");
    fProcedures_["connect"] = GetProcAddressManual(hModules_[L"ws2_32.dll"], "connect");
    fProcedures_["closesocket"] = GetProcAddressManual(hModules_[L"ws2_32.dll"], "closesocket");
    fProcedures_["recv"] = GetProcAddressManual(hModules_[L"ws2_32.dll"], "recv");

    // WinApi functions

    // LDRLOADDLL_HASH
}
