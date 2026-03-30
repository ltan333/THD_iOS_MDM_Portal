"use client";

import React, { useState } from "react";
import { Table, Input, Button, Tag, Drawer, Select, Dropdown, MenuProps, Modal, Form } from "antd";
import { 
    Search, 
    Plus, 
    Smartphone, 
    Monitor,
    Users,
    ChevronDown,
    MoreVertical,
    FolderOpen,
    FolderPlus,
    FilePlus2,
    Trash2,
    Settings
} from "lucide-react";
import type { ColumnsType } from "antd/es/table";

interface Device {
    id: string;
    name: string;
    os: "iOS" | "iPadOS" | "macOS";
    model: string;
    status: "online" | "offline";
    lastSeen: string;
}

interface DeviceGroup {
    id: string;
    name: string;
    description: string;
    deviceCount: number;
    createdDate: string;
    devices: Device[];
}

const mockGroups: DeviceGroup[] = [
    {
        id: "g1",
        name: "Executive Team",
        description: "Devices belonging to the executive team",
        deviceCount: 2,
        createdDate: "2026-01-15",
        devices: [
            { id: "d1", name: "CEO's iPhone", os: "iOS", model: "iPhone 16", status: "online", lastSeen: "2026-03-26 02:30 PM" },
            { id: "d2", name: "CEO's iPad", os: "iPadOS", model: "iPad Pro 13-inch", status: "offline", lastSeen: "2026-03-25 10:15 AM" }
        ]
    },
    {
        id: "g2",
        name: "Development Team",
        description: "All devices used by developers",
        deviceCount: 3,
        createdDate: "2026-02-10",
        devices: [
            { id: "d3", name: "Dev MacBook Pro", os: "macOS", model: "MacBook Pro 14-inch", status: "online", lastSeen: "2026-03-26 09:15 AM" },
            { id: "d4", name: "Test iPhone 15", os: "iOS", model: "iPhone 15", status: "online", lastSeen: "2026-03-26 11:45 AM" },
            { id: "d5", name: "Test iPad", os: "iPadOS", model: "iPad Air", status: "online", lastSeen: "2026-03-26 08:30 AM" }
        ]
    },
    {
        id: "g3",
        name: "Marketing Department",
        description: "Devices for marketing staff",
        deviceCount: 1,
        createdDate: "2026-03-05",
        devices: [
            { id: "d6", name: "Marketing iPad", os: "iPadOS", model: "iPad Pro 11-inch", status: "offline", lastSeen: "2026-03-25 04:12 PM" }
        ]
    }
];

export function DeviceGroupsList() {
    const [selectedGroup, setSelectedGroup] = useState<DeviceGroup | null>(null);
    const [isDrawerVisible, setIsDrawerVisible] = useState(false);
    const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [form] = Form.useForm();

    const handleGroupClick = (group: DeviceGroup) => {
        setSelectedGroup(group);
        setIsDrawerVisible(true);
    };

    const columns: ColumnsType<DeviceGroup> = [
        {
            title: "GROUP NAME",
            dataIndex: "name",
            key: "name",
            render: (text, record) => (
                <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-indigo-50 flex items-center justify-center border border-indigo-100">
                        <FolderOpen className="w-4 h-4 text-indigo-600" />
                    </div>
                    <div className="flex flex-col">
                        <a 
                            href="#" 
                            className="text-[#de2a15] hover:text-[#c22412] font-medium transition-colors"
                            onClick={(e) => {
                                e.preventDefault();
                                handleGroupClick(record);
                            }}
                        >
                            {text}
                        </a>
                        <span className="text-xs text-slate-500">{record.description}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "TOTAL DEVICES",
            dataIndex: "deviceCount",
            key: "deviceCount",
            render: (count) => (
                <div className="flex items-center gap-2">
                    <Users className="w-4 h-4 text-slate-400" />
                    <span className="font-medium text-slate-700">{count} devices</span>
                </div>
            ),
        },
        {
            title: "CREATED DATE",
            dataIndex: "createdDate",
            key: "createdDate",
            render: (date) => <span className="text-slate-600">{date}</span>,
        },
        {
            title: "ACTIONS",
            key: "actions",
            width: 100,
            render: (_, record) => {
                const actionMenu: MenuProps['items'] = [
                    {
                        key: 'assign-profile',
                        icon: <FilePlus2 className="w-4 h-4" />,
                        label: 'Assign Profile',
                        onClick: () => console.log('Assign profile to group', record.id)
                    },
                    {
                        type: 'divider',
                    },
                    {
                        key: 'delete-group',
                        icon: <Trash2 className="w-4 h-4 text-red-500" />,
                        label: <span className="text-red-500">Delete Group</span>,
                        onClick: () => {
                            Modal.confirm({
                                title: 'Delete Device Group',
                                content: `Are you sure you want to delete the group "${record.name}"? Devices in this group will not be deleted.`,
                                okText: 'Delete',
                                okType: 'danger',
                                cancelText: 'Cancel',
                                onOk: () => console.log('Delete group:', record.id)
                            });
                        }
                    }
                ];

                return (
                    <Dropdown menu={{ items: actionMenu }} trigger={['click']} placement="bottomRight">
                        <Button type="text" icon={<MoreVertical className="w-4 h-4 text-slate-500" />} />
                    </Dropdown>
                );
            },
        }
    ];

    const deviceColumns: ColumnsType<Device> = [
        {
            title: "DEVICE",
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
                        <span className="font-medium text-slate-800">{text}</span>
                        <span className="text-xs text-slate-500">{record.model}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "STATUS",
            dataIndex: "status",
            key: "status",
            render: (status) => (
                <Tag color={status === "online" ? "success" : "default"} className="rounded-full px-2">
                    {status.toUpperCase()}
                </Tag>
            ),
        }
    ];

    const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
        setSelectedRowKeys(newSelectedRowKeys);
    };

    const rowSelection = {
        selectedRowKeys,
        onChange: onSelectChange,
    };

    const actionMenu: MenuProps['items'] = [
        {
            key: 'assign-profile',
            icon: <FilePlus2 className="w-4 h-4" />,
            label: 'Assign Profile to Group',
            onClick: () => console.log('Assign profile', selectedRowKeys)
        }
    ];

    return (
        <div className="flex flex-col h-[calc(100vh-64px)] bg-slate-50 relative border-none overflow-hidden rounded-none shadow-none z-0">
            {/* Top Toolbar */}
            <div className="flex flex-wrap items-center justify-between p-4 gap-4 bg-white border-b border-slate-200 z-10 shadow-sm">
                <div className="flex items-center gap-3">
                    <div className="flex group">
                        <Input
                            placeholder="Search group name..."
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
                
                <Button 
                    type="primary" 
                    icon={<Plus className="w-4 h-4" />}
                    className="bg-[#de2a15] hover:bg-[#c22412] text-white font-medium px-5 h-8 border-none shadow-sm transition-colors rounded-md"
                    onClick={() => setIsCreateModalVisible(true)}
                >
                    CREATE GROUP
                </Button>
            </div>

            {/* Sub Toolbar */}
            <div className="flex items-center justify-between px-4 py-3 bg-slate-50 border-b border-slate-200 z-10">
                <div className="flex items-center gap-4 text-sm text-slate-600">
                    <span className="font-bold text-slate-800 tracking-wide uppercase">
                        DEVICE GROUPS <span className="font-normal text-slate-500">(1 - 3 of 3)</span>
                    </span>
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
                </div>
            </div>

            {/* Table */}
            <div className="flex-1 overflow-auto border-t border-slate-200 z-10 relative scrollbar-hide">
                <Table
                    rowSelection={rowSelection}
                    columns={columns}
                    dataSource={mockGroups}
                    pagination={false}
                    rowKey="id"
                    className="custom-data-table"
                    rowClassName="hover:bg-slate-50 transition-colors cursor-pointer"
                    onRow={(record) => ({
                        onClick: (e) => {
                            // Ngăn không cho sự kiện click lan truyền nếu click vào các phần tử tương tác khác (như Dropdown)
                            const target = e.target as HTMLElement;
                            const isActionArea = target.closest('.ant-dropdown-trigger') || target.closest('.ant-btn');
                            
                            if (!isActionArea) {
                                handleGroupClick(record);
                            }
                        },
                    })}
                />
            </div>

            {/* Group Detail Drawer */}
            <Drawer
                title={
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-lg bg-indigo-50 flex items-center justify-center border border-indigo-100">
                            <FolderOpen className="w-5 h-5 text-indigo-600" />
                        </div>
                        <div>
                            <div className="font-bold text-slate-800 text-lg">{selectedGroup?.name}</div>
                            <div className="text-xs text-slate-500">
                                Created on: {selectedGroup?.createdDate} • {selectedGroup?.deviceCount} devices
                            </div>
                        </div>
                    </div>
                }
                placement="right"
                width={500}
                onClose={() => setIsDrawerVisible(false)}
                open={isDrawerVisible}
                className="custom-drawer"
                styles={{
                    header: { padding: '20px 24px', borderBottom: '1px solid #e2e8f0' },
                    body: { padding: '24px', backgroundColor: '#f8fafc' }
                }}
            >
                {selectedGroup && (
                    <div className="flex flex-col h-full">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide">Devices in Group</h3>
                        </div>
                        
                        <div className="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm flex-1">
                            <Table
                                columns={deviceColumns}
                                dataSource={selectedGroup.devices}
                                pagination={false}
                                rowKey="id"
                                size="small"
                                className="border-none"
                            />
                        </div>
                    </div>
                )}
            </Drawer>

            {/* Create Group Modal */}
            <Modal
                title="Create Device Group"
                open={isCreateModalVisible}
                onOk={() => {
                    form.validateFields().then(values => {
                        console.log('Create group:', values);
                        setIsCreateModalVisible(false);
                        form.resetFields();
                    });
                }}
                onCancel={() => {
                    setIsCreateModalVisible(false);
                    form.resetFields();
                }}
                okText="Create Group"
                okButtonProps={{ className: "bg-[#de2a15] hover:bg-[#c22412]" }}
            >
                <div className="py-4">
                    <Form form={form} layout="vertical">
                        <Form.Item 
                            name="name" 
                            label="Group Name" 
                            rules={[{ required: true, message: 'Please enter a group name' }]}
                        >
                            <Input placeholder="e.g. Marketing Department" />
                        </Form.Item>
                        <Form.Item 
                            name="description" 
                            label="Description"
                        >
                            <Input.TextArea placeholder="Enter group description..." rows={4} />
                        </Form.Item>
                    </Form>
                </div>
            </Modal>

            <style jsx global>{`
                .scrollbar-hide::-webkit-scrollbar {
                    display: none;
                }
                .scrollbar-hide {
                    -ms-overflow-style: none;
                    scrollbar-width: none;
                }
                
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
            `}</style>
        </div>
    );
}
