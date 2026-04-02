"use client";

import React, { useState, useEffect, useCallback } from "react";
import { Table, Input, Button, Tag, Drawer, Select, Dropdown, MenuProps, Modal, Form, App, Tabs } from "antd";
import { 
    Search, 
    Plus, 
    Smartphone, 
    Monitor,
    Users,
    ChevronDown,
    MoreVertical,
    FolderOpen,
    FilePlus2,
    Trash2
} from "lucide-react";
import type { ColumnsType } from "antd/es/table";
import { deviceGroupService } from "@/services/device-group.service";
import { deviceService } from "@/services/device.service";
import { profileService } from "@/services/profile.service";
import { DeviceGroupResponse } from "@/types/device-group.type";
import { DeviceResponse } from "@/types/device.type";
import { ProfileResponse, ProfileAssignmentResponse } from "@/types/profile.type";
import dayjs from "dayjs";

interface DeviceProfileView {
    id: number;
    name: string;
    status: string;
    configurations: string[];
}

export function DeviceGroupsList() {
    const { message: antdMessage, modal } = App.useApp();
    const [groups, setGroups] = useState<DeviceGroupResponse[]>([]);
    const [loading, setLoading] = useState(false);
    const [selectedGroup, setSelectedGroup] = useState<DeviceGroupResponse | null>(null);
    const [profiles, setProfiles] = useState<ProfileResponse[]>([]);
    const [isDrawerVisible, setIsDrawerVisible] = useState(false);
    const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
    const [isAssignProfileModalVisible, setIsAssignProfileModalVisible] = useState(false);
    const [selectedProfileId, setSelectedProfileId] = useState<number | null>(null);
    const [targetGroupIds, setTargetGroupIds] = useState<number[]>([]);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [groupAssignedProfiles, setGroupAssignedProfiles] = useState<ProfileResponse[]>([]);
    const [isDeviceDetailModalVisible, setIsDeviceDetailModalVisible] = useState(false);
    const [selectedDevice, setSelectedDevice] = useState<DeviceResponse | null>(null);
    const [selectedDeviceProfiles, setSelectedDeviceProfiles] = useState<DeviceProfileView[]>([]);
    const [loadingDeviceDetail, setLoadingDeviceDetail] = useState(false);
    const [form] = Form.useForm();

    const fetchGroups = useCallback(async () => {
        try {
            setLoading(true);
            const response = await deviceGroupService.getGroups({ search: searchQuery });
            if (response.is_success) {
                setGroups(response.data.items || []);
            } else {
                antdMessage.error(response.message || "Failed to fetch device groups");
            }
        } catch (error) {
            console.error("Fetch groups error:", error);
            antdMessage.error("An error occurred while fetching device groups");
        } finally {
            setLoading(false);
        }
    }, [searchQuery, antdMessage]);

    useEffect(() => {
        fetchGroups();
    }, [fetchGroups]);

    const fetchProfiles = async () => {
        try {
            const response = await profileService.getProfiles({ limit: 100, status: "active" });
            if (response.is_success && response.data) {
                setProfiles(response.data.items || []);
            }
        } catch (error) {
            console.error("Fetch profiles error:", error);
        }
    };

    useEffect(() => {
        fetchProfiles();
    }, []);

    const openAssignProfileModal = (groupIds: number[]) => {
        setTargetGroupIds(groupIds);
        setSelectedProfileId(null);
        fetchProfiles();
        setIsAssignProfileModalVisible(true);
    };

    const getProfileConfigurations = (profile: ProfileResponse): string[] => {
        const configurations: string[] = [];
        if (profile.network_config && Object.keys(profile.network_config).length > 0) configurations.push("Wi-Fi / Network");
        if (profile.restrictions && Object.keys(profile.restrictions).length > 0) configurations.push("Restrictions");
        if (profile.security_settings && Object.keys(profile.security_settings).length > 0) configurations.push("Passcode / Security");
        if (profile.content_filter && Object.keys(profile.content_filter).length > 0) configurations.push("Content Filter");
        if (profile.payloads && Object.keys(profile.payloads).length > 0) configurations.push("Custom Payloads");
        if (profile.compliance_rules && Object.keys(profile.compliance_rules).length > 0) configurations.push("Compliance Rules");
        if (configurations.length === 0) configurations.push("General");
        return configurations;
    };

    const getGroupAssignedProfiles = async (groupId: number): Promise<ProfileResponse[]> => {
        const profilesRes = await profileService.getProfiles({ limit: 100, status: "active" });
        if (!profilesRes.is_success || !profilesRes.data) {
            return [];
        }

        const candidates = profilesRes.data.items || [];
        const assignmentChecks = await Promise.all(
            candidates.map(async (profile) => {
                try {
                    const assignmentRes = await profileService.getProfileAssignments(profile.id);
                    if (!assignmentRes.is_success || !assignmentRes.data) {
                        return null;
                    }

                    const assignedToGroup = assignmentRes.data.some(
                        (assignment) => assignment.target_type === "group" && assignment.group_id === groupId
                    );

                    return assignedToGroup ? profile : null;
                } catch {
                    return null;
                }
            })
        );

        return assignmentChecks.filter((item): item is ProfileResponse => item !== null);
    };

    const getProfilesForDeviceInGroup = async (deviceId: string, groupId: number): Promise<DeviceProfileView[]> => {
        const profilesRes = await profileService.getProfiles({ limit: 100, status: "active" });
        if (!profilesRes.is_success || !profilesRes.data) {
            return [];
        }

        const candidates = profilesRes.data.items || [];
        const checks = await Promise.all(
            candidates.map(async (profile) => {
                try {
                    const assignmentRes = await profileService.getProfileAssignments(profile.id);
                    if (!assignmentRes.is_success || !assignmentRes.data) {
                        return null;
                    }

                    const assignedToDevice = assignmentRes.data.some(
                        (assignment) => assignment.target_type === "device" && assignment.device_id === deviceId
                    );
                    const assignedToGroup = assignmentRes.data.some(
                        (assignment) => assignment.target_type === "group" && assignment.group_id === groupId
                    );

                    if (!assignedToDevice && !assignedToGroup) {
                        return null;
                    }

                    return {
                        id: profile.id,
                        name: profile.name,
                        status: profile.status || "active",
                        configurations: getProfileConfigurations(profile),
                    } as DeviceProfileView;
                } catch {
                    return null;
                }
            })
        );

        return checks.filter((item): item is DeviceProfileView => item !== null);
    };

    const handleAssignProfileToGroups = async () => {
        if (!selectedProfileId) {
            antdMessage.warning("Please select a profile");
            return;
        }

        if (targetGroupIds.length === 0) {
            antdMessage.warning("Please select at least one group");
            return;
        }

        try {
            const assignmentRes = await profileService.getProfileAssignments(selectedProfileId);
            const existingAssignments = assignmentRes.is_success && assignmentRes.data ? assignmentRes.data : [];
            const existingGroupIds = new Set(
                existingAssignments.filter((item) => item.target_type === "group" && item.group_id).map((item) => item.group_id as number)
            );
            const existingDeviceIds = new Set(
                existingAssignments.filter((item) => item.target_type === "device" && item.device_id).map((item) => item.device_id as string)
            );

            const affectedDevices = new Set<string>();

            for (const groupId of targetGroupIds) {
                if (!existingGroupIds.has(groupId)) {
                    await profileService.assignProfile(selectedProfileId, {
                        target_type: "group",
                        group_id: groupId,
                        schedule_type: "immediate",
                    });
                    existingGroupIds.add(groupId);
                }

                const groupDetailRes = await deviceGroupService.getGroupById(groupId);
                const devices = groupDetailRes.is_success && groupDetailRes.data?.devices ? groupDetailRes.data.devices : [];
                devices.forEach((device) => {
                    if (device.id) {
                        affectedDevices.add(device.id);
                    }
                });
            }

            for (const deviceId of Array.from(affectedDevices)) {
                if (existingDeviceIds.has(deviceId)) {
                    continue;
                }
                await profileService.assignProfile(selectedProfileId, {
                    target_type: "device",
                    device_id: deviceId,
                    schedule_type: "immediate",
                });
                existingDeviceIds.add(deviceId);
            }

            antdMessage.success("Đã gán profile cho group và toàn bộ thiết bị trong group.");
            setIsAssignProfileModalVisible(false);
            setSelectedProfileId(null);
            setTargetGroupIds([]);
            setSelectedRowKeys([]);
        } catch (error) {
            console.error("Assign profile error:", error);
            antdMessage.error("Failed to assign profile to group");
        }
    };

    const handleGroupClick = async (group: DeviceGroupResponse) => {
        try {
            // Fetch detail to get devices in group
            const [groupRes, assignedProfiles] = await Promise.all([
                deviceGroupService.getGroupById(group.id),
                getGroupAssignedProfiles(group.id),
            ]);
            if (groupRes.is_success) {
                setSelectedGroup(groupRes.data);
                setGroupAssignedProfiles(assignedProfiles);
                setIsDrawerVisible(true);
            } else {
                antdMessage.error(groupRes.message || "Failed to fetch group details");
            }
        } catch (error) {
            console.error("Fetch group detail error:", error);
            antdMessage.error("An error occurred while fetching group details");
        }
    };

    const handleGroupDeviceClick = async (device: DeviceResponse) => {
        if (!selectedGroup?.id || !device.id) {
            return;
        }

        setLoadingDeviceDetail(true);
        setSelectedDeviceProfiles([]);
        setIsDeviceDetailModalVisible(true);
        try {
            const [deviceRes, profilesInScope] = await Promise.all([
                deviceService.getDeviceById(device.id),
                getProfilesForDeviceInGroup(device.id, selectedGroup.id),
            ]);

            if (deviceRes.is_success && deviceRes.data) {
                setSelectedDevice(deviceRes.data);
            } else {
                setSelectedDevice(device);
            }
            setSelectedDeviceProfiles(profilesInScope);
        } catch (error) {
            console.error("Fetch device detail from group error:", error);
            antdMessage.error("Failed to load device details");
            setSelectedDevice(device);
        } finally {
            setLoadingDeviceDetail(false);
        }
    };

    const handleDeleteGroup = async (groupId: number) => {
        try {
            console.log("Calling delete API for group ID:", groupId);
            const response = await deviceGroupService.deleteGroup(groupId);
            
            // Log response for debugging
            console.log("Delete response:", response);

            // Delete group API does not always return standard response structure
            // If the call succeeds without throwing, we can assume it worked.
            antdMessage.success("Group deleted successfully");
            fetchGroups();
            return true;
        } catch (error: any) {
            console.error("Delete group error details:", error);
            // Some backends return success responses inside catch if status is 200 but parsing fails
            if (error?.response?.status === 200 || error?.response?.status === 204) {
                antdMessage.success("Group deleted successfully");
                fetchGroups();
                return true;
            } else {
                antdMessage.error(error?.response?.data?.message || "Failed to delete group");
                return false;
            }
        }
    };

    const handleCreateGroup = async (values: { name: string; description?: string }) => {
        try {
            const response = await deviceGroupService.createGroup(values);
            if (response.is_success) {
                antdMessage.success("Group created successfully");
                setIsCreateModalVisible(false);
                form.resetFields();
                fetchGroups();
            } else {
                antdMessage.error(response.message || "Failed to create group");
            }
        } catch (error) {
            console.error("Create group error:", error);
            antdMessage.error("An error occurred while creating group");
        }
    };

    const columns: ColumnsType<DeviceGroupResponse> = [
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
            dataIndex: "device_count",
            key: "device_count",
            render: (count) => (
                <div className="flex items-center gap-2">
                    <Users className="w-4 h-4 text-slate-400" />
                    <span className="font-medium text-slate-700">{count || 0} devices</span>
                </div>
            ),
        },
        {
            title: "CREATED DATE",
            dataIndex: "created_at",
            key: "created_at",
            render: (date) => <span className="text-slate-600">{date ? dayjs(date).format("YYYY-MM-DD HH:mm") : "-"}</span>,
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
                        onClick: (e) => {
                            e.domEvent.stopPropagation();
                            openAssignProfileModal([record.id]);
                        }
                    },
                    {
                        type: 'divider',
                    },
                    {
                        key: 'delete-group',
                        icon: <Trash2 className="w-4 h-4 text-red-500" />,
                        label: <span className="text-red-500">Delete Group</span>,
                        onClick: (e) => {
                            e.domEvent.stopPropagation();
                            
                            // Sử dụng setTimeout để đảm bảo Dropdown đóng trước khi hiện Modal
                            // và tránh xung đột event loop
                            setTimeout(() => {
                                modal.confirm({
                                    title: 'Delete Device Group',
                                    content: `Are you sure you want to delete the group "${record.name}"? Devices in this group will not be deleted.`,
                                    okText: 'Delete',
                                    okType: 'danger',
                                    cancelText: 'Cancel',
                                    onOk: () => {
                                        return new Promise<void>((resolve) => {
                                            handleDeleteGroup(record.id).then(() => {
                                                resolve();
                                            }).catch(() => {
                                                resolve();
                                            });
                                        });
                                    }
                                });
                            }, 10);
                        }
                    }
                ];

                return (
                    <Dropdown menu={{ items: actionMenu }} trigger={['click']} placement="bottomRight">
                        <Button 
                            type="text" 
                            icon={<MoreVertical className="w-4 h-4 text-slate-500" />} 
                            onClick={(e) => e.stopPropagation()}
                        />
                    </Dropdown>
                );
            },
        }
    ];

    const deviceColumns: ColumnsType<DeviceResponse> = [
        {
            title: "DEVICE",
            dataIndex: "name",
            key: "name",
            render: (text, record) => (
                <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                        {record.platform?.toLowerCase() === "macos" ? (
                            <Monitor className="w-4 h-4 text-slate-600" />
                        ) : (
                            <Smartphone className="w-4 h-4 text-slate-600" />
                        )}
                    </div>
                    <div className="flex flex-col">
                        <button
                            type="button"
                            className="text-left font-medium text-slate-800 hover:text-[#de2a15] transition-colors"
                            onClick={(e) => {
                                e.stopPropagation();
                                handleGroupDeviceClick(record);
                            }}
                        >
                            {text || record.model || "Unknown Device"}
                        </button>
                        <span className="text-xs text-slate-500">{record.model || "Unknown Model"}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "STATUS",
            dataIndex: "status",
            key: "status",
            render: (status) => (
                <Tag color={status?.toLowerCase() === "active" ? "success" : "default"} className="rounded-full px-2">
                    {status?.toUpperCase() || "UNKNOWN"}
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
            onClick: () => openAssignProfileModal((selectedRowKeys as number[]).map((key) => Number(key)))
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
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            onPressEnter={() => fetchGroups()}
                            prefix={<Search className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" />}
                            className="w-64 h-8 rounded-r-none border-r-0 hover:border-[#de2a15] focus:border-[#de2a15] focus:shadow-none transition-colors"
                        />
                        <Button 
                            type="primary" 
                            onClick={() => fetchGroups()}
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
                        DEVICE GROUPS <span className="font-normal text-slate-500">(1 - {groups.length} of {groups.length})</span>
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
                    dataSource={groups}
                    loading={loading}
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
                                Created on: {selectedGroup?.created_at ? dayjs(selectedGroup.created_at).format("YYYY-MM-DD") : "-"} • {selectedGroup?.device_count || 0} devices
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
                        <div className="mb-4">
                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-2">Profiles Assigned to Group</h3>
                            <div className="bg-white rounded-xl border border-slate-200 p-3 shadow-sm">
                                {groupAssignedProfiles.length > 0 ? (
                                    <div className="flex flex-wrap gap-2">
                                        {groupAssignedProfiles.map((profile) => (
                                            <Tag key={profile.id} className="rounded-full px-3 py-1 m-0 bg-indigo-50 text-indigo-700 border-indigo-200">
                                                {profile.name}
                                            </Tag>
                                        ))}
                                    </div>
                                ) : (
                                    <div className="text-sm text-slate-500">No profile is assigned directly to this group.</div>
                                )}
                            </div>
                        </div>

                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide">Devices in Group</h3>
                            <span className="text-xs text-slate-500">Click a device to view details and profile configurations</span>
                        </div>
                        
                        <div className="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm flex-1">
                            <Table
                                columns={deviceColumns}
                                dataSource={selectedGroup.devices}
                                pagination={false}
                                rowKey="id"
                                size="small"
                                className="border-none"
                                rowClassName="cursor-pointer hover:bg-slate-50"
                                onRow={(record) => ({
                                    onClick: () => handleGroupDeviceClick(record),
                                })}
                            />
                        </div>
                    </div>
                )}
            </Drawer>

            <Modal
                title={selectedDevice?.name || selectedDevice?.model || "Device Detail"}
                open={isDeviceDetailModalVisible}
                onCancel={() => {
                    setIsDeviceDetailModalVisible(false);
                    setSelectedDevice(null);
                    setSelectedDeviceProfiles([]);
                }}
                footer={null}
                width={760}
            >
                <Tabs
                    defaultActiveKey="info"
                    items={[
                        {
                            key: "info",
                            label: "Device Info",
                            children: (
                                <div className="grid grid-cols-2 gap-y-4 gap-x-6 py-2">
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">Model</div>
                                        <div className="text-sm text-slate-800 font-medium">{selectedDevice?.model || "-"}</div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">Platform</div>
                                        <div className="text-sm text-slate-800 font-medium">{selectedDevice?.platform || "-"}</div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">OS Version</div>
                                        <div className="text-sm text-slate-800 font-medium">{selectedDevice?.os_version || "-"}</div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">Status</div>
                                        <div className="text-sm text-slate-800 font-medium">{selectedDevice?.status || "-"}</div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">Serial Number</div>
                                        <div className="text-sm text-slate-800 font-mono">{selectedDevice?.serial_number || "-"}</div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-slate-500 font-medium mb-1">UDID</div>
                                        <div className="text-sm text-slate-800 font-mono break-all">{selectedDevice?.udid || "-"}</div>
                                    </div>
                                </div>
                            ),
                        },
                        {
                            key: "profiles",
                            label: "Profiles Configuration",
                            children: (
                                <div className="py-2">
                                    {loadingDeviceDetail ? (
                                        <div className="text-sm text-slate-500">Loading profile configurations...</div>
                                    ) : selectedDeviceProfiles.length > 0 ? (
                                        <div className="space-y-3">
                                            {selectedDeviceProfiles.map((profile) => (
                                                <div key={profile.id} className="rounded-lg border border-slate-200 p-3 bg-slate-50">
                                                    <div className="flex items-center justify-between">
                                                        <div className="text-sm font-semibold text-slate-800">{profile.name}</div>
                                                        <Tag color={profile.status === "active" ? "success" : "default"} className="rounded-full px-2">
                                                            {profile.status.toUpperCase()}
                                                        </Tag>
                                                    </div>
                                                    <div className="text-xs text-slate-600 mt-2">{profile.configurations.join(" • ")}</div>
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <div className="text-sm text-slate-500">No profile configuration recorded for this device.</div>
                                    )}
                                </div>
                            ),
                        },
                    ]}
                />
            </Modal>

            {/* Create Group Modal */}
            <Modal
                title="Create Device Group"
                open={isCreateModalVisible}
                onOk={() => {
                    form.validateFields().then(values => {
                        handleCreateGroup(values);
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

            <Modal
                title="Assign Profile to Group"
                open={isAssignProfileModalVisible}
                onOk={handleAssignProfileToGroups}
                onCancel={() => {
                    setIsAssignProfileModalVisible(false);
                    setSelectedProfileId(null);
                    setTargetGroupIds([]);
                }}
                okText="Assign Profile"
                okButtonProps={{ className: "bg-[#de2a15] hover:bg-[#c22412]", disabled: !selectedProfileId }}
            >
                <div className="py-4">
                    <p className="mb-4 text-slate-600">
                        Select a profile to assign to {targetGroupIds.length} group(s). Devices in selected groups will apply this common profile.
                    </p>
                    <Select
                        className="w-full"
                        placeholder="Select a profile"
                        value={selectedProfileId ?? undefined}
                        onChange={(value) => setSelectedProfileId(value)}
                        options={profiles.map((profile) => ({ value: profile.id, label: profile.name }))}
                    />
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
