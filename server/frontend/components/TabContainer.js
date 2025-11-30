// Main tab container that renders different C2 interface components
// depending on which tab is currently active.
function TabContainer({ activeTab, setActiveTab, selectedClient, clients, onClientSelect, setStatus, refreshMessages }) {
    return (
        <div className="section">
            <div className="tab-container">

                {/* ---- TAB BUTTONS (Navigation bar) ---- */}
                <div className="tab-buttons">

                    {/* Commands tab button */}
                    <button 
                        className={`tab-button ${activeTab === 'commands' ? 'active' : ''}`}
                        onClick={() => setActiveTab('commands')}
                    >
                        📝 Commands
                    </button>

                    {/* TODO: MyCommands tab button */}
                    <button 
                        className={`tab-button ${activeTab === 'mycommands' ? 'active' : ''}`}
                        onClick={() => setActiveTab('mycommands')}
                    >
                        📝 MyCommands
                    </button>

                    {/* File transfer tab button */}
                    <button 
                        className={`tab-button ${activeTab === 'files' ? 'active' : ''}`}
                        onClick={() => setActiveTab('files')}
                    >
                        📦 File Transfer
                    </button>

                    {/* Shellcode execution tab */}
                    <button 
                        className={`tab-button ${activeTab === 'shellcode' ? 'active' : ''}`}
                        onClick={() => setActiveTab('shellcode')}
                    >
                        ⚡ Shellcode
                    </button>
                </div>

                {/* ---- TAB CONTENT AREA ---- */}
                <div className="tab-content">

                    {/* Commands tab content */}
                    {activeTab === 'commands' && (
                        <CommandTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                            refreshMessages={refreshMessages}
                        />
                    )}

                    {/* TODO: MyCommands tab content */}
                    {activeTab === 'mycommands' && (
                        <MyCommandTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                            refreshMessages={refreshMessages}
                        />
                    )}
                    
                    {/* File Transfer tab content */}
                    {activeTab === 'files' && (
                        <FileTransferTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                        />
                    )}
                    
                    {/* Shellcode execution tab content */}
                    {activeTab === 'shellcode' && (
                        <ShellcodeTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                        />
                    )}
                </div>
            </div>
        </div>
    );
}
