function TabContainer({ activeTab, setActiveTab, selectedClient, clients, onClientSelect, setStatus, refreshMessages }) {
    return (
        <div className="section">
            <div className="tab-container">
                <div className="tab-buttons">
                    <button 
                        className={`tab-button ${activeTab === 'commands' ? 'active' : ''}`}
                        onClick={() => setActiveTab('commands')}
                    >
                        📝 Commands
                    </button>
                    <button 
                        className={`tab-button ${activeTab === 'files' ? 'active' : ''}`}
                        onClick={() => setActiveTab('files')}
                    >
                        📦 File Transfer
                    </button>
                    <button 
                        className={`tab-button ${activeTab === 'shellcode' ? 'active' : ''}`}
                        onClick={() => setActiveTab('shellcode')}
                    >
                        ⚡ Shellcode
                    </button>
                    <button 
                        className={`tab-button ${activeTab === 'config' ? 'active' : ''}`}
                        onClick={() => setActiveTab('config')}
                    >
                        ⚙️ Configuration
                    </button>
                </div>
                
                <div className="tab-content">
                    {activeTab === 'commands' && (
                        <CommandTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                            refreshMessages={refreshMessages}
                        />
                    )}
                    
                    {activeTab === 'files' && (
                        <FileTransferTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                        />
                    )}
                    
                    {activeTab === 'shellcode' && (
                        <ShellcodeTab 
                            selectedClient={selectedClient}
                            clients={clients}
                            onClientSelect={onClientSelect}
                            setStatus={setStatus}
                        />
                    )}
                    
                    {activeTab === 'config' && (
                        <ConfigTab 
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