import React from 'react';
import { formatMemory } from '../lib/sessions.js';

export default function ClientInfo({ clientInfo }) {
    if (!clientInfo) return <div>No client information available.</div>;
    const productTypes = {
        1: 'Workstation',
        2: 'Domain controller',
        3: 'Server',
    };
    return (
        <div className="system-info">
            <strong>System information</strong>
            <div>Hostname: {clientInfo.hostname || 'Unknown'}</div>
            <div>User: {clientInfo.user || 'Unknown'}</div>
            <div>Domain: {clientInfo.domain || 'Unknown'}</div>
            <div>
                Reported IPs: {clientInfo.ip_addresses?.join(', ') || 'Unknown'}
            </div>
            <div>
                MAC addresses:{' '}
                {clientInfo.mac_addresses?.join(', ') || 'Unknown'}
            </div>
            <div>Manufacturer: {clientInfo.manufacturer || 'Unknown'}</div>
            <div>Model: {clientInfo.model || 'Unknown'}</div>
            <div>Boot time: {clientInfo.boot_time || 'Unknown'}</div>
            <div>OS: {clientInfo.os || 'Unknown'}</div>
            <div>Architecture: {clientInfo.arch || 'Unknown'}</div>
            <div>Build: {clientInfo.build_number ?? 'Unknown'}</div>
            <div>
                Product type:{' '}
                {productTypes[clientInfo.product_type] || 'Unknown'}
            </div>
            <div>
                Logical processors:{' '}
                {clientInfo.number_of_processors ?? 'Unknown'}
            </div>
            <div>
                Physical memory:{' '}
                {formatMemory(clientInfo.total_physical_memory)}
            </div>
            <div>
                Available memory:{' '}
                {formatMemory(clientInfo.available_physical_memory)}
            </div>
        </div>
    );
}
