"use client";

import React, { useState, useEffect, useCallback } from "react";
import { Table, Select, Input, Button, Tag, Modal, Tabs, Dropdown, App } from "antd";
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
    MoreVertical,
    Lock,
    LockOpen,
    Trash2,
    Power,
    PowerOff,
    AlertCircle,
    MessageSquare,
    Contact
} from "lucide-react";
import type { ColumnsType } from "antd/es/table";
import { deviceService } from "@/services/device.service";
import { deviceGroupService } from "@/services/device-group.service";
import { profileService } from "@/services/profile.service";
import { DeviceResponse, DeviceWipeRequest } from "@/types/device.type";
import { DeviceGroupResponse } from "@/types/device-group.type";
import { ProfileResponse } from "@/types/profile.type";
import { useDebounce } from "@/hooks";

interface DeviceProfileView {
    id: number;
    name: string;
    status: string;
    configurations: string[];
}


export function DevicesList() {
    const { message: antdMessage, modal } = App.useApp();
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [selectedDevice, setSelectedDevice] = useState<DeviceResponse | null>(null);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isGroupModalVisible, setIsGroupModalVisible] = useState(false);
    const [isProfileModalVisible, setIsProfileModalVisible] = useState(false);

    // Data states
    const [devices, setDevices] = useState<DeviceResponse[]>([]);
    const [groups, setGroups] = useState<DeviceGroupResponse[]>([]);
    const [profiles, setProfiles] = useState<ProfileResponse[]>([]);
    const [deviceProfiles, setDeviceProfiles] = useState<DeviceProfileView[]>([]);
    const [loading, setLoading] = useState(false);
    const [pagination, setPagination] = useState({ current: 1, pageSize: 50, total: 0 });
    const [selectedGroupToAdd, setSelectedGroupToAdd] = useState<number | null>(null);
    const [selectedProfileToAdd, setSelectedProfileToAdd] = useState<number | null>(null);

    // Filters
    const [search, setSearch] = useState("");
    const [statusFilter, setStatusFilter] = useState("all");
    const [platformFilter, setPlatformFilter] = useState("all");
    const debouncedSearch = useDebounce(search, 320);

    const fetchDevices = useCallback(async () => {
        setLoading(true);
        try {
            const params: any = {
                page: pagination.current,
                limit: pagination.pageSize,
            };
            if (debouncedSearch) params.search = debouncedSearch;
            if (statusFilter !== "all") params.status = statusFilter;
            if (platformFilter !== "all") params.platform = platformFilter;

            const res = await deviceService.getDevices(params);
            if (res.is_success && res.data) {
                setDevices(res.data.items || []);
                setPagination(prev => ({ ...prev, total: res.data?.pagination?.total || 0 }));
            }
        } catch (error) {
            console.error("Failed to fetch devices", error);
            antdMessage.error("Failed to fetch devices");
        } finally {
            setLoading(false);
        }
    }, [pagination.current, pagination.pageSize, debouncedSearch, statusFilter, platformFilter, antdMessage]);

    const fetchGroups = async () => {
        try {
            const res = await deviceGroupService.getGroups({ limit: 100 });
            if (res.is_success && res.data) {
                setGroups(res.data.items || []);
            }
        } catch (error) {
            console.error("Failed to fetch groups", error);
        }
    };

    const fetchProfiles = async () => {
        try {
            const res = await profileService.getProfiles({ limit: 100, status: "active" });
            if (res.is_success && res.data) {
                setProfiles(res.data.items || []);
            }
        } catch (error) {
            console.error("Failed to fetch profiles", error);
        }
    };

    useEffect(() => {
        fetchDevices();
    }, [fetchDevices]);

    useEffect(() => {
        fetchProfiles();
    }, []);

    const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
        setSelectedRowKeys(newSelectedRowKeys);
    };

    const handleDeviceClick = async (record: DeviceResponse) => {
        setSelectedDevice(record);
        setDeviceProfiles([]);
        setIsModalVisible(true);
        try {
            const res = await deviceService.getDeviceById(record.id);
            if (res.is_success && res.data) {
                setSelectedDevice(res.data);
            }
            await fetchDeviceProfiles(record.id);
        } catch (error) {
            console.error("Failed to fetch device detail", error);
            antdMessage.error("Failed to fetch device detail");
        }
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

    const fetchDeviceProfiles = async (deviceId: string) => {
        try {
            const profilesRes = await profileService.getProfiles({ limit: 100 });
            if (!profilesRes.is_success || !profilesRes.data) {
                setDeviceProfiles([]);
                return;
            }

            const candidates = profilesRes.data.items || [];
            const assignedResults = await Promise.all(
                candidates.map(async (profile) => {
                    try {
                        const assignmentRes = await profileService.getProfileAssignments(profile.id);
                        if (!assignmentRes.is_success || !assignmentRes.data) {
                            return null;
                        }

                        const assignedToDevice = assignmentRes.data.some(
                            (assignment) => assignment.target_type === "device" && assignment.device_id === deviceId
                        );

                        if (!assignedToDevice) {
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

            setDeviceProfiles(assignedResults.filter((item): item is DeviceProfileView => item !== null));
        } catch (error) {
            console.error("Failed to fetch assigned profiles", error);
            setDeviceProfiles([]);
        }
    };

    const rowSelection = {
        selectedRowKeys,
        onChange: onSelectChange,
    };

    const getApiErrorMessage = (error: unknown, fallback: string) => {
        if (error && typeof error === "object" && "message" in error) {
            const message = (error as { message?: string }).message;
            if (typeof message === "string" && message.trim()) {
                return message;
            }
        }
        return fallback;
    };

    const executeDeviceAction = async (
        action: () => Promise<unknown>,
        successMessage: string,
        fallbackErrorMessage: string
    ) => {
        try {
            await action();
            antdMessage.success(successMessage);
        } catch (error) {
            antdMessage.error(getApiErrorMessage(error, fallbackErrorMessage));
        }
    };

    const actionButtonBaseClass = "h-9 w-full px-2 text-xs font-medium whitespace-nowrap justify-center bg-white border-slate-200 text-slate-700 hover:bg-slate-50 hover:text-primary-600 hover:border-primary-300 transition-all duration-200 shadow-sm rounded-lg";

    const actionMenu: MenuProps['items'] = [
        {
            key: 'add-to-group',
            icon: <FolderPlus className="w-4 h-4" />,
            label: 'Add to Group',
            onClick: () => {
                fetchGroups();
                setIsGroupModalVisible(true);
            }
        },
        {
            key: 'assign-profile',
            icon: <FilePlus2 className="w-4 h-4" />,
            label: 'Assign Profile',
            onClick: () => {
                setSelectedProfileToAdd(null);
                fetchProfiles();
                setIsProfileModalVisible(true);
            }
        }
    ];

    const columns: ColumnsType<DeviceResponse> = [
        {
            title: "DEVICE NAME",
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
                        <a 
                            href="#" 
                            className="text-primary-600 hover:text-primary-700 font-medium transition-colors"
                            onClick={(e) => {
                                e.preventDefault();
                                handleDeviceClick(record);
                            }}
                        >
                            {text || record.model || "Unknown Device"}
                        </a>
                        <span className="text-xs text-slate-500">{record.model}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "OS",
            dataIndex: "platform",
            key: "platform",
            render: (text) => (
                <div className="flex items-center gap-2 text-slate-700">
                    <Apple className="w-4 h-4" fill="currentColor" />
                    <span className="font-medium">{text || "iOS"}</span>
                </div>
            ),
        },
        {
            title: "VERSION",
            dataIndex: "os_version",
            key: "os_version",
            render: (text) => <span className="text-slate-700 font-mono text-sm">{text}</span>,
        },
        {
            title: "BATTERY",
            dataIndex: "battery_level",
            key: "battery_level",
            render: (percent) => {
                if (percent == null) return <span className="text-slate-400">-</span>;
                return (
                    <div className="flex items-center gap-2">
                        <Battery className={`w-4 h-4 ${percent < 20 ? 'text-red-500' : percent < 50 ? 'text-orange-500' : 'text-emerald-500'}`} />
                        <span className="text-slate-700">{Math.round(percent * 1)}%</span>
                    </div>
                );
            }
        },
        {
            title: "STORAGE",
            dataIndex: "storage_capacity",
            key: "storage_capacity",
            render: (_, record) => {
                if (!record.storage_capacity) return <span className="text-slate-400">-</span>;
                const totalGB = Math.round(record.storage_capacity / (1024 * 1024 * 1024));
                const usedGB = record.storage_used ? Math.round(record.storage_used / (1024 * 1024 * 1024)) : 0;
                const availableGB = totalGB - usedGB;
                
                return (
                    <div className="flex items-center gap-2">
                        <HardDrive className="w-4 h-4 text-slate-400" />
                        <span className="text-slate-700">{availableGB} GB </span>
                        <span className="text-slate-400 text-xs">/ {totalGB} GB</span>
                    </div>
                );
            }
        },
        {
            title: "STATUS",
            dataIndex: "status",
            key: "status",
            render: (status) => (
                <Tag color={status?.toLowerCase() === "active" ? "success" : "default"} className="rounded-full px-3">
                    {status?.toUpperCase() || "UNKNOWN"}
                </Tag>
            ),
        },
        {
            title: "LAST SEEN",
            dataIndex: "last_seen",
            key: "last_seen",
            render: (text) => {
                if (!text) return <span className="text-slate-400">-</span>;
                return <span className="text-slate-500 text-sm">{new Date(text).toLocaleString()}</span>;
            }
        },
    ];

    return (
        <div className="motion-safe-fade flex flex-col h-[calc(100vh-64px)] bg-transparent relative border-none overflow-hidden rounded-none shadow-none z-0">
            {/* Top Toolbar */}
            <div className="liquid-glass liquid-layout-toolbar motion-safe-fade flex flex-wrap items-center justify-between p-4 gap-4 bg-white border-b border-slate-200 z-10 shadow-sm">
                <div className="flex items-center gap-6">
                    <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-700">OS:</span>
                        <Select className="cursor-pointer" value={platformFilter}
                            onChange={(val) => {
                                setPlatformFilter(val);
                                setPagination(prev => ({ ...prev, current: 1 }));
                            }}
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
                        <Select className="cursor-pointer" value={statusFilter}
                            onChange={(val) => {
                                setStatusFilter(val);
                                setPagination(prev => ({ ...prev, current: 1 }));
                            }}
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
                            value={search}
                            onChange={(e) => {
                                setSearch(e.target.value);
                                setPagination(prev => ({ ...prev, current: 1 }));
                            }}
                            prefix={<Search className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" />}
                            className="w-64 h-8 rounded-r-none border-r-0 hover:border-primary-600 focus:border-primary-600 focus:shadow-none transition-colors"
                        />
                        <Button 
                            type="primary" 
                            className="liquid-btn bg-primary-600 hover:bg-primary-700 rounded-l-none h-8 w-10 px-0 flex items-center justify-center border-none shadow-sm transition-colors"
                            icon={<Search className="w-4 h-4 text-white" strokeWidth={2.5} />}
                        />
                    </div>
                </div>
            </div>

            {/* Sub Toolbar / Table Controls */}
            <div className="liquid-glass liquid-layout-toolbar motion-safe-fade flex items-center justify-between px-4 py-3 bg-slate-50 border-b border-slate-200 z-10">
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
                        <span className="text-primary-600 font-bold">1 of 1</span>
                        <Button type="text" size="small" disabled className="text-slate-400">&rarr;</Button>
                    </div>

                    <div className="border-l border-slate-300 pl-4">
                        <Button type="text" size="small" icon={<RefreshCcw className="w-4 h-4 text-slate-500 hover:text-slate-800" />} />
                    </div>
                </div>

                <div className="flex items-center gap-4 text-sm">
                    {selectedRowKeys.length > 0 && (
                        <div className="flex items-center gap-2 mr-4 border-r border-slate-300 pr-4">
                            <span className="text-primary-600 font-medium">{selectedRowKeys.length} selected</span>
                            <Dropdown menu={{ items: actionMenu }} trigger={['click']} placement="bottomRight">
                                <Button size="small" type="primary" className="liquid-btn bg-primary-600 hover:bg-primary-700 flex items-center gap-1">
                                    Actions <ChevronDown className="w-3 h-3" />
                                </Button>
                            </Dropdown>
                        </div>
                    )}
                    <span className="flex items-center gap-1 cursor-pointer text-slate-700 hover:text-slate-900 font-medium">
                        Columns (7) <ChevronDown className="w-4 h-4" />
                    </span>
                    <Button type="text" icon={<Filter className="w-4 h-4 text-primary-600" />} className="text-primary-600 bg-primary-50 hover:bg-primary-100 rounded-full w-8 h-8 flex items-center justify-center p-0 border border-primary-200 transition-colors" />
                </div>
            </div>

            {/* Table */}
            <div className="liquid-glass liquid-layout-content motion-safe-fade flex-1 overflow-auto border-t border-slate-200 z-10 relative scrollbar-hide">
                <Table
                    rowSelection={rowSelection}
                    columns={columns}
                    dataSource={devices}
                    rowKey="id"
                    loading={loading}
                    pagination={false}
                    className="custom-data-table"
                    rowClassName="motion-safe-lift hover:bg-slate-50 transition-colors cursor-pointer"
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
                            {selectedDevice?.platform?.toLowerCase() === "macos" ? (
                                <Monitor className="w-5 h-5 text-slate-700" />
                            ) : (
                                <Smartphone className="w-5 h-5 text-slate-700" />
                            )}
                        </div>
                        <div>
                            <div className="font-bold text-slate-800 text-lg">{selectedDevice?.name || selectedDevice?.model || "Unknown Device"}</div>
                            <div className="text-xs text-slate-500 flex items-center gap-2">
                                <span className={selectedDevice?.status?.toLowerCase() === "active" ? "text-emerald-500 font-medium" : "text-slate-500"}>
                                    ● {selectedDevice?.status?.toUpperCase() || "UNKNOWN"}
                                </span>
                                <span>|</span>
                                <span>Last seen: {selectedDevice?.last_seen ? new Date(selectedDevice.last_seen).toLocaleString() : "-"}</span>
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
                                                <AppWindow className="w-4 h-4 text-slate-600" /> Device Actions
                                            </h3>
                                            <div className="grid grid-cols-5 gap-2">
                                                <Button 
                                                    icon={<Lock className="w-4 h-4" />}
                                                    className={actionButtonBaseClass}
                                                    onClick={() => {
                                                        let lockMessage = "Thiết bị bị khóa bởi Quản trị viên";
                                                        let lockPhoneNumber = "";
                                                        let lockFootnote = "";
                                                        modal.confirm({
                                                            title: (
                                                                <div className="flex items-center gap-3 pb-3 border-b border-gray-200">
                                                                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-error-500 to-error-600 flex items-center justify-center shadow-md">
                                                                        <Lock className="w-5 h-5 text-white" />
                                                                    </div>
                                                                    <div>
                                                                        <h3 className="text-lg font-semibold text-gray-900 m-0">Khóa Thiết Bị</h3>
                                                                        <p className="text-sm text-gray-500 m-0">Cấu hình hiển thị màn hình khóa</p>
                                                                    </div>
                                                                </div>
                                                            ),
                                                            width: 540,
                                                            content: (
                                                                <div className="space-y-4 pt-2">
                                                                    <div className="bg-warning-50 border border-warning-200 rounded-lg p-3 flex items-start gap-3">
                                                                        <AlertCircle className="w-5 h-5 text-warning-600 flex-shrink-0 mt-0.5" />
                                                                        <div className="flex-1">
                                                                            <p className="text-sm font-medium text-warning-900 m-0">Thiết bị sẽ bị khóa ngay lập tức</p>
                                                                            <p className="text-xs text-warning-700 m-0 mt-1">Người dùng sẽ thấy thông tin này trên màn hình khóa</p>
                                                                        </div>
                                                                    </div>
                                                                    
                                                                    <div className="space-y-1">
                                                                        <label className="text-sm font-semibold text-gray-700 flex items-center gap-2">
                                                                            <MessageSquare className="w-4 h-4 text-gray-500" />
                                                                            Thông Báo Màn Hình Khóa
                                                                        </label>
                                                                        <Input.TextArea
                                                                            rows={3}
                                                                            defaultValue={lockMessage}
                                                                            maxLength={200}
                                                                            placeholder="Nhập thông báo hiển thị trên màn hình khóa"
                                                                            showCount
                                                                            className="rounded-lg"
                                                                            onChange={(e) => {
                                                                                lockMessage = e.target.value;
                                                                            }}
                                                                        />
                                                                        <p className="text-xs text-gray-500 mt-1">Thông báo này sẽ được hiển thị nổi bật trên thiết bị bị khóa</p>
                                                                    </div>
                                                                    
                                                                    <div className="space-y-1">
                                                                        <label className="text-sm font-semibold text-gray-700 flex items-center gap-2">
                                                                            <Contact className="w-4 h-4 text-gray-500" />
                                                                            Số Điện Thoại Liên Hệ
                                                                            <span className="text-xs font-normal text-gray-500">(Tùy chọn)</span>
                                                                        </label>
                                                                        <Input
                                                                            prefix={<Contact className="w-4 h-4 text-gray-400" />}
                                                                            placeholder="Ví dụ: +84 123 456 789"
                                                                            maxLength={30}
                                                                            className="rounded-lg h-10"
                                                                            onChange={(e) => {
                                                                                lockPhoneNumber = e.target.value;
                                                                            }}
                                                                        />
                                                                        <p className="text-xs text-gray-500 mt-1">Hiển thị số điện thoại liên hệ trên màn hình khóa</p>
                                                                    </div>
                                                                    
                                                                    <div className="space-y-1">
                                                                        <label className="text-sm font-semibold text-gray-700 flex items-center gap-2">
                                                                            <MessageSquare className="w-4 h-4 text-gray-500" />
                                                                            Ghi Chú Bổ Sung
                                                                            <span className="text-xs font-normal text-gray-500">(Tùy chọn)</span>
                                                                        </label>
                                                                        <Input
                                                                            prefix={<MessageSquare className="w-4 h-4 text-gray-400" />}
                                                                            placeholder="Ví dụ: Vui lòng trả lại phòng IT"
                                                                            maxLength={100}
                                                                            className="rounded-lg h-10"
                                                                            onChange={(e) => {
                                                                                lockFootnote = e.target.value;
                                                                            }}
                                                                        />
                                                                        <p className="text-xs text-gray-500 mt-1">Thêm hướng dẫn hoặc thông tin bổ sung</p>
                                                                    </div>
                                                                </div>
                                                            ),
                                                            okText: 'Khóa Thiết Bị',
                                                            cancelText: 'Hủy',
                                                            okButtonProps: { 
                                                                danger: true,
                                                                icon: <Lock className="w-4 h-4" />,
                                                                className: 'h-10 px-5 rounded-lg font-semibold'
                                                            },
                                                            cancelButtonProps: {
                                                                className: 'h-10 px-5 rounded-lg font-medium'
                                                            },
                                                            onOk: async () => {
                                                                if(!selectedDevice?.id) return;
                                                                await executeDeviceAction(
                                                                    () => deviceService.lockDevice(selectedDevice.id, {
                                                                        message: lockMessage.trim(),
                                                                        phone_number: lockPhoneNumber.trim() || undefined,
                                                                        footnote: lockFootnote.trim() || undefined,
                                                                    }),
                                                                    "Lệnh khóa đã được gửi thành công",
                                                                    "Không thể gửi lệnh khóa"
                                                                );
                                                            }
                                                        });
                                                    }}
                                                >
                                                    Khóa
                                                </Button>
                                                <Button 
                                                    icon={<LockOpen className="w-4 h-4" />}
                                                    className={actionButtonBaseClass}
                                                    onClick={() => {
                                                        modal.confirm({
                                                            title: 'Unlock Device',
                                                            content: 'This will disable Lost Mode on the device. Continue?',
                                                            okText: 'Unlock',
                                                            onOk: async () => {
                                                                if(!selectedDevice?.id) return;
                                                                await executeDeviceAction(
                                                                    () => deviceService.unlockDevice(selectedDevice.id),
                                                                    "Unlock command sent",
                                                                    "Failed to send unlock command"
                                                                );
                                                            }
                                                        });
                                                    }}
                                                >
                                                    Unlock
                                                </Button>
                                                <Button 
                                                    icon={<Trash2 className="w-4 h-4" />}
                                                    danger
                                                    className={`${actionButtonBaseClass} border-slate-300`}
                                                    onClick={() => {
                                                        modal.confirm({
                                                            title: 'Wipe Device',
                                                            content: 'WARNING: This will factory reset the device. All data will be lost. Are you sure?',
                                                            okText: 'Continue',
                                                            okButtonProps: { danger: true },
                                                            onOk: async () => {
                                                                if(!selectedDevice?.id) return;
                                                                modal.confirm({
                                                                    title: 'Confirm Device Wipe',
                                                                    content: 'This action is irreversible. Confirm to enqueue wipe command.',
                                                                    okText: 'Wipe',
                                                                    okButtonProps: { danger: true },
                                                                    onOk: async () => {
                                                                        const wipePayload: DeviceWipeRequest = {
                                                                            obliteration_behavior: "Always"
                                                                        };
                                                                        await executeDeviceAction(
                                                                            () => deviceService.wipeDevice(selectedDevice.id, wipePayload),
                                                                            "Wipe command sent",
                                                                            "Failed to send wipe command"
                                                                        );
                                                                    }
                                                                });
                                                            }
                                                        });
                                                    }}
                                                >
                                                    Erase
                                                </Button>
                                                <Button 
                                                    icon={<Power className="w-4 h-4" />}
                                                    className={actionButtonBaseClass}
                                                    onClick={async () => {
                                                        if(!selectedDevice?.id) return;
                                                        await executeDeviceAction(
                                                            () => deviceService.restartDevice(selectedDevice.id, { notify_user: true }),
                                                            "Restart command sent",
                                                            "Failed to send restart command"
                                                        );
                                                    }}
                                                >
                                                    Restart
                                                </Button>
                                                <Button 
                                                    icon={<PowerOff className="w-4 h-4" />}
                                                    danger
                                                    className={`${actionButtonBaseClass} border-slate-300`}
                                                    onClick={() => {
                                                        modal.confirm({
                                                            title: 'Shutdown Device',
                                                            content: 'Device will power off remotely. Do you want to continue?',
                                                            okText: 'Continue',
                                                            okButtonProps: { danger: true },
                                                            onOk: async () => {
                                                                if(!selectedDevice?.id) return;
                                                                modal.confirm({
                                                                    title: 'Confirm Shutdown',
                                                                    content: 'Please confirm again to send shutdown command.',
                                                                    okText: 'Shutdown',
                                                                    okButtonProps: { danger: true },
                                                                    onOk: async () => {
                                                                        await executeDeviceAction(
                                                                            () => deviceService.shutdownDevice(selectedDevice.id),
                                                                            "Shutdown command sent",
                                                                            "Failed to send shutdown command"
                                                                        );
                                                                    }
                                                                });
                                                            }
                                                        });
                                                    }}
                                                >
                                                    Shutdown
                                                </Button>
                                            </div>
                                        </div>
                                        
                                        <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
                                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-4 border-b border-slate-100 pb-2 flex items-center gap-2">
                                                <Apple className="w-4 h-4 text-slate-600" fill="currentColor" /> System Details
                                            </h3>
                                            <div className="grid grid-cols-2 gap-y-4 gap-x-6">
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Owner ID</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.owner_id || "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Enrollment Date</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.enrolled_at ? new Date(selectedDevice.enrolled_at).toLocaleString() : "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Model</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.model || "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">OS Version</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.platform} {selectedDevice.os_version}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Serial Number</div>
                                                    <div className="text-sm text-slate-800 font-mono">{selectedDevice.serial_number || "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">UDID</div>
                                                    <div className="text-sm text-slate-800 font-mono text-[11px] truncate" title={selectedDevice.udid}>{selectedDevice.udid || "-"}</div>
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
                                                        <Battery className={`w-4 h-4 ${!selectedDevice.battery_level ? 'text-slate-400' : selectedDevice.battery_level < 0.2 ? 'text-red-500' : selectedDevice.battery_level < 0.5 ? 'text-orange-500' : 'text-emerald-500'}`} />
                                                        <span className="text-sm text-slate-800 font-medium">{selectedDevice.battery_level ? Math.round(selectedDevice.battery_level * 1) + '%' : "-"}</span>
                                                    </div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Storage</div>
                                                    <div className="text-sm text-slate-800 font-medium">
                                                        {selectedDevice.storage_capacity && selectedDevice.storage_used 
                                                            ? `${Math.round((selectedDevice.storage_capacity - selectedDevice.storage_used) / (1024*1024*1024))} GB free of ${Math.round(selectedDevice.storage_capacity / (1024*1024*1024))} GB`
                                                            : "-"}
                                                    </div>
                                                    {selectedDevice.storage_capacity && selectedDevice.storage_used && (
                                                        <div className="w-full bg-slate-100 rounded-full h-1.5 mt-2">
                                                            <div 
                                                                className="bg-blue-500 h-1.5 rounded-full" 
                                                                style={{ width: `${(selectedDevice.storage_used / selectedDevice.storage_capacity) * 100}%` }}
                                                            ></div>
                                                        </div>
                                                    )}
                                                </div>
                                            </div>
                                        </div>
                                        
                                        <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
                                            <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wide mb-4 border-b border-slate-100 pb-2 flex items-center gap-2">
                                                <Wifi className="w-4 h-4 text-slate-600" /> Network & Compliance
                                            </h3>
                                            <div className="grid grid-cols-2 gap-y-4 gap-x-6">
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">IP Address</div>
                                                    <div className="text-sm text-slate-800 font-mono">{selectedDevice.ip_address || "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">MAC Address</div>
                                                    <div className="text-sm text-slate-800 font-mono">{selectedDevice.mac_address || "-"}</div>
                                                </div>
                                                <div>
                                                    <div className="text-xs text-slate-500 font-medium mb-1">Compliance</div>
                                                    <div className="text-sm text-slate-800 font-medium">{selectedDevice.compliance_status || "-"}</div>
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
                                                <Tag className="m-0 bg-blue-50 text-blue-700 border-blue-200 font-medium">{deviceProfiles.length} Profiles</Tag>
                                            </div>
                                            <div className="divide-y divide-slate-100">
                                                {deviceProfiles.length > 0 ? deviceProfiles.map((profile) => (
                                                    <div key={profile.id} className="p-4 flex items-center justify-between hover:bg-slate-50 transition-colors">
                                                        <div className="flex items-center gap-3">
                                                            <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center border border-slate-200">
                                                                <Shield className="w-4 h-4 text-slate-500" />
                                                            </div>
                                                            <div>
                                                                <div className="text-sm font-medium text-slate-800">{profile.name}</div>
                                                                <div className="text-xs text-slate-500 mt-0.5">{profile.configurations.join(" • ")}</div>
                                                            </div>
                                                        </div>
                                                        <Tag color={profile.status === 'active' ? 'success' : profile.status === 'failed' ? 'error' : 'warning'} className="rounded-full px-3">
                                                            {profile.status.toUpperCase()}
                                                        </Tag>
                                                    </div>
                                                )) : (
                                                    <div className="p-6 text-center text-slate-500">
                                                        No profiles installed on this device.
                                                    </div>
                                                )}
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
                onOk={async () => {
                    if (!selectedGroupToAdd) {
                        antdMessage.warning("Please select a group");
                        return;
                    }
                    try {
                        await deviceGroupService.addDevicesToGroup(selectedGroupToAdd, {
                            device_ids: selectedRowKeys as string[]
                        });
                        antdMessage.success("Devices added to group successfully");
                        setIsGroupModalVisible(false);
                        setSelectedRowKeys([]);
                    } catch (error) {
                        antdMessage.error("Failed to add devices to group");
                    }
                }}
                onCancel={() => setIsGroupModalVisible(false)}
                okText="Add to Group"
                okButtonProps={{ className: "bg-primary-600 hover:bg-primary-700", disabled: !selectedGroupToAdd }}
            >
                <div className="py-4">
                    <p className="mb-4 text-slate-600">
                        Select a group to add the {selectedRowKeys.length} selected device(s) to:
                    </p>
                    <Select
                        className="w-full"
                        placeholder="Select a group"
                        onChange={(val) => setSelectedGroupToAdd(val)}
                        options={groups.map(g => ({ value: g.id, label: g.name }))}
                    />
                </div>
            </Modal>

            {/* Assign Profile Modal */}
            <Modal
                title="Assign Configuration Profile"
                open={isProfileModalVisible}
                onOk={async () => {
                    if (!selectedProfileToAdd) {
                        antdMessage.warning("Please select a profile");
                        return;
                    }
                    try {
                        const profileId = Number(selectedProfileToAdd);
                        for (const deviceId of selectedRowKeys) {
                            await profileService.assignProfile(profileId, {
                                target_type: "device",
                                device_id: deviceId as string,
                                schedule_type: "immediate",
                            });
                        }
                        antdMessage.success("Profiles assigned successfully");
                        setIsProfileModalVisible(false);
                        setSelectedRowKeys([]);
                        setSelectedProfileToAdd(null);
                        if (selectedDevice?.id && selectedRowKeys.includes(selectedDevice.id)) {
                            await fetchDeviceProfiles(selectedDevice.id);
                        }
                    } catch (error) {
                        antdMessage.error("Failed to assign profile");
                    }
                }}
                onCancel={() => {
                    setIsProfileModalVisible(false);
                    setSelectedProfileToAdd(null);
                }}
                okText="Assign Profile"
                okButtonProps={{ className: "bg-primary-600 hover:bg-primary-700", disabled: !selectedProfileToAdd }}
            >
                <div className="py-4">
                    <p className="mb-4 text-slate-600">
                        Select a configuration profile to assign to the {selectedRowKeys.length} selected device(s):
                    </p>
                    <Select
                        className="w-full"
                        placeholder="Select a profile"
                        onChange={(val) => setSelectedProfileToAdd(val)}
                        options={profiles.map((profile) => ({ value: profile.id, label: profile.name }))}
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
                    color: #1890ff !important;
                    font-weight: 600 !important;
                }
                .custom-tabs .ant-tabs-ink-bar {
                    background: #1890ff !important;
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
