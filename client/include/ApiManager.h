/* ApiManager */
// this files manages all of the resolve apis
// the resolve apis are called through this class
//     remove ApiLoader.h? or maybe just have api manager call it?

#ifndef API_MANAGER_H
#define API_MANAGER_H

#include <unordered_map>

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <string> // TODO: remove this?

/*
// TODO: uncomment when im ready to remove windows.h

typedef unsigned short USHORT;
*/

typedef struct _UNICODE_STRING
{
    USHORT Length;
    USHORT MaximumLength;
    PWSTR Buffer;
} UNICODE_STRING;
typedef UNICODE_STRING *PUNICODE_STRING;


typedef struct _LDR_DATA_TABLE_ENTRY {
    LIST_ENTRY InLoadOrderLinks; // PVOID Reserved1[2];
    LIST_ENTRY InMemoryOrderLinks;
    LIST_ENTRY InInitializationOrderLinks; // PVOID Reserved2[2];
    PVOID DllBase; // pointer to the base address of where the DLL/module is loaded
    PVOID EntryPoint;
	ULONG SizeOfImage;
    UNICODE_STRING FullDllName;
    UNICODE_STRING BaseDllName;
    PVOID Reserved5[3];
    union
    {
        ULONG CheckSum;
        PVOID Reserved6;
    };
    ULONG TimeDateStamp;
} LDR_DATA_TABLE_ENTRY, *PLDR_DATA_TABLE_ENTRY;

typedef struct _PEB_LDR_DATA {
    ULONG Length;
    BOOLEAN Initialized;
    PVOID SsHandle;
    LIST_ENTRY InLoadOrderModuleList;
    LIST_ENTRY InMemoryOrderModuleList;
    LIST_ENTRY InInitializationOrderModuleList;
    //PVOID EntryInProgress;
	//#if (NTDDI_VERSION >= NTDDI_WIN7)
    //UCHAR ShutdownInProgress;
    //PVOID ShutdownThreadId;
} PEB_LDR_DATA,*PPEB_LDR_DATA;

// TODO: this struct isn't needed? we only need Ldr anyways... might have to uncomment to avoid using certain #includes tho
// If it matters: https://github.com/HavocFramework/Havoc/blob/main/payloads/DllLdr/Include/Native.h#L41
typedef struct _PEB {
	BYTE Reserved1[2];
	BYTE BeingDebugged;
	BYTE Reserved2[1];
	PVOID Reserved3[2];
	PPEB_LDR_DATA Ldr;
    /*
	PRTL_USER_PROCESS_PARAMETERS ProcessParameters;
	PVOID Reserved4[3];
	PVOID AtlThunkSListPtr;
	PVOID Reserved5;
	ULONG Reserved6;
	PVOID Reserved7;
	ULONG Reserved8;
	ULONG AtlThunkSListPtr32;
	PVOID Reserved9[45];
	BYTE Reserved10[96];
	PPS_POST_PROCESS_INIT_ROUTINE PostProcessInitRoutine;
	BYTE Reserved11[128];
	PVOID Reserved12[1];
	ULONG SessionId;
    */
} PEB,*PPEB;

/* 64bit
typedef struct _PEB {
    BYTE Reserved1[2];
    BYTE BeingDebugged;
    BYTE Reserved2[21];
    PPEB_LDR_DATA LoaderData;
    PRTL_USER_PROCESS_PARAMETERS ProcessParameters;
    BYTE Reserved3[520];
    PPS_POST_PROCESS_INIT_ROUTINE PostProcessInitRoutine;
    BYTE Reserved4[136];
    ULONG SessionId;
} PEB;
*/






// TODO: maybe rename to api loader? or libLoaderApi?
class ApiManager
{
public:
    //ApiManager(); // TODO: loads specific libs and apis based off of config??
	void loadAPIs(); // TODO: this should load all apis needed by everything in ClientSubsystem?

	// TODO: add a template return function that casts the stored 
	//		 procedure to the correct type. after creating this getter method,
	//		 i should be able to move the maps back to private

	// TODOs:
	// decltype(expr) just get the type of the variable which you can use a a type specifier for another variable
	// #define DECL_API(x) decltype(x) * x ... is just a pointer to a function (or variable) 
	std::unordered_map<LPCWSTR, HMODULE> hModules_; // TODO: use hashes eventually
	std::unordered_map<LPCSTR, FARPROC> fProcedures_; // TODO: use hashes eventually

private:


	// TODO: why is this static? change to inline? 
	static PPEB getPeb()
	{

		#if defined(_WIN64)
    		//#define GET_PEB() (PPEB) __readgsqword(0x60)
			return reinterpret_cast<PPEB>(__readgsqword(0x60));
		#else 
    		//#define GET_PEB() (PPEB) __readfsdword(0x30)
			return reinterpret_cast<PPEB>(__readgsqword(0x30));
		#endif
	}

	//PebWalker (OWNS PEB traversal logic)
	//HashManager (OWNS API hashing algorithms)... hashes should be predefined?

	// TODO: uint32_t should be the hash? initialize the unorder map with all hashes i want to 
	// define and set their associated FARPROC to nullptr? 
	//std::unordered_map<uint32_t, FARPROC> (OWNS API address cache)... maps api to resolved version?

	//
	// 1. RESOLVE phase (finding)
		//FARPROC resolveAddress(const std::string& function_name);
	// Walk PEB, find module, parse exports, get address

	// 2. LOAD phase (preparing)  
	/*
	template<typename FuncType>
	FuncType loadFunction(const string& function_name) 
	{
		FARPROC address = resolveAddress(function_name);
		return reinterpret_cast<FuncType>(address);  // Ready to call!
	}*/

	// HMODULE GetModuleHandleManualHashO(DWORD moduleHash); // using offsets
	// GetModuleHandle by literal name
	HMODULE GetModuleHandleManual(LPCWSTR lpModuleName);
	// GetModuleHandle by hash
	HMODULE GetModuleHandleManual(DWORD moduleHash);
	
	// FARPROC GetProcAddressManual(LPCSTR lpModuleName, LPCSTR lpProcName); remove version? it does GetModuleHandle inside of it
	// GetProcAddress by literal name
	FARPROC GetProcAddressManual(PVOID pModule, LPCSTR lpProcName);
	// GetProcAddress by hash
	FARPROC GetProcAddressManual(HMODULE hModule, DWORD procHash);

	// TODO: move these outside of class?
	CHAR HideChar(CHAR c) { return c; }
	WCHAR HideChar(WCHAR c) { return c; }



};

// TODO: define hashed apis here?

#endif
