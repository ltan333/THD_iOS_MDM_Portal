"use client";

import React, { useState } from "react";
import { Table, Select, Input, Button, Tag, Modal, Tabs, Dropdown } from "antd";
import type { MenuProps } from "antd";
import { 
    Search, 
    Apple, 
    Smartphone, 
    Monitor,
    Battery,
    HardDrive,
    User,
    RefreshCcw,
    ChevronDown,
    Filter,
    Shield,
    Info,
    Wifi,
    Settings2,
    AppWindow,
    FileText,
    MapPin,
    FolderPlus,
    FilePlus2,
    MoreVertical
} from "lucide-react";
import type { ColumnsType } from "antd/es/table";

interface DeviceProfile {
    id: string;
    name: string;
    status: "active" | "pending" | "failed";
}

interface DeviceType {
    key: string;
    name: string;
    owner: string;
    enrollmentDate: string;
    os: "iOS" | "iPadOS" | "macOS";
    model: string;
    version: string;
    batteryPercent: number;
    availableMemory: string;
    totalMemory: string;
    addedBy: string;
    lastSeen: string;
    status: "online" | "offline";
    profiles: DeviceProfile[];
    serialNumber: string;
    ipAddress: string;
}

const mockData: DeviceType[] = [
    {
        key: "1",
        name: "iPhone",
        owner: "Nguyễn Văn A",
        enrollmentDate: "2026-03-20 08:00 AM",
        os: "iOS",
        model: "iPhone 15 Pro Max",
        version: "17.4.1",
        batteryPercent: 85,
        availableMemory: "124 GB",
        totalMemory: "256 GB",
        addedBy: "Admin",
        lastSeen: "2026-03-26 10:23 AM",
        status: "online",
        serialNumber: "F1234567890",
        ipAddress: "192.168.1.101",
        profiles: [
            { id: "p1", name: "Allow SOTI MobiControl", status: "active" },
            { id: "p2", name: "Corporate Wi-Fi", status: "active" }
        ]
    },
    {
        key: "2",
        name: "iPad",
        owner: "Marketing Dept",
        enrollmentDate: "2026-03-21 09:30 AM",
        os: "iPadOS",
        model: "iPad Pro 11-inch (M4)",
        version: "17.5",
        batteryPercent: 42,
        availableMemory: "45 GB",
        totalMemory: "128 GB",
        addedBy: "Le An",
        lastSeen: "2026-03-25 04:12 PM",
        status: "offline",
        serialNumber: "F0987654321",
        ipAddress: "192.168.1.105",
        profiles: [
            { id: "p1", name: "Allow SOTI MobiControl", status: "active" },
            { id: "p3", name: "Disable Camera", status: "pending" }
        ]
    },
    {
        key: "3",
        name: "MacBook Pro",
        owner: "Dev Team",
        enrollmentDate: "2026-03-22 10:15 AM",
        os: "macOS",
        model: "MacBook Pro 14-inch (M3)",
        version: "14.4.1",
        batteryPercent: 100,
        availableMemory: "312 GB",
        totalMemory: "512 GB",
        addedBy: "Huy",
        lastSeen: "2026-03-26 09:15 AM",
        status: "online",
        serialNumber: "M1234567890",
        ipAddress: "192.168.1.110",
        profiles: [
            { id: "p1", name: "Allow SOTI MobiControl", status: "active" },
            { id: "p4", name: "Developer Tools Config", status: "active" }
        ]
    },
    {
        key: "4",
        name: "iPhone",
        owner: "CEO",
        enrollmentDate: "2026-03-23 11:00 AM",
        os: "iOS",
        model: "iPhone 16",
        version: "18.0",
        batteryPercent: 92,
        availableMemory: "400 GB",
        totalMemory: "512 GB",
        addedBy: "Admin",
        lastSeen: "2026-03-26 02:30 PM",
        status: "online",
        serialNumber: "F5555555555",
        ipAddress: "192.168.1.102",
        profiles: [
            { id: "p1", name: "Allow SOTI MobiControl", status: "active" },
            { id: "p2", name: "Corporate Wi-Fi", status: "active" },
            { id: "p5", name: "Executive VPN", status: "active" }
        ]
    }
];

export function DevicesList() {
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [selectedDevice, setSelectedDevice] = useState<DeviceType | null>(null);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isGroupModalVisible, setIsGroupModalVisible] = useState(false);
    const [isProfileModalVisible, setIsProfileModalVisible] = useState(false);

    const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
        setSelectedRowKeys(newSelectedRowKeys);
    };

    const handleDeviceClick = (record: DeviceType) => {
        setSelectedDevice(record);
        setIsModalVisible(true);
    };

    const rowSelection = {
        selectedRowKeys,
        onChange: onSelectChange,
    };

    const actionMenu: MenuProps['items'] = [
        {
            key: 'add-to-group',
            icon: <FolderPlus className="w-4 h-4" />,
            label: 'Add to Group',
            onClick: () => setIsGroupModalVisible(true)
        },
        {
            key: 'assign-profile',
            icon: <FilePlus2 className="w-4 h-4" />,
            label: 'Assign Profile',
            onClick: () => setIsProfileModalVisible(true)
        }
    ];

    const columns: ColumnsType<DeviceType> = [
        {
            title: "DEVICE NAME",
            dataIndex: "name",
            key: "name",
            render: (text, record) => (
                <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                        {record.os === "macOS" ? (
                            <Monitor className="w-4 h-4 text-slate-600" />
                        ) : (
                            <Smartphone className="w-4 h-4 text-slate-600" />
                        )}
                    </div>
                    <div className="flex flex-col">
                        <a 
                            href="#" 
                            className="text-[#de2a15] hover:text-[#c22412] font-medium transition-colors"
                            onClick={(e) => {
                                e.preventDefault();
                                handleDeviceClick(record);
                            }}
                        >
                            {text}
                        </a>
                        <span className="text-xs text-slate-500">{record.model}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "OS",
            dataIndex: "os",
            key: "os",
            render: (text) => (
                <div className="flex items-center gap-2 text-slate-700">
                    <Apple className="w-4 h-4" fill="currentColor" />
                    <span className="font-medium">{text}</span>
                </div>
            ),
        },
        {
            title: "VERSION",
            dataIndex: "version",
            key: "version",
            render: (text) => <span className="text-slate-700 font-mono text-sm">{text}</span>,
        },
        {
            title: "BATTERY",
            dataIndex: "batteryPercent",
            key: "batteryPercent",
            render: (percent) => (
                <div className="flex items-center gap-2">
                    <Battery className={`w-4 h-4 ${percent < 20 ? 'text-red-500' : percent < 50 ? 'text-orange-500' : 'text-emerald-500'}`} />
                    <span className="text-slate-700">{percent}%</span>
                </div>
            ),
        },
        {
            title: "AVAILABLE MEMORY",
            dataIndex: "availableMemory",
            key: "availableMemory",
            render: (text, record) => (
                <div className="flex items-center gap-2">
                    <HardDrive className="w-4 h-4 text-slate-400" />
                    <span className="text-slate-700">{text}</span>
                    <span className="text-slate-400 text-xs">/ {record.totalMemory}</span>
                </div>
            ),
        },
        {
            title: "STATUS",
            dataIndex: "status",
            key: "status",
            render: (status) => (
                <Tag color={status === "online" ? "success" : "default"} className="rounded-full px-3">
                    {status.toUpperCase()}
                </Tag>
            ),
        },
        {
            title: "LAST SEEN",
            dataIndex: "lastSeen",
            key: "lastSeen",
            render: (text) => <span className="text-slate-500 text-sm">{text}</span>,
        },
    ];

    return (
        <div className="flex flex-col h-[calc(100vh-64px)] bg-slate-50 relative border-none overflow-hidden rounded-none shadow-none z-0">
            {/* Top Toolbar */}
            <div className="flex flex-wrap items-center justify-between p-4 gap-4 bg-white border-b border-slate-200 z-10 shadow-sm">
                <div className="flex items-center gap-6">
                    <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-700">OS:</span>
                        <Select className="cursor-pointer" defaultValue="all"
                            style={{ width: 120 }}
                            options={[
                                { value: "all", label: "All" },
                                { value: "ios", label: "iOS" },
                                { value: "ipados", label: "iPadOS" },
                                { value: "macos", label: "macOS" },
                            ]}
                        />
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-700">Status:</span>
                        <Select className="cursor-pointer" defaultValue="all"
                            style={{ width: 120 }}
                            options={[
                                { value: "all", label: "All" },
                                { value: "online", label: "Online" },
                                { value: "offline", label: "Offline" },
                            ]}
                        />
                    </div>
                </div>

                <div className="flex items-center gap-3">
                    <div className="flex group">
                        <Input
                            placeholder="Search device name, model..."
                            prefix={<Search className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" />}
                            className="w-64 h-8 rounded-r-none border-r-0 hover:border-[#de2a15] focus:border-[#de2a15] focus:shadow-none transition-colors"
                        />
                        <Button 
                            type="primary" 
                            className="bg-[#de2a15] hover:bg-[#c22412] rounded-l-none h-8 w-10 px-0 flex items-center justify-center border-none shadow-sm transition-colors"
                            icon={<Search className="w-4 h-4 text-white" strokeWidth={2.5} />}
                        />
                    </div>
                </div>
            </div>

            {/* Sub Toolbar / Table Controls */}
            <div className="flex items-center justify-between px-4 py-3 bg-slate-50 border-b border-slate-200 z-10">
                <div className="flex items-center gap-4 text-sm text-slate-600">
                    <span className="font-bold text-slate-800 tracking-wide uppercase">
                        DEVICES <span className="font-normal text-slate-500">(1 - 4 of 4)</span>
                    </span>
                    
                    <div className="flex items-center gap-2 border-l border-slate-300 pl-4">
                        <Select className="cursor-pointer" defaultValue="50"
                            size="small"
                            style={{ width: 70 }}
                            options={[
                                { value: "25", label: "25" },
                                { value: "50", label: "50" },
                                { value: "100", label: "100" },
                            ]}
                        />
                        <span>Per Page</span>
                    </div>

                    <div className="flex items-center gap-2 border-l border-slate-300 pl-4">
                        <Button type="text" size="small" disabled className="text-slate-400">&larr;</Button>
                        <span className="text-[#de2a15] font-bold">1 of 1</span>
                        <Button type="text" size="small" disabled className="text-slate-400">&rarr;</Button>
                    </div>

                    <div className="border-l border-slate-300 pl-4">
                        <Button type="text" size="small" icon={<RefreshCcw className="w-4 h-4 text-slate-500 hover:text-slate-800" />} />
                    </div>
                </div>

                <div className="flex items-center gap-4 text-sm">
                    {selectedRowKeys.length > 0 && (
                        <div className="flex items-center gap-2 mr-4 border-r border-slate-300 pr-4">
                            <span className="text-[#de2a15] font-medium">{selectedRowKeys.length} selected</span>
                            <Dropdown menu={{ items: actionMenu }} trigger={['click']} placement="bottomRight">
                                <Button size="small" type="primary" className="bg-[#de2a15] hover:bg-[#c22412] flex items-center gap-1">
                                    Actions <ChevronDown className="w-3 h-3" />
                                </Button>
                            </Dropdown>
                        </div>
                    )}
                    <span className="flex items-center gap-1 cursor-pointer text-slate-700 hover:text-slate-900 font-medium">
                        Columns (7) <ChevronDown className="w-4 h-4" />
                    </span>
                    <Button type="text" icon={<Filter className="w-4 h-4 text-[#de2a15]" />} className="text-[#de2a15] bg-red-50 hover:bg-red-100 rounded-full w-8 h-8 flex items-center justify-center p-0 border border-red-200 transition-colors" />
                </div>
            </div>

            {/* Table */}
            <div className="flex-1 overflow-auto border-t border-slate-200 z-10 relative scrollbar-hide">
                <Table
                    rowSelection={rowSelection}
                    columns={columns}
                    dataSource={mockData}
                    pagination={false}
                    className="custom-data-table"
                    rowClassName="hover:bg-slate-50 transition-colors cursor-pointer"
                    onRow={(record) => ({
                        onClick: () => handleDeviceClick(record),
                    })}
                />
            </div>

            {/* Device Detail Modal */}
            <Modal
                title={
                    <div className="flex items-center gap-3 pb-2">
                        <div className="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center">
                            {selectedDevice?.os === "macOS" ? (
                                <Monitor className="w-5 h-5 text-slate-700" />
                            ) : (
                                <Smartphone className="w-5 h-5 text-slate-700" />
                            )}
                        </div>
                        <div>
                            <div className="font-bold text-slate-800 text-lg">{selectedDevice?.name}</div>
                            <div className="text-xs text-slate-500 flex items-center gap-2">
                                <span className={selectedDevice?.status === "online" ? "text-emerald-500 font-medium" : "text-slate-500"}>
                                    ● {selectedDevice?.status === "online" ? "Online" : "Offline"}
                                </span>
                                <span>|</span>
                                <span>Last seen: {selectedDevice?.lastSeen}</span>
                            </div>
                        </div>
                    </div>
                }
                open={isModalVisible}
                onCancel={() => setIsModalVisible(false)}
                footer={null}
                width={800}
                className="custom-modal"
                styles={{
                    body: { padding: 0, backgroundColor: '#f8fafc', height: '70vh' },
                    content: { padding: 0, overflow: 'hidden' },
                    header: { padding: '20px 24px 0', borderBottom: '1px solid #e2e8f0', margin: 0 }
                }}
            >
                {selectedDevice && (
                    <Tabs
                        defaultActiveKey="info"
                        className="custom-tabs h-full flex flex-col"
                        items={[
                            {
                                key: "info",
                                label: (
                                    <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
                                        <Info className="w-4 h-4 text-slate-500" />
                                        DEVICE INFO
                                    </div>
                                ),
                                children: (
                                    <div className="p-6 overflow-y-auto h-full scrollbar-hide space-y-6">
                                        <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
                                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-4 border-b border-slate-100 pb-2 flex items-center gap-2">
                                                <Apple className="w-4 h-4 text-slate-600" fill="currentColor" /> System Details
                                            </h3>
                                            <div className="grid grid-cols-2 gap-y-4 gap-x-6">
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Owner</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.owner}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Enrollment Date</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.enrollmentDate}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Model</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.model}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">OS Version</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.os} {selectedDevice.version}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Serial Number</div>
                                                    <div className="text-sm text-slate-800 font-mono">{selectedDevice.serialNumber}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Added By</div>
                                                    <div className="flex items-center gap-1.5 text-sm text-slate-800 font-medium">
                                                        <User className="w-3.5 h-3.5 text-slate-400" /> {selectedDevice.addedBy}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
                                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-4 border-b border-slate-100 pb-2 flex items-center gap-2">
                                                <HardDrive className="w-4 h-4 text-slate-600" /> Hardware Status
                                            </h3>
                                            <div className="grid grid-cols-2 gap-y-4 gap-x-6">
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Battery Level</div>
                                                    <div className="flex items-center gap-2">
                                                        <Battery className={`w-4 h-4 ${selectedDevice.batteryPercent < 20 ? 'text-red-500' : selectedDevice.batteryPercent < 50 ? 'text-orange-500' : 'text-emerald-500'}`} />
                                                        <span className="text-sm text-slate-800 font-medium">{selectedDevice.batteryPercent}%</span>
                                                    </div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Storage</div>
                                                    <div className="text-sm text-slate-800 font-medium">
                                                        {selectedDevice.availableMemory} free of {selectedDevice.totalMemory}
                                                    </div>
                                                    <div className="w-full bg-slate-100 rounded-full h-1.5 mt-2">
                                                        <div 
                                                            className="bg-blue-500 h-1.5 rounded-full" 
                                                            style={{ width: `${100 - (parseInt(selectedDevice.availableMemory) / parseInt(selectedDevice.totalMemory) * 100)}%` }}
                                                        ></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        
                                        <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
                                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-4 border-b border-slate-100 pb-2 flex items-center gap-2">
                                                <Wifi className="w-4 h-4 text-slate-600" /> Network
                                            </h3>
                                            <div className="grid grid-cols-2 gap-y-4 gap-x-6">
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">IP Address</div>
                                                    <div className="text-sm text-slate-800 font-mono">{selectedDevice.ipAddress}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                )
                            },
                            {
                                key: "profiles",
                                label: (
                                    <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
                                        <Settings2 className="w-4 h-4 text-slate-500" />
                                        CONFIGURATIONS
                                    </div>
                                ),
                                children: (
                                    <div className="p-6 overflow-y-auto h-full scrollbar-hide">
                                        <div className="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm">
                                            <div className="px-5 py-4 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
                                                <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide m-0">Installed Profiles</h3>
                                                <Tag className="m-0 bg-blue-50 text-blue-700 border-blue-200 font-medium">{selectedDevice.profiles.length} Profiles</Tag>
                                            </div>
                                            <div className="divide-y divide-slate-100">
                                                {selectedDevice.profiles.map(profile => (
                                                    <div key={profile.id} className="p-4 flex items-center justify-between hover:bg-slate-50 transition-colors">
                                                        <div className="flex items-center gap-3">
                                                            <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center border border-slate-200">
                                                                <Shield className="w-4 h-4 text-slate-500" />
                                                            </div>
                                                            <div>
                                                                <div className="text-sm font-medium text-slate-800">{profile.name}</div>
                                                                <div className="text-xs text-slate-500 mt-0.5">Profile Identifier: com.company.{profile.id}</div>
                                                            </div>
                                                        </div>
                                                        <Tag color={profile.status === 'active' ? 'success' : profile.status === 'failed' ? 'error' : 'warning'} className="rounded-full px-3">
                                                            {profile.status.toUpperCase()}
                                                        </Tag>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                )
                            },
                            {
                                key: "applications",
                                label: (
                                    <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
                                        <AppWindow className="w-4 h-4 text-slate-500" />
                                        APPLICATIONS
                                    </div>
                                ),
                                children: (
                                    <div className="p-6 overflow-y-auto h-full scrollbar-hide flex items-center justify-center">
                                        <div className="text-center">
                                            <AppWindow className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                                            <h3 className="text-lg font-medium text-slate-700">Applications History</h3>
                                            <p className="text-slate-500 max-w-sm mt-2">Data will be fetched from API and displayed here in future updates.</p>
                                        </div>
                                    </div>
                                )
                            },
                            {
                                key: "logs",
                                label: (
                                    <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
                                        <FileText className="w-4 h-4 text-slate-500" />
                                        LOGS
                                    </div>
                                ),
                                children: (
                                    <div className="p-6 overflow-y-auto h-full scrollbar-hide flex items-center justify-center">
                                        <div className="text-center">
                                            <FileText className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                                            <h3 className="text-lg font-medium text-slate-700">Device Logs</h3>
                                            <p className="text-slate-500 max-w-sm mt-2">Data will be fetched from API and displayed here in future updates.</p>
                                        </div>
                                    </div>
                                )
                            },
                            {
                                key: "location",
                                label: (
                                    <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
                                        <MapPin className="w-4 h-4 text-slate-500" />
                                        LOCATION
                                    </div>
                                ),
                                children: (
                                    <div className="p-6 overflow-y-auto h-full scrollbar-hide flex items-center justify-center">
                                        <div className="text-center">
                                            <MapPin className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                                            <h3 className="text-lg font-medium text-slate-700">Device Location</h3>
                                            <p className="text-slate-500 max-w-sm mt-2">Data will be fetched from API and displayed here in future updates.</p>
                                        </div>
                                    </div>
                                )
                            }
                        ]}
                    />
                )}
            </Modal>

            {/* Add to Group Modal */}
            <Modal
                title="Add Devices to Group"
                open={isGroupModalVisible}
                onOk={() => {
                    setIsGroupModalVisible(false);
                    // Add logic to save devices to group
                }}
                onCancel={() => setIsGroupModalVisible(false)}
                okText="Add to Group"
                okButtonProps={{ className: "bg-[#de2a15] hover:bg-[#c22412]" }}
            >
                <div className="py-4">
                    <p className="mb-4 text-slate-600">
                        Select a group to add the {selectedRowKeys.length} selected device(s) to:
                    </p>
                    <Select
                        className="w-full"
                        placeholder="Select a group"
                        options={[
                            { value: 'g1', label: 'Executive Team' },
                            { value: 'g2', label: 'Development Team' },
                            { value: 'g3', label: 'Marketing Department' }
                        ]}
                    />
                </div>
            </Modal>

            {/* Assign Profile Modal */}
            <Modal
                title="Assign Configuration Profile"
                open={isProfileModalVisible}
                onOk={() => {
                    setIsProfileModalVisible(false);
                    // Add logic to assign profile
                }}
                onCancel={() => setIsProfileModalVisible(false)}
                okText="Assign Profile"
                okButtonProps={{ className: "bg-[#de2a15] hover:bg-[#c22412]" }}
            >
                <div className="py-4">
                    <p className="mb-4 text-slate-600">
                        Select a configuration profile to assign to the {selectedRowKeys.length} selected device(s):
                    </p>
                    <Select
                        className="w-full"
                        placeholder="Select a profile"
                        options={[
                            { value: 'p1', label: 'Allow SOTI MobiControl' },
                            { value: 'p2', label: 'Corporate Wi-Fi' },
                            { value: 'p3', label: 'Disable Camera' },
                            { value: 'p4', label: 'Developer Tools Config' },
                            { value: 'p5', label: 'Executive VPN' }
                        ]}
                    />
                </div>
            </Modal>

            {/* Custom Styles */}
            <style jsx global>{`
                /* Hide scrollbar for a cleaner look */
                .scrollbar-hide::-webkit-scrollbar {
                    display: none;
                }
                .scrollbar-hide {
                    -ms-overflow-style: none;
                    scrollbar-width: none;
                }
                
                /* Table modifications */
                .custom-data-table {
                    background: #ffffff !important;
                }
                .custom-data-table .ant-table {
                    background: #ffffff !important;
                }
                .custom-data-table .ant-table-container {
                    background: #ffffff !important;
                }
                .custom-data-table .ant-table-thead > tr > th {
                    background: #f8fafc !important;
                    color: #334155 !important;
                    font-weight: 600 !important;
                    font-size: 13px !important;
                    border-bottom: 1px solid #e2e8f0 !important;
                    padding: 12px 16px !important;
                }
                .custom-data-table .ant-table-tbody > tr > td {
                    background: #ffffff !important;
                    padding: 12px 16px !important;
                    border-bottom: 1px solid #f1f5f9 !important;
                    font-size: 14px !important;
                    transition: background-color 0.2s ease;
                }
                .custom-data-table .ant-table-tbody > tr:hover > td {
                    background: #f8fafc !important;
                }
                
                /* Modal Tabs Customization */
                .custom-modal .ant-modal-content {
                    border-radius: 12px;
                    overflow: hidden;
                }
                
                .custom-tabs {
                    height: 100%;
                    display: flex;
                    flex-direction: column;
                }
                .custom-tabs .ant-tabs-nav {
                    margin-bottom: 0 !important;
                    padding: 0 16px;
                    border-bottom: 1px solid #e2e8f0;
                    background: #ffffff;
                }
                .custom-tabs .ant-tabs-tab {
                    padding: 16px 0 !important;
                    margin: 0 16px 0 0 !important;
                }
                .custom-tabs .ant-tabs-tab-active .ant-tabs-tab-btn {
                    color: #de2a15 !important;
                    font-weight: 600 !important;
                }
                .custom-tabs .ant-tabs-ink-bar {
                    background: #de2a15 !important;
                    height: 3px !important;
                    border-radius: 3px 3px 0 0;
                }
                .custom-tabs .ant-tabs-content-holder {
                    flex: 1;
                    overflow: hidden;
                }
                .custom-tabs .ant-tabs-content {
                    height: 100%;
                }
                .custom-tabs .ant-tabs-tabpane {
                    height: 100%;
                }
            `}</style>
        </div>
    );
}
