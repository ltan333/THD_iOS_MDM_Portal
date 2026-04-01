"use client";

import React, { useState, useEffect, useCallback } from "react";
import {
    Table, Select, Input, Button, Modal, Tabs, Form, Switch,
    InputNumber, App, Tag, Tooltip, Dropdown, Popconfirm, Divider
} from "antd";
import type { MenuProps } from "antd";
import type { ColumnsType } from "antd/es/table";
import {
    Search, Plus, Apple, Smartphone, Monitor, Settings2, RefreshCcw,
    ChevronDown, Filter, CheckCircle2, PenSquare, X, MoreVertical,
    Shield, Lock, Wifi, Server, Globe, Trash2, Copy, RotateCcw,
    AlertCircle, Archive, Eye
} from "lucide-react";
import { profileService } from "@/services/profile.service";
import { ProfileResponse, CreateProfileRequest, UpdateProfileRequest } from "@/types/profile.type";

// ─── helpers ────────────────────────────────────────────────────────────────

const PLATFORM_OPTIONS = [
    { value: "all", label: "All Platforms" },
    { value: "ios", label: "iOS / iPadOS" },
    { value: "android", label: "Android" },
    { value: "macos", label: "macOS" },
    { value: "windows", label: "Windows" },
];

const STATUS_OPTIONS = [
    { value: "all", label: "All Statuses" },
    { value: "active", label: "Active" },
    { value: "draft", label: "Draft" },
    { value: "archived", label: "Archived" },
];

const SCOPE_OPTIONS = [
    { value: "device", label: "Device" },
    { value: "user", label: "User" },
    { value: "group", label: "Group" },
];

const WIFI_SECURITY_OPTIONS = [
    { value: "WPA2", label: "WPA2" },
    { value: "WPA3", label: "WPA3" },
    { value: "WPA", label: "WPA" },
    { value: "WEP", label: "WEP" },
    { value: "None", label: "None (Open)" },
];

function platformIcon(platform: string) {
    switch (platform) {
        case "ios": return <Smartphone className="w-4 h-4 shrink-0" />;
        case "android": return <Smartphone className="w-4 h-4 shrink-0" />;
        case "macos": return <Monitor className="w-4 h-4 shrink-0" />;
        case "windows": return <Monitor className="w-4 h-4 shrink-0" />;
        default: return <Settings2 className="w-4 h-4 shrink-0" />;
    }
}

function platformLabel(platform: string) {
    return PLATFORM_OPTIONS.find(p => p.value === platform)?.label ?? platform;
}

function statusTag(status: string) {
    switch (status) {
        case "active": return <Tag color="success" className="font-medium">Active</Tag>;
        case "draft": return <Tag color="default" className="font-medium">Draft</Tag>;
        case "archived": return <Tag color="warning" className="font-medium">Archived</Tag>;
        default: return <Tag>{status}</Tag>;
    }
}

// ─── blocked-websites list sub-component ────────────────────────────────────

function BlockedWebsitesList({
    value = [],
    onChange,
}: {
    value?: string[];
    onChange?: (v: string[]) => void;
}) {
    const [input, setInput] = useState("");
    const add = () => {
        const url = input.trim();
        if (url && !value.includes(url)) {
            onChange?.([...value, url]);
        }
        setInput("");
    };
    const remove = (url: string) => onChange?.(value.filter(u => u !== url));
    return (
        <div className="flex flex-col gap-2">
            <div className="flex gap-2">
                <Input
                    placeholder="https://example.com"
                    value={input}
                    onChange={e => setInput(e.target.value)}
                    onPressEnter={add}
                    className="flex-1"
                />
                <Button onClick={add} type="default">Add</Button>
            </div>
            {value.length > 0 && (
                <div className="flex flex-col gap-1 max-h-40 overflow-y-auto border border-slate-200 rounded-md p-2">
                    {value.map(url => (
                        <div key={url} className="flex items-center justify-between text-sm bg-slate-50 rounded px-2 py-1">
                            <span className="text-slate-700 truncate">{url}</span>
                            <Button
                                type="text" size="small"
                                icon={<X className="w-3 h-3" />}
                                onClick={() => remove(url)}
                                className="text-slate-400 hover:text-red-500 shrink-0"
                            />
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

// ─── profile form ────────────────────────────────────────────────────────────

interface ProfileFormValues {
    name: string;
    platform: string;
    scope: string;
    // security
    passcode_required: boolean;
    min_passcode_length: number;
    max_failed_attempts: number;
    screen_lock_timeout: number;
    // restrictions
    camera_disabled: boolean;
    bluetooth_disabled: boolean;
    airdrop_disabled: boolean;
    app_store_disabled: boolean;
    // wifi
    wifi_enabled: boolean;
    wifi_ssid: string;
    wifi_security_type: string;
    wifi_password: string;
    // vpn
    vpn_enabled: boolean;
    vpn_server: string;
    vpn_username: string;
    vpn_password: string;
    // content filter
    safe_browsing_enabled: boolean;
    blocked_websites: string[];
}

function buildPayload(values: ProfileFormValues): CreateProfileRequest {
    const payload: CreateProfileRequest = {
        name: values.name,
        platform: values.platform,
        scope: values.scope,
    };

    // security_settings
    const sec: Record<string, unknown> = {};
    if (values.passcode_required) sec.passcode_required = true;
    if (values.min_passcode_length > 0) sec.min_passcode_length = values.min_passcode_length;
    if (values.max_failed_attempts > 0) sec.max_failed_attempts = values.max_failed_attempts;
    if (values.screen_lock_timeout > 0) sec.screen_lock_timeout = values.screen_lock_timeout;
    if (Object.keys(sec).length) payload.security_settings = sec;

    // restrictions
    const res: Record<string, unknown> = {};
    if (values.camera_disabled) res.camera_disabled = true;
    if (values.bluetooth_disabled) res.bluetooth_disabled = true;
    if (values.airdrop_disabled) res.airdrop_disabled = true;
    if (values.app_store_disabled) res.app_store_disabled = true;
    if (Object.keys(res).length) payload.restrictions = res;

    // network_config
    const net: Record<string, unknown> = {};
    if (values.wifi_enabled && values.wifi_ssid) {
        net.wifi = {
            ssid: values.wifi_ssid,
            security_type: values.wifi_security_type || "WPA2",
            password: values.wifi_password || "",
        };
    }
    if (values.vpn_enabled && values.vpn_server) {
        net.vpn = {
            enabled: true,
            server: values.vpn_server,
            username: values.vpn_username || "",
            password: values.vpn_password || "",
        };
    }
    if (Object.keys(net).length) payload.network_config = net;

    // content_filter
    const cf: Record<string, unknown> = {};
    if (values.safe_browsing_enabled) cf.safe_browsing_enabled = true;
    if (values.blocked_websites?.length) cf.blocked_websites = values.blocked_websites;
    if (Object.keys(cf).length) payload.content_filter = cf;

    return payload;
}

function profileToFormValues(p: ProfileResponse): ProfileFormValues {
    const sec = (p.security_settings ?? {}) as Record<string, unknown>;
    const res = (p.restrictions ?? {}) as Record<string, unknown>;
    const net = (p.network_config ?? {}) as Record<string, unknown>;
    const wifi = (net.wifi ?? {}) as Record<string, unknown>;
    const vpn = (net.vpn ?? {}) as Record<string, unknown>;
    const cf = (p.content_filter ?? {}) as Record<string, unknown>;
    return {
        name: p.name,
        platform: p.platform || "ios",
        scope: p.scope || "device",
        passcode_required: !!sec.passcode_required,
        min_passcode_length: (sec.min_passcode_length as number) || 0,
        max_failed_attempts: (sec.max_failed_attempts as number) || 0,
        screen_lock_timeout: (sec.screen_lock_timeout as number) || 0,
        camera_disabled: !!res.camera_disabled,
        bluetooth_disabled: !!res.bluetooth_disabled,
        airdrop_disabled: !!res.airdrop_disabled,
        app_store_disabled: !!res.app_store_disabled,
        wifi_enabled: !!wifi.ssid,
        wifi_ssid: (wifi.ssid as string) || "",
        wifi_security_type: (wifi.security_type as string) || "WPA2",
        wifi_password: (wifi.password as string) || "",
        vpn_enabled: !!vpn.enabled,
        vpn_server: (vpn.server as string) || "",
        vpn_username: (vpn.username as string) || "",
        vpn_password: (vpn.password as string) || "",
        safe_browsing_enabled: !!cf.safe_browsing_enabled,
        blocked_websites: (cf.blocked_websites as string[]) || [],
    };
}

// ─── main component ───────────────────────────────────────────────────────────

export function ProfilesList() {
    const { message: antdMessage } = App.useApp();
    const [form] = Form.useForm<ProfileFormValues>();

    // list state
    const [profiles, setProfiles] = useState<ProfileResponse[]>([]);
    const [loading, setLoading] = useState(false);
    const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 });

    // filters
    const [search, setSearch] = useState("");
    const [searchInput, setSearchInput] = useState("");
    const [platformFilter, setPlatformFilter] = useState("all");
    const [statusFilter, setStatusFilter] = useState("all");

    // modal
    const [modalOpen, setModalOpen] = useState(false);
    const [editingProfile, setEditingProfile] = useState<ProfileResponse | null>(null);
    const [submitting, setSubmitting] = useState(false);
    const [wifiEnabled, setWifiEnabled] = useState(false);
    const [vpnEnabled, setVpnEnabled] = useState(false);

    // ── fetch ─────────────────────────────────────────────────────────────────
    const fetchProfiles = useCallback(async (page = pagination.current, pageSize = pagination.pageSize) => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = { page, limit: pageSize };
            if (search) params.search = search;
            if (platformFilter !== "all") params.platform = platformFilter;
            if (statusFilter !== "all") params.status = statusFilter;

            const res = await profileService.getProfiles(params as Parameters<typeof profileService.getProfiles>[0]);
            if (res.is_success && res.data) {
                setProfiles(res.data.items || []);
                setPagination(prev => ({ ...prev, current: page, pageSize, total: res.data?.total ?? 0 }));
            }
        } catch {
            antdMessage.error("Failed to fetch profiles");
        } finally {
            setLoading(false);
        }
    }, [search, platformFilter, statusFilter, pagination.current, pagination.pageSize, antdMessage]);

    useEffect(() => {
        fetchProfiles(1, pagination.pageSize);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [search, platformFilter, statusFilter]);

    // ── actions ───────────────────────────────────────────────────────────────
    const handleRepush = async (profile: ProfileResponse) => {
        try {
            const res = await profileService.repush(profile.id);
            if (res.is_success) antdMessage.success(`Repush queued for "${profile.name}"`);
            else antdMessage.error(res.message || "Repush failed");
        } catch {
            antdMessage.error("Repush failed");
        }
    };

    const handleDelete = async (profile: ProfileResponse) => {
        try {
            const res = await profileService.deleteProfile(profile.id);
            if (res.is_success) {
                antdMessage.success(`Profile "${profile.name}" deleted`);
                fetchProfiles();
            } else {
                antdMessage.error(res.message || "Delete failed");
            }
        } catch {
            antdMessage.error("Delete failed");
        }
    };

    const handleStatusChange = async (profile: ProfileResponse, status: "active" | "draft" | "archived") => {
        try {
            const res = await profileService.updateStatus(profile.id, { status });
            if (res.is_success) {
                antdMessage.success(`Status updated to ${status}`);
                fetchProfiles();
            } else {
                antdMessage.error(res.message || "Status update failed");
            }
        } catch {
            antdMessage.error("Status update failed");
        }
    };

    const handleDuplicate = async (profile: ProfileResponse) => {
        try {
            const res = await profileService.duplicate(profile.id);
            if (res.is_success) {
                antdMessage.success(`Profile "${profile.name}" duplicated`);
                fetchProfiles();
            } else {
                antdMessage.error(res.message || "Duplicate failed");
            }
        } catch {
            antdMessage.error("Duplicate failed");
        }
    };

    // ── modal helpers ─────────────────────────────────────────────────────────
    const openCreate = () => {
        setEditingProfile(null);
        form.resetFields();
        form.setFieldsValue({ platform: "ios", scope: "device", blocked_websites: [] });
        setWifiEnabled(false);
        setVpnEnabled(false);
        setModalOpen(true);
    };

    const openEdit = (profile: ProfileResponse) => {
        setEditingProfile(profile);
        const vals = profileToFormValues(profile);
        form.setFieldsValue(vals);
        setWifiEnabled(vals.wifi_enabled);
        setVpnEnabled(vals.vpn_enabled);
        setModalOpen(true);
    };

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            const payload = buildPayload(values);
            let res;
            if (editingProfile) {
                res = await profileService.updateProfile(editingProfile.id, payload as UpdateProfileRequest);
            } else {
                res = await profileService.createProfile(payload);
            }
            if (res.is_success) {
                antdMessage.success(editingProfile ? "Profile updated" : "Profile created");
                setModalOpen(false);
                fetchProfiles();
            } else {
                antdMessage.error(res.message || "Operation failed");
            }
        } catch (err: unknown) {
            // form validation error — antd shows inline; network error
            if (err && typeof err === "object" && "errorFields" in err) return;
            antdMessage.error("Operation failed");
        } finally {
            setSubmitting(false);
        }
    };

    // ── row actions dropdown ──────────────────────────────────────────────────
    const rowMenu = (profile: ProfileResponse): MenuProps => ({
        items: [
            {
                key: "view",
                label: "View details",
                icon: <Eye className="w-4 h-4" />,
                onClick: () => openEdit(profile),
            },
            {
                key: "repush",
                label: "Repush to devices",
                icon: <RotateCcw className="w-4 h-4" />,
                onClick: () => handleRepush(profile),
                disabled: profile.status !== "active",
            },
            {
                key: "duplicate",
                label: "Duplicate",
                icon: <Copy className="w-4 h-4" />,
                onClick: () => handleDuplicate(profile),
            },
            { type: "divider" },
            ...(profile.status !== "active" ? [{
                key: "activate",
                label: "Set Active",
                icon: <CheckCircle2 className="w-4 h-4" />,
                onClick: () => handleStatusChange(profile, "active"),
            }] : []),
            ...(profile.status !== "draft" ? [{
                key: "draft",
                label: "Move to Draft",
                icon: <PenSquare className="w-4 h-4" />,
                onClick: () => handleStatusChange(profile, "draft"),
            }] : []),
            ...(profile.status !== "archived" ? [{
                key: "archive",
                label: "Archive",
                icon: <Archive className="w-4 h-4" />,
                onClick: () => handleStatusChange(profile, "archived"),
            }] : []),
            { type: "divider" },
            {
                key: "delete",
                label: (
                    <Popconfirm
                        title="Delete profile?"
                        description="This action cannot be undone."
                        onConfirm={() => handleDelete(profile)}
                        okText="Delete"
                        okButtonProps={{ danger: true }}
                    >
                        <span className="text-red-500">Delete</span>
                    </Popconfirm>
                ),
                icon: <Trash2 className="w-4 h-4 text-red-500" />,
            },
        ],
    });

    // ── table columns ─────────────────────────────────────────────────────────
    const columns: ColumnsType<ProfileResponse> = [
        {
            title: "PROFILE NAME",
            dataIndex: "name",
            key: "name",
            render: (text, record) => (
                <div className="flex items-center gap-3">
                    {record.status === "active" ? (
                        <CheckCircle2 className="w-5 h-5 text-emerald-500 shrink-0" strokeWidth={1.5} />
                    ) : record.status === "archived" ? (
                        <Archive className="w-5 h-5 text-amber-400 shrink-0" strokeWidth={1.5} />
                    ) : (
                        <div className="w-5 h-5 rounded-full border border-slate-300 flex items-center justify-center shrink-0">
                            <PenSquare className="w-3 h-3 text-slate-500" strokeWidth={2} />
                        </div>
                    )}
                    <button
                        className="text-slate-700 hover:text-[#de2a15] font-medium text-left"
                        onClick={() => openEdit(record)}
                    >
                        {text}
                    </button>
                </div>
            ),
        },
        {
            title: "PLATFORM",
            dataIndex: "platform",
            key: "platform",
            render: (platform) => (
                <div className="flex items-center gap-2 text-slate-700">
                    {platformIcon(platform)}
                    <span>{platformLabel(platform)}</span>
                </div>
            ),
        },
        {
            title: "SCOPE",
            dataIndex: "scope",
            key: "scope",
            render: (scope) => <span className="capitalize text-slate-700">{scope}</span>,
        },
        {
            title: "VERSION",
            dataIndex: "version",
            key: "version",
            render: (v) => <span className="text-slate-700">{v}.0</span>,
        },
        {
            title: "STATUS",
            dataIndex: "status",
            key: "status",
            render: statusTag,
        },
        {
            title: "CREATED AT",
            dataIndex: "created_at",
            key: "created_at",
            render: (v) => (
                <span className="text-slate-500 text-sm">
                    {v ? new Date(v).toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric" }) : "—"}
                </span>
            ),
        },
        {
            title: "",
            key: "actions",
            width: 48,
            render: (_, record) => (
                <Dropdown menu={rowMenu(record)} trigger={["click"]} placement="bottomRight">
                    <Button
                        type="text"
                        size="small"
                        icon={<MoreVertical className="w-4 h-4 text-slate-500" />}
                        className="hover:bg-slate-100"
                    />
                </Dropdown>
            ),
        },
    ];

    // ── totals display ────────────────────────────────────────────────────────
    const start = (pagination.current - 1) * pagination.pageSize + 1;
    const end = Math.min(pagination.current * pagination.pageSize, pagination.total);
    const totalPages = Math.ceil(pagination.total / pagination.pageSize) || 1;

    // ─────────────────────────────────────────────────────────────────────────
    return (
        <div className="flex flex-col h-[calc(100vh-64px)] bg-slate-50 relative border-none overflow-hidden rounded-none shadow-none z-0">

            {/* ── Top Toolbar ── */}
            <div className="flex flex-wrap items-center justify-between p-4 gap-4 bg-white border-b border-slate-200 z-10 shadow-sm">
                <div className="flex items-center gap-6">
                    <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-700">Platform:</span>
                        <Select
                            value={platformFilter}
                            onChange={(v: string) => { setPlatformFilter(v); }}
                            style={{ width: 150 }}
                            options={PLATFORM_OPTIONS}
                        />
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-700">Status:</span>
                        <Select
                            value={statusFilter}
                            onChange={(v: string) => { setStatusFilter(v); }}
                            style={{ width: 140 }}
                            options={STATUS_OPTIONS}
                        />
                    </div>
                </div>

                <div className="flex items-center gap-3">
                    <div className="flex group">
                        <Input
                            placeholder="Search profile name..."
                            prefix={<Search className="w-4 h-4 text-slate-400" />}
                            className="w-64 h-8 rounded-r-none border-r-0"
                            value={searchInput}
                            onChange={e => setSearchInput(e.target.value)}
                            onPressEnter={() => setSearch(searchInput)}
                        />
                        <Button
                            type="primary"
                            className="bg-[#de2a15] hover:bg-[#c22412] rounded-l-none h-8 w-10 px-0 flex items-center justify-center border-none"
                            icon={<Search className="w-4 h-4 text-white" strokeWidth={2.5} />}
                            onClick={() => setSearch(searchInput)}
                        />
                    </div>
                    <Button
                        type="primary"
                        icon={<Plus className="w-4 h-4" />}
                        className="bg-[#de2a15] hover:bg-[#c22412] text-white font-medium px-5 h-8 border-none rounded-md"
                        onClick={openCreate}
                    >
                        ADD PROFILE
                    </Button>
                </div>
            </div>

            {/* ── Sub Toolbar ── */}
            <div className="flex items-center justify-between px-4 py-3 bg-slate-50 border-b border-slate-200 z-10">
                <div className="flex items-center gap-4 text-sm text-slate-600">
                    <span className="font-bold text-slate-800 tracking-wide uppercase">
                        PROFILES{" "}
                        <span className="font-normal text-slate-500">
                            ({pagination.total > 0 ? `${start} - ${end} of ${pagination.total}` : "0"})
                        </span>
                    </span>

                    <div className="flex items-center gap-2 border-l border-slate-300 pl-4">
                        <Select
                            value={String(pagination.pageSize)}
                            size="small"
                            style={{ width: 70 }}
                            onChange={(v: string) => fetchProfiles(1, Number(v))}
                            options={[
                                { value: "20", label: "20" },
                                { value: "50", label: "50" },
                                { value: "100", label: "100" },
                            ]}
                        />
                        <span>Per Page</span>
                    </div>

                    <div className="flex items-center gap-2 border-l border-slate-300 pl-4">
                        <Button
                            type="text" size="small"
                            disabled={pagination.current <= 1}
                            onClick={() => fetchProfiles(pagination.current - 1)}
                        >&larr;</Button>
                        <span className="text-[#de2a15] font-bold">
                            {pagination.current} of {totalPages}
                        </span>
                        <Button
                            type="text" size="small"
                            disabled={pagination.current >= totalPages}
                            onClick={() => fetchProfiles(pagination.current + 1)}
                        >&rarr;</Button>
                    </div>

                    <div className="border-l border-slate-300 pl-4">
                        <Tooltip title="Refresh">
                            <Button
                                type="text" size="small"
                                icon={<RefreshCcw className="w-4 h-4 text-slate-500" />}
                                onClick={() => fetchProfiles()}
                                loading={loading}
                            />
                        </Tooltip>
                    </div>
                </div>

                <div className="flex items-center gap-4 text-sm">
                    <span className="flex items-center gap-1 text-slate-700 font-medium">
                        Columns (6) <ChevronDown className="w-4 h-4" />
                    </span>
                    <Button
                        type="text"
                        icon={<Filter className="w-4 h-4 text-[#de2a15]" />}
                        className="text-[#de2a15] bg-red-50 hover:bg-red-100 rounded-full w-8 h-8 flex items-center justify-center p-0 border border-red-200"
                    />
                </div>
            </div>

            {/* ── Table ── */}
            <div className="flex-1 overflow-auto border-t border-slate-200 z-10 relative scrollbar-hide">
                <Table
                    columns={columns}
                    dataSource={profiles}
                    rowKey="id"
                    loading={loading}
                    pagination={false}
                    className="custom-data-table"
                    rowClassName="hover:bg-red-50 transition-colors"
                    locale={{ emptyText: "No profiles found" }}
                />
            </div>

            {/* ── Create / Edit Modal ── */}
            <Modal
                title={null}
                open={modalOpen}
                onCancel={() => setModalOpen(false)}
                footer={null}
                width={860}
                className="custom-modal form-modal"
                styles={{ body: { padding: 0 } }}
                centered
                closeIcon={
                    <div className="bg-white hover:bg-slate-100 p-2 rounded-full cursor-pointer">
                        <X className="w-5 h-5 text-slate-700" />
                    </div>
                }
            >
                <div className="flex flex-col bg-white">
                    {/* header */}
                    <div className="px-6 py-4 border-b border-slate-200 flex items-center gap-3 bg-slate-50">
                        <div className="w-10 h-10 rounded-xl bg-white flex items-center justify-center shadow-sm border border-slate-200">
                            <Apple className="w-6 h-6 text-slate-800" fill="currentColor" />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-slate-800 m-0 tracking-wide">
                                {editingProfile ? "EDIT PROFILE" : "CREATE PROFILE"}
                            </h2>
                            <p className="text-xs text-slate-500 font-medium m-0">
                                Configure settings and restrictions for Apple devices
                            </p>
                        </div>
                    </div>

                    <Form
                        form={form}
                        layout="vertical"
                        className="custom-form"
                        initialValues={{
                            platform: "ios",
                            scope: "device",
                            passcode_required: false,
                            min_passcode_length: 0,
                            max_failed_attempts: 0,
                            screen_lock_timeout: 0,
                            camera_disabled: false,
                            bluetooth_disabled: false,
                            airdrop_disabled: false,
                            app_store_disabled: false,
                            wifi_enabled: false,
                            wifi_security_type: "WPA2",
                            vpn_enabled: false,
                            safe_browsing_enabled: false,
                            blocked_websites: [],
                        }}
                    >
                        <Tabs
                            defaultActiveKey="general"
                            className="custom-tabs"
                            items={[
                                // ── General ──────────────────────────────────────────────────────
                                {
                                    key: "general",
                                    label: (
                                        <span className="flex items-center gap-2 px-3 font-semibold uppercase tracking-wider text-[13px]">
                                            <AlertCircle className="w-4 h-4 text-red-500" fill="currentColor" stroke="white" />
                                            General
                                        </span>
                                    ),
                                    children: (
                                        <div className="p-8 overflow-y-auto max-h-[60vh]">
                                            <Form.Item
                                                name="name"
                                                label={<span className="font-semibold text-slate-700">Profile Name <span className="text-red-500">*</span></span>}
                                                rules={[{ required: true, message: "Profile name is required" }]}
                                            >
                                                <Input placeholder="Enter profile name..." className="h-10" />
                                            </Form.Item>

                                            <div className="grid grid-cols-2 gap-4">
                                                <Form.Item
                                                    name="platform"
                                                    label={<span className="font-semibold text-slate-700">Platform</span>}
                                                >
                                                    <Select options={PLATFORM_OPTIONS.filter(p => p.value !== "all")} />
                                                </Form.Item>
                                                <Form.Item
                                                    name="scope"
                                                    label={<span className="font-semibold text-slate-700">Scope</span>}
                                                >
                                                    <Select options={SCOPE_OPTIONS} />
                                                </Form.Item>
                                            </div>

                                            {editingProfile && (
                                                <div className="grid grid-cols-2 gap-4 mt-2">
                                                    <div className="glass-card p-4 rounded-lg">
                                                        <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Status</div>
                                                        <div className="font-medium text-slate-800">{statusTag(editingProfile.status)}</div>
                                                    </div>
                                                    <div className="glass-card p-4 rounded-lg">
                                                        <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Version</div>
                                                        <div className="font-medium text-slate-800">{editingProfile.version}.0</div>
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    ),
                                },

                                // ── Security ──────────────────────────────────────────────────────
                                {
                                    key: "security",
                                    label: (
                                        <span className="flex items-center gap-2 px-3 font-semibold uppercase tracking-wider text-[13px]">
                                            <Shield className="w-4 h-4 text-blue-500" />
                                            Security
                                        </span>
                                    ),
                                    children: (
                                        <div className="p-8 overflow-y-auto max-h-[60vh] flex flex-col gap-4">
                                            <div className="flex items-center justify-between py-3 border-b border-slate-100">
                                                <div>
                                                    <div className="font-semibold text-slate-700">Require Passcode</div>
                                                    <div className="text-xs text-slate-500 mt-0.5">Force device to have a passcode set</div>
                                                </div>
                                                <Form.Item name="passcode_required" valuePropName="checked" className="mb-0">
                                                    <Switch />
                                                </Form.Item>
                                            </div>

                                            <div className="grid grid-cols-3 gap-4">
                                                <Form.Item
                                                    name="min_passcode_length"
                                                    label={<span className="text-sm font-medium text-slate-600">Min Passcode Length</span>}
                                                    tooltip="0 = no minimum"
                                                >
                                                    <InputNumber min={0} max={16} className="w-full" />
                                                </Form.Item>
                                                <Form.Item
                                                    name="max_failed_attempts"
                                                    label={<span className="text-sm font-medium text-slate-600">Max Failed Attempts</span>}
                                                    tooltip="0 = no limit"
                                                >
                                                    <InputNumber min={0} max={20} className="w-full" />
                                                </Form.Item>
                                                <Form.Item
                                                    name="screen_lock_timeout"
                                                    label={<span className="text-sm font-medium text-slate-600">Screen Lock Timeout (s)</span>}
                                                    tooltip="Inactivity seconds before screen locks. 0 = disabled"
                                                >
                                                    <InputNumber min={0} step={60} className="w-full" />
                                                </Form.Item>
                                            </div>

                                            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 text-xs text-blue-700">
                                                <Lock className="w-3.5 h-3.5 inline mr-1.5" />
                                                On iOS, enabling a passcode automatically enables device encryption.
                                            </div>
                                        </div>
                                    ),
                                },

                                // ── Restrictions ───────────────────────────────────────────────────
                                {
                                    key: "restrictions",
                                    label: (
                                        <span className="flex items-center gap-2 px-3 font-semibold uppercase tracking-wider text-[13px]">
                                            <Lock className="w-4 h-4 text-orange-500" />
                                            Restrictions
                                        </span>
                                    ),
                                    children: (
                                        <div className="p-8 overflow-y-auto max-h-[60vh] flex flex-col gap-1">
                                            {([
                                                { name: "camera_disabled", label: "Disable Camera", desc: "Blocks access to the device camera app" },
                                                { name: "bluetooth_disabled", label: "Disable Bluetooth", desc: "Prevents users from enabling Bluetooth" },
                                                { name: "airdrop_disabled", label: "Disable AirDrop", desc: "Blocks AirDrop file sharing (iOS only)" },
                                                { name: "app_store_disabled", label: "Disable App Store", desc: "Prevents installing apps from the App Store" },
                                            ] as const).map(r => (
                                                <div key={r.name} className="flex items-center justify-between py-3.5 border-b border-slate-100 last:border-0">
                                                    <div>
                                                        <div className="font-semibold text-slate-700">{r.label}</div>
                                                        <div className="text-xs text-slate-500 mt-0.5">{r.desc}</div>
                                                    </div>
                                                    <Form.Item name={r.name} valuePropName="checked" className="mb-0">
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                            ))}
                                        </div>
                                    ),
                                },

                                // ── Network ────────────────────────────────────────────────────────
                                {
                                    key: "network",
                                    label: (
                                        <span className="flex items-center gap-2 px-3 font-semibold uppercase tracking-wider text-[13px]">
                                            <Wifi className="w-4 h-4 text-green-500" />
                                            Network
                                        </span>
                                    ),
                                    children: (
                                        <div className="p-8 overflow-y-auto max-h-[60vh] flex flex-col gap-6">
                                            {/* WiFi */}
                                            <div className="border border-slate-200 rounded-xl p-5">
                                                <div className="flex items-center justify-between mb-4">
                                                    <div className="flex items-center gap-2">
                                                        <Wifi className="w-5 h-5 text-green-500" />
                                                        <span className="font-semibold text-slate-700">Wi-Fi Configuration</span>
                                                    </div>
                                                    <Form.Item name="wifi_enabled" valuePropName="checked" className="mb-0">
                                                        <Switch onChange={setWifiEnabled} />
                                                    </Form.Item>
                                                </div>
                                                {wifiEnabled && (
                                                    <div className="flex flex-col gap-3">
                                                        <div className="grid grid-cols-2 gap-3">
                                                            <Form.Item name="wifi_ssid" label="SSID" rules={[{ required: wifiEnabled, message: "SSID required" }]} className="mb-0">
                                                                <Input placeholder="Network name" />
                                                            </Form.Item>
                                                            <Form.Item name="wifi_security_type" label="Security" className="mb-0">
                                                                <Select options={WIFI_SECURITY_OPTIONS} />
                                                            </Form.Item>
                                                        </div>
                                                        <Form.Item name="wifi_password" label="Password" className="mb-0">
                                                            <Input.Password placeholder="Leave blank for open network" />
                                                        </Form.Item>
                                                    </div>
                                                )}
                                            </div>

                                            {/* VPN */}
                                            <div className="border border-slate-200 rounded-xl p-5">
                                                <div className="flex items-center justify-between mb-4">
                                                    <div className="flex items-center gap-2">
                                                        <Server className="w-5 h-5 text-purple-500" />
                                                        <span className="font-semibold text-slate-700">VPN Configuration (IKEv2)</span>
                                                    </div>
                                                    <Form.Item name="vpn_enabled" valuePropName="checked" className="mb-0">
                                                        <Switch onChange={setVpnEnabled} />
                                                    </Form.Item>
                                                </div>
                                                {vpnEnabled && (
                                                    <div className="flex flex-col gap-3">
                                                        <Form.Item name="vpn_server" label="Server Address" rules={[{ required: vpnEnabled, message: "Server required" }]} className="mb-0">
                                                            <Input placeholder="vpn.example.com" />
                                                        </Form.Item>
                                                        <div className="grid grid-cols-2 gap-3">
                                                            <Form.Item name="vpn_username" label="Username" className="mb-0">
                                                                <Input placeholder="VPN username" />
                                                            </Form.Item>
                                                            <Form.Item name="vpn_password" label="Password" className="mb-0">
                                                                <Input.Password placeholder="VPN password" />
                                                            </Form.Item>
                                                        </div>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ),
                                },

                                // ── Content Filter ─────────────────────────────────────────────────
                                {
                                    key: "content_filter",
                                    label: (
                                        <span className="flex items-center gap-2 px-3 font-semibold uppercase tracking-wider text-[13px]">
                                            <Globe className="w-4 h-4 text-red-500" />
                                            Content Filter
                                        </span>
                                    ),
                                    children: (
                                        <div className="p-8 overflow-y-auto max-h-[60vh] flex flex-col gap-4">
                                            <div className="flex items-center justify-between py-3 border-b border-slate-100">
                                                <div>
                                                    <div className="font-semibold text-slate-700">Safe Browsing</div>
                                                    <div className="text-xs text-slate-500 mt-0.5">Auto-filter adult content (requires supervised device)</div>
                                                </div>
                                                <Form.Item name="safe_browsing_enabled" valuePropName="checked" className="mb-0">
                                                    <Switch />
                                                </Form.Item>
                                            </div>

                                            <div>
                                                <div className="font-semibold text-slate-700 mb-1">Blocked Websites</div>
                                                <div className="text-xs text-slate-500 mb-3">Add URLs to block on the device browser (requires supervised device)</div>
                                                <Form.Item name="blocked_websites" className="mb-0">
                                                    <BlockedWebsitesList />
                                                </Form.Item>
                                            </div>

                                            <div className="bg-amber-50 border border-amber-200 rounded-lg p-3 text-xs text-amber-700">
                                                <AlertCircle className="w-3.5 h-3.5 inline mr-1.5" />
                                                Content filtering (com.apple.webcontent-filter) requires supervised devices.
                                            </div>
                                        </div>
                                    ),
                                },
                            ]}
                        />

                        {/* footer */}
                        <Divider className="my-0" />
                        <div className="flex justify-end gap-3 px-6 py-4">
                            <Button onClick={() => setModalOpen(false)}>Cancel</Button>
                            <Button
                                type="primary"
                                className="bg-[#de2a15] hover:bg-[#c22412] border-none"
                                onClick={handleSubmit}
                                loading={submitting}
                            >
                                {editingProfile ? "Save Changes" : "Create Profile"}
                            </Button>
                        </div>
                    </Form>
                </div>
            </Modal>
        </div>
    );
}
