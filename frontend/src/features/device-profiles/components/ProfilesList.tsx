"use client";

import React, { useState, useEffect, useCallback } from "react";
import { Table, Select, Input, Button, Dropdown, MenuProps, Tag, Modal, Tabs, Form, Checkbox, InputNumber, App } from "antd";
import { 
 Search, 
 Plus, 
 Apple, 
 Smartphone, 
 CheckCircle2, 
 Settings2, 
 RefreshCcw,
 ChevronDown,
 Filter,
 PenSquare,
 Monitor,
 MonitorPlay,
 Users,
 AlertCircle,
 Shield,
 Lock,
 Globe,
 Wifi,
 Server,
 Mail,
 Calendar,
 Contact,
 Radio,
 Printer,
 Key,
 Tv,
 MessageSquare,
 Link as LinkIcon,
 Trash2,
 Minus,
 MonitorUp,
 X
} from "lucide-react";
import type { ColumnsType } from "antd/es/table";
import { profileService } from "@/services/profile.service";
import { AssignProfileRequest, CreateProfileRequest, ProfileResponse, UpdateProfileRequest } from "@/types/profile.type";
import dayjs from "dayjs";

interface ProfileType {
 key: string;
 id: number;
 name: string;
 status: "active" | "draft" | "archived" | string;
 family: "Apple" | "Android Plus" | "Windows Modern" | string;
 installMethod: string;
 version: string;
 hasDraft?: boolean;
 configurations: number;
 activeConfigs: {key: string, name: string, icon: React.ReactNode}[];
 packages: number;
 assignedDate: string;
 assignedBy: string;
}

export function ProfilesList() {
 const { message } = App.useApp();

 // Form instances for configs that map to structured backend fields
 const [passcodeForm] = Form.useForm();
 const [restrictionsForm] = Form.useForm();
 const [wifiForm] = Form.useForm();

 // Platform selection
 const [selectedPlatform, setSelectedPlatform] = useState<string>("ios");

 // Config data captured when each modal's SAVE is clicked
 const [passcodeData, setPasscodeData] = useState<any>({});
 const [restrictionsData, setRestrictionsData] = useState<any>({});
 const [wifiData, setWifiData] = useState<any>({});

 // WiFi controlled states (fields not in Form.Item)
 const [wifiSsid, setWifiSsid] = useState("");
 const [wifiSecurityType, setWifiSecurityType] = useState("none");
 const [wifiProxySetup, setWifiProxySetup] = useState("none");

 // VPN controlled states
 const [vpnConnectionName, setVpnConnectionName] = useState("");
 const [vpnServer, setVpnServer] = useState("");
 const [vpnConnectionType, setVpnConnectionType] = useState("l2tp");
 const [vpnData, setVpnData] = useState<any>({});

 // Assignment modal state
 const [isAssignModalVisible, setIsAssignModalVisible] = useState(false);
 const [assignTargetType, setAssignTargetType] = useState<"device" | "group">("device");
 const [assignDeviceId, setAssignDeviceId] = useState("");
 const [assignGroupId, setAssignGroupId] = useState("");
 const [assignScheduleType, setAssignScheduleType] = useState<"immediate" | "scheduled">("immediate");
 const [assignLoading, setAssignLoading] = useState(false);

 const [profiles, setProfiles] = useState<ProfileType[]>([]);
 const [loading, setLoading] = useState(false);
 const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 50,
    total: 0
 });
 const [searchQuery, setSearchQuery] = useState("");
 const [familyFilter, setFamilyFilter] = useState("all");
 const [statusFilter, setStatusFilter] = useState("none");
 const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
 const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
 const [isProfileFormModalVisible, setIsProfileFormModalVisible] = useState(false);
 const [isAddConfigModalVisible, setIsAddConfigModalVisible] = useState(false);
 const [isAppleExpanded, setIsAppleExpanded] = useState(true);
 const [isPasscodeConfigVisible, setIsPasscodeConfigVisible] = useState(false);
 const [isRestrictionsConfigVisible, setIsRestrictionsConfigVisible] = useState(false);
 const [isDomainsConfigVisible, setIsDomainsConfigVisible] = useState(false);
 const [isHttpProxyConfigVisible, setIsHttpProxyConfigVisible] = useState(false);
 const [isDnsProxyConfigVisible, setIsDnsProxyConfigVisible] = useState(false);
 const [isContentFilterConfigVisible, setIsContentFilterConfigVisible] = useState(false);
 const [isCertificateTransparencyConfigVisible, setIsCertificateTransparencyConfigVisible] = useState(false);
 const [isWifiConfigVisible, setIsWifiConfigVisible] = useState(false);
 const [isVpnConfigVisible, setIsVpnConfigVisible] = useState(false);
 const [isAirPlayConfigVisible, setIsAirPlayConfigVisible] = useState(false);
 const [isAirPlaySecurityConfigVisible, setIsAirPlaySecurityConfigVisible] = useState(false);
 const [isAirPrintConfigVisible, setIsAirPrintConfigVisible] = useState(false);
 const [isCalendarConfigVisible, setIsCalendarConfigVisible] = useState(false);
 const [isContactsConfigVisible, setIsContactsConfigVisible] = useState(false);
 const [isExchangeConfigVisible, setIsExchangeConfigVisible] = useState(false);
 const [isGoogleAccountConfigVisible, setIsGoogleAccountConfigVisible] = useState(false);
 const [isLdapConfigVisible, setIsLdapConfigVisible] = useState(false);
 const [isMailConfigVisible, setIsMailConfigVisible] = useState(false);
 const [isMacOsServerConfigVisible, setIsMacOsServerConfigVisible] = useState(false);
 const [isScepConfigVisible, setIsScepConfigVisible] = useState(false);
 const [isCellularConfigVisible, setIsCellularConfigVisible] = useState(false);
 const [isNotificationsConfigVisible, setIsNotificationsConfigVisible] = useState(false);
 const [isConferenceRoomConfigVisible, setIsConferenceRoomConfigVisible] = useState(false);
 const [isTvRemoteConfigVisible, setIsTvRemoteConfigVisible] = useState(false);
 const [isLockScreenMessageConfigVisible, setIsLockScreenMessageConfigVisible] = useState(false);
 const [isWebClipConfigVisible, setIsWebClipConfigVisible] = useState(false);
 const [isSubscribedCalendarConfigVisible, setIsSubscribedCalendarConfigVisible] = useState(false);
 const [hasPasscodeConfig, setHasPasscodeConfig] = useState(false);
 const [hasRestrictionsConfig, setHasRestrictionsConfig] = useState(false);
 const [hasDomainsConfig, setHasDomainsConfig] = useState(false);
 const [hasHttpProxyConfig, setHasHttpProxyConfig] = useState(false);
 const [hasDnsProxyConfig, setHasDnsProxyConfig] = useState(false);
 const [hasContentFilterConfig, setHasContentFilterConfig] = useState(false);
 const [hasCertificateTransparencyConfig, setHasCertificateTransparencyConfig] = useState(false);
 const [hasWifiConfig, setHasWifiConfig] = useState(false);
 const [hasVpnConfig, setHasVpnConfig] = useState(false);
 const [hasAirPlayConfig, setHasAirPlayConfig] = useState(false);
 const [hasAirPlaySecurityConfig, setHasAirPlaySecurityConfig] = useState(false);
 const [hasAirPrintConfig, setHasAirPrintConfig] = useState(false);
 const [hasCalendarConfig, setHasCalendarConfig] = useState(false);
 const [hasContactsConfig, setHasContactsConfig] = useState(false);
 const [hasExchangeConfig, setHasExchangeConfig] = useState(false);
 const [hasGoogleAccountConfig, setHasGoogleAccountConfig] = useState(false);
 const [hasLdapConfig, setHasLdapConfig] = useState(false);
 const [hasMailConfig, setHasMailConfig] = useState(false);
 const [hasMacOsServerConfig, setHasMacOsServerConfig] = useState(false);
 const [hasScepConfig, setHasScepConfig] = useState(false);
 const [hasCellularConfig, setHasCellularConfig] = useState(false);
 const [hasNotificationsConfig, setHasNotificationsConfig] = useState(false);
 const [hasConferenceRoomConfig, setHasConferenceRoomConfig] = useState(false);
 const [hasTvRemoteConfig, setHasTvRemoteConfig] = useState(false);
 const [hasLockScreenMessageConfig, setHasLockScreenMessageConfig] = useState(false);
 const [hasWebClipConfig, setHasWebClipConfig] = useState(false);
 const [hasSubscribedCalendarConfig, setHasSubscribedCalendarConfig] = useState(false);

 // State cho TV Remote
 const [allowedRemotes, setAllowedRemotes] = useState<string[]>([]);
 const [allowedTvs, setAllowedTvs] = useState<string[]>([]);
 const [selectedAllowedRemoteIdx, setSelectedAllowedRemoteIdx] = useState<number | null>(null);
 const [selectedAllowedTvIdx, setSelectedAllowedTvIdx] = useState<number | null>(null);

 // State cho Notifications
 const [notificationSettings, setNotificationSettings] = useState<{appBundleId: string}[]>([]);
 const [selectedNotificationSettingIdx, setSelectedNotificationSettingIdx] = useState<number | null>(null);

 // State cho LDAP Search Settings
 const [ldapSearchSettings, setLdapSearchSettings] = useState<{description: string, scope: string, searchBase: string}[]>([]);
 const [selectedLdapSearchIdx, setSelectedLdapSearchIdx] = useState<number | null>(null);

 // State cho AirPrint
 const [airPrintPrinters, setAirPrintPrinters] = useState<{host: string, useTls: boolean, port: string, resourcePath: string}[]>([]);
 const [selectedAirPrintIdx, setSelectedAirPrintIdx] = useState<number | null>(null);

 // State cho AirPlay
 const [airPlayPasswords, setAirPlayPasswords] = useState<{deviceName: string, password: string}[]>([]);
 const [airPlayAllowed, setAirPlayAllowed] = useState<string[]>([]);
 const [selectedAirPlayPasswordIdx, setSelectedAirPlayPasswordIdx] = useState<number | null>(null);
 const [selectedAirPlayAllowedIdx, setSelectedAirPlayAllowedIdx] = useState<number | null>(null);

 // State cho Certificate Transparency
 const [excludedCertificates, setExcludedCertificates] = useState<string[]>([]);
 const [excludedDomains, setExcludedDomains] = useState<string[]>([]);
 const [selectedExcludedCertificateIdx, setSelectedExcludedCertificateIdx] = useState<number | null>(null);
 const [selectedExcludedDomainIdx, setSelectedExcludedDomainIdx] = useState<number | null>(null);

 // State cho Content Filter
 const [contentFilterType, setContentFilterType] = useState<string>("limit-adult");
 
 // Arrays cho Limit Adult Content
 const [allowedUrls, setAllowedUrls] = useState<string[]>([]);
 const [unallowedUrls, setUnallowedUrls] = useState<string[]>([]);
 const [selectedAllowedUrlIdx, setSelectedAllowedUrlIdx] = useState<number | null>(null);
 const [selectedUnallowedUrlIdx, setSelectedUnallowedUrlIdx] = useState<number | null>(null);

 // Arrays cho Specific Websites Only
 const [specificWebsites, setSpecificWebsites] = useState<{url: string, name: string}[]>([]);
 const [selectedSpecificWebsiteIdx, setSelectedSpecificWebsiteIdx] = useState<number | null>(null);

 // Arrays cho Plugin Custom Data
 const [pluginCustomData, setPluginCustomData] = useState<{key: string, value: string}[]>([]);
 const [selectedPluginDataIdx, setSelectedPluginDataIdx] = useState<number | null>(null);

 // State cho danh sách Domains
 const [unmarkedEmailDomains, setUnmarkedEmailDomains] = useState<string[]>([]);
 const [managedSafariDomains, setManagedSafariDomains] = useState<string[]>([]);
 const [safariPasswordDomains, setSafariPasswordDomains] = useState<string[]>([]);
 
 // State lưu index đang được chọn
 const [selectedEmailDomainIdx, setSelectedEmailDomainIdx] = useState<number | null>(null);
 const [selectedSafariDomainIdx, setSelectedSafariDomainIdx] = useState<number | null>(null);
 const [selectedPasswordDomainIdx, setSelectedPasswordDomainIdx] = useState<number | null>(null);

 const [isProfileDetailModalVisible, setIsProfileDetailModalVisible] = useState(false);
 const [selectedProfile, setSelectedProfile] = useState<ProfileType | null>(null);

 const handleProfileClick = (profile: ProfileType) => {
  setSelectedProfile(profile);
  setIsProfileDetailModalVisible(true);
 };

 const handleDeleteProfile = async () => {
  if (!selectedProfile) return;
  
  Modal.confirm({
   title: 'Delete Profile',
   content: `Are you sure you want to delete profile "${selectedProfile.name}"?`,
   okText: 'Delete',
   okButtonProps: { danger: true },
   cancelText: 'Cancel',
   onOk: () => new Promise(async (resolve, reject) => {
    try {
     const response = await profileService.deleteProfile(selectedProfile.id);
     // Coi 204 (No Content) hoặc 200 là thành công
     if (response.is_success || (response as any).status === 204) {
      message.success('Profile deleted successfully');
      setIsProfileDetailModalVisible(false);
      fetchProfiles();
      resolve(true);
     } else {
      message.error(response.message || 'Failed to delete profile');
      reject(new Error(response.message || 'Failed to delete profile'));
     }
    } catch (error: any) {
     // Xử lý case axios throw lỗi nhưng HTTP status là 204 hoặc 200
     if (error?.response?.status === 204 || error?.response?.status === 200) {
      message.success('Profile deleted successfully');
      setIsProfileDetailModalVisible(false);
      fetchProfiles();
      resolve(true);
     } else {
      message.error('An error occurred while deleting profile');
      reject(error);
     }
    }
   }),
  });
 };

 const handleUpdateStatus = async (newStatus: "active" | "draft" | "archived") => {
  if (!selectedProfile) return;
  try {
   const response = await profileService.updateProfileStatus(selectedProfile.id, newStatus);
   if (response.is_success) {
    message.success(`Profile status updated to ${newStatus}`);
    setIsProfileDetailModalVisible(false);
    fetchProfiles();
   } else {
    message.error(response.message || "Failed to update status");
   }
  } catch (error) {
   message.error("An error occurred while updating status");
  }
 };

 const handleAssignProfile = async () => {
  if (!selectedProfile) return;
  if (assignTargetType === "device" && !assignDeviceId.trim()) {
   message.error("Device ID is required");
   return;
  }
  if (assignTargetType === "group" && !assignGroupId.trim()) {
   message.error("Group ID is required");
   return;
  }
  setAssignLoading(true);
  try {
   const payload: AssignProfileRequest = {
    target_type: assignTargetType,
    device_id: assignTargetType === "device" ? assignDeviceId.trim() : undefined,
    group_id: assignTargetType === "group" ? Number(assignGroupId) : undefined,
    schedule_type: assignScheduleType,
   };
   const response = await profileService.assignProfile(selectedProfile.id, payload);
   if (response.is_success) {
    message.success("Profile assigned successfully");
    setIsAssignModalVisible(false);
    setAssignDeviceId("");
    setAssignGroupId("");
   } else {
    message.error(response.message || "Failed to assign profile");
   }
  } catch (error) {
   message.error("An error occurred while assigning profile");
  } finally {
   setAssignLoading(false);
  }
 };

 const [newProfileName, setNewProfileName] = useState("");
 const [newProfileDesc, setNewProfileDesc] = useState("");
const [editingProfileId, setEditingProfileId] = useState<number | null>(null);
const [editingProfileSnapshot, setEditingProfileSnapshot] = useState<ProfileResponse | null>(null);

const handleSaveProfile = async () => {
  if (!newProfileName.trim()) {
   message.error("Tên cấu hình là bắt buộc.");
   return;
  }

  try {
   // Build security_settings from passcode form
   const security_settings: any = hasPasscodeConfig
    ? { passcode: passcodeData }
    : (editingProfileSnapshot?.security_settings || {});

   // Build network_config from wifi/vpn/proxy data
   const network_config: any = {};
   if (hasWifiConfig) {
    network_config.wifi = { ssid: wifiSsid, ...wifiData };
   }
   if (hasVpnConfig) {
    network_config.vpn = { connection_name: vpnConnectionName, server: vpnServer, type: vpnConnectionType, ...vpnData };
   }
   if (hasHttpProxyConfig) {
    network_config.proxy = { enabled: true };
   }
   // Preserve existing network config if nothing changed
   if (!hasWifiConfig && !hasVpnConfig && !hasHttpProxyConfig && editingProfileSnapshot?.network_config) {
    Object.assign(network_config, editingProfileSnapshot.network_config);
   }

   // Build restrictions from restrictions form
   const restrictions: any = hasRestrictionsConfig
    ? restrictionsData
    : (editingProfileSnapshot?.restrictions || {});

   // Build content_filter
   const content_filter: any = hasContentFilterConfig
    ? {
       type: contentFilterType,
       safe_browsing: contentFilterType === "limit-adult",
       allowed_domains: allowedUrls,
       blocked_websites: unallowedUrls,
       specific_websites: specificWebsites,
       plugin_data: pluginCustomData,
      }
    : (editingProfileSnapshot?.content_filter || {});

   // Build payloads for configs not in structured backend fields
   const payloads: any = { ...(editingProfileSnapshot?.payloads || {}) };
   if (hasDnsProxyConfig) payloads.dns_proxy = { enabled: true };
   if (hasCellularConfig) payloads.cellular = { enabled: true };
   if (hasAirPlayConfig) payloads.airplay = { passwords: airPlayPasswords, allowed: airPlayAllowed };
   if (hasAirPlaySecurityConfig) payloads.airplay_security = { enabled: true };
   if (hasAirPrintConfig) payloads.airprint = { printers: airPrintPrinters };
   if (hasCertificateTransparencyConfig) payloads.certificate_transparency = { excluded_certificates: excludedCertificates, excluded_domains: excludedDomains };
   if (hasDomainsConfig) payloads.domains = { unmarked_email: unmarkedEmailDomains, managed_safari: managedSafariDomains, safari_password: safariPasswordDomains };
   if (hasCalendarConfig) payloads.calendar = { enabled: true };
   if (hasContactsConfig) payloads.contacts = { enabled: true };
   if (hasExchangeConfig) payloads.exchange = { enabled: true };
   if (hasGoogleAccountConfig) payloads.google_account = { enabled: true };
   if (hasLdapConfig) payloads.ldap = { search_settings: ldapSearchSettings };
   if (hasMailConfig) payloads.mail = { enabled: true };
   if (hasMacOsServerConfig) payloads.macos_server = { enabled: true };
   if (hasScepConfig) payloads.scep = { enabled: true };
   if (hasNotificationsConfig) payloads.notifications = { settings: notificationSettings };
   if (hasConferenceRoomConfig) payloads.conference_room = { enabled: true };
   if (hasTvRemoteConfig) payloads.tv_remote = { allowed_remotes: allowedRemotes, allowed_tvs: allowedTvs };
   if (hasLockScreenMessageConfig) payloads.lock_screen_message = { enabled: true };
   if (hasWebClipConfig) payloads.web_clip = { enabled: true };
   if (hasSubscribedCalendarConfig) payloads.subscribed_calendar = { enabled: true };

   const payload: UpdateProfileRequest | CreateProfileRequest = {
    name: newProfileName.trim(),
    platform: selectedPlatform || editingProfileSnapshot?.platform || "ios",
    scope: editingProfileSnapshot?.scope || "device",
    compliance_rules: editingProfileSnapshot?.compliance_rules || {},
    content_filter,
    network_config,
    payloads,
    restrictions,
    security_settings,
   };

   const response = editingProfileId
    ? await profileService.updateProfile(editingProfileId, payload as UpdateProfileRequest)
    : await profileService.createProfile(payload as CreateProfileRequest);

   if (response.is_success) {
    message.success(editingProfileId ? "Profile saved successfully" : "Profile created successfully");
    setIsProfileFormModalVisible(false);
    setNewProfileName("");
    setNewProfileDesc("");
    setEditingProfileId(null);
    setEditingProfileSnapshot(null);
    fetchProfiles();
   } else {
    message.error(response.message || "Failed to save profile");
   }
  } catch (error) {
   message.error("An error occurred while saving profile");
  }
 };

 const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
  setSelectedRowKeys(newSelectedRowKeys);
 };

 const fetchProfiles = useCallback(async () => {
  setLoading(true);
  try {
   const params = {
    page: pagination.current,
    limit: pagination.pageSize,
    search: searchQuery || undefined,
    platform: familyFilter !== "all" ? familyFilter : undefined,
    status: statusFilter !== "none" ? statusFilter : undefined,
   };

   const response = await profileService.getProfiles(params);
    if (response.is_success && response.data && response.data.items) {
     const data = response.data.items.map((item: ProfileResponse) => {
      const activeConfigs = [];
      if (item.network_config && Object.keys(item.network_config).length > 0) activeConfigs.push({key: 'network_config', name: 'Wi-Fi / Network', icon: <Wifi className="w-4 h-4" />});
      if (item.restrictions && Object.keys(item.restrictions).length > 0) activeConfigs.push({key: 'restrictions', name: 'Restrictions', icon: <Shield className="w-4 h-4" />});
      if (item.security_settings && Object.keys(item.security_settings).length > 0) activeConfigs.push({key: 'security_settings', name: 'Passcode / Security', icon: <Lock className="w-4 h-4" />});
      if (item.content_filter && Object.keys(item.content_filter).length > 0) activeConfigs.push({key: 'content_filter', name: 'Content Filter', icon: <Globe className="w-4 h-4" />});
      if (item.payloads && Object.keys(item.payloads).length > 0) activeConfigs.push({key: 'payloads', name: 'Custom Payloads', icon: <Settings2 className="w-4 h-4" />});
      if (item.compliance_rules && Object.keys(item.compliance_rules).length > 0) activeConfigs.push({key: 'compliance_rules', name: 'Compliance Rules', icon: <CheckCircle2 className="w-4 h-4" />});

      return {
       key: item.id.toString(),
       id: item.id,
       name: item.name,
       status: item.status,
       family: item.platform === 'ios' || item.platform === 'macos' ? 'Apple' : item.platform === 'android' ? 'Android Plus' : item.platform === 'windows' ? 'Windows Modern' : item.platform,
       installMethod: "Automatic",
       version: `${item.version}.0`,
       hasDraft: item.status === 'draft',
       configurations: activeConfigs.length || 1,
       activeConfigs: activeConfigs.length > 0 ? activeConfigs : [{key: 'general', name: 'General Information', icon: <AlertCircle className="w-4 h-4" />}],
       packages: 0,
       assignedDate: item.updated_at ? dayjs(item.updated_at).format('YYYY-MM-DD hh:mm:ss A') : "N/A",
       assignedBy: "System"
      };
     });
    
    setProfiles(data);
    setPagination(prev => ({
     ...prev,
     total: response.data?.pagination?.total || 0
    }));
   } else {
    message.error(response.message || "Failed to fetch profiles");
   }
  } catch (error) {
   console.error("Failed to fetch profiles:", error);
   message.error("An error occurred while fetching profiles");
  } finally {
   setLoading(false);
  }
 }, [pagination.current, pagination.pageSize, searchQuery, familyFilter, statusFilter, message]);

 useEffect(() => {
  fetchProfiles();
 }, [fetchProfiles]);

 const handleAddProfileClick = () => {
 setIsAppleExpanded(false); // Reset expansion state when opening modal
 setIsCreateModalVisible(true);
 };

 const handleAppleClick = (platform: string = "ios") => {
 setSelectedPlatform(platform);
 setEditingProfileId(null);
 setEditingProfileSnapshot(null);
 setNewProfileName("");
 setNewProfileDesc("");
 setIsCreateModalVisible(false);
 setIsProfileFormModalVisible(true);
 };

 const rowSelection = {
 selectedRowKeys,
 onChange: onSelectChange,
 };

 const columns: ColumnsType<ProfileType> = [
 {
 title: "PROFILE NAME",
 dataIndex: "name",
 key: "name",
 render: (text, record) => (
 <div className="flex items-center gap-3">
 {record.status === "active" ? (
 <CheckCircle2 className="w-5 h-5 text-emerald-500" strokeWidth={1.5} />
 ) : (
 <div className="w-5 h-5 rounded-full border border-slate-300 flex items-center justify-center">
 <PenSquare className="w-3 h-3 text-slate-500" strokeWidth={2} />
 </div>
 )}
 <a 
      href="#" 
      className="text-slate-700 hover:text-[#de2a15] font-medium"
      onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        handleProfileClick(record);
      }}
    >
      {text}
    </a>
 </div>
 ),
 },
 {
 title: "FAMILY",
 dataIndex: "family",
 key: "family",
 render: (text) => (
 <div className="flex items-center gap-2 text-slate-700">
 {text === "Apple" ? (
 <Apple className="w-5 h-5" strokeWidth={1.5} />
 ) : text === "Android Plus" ? (
 <Smartphone className="w-5 h-5" strokeWidth={1.5} />
 ) : (
 <Settings2 className="w-5 h-5" strokeWidth={1.5} />
 )}
 <span>{text}</span>
 </div>
 ),
 },
 {
 title: "INSTALL METHOD",
 dataIndex: "installMethod",
 key: "installMethod",
 render: (text) => <span className="text-slate-700">{text}</span>,
 },
 {
 title: "VERSION",
 dataIndex: "version",
 key: "version",
 render: (text, record) => (
 <div className="flex items-center gap-2">
 <span className="text-slate-700">{text}</span>
 {record.hasDraft && (
 <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border border-slate-200 text-slate-700 bg-slate-50">
 <PenSquare className="w-3 h-3" />
 Draft available
 </span>
 )}
 </div>
 ),
 },
 {
 title: "CONFIGURATIONS",
 dataIndex: "configurations",
 key: "configurations",
 render: (text) => <span className="text-slate-700">{text}</span>,
 },
 {
 title: "PACKAGES",
 dataIndex: "packages",
 key: "packages",
 render: (text) => <span className="text-slate-700">{text}</span>,
 },
 {
 title: "ASSIGNED DATE",
 dataIndex: "assignedDate",
 key: "assignedDate",
 render: (text) => <span className="text-slate-500">{text}</span>,
 },
 {
 title: "ASSIGNED BY",
 dataIndex: "assignedBy",
 key: "assignedBy",
 render: (text) => <span className="text-slate-500">{text}</span>,
 },
 ];

 return (
 <div className="flex flex-col h-[calc(100vh-64px)] bg-slate-50 relative border-none overflow-hidden rounded-none shadow-none z-0">
 {/* Top Toolbar */}
        <div className="flex flex-wrap items-center justify-between p-4 gap-4 bg-white border-b border-slate-200 z-10 shadow-sm">
        <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-slate-700 dark:text-slate-300">Family:</span>
        <Select className="cursor-pointer" value={familyFilter}
        onChange={(val) => { setFamilyFilter(val); setPagination(prev => ({ ...prev, current: 1 })); }}
        style={{ width: 120 }}
        options={[
        { value: "all", label: "All" },
        { value: "apple", label: "Apple" },
        { value: "android", label: "Android Plus" },
        ]}
        />
        </div>
        <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-slate-700 dark:text-slate-300">Filters:</span>
        <Select className="cursor-pointer" value={statusFilter}
        onChange={(val) => { setStatusFilter(val); setPagination(prev => ({ ...prev, current: 1 })); }}
        style={{ width: 120 }}
        options={[
        { value: "none", label: "None" },
        { value: "active", label: "Active" },
        { value: "draft", label: "Draft" },
        ]}
        />
        </div>
        </div>

        <div className="flex items-center gap-3">
        <div className="flex group">
        <Input
        placeholder="Search profile name"
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        onPressEnter={() => setPagination(prev => ({ ...prev, current: 1 }))}
        prefix={<Search className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" />}
        className="w-64 h-8 rounded-r-none border-r-0 hover:border-[#de2a15] focus:border-[#de2a15] focus:shadow-none transition-colors"
        />
        <Button 
        type="primary" 
        onClick={() => setPagination(prev => ({ ...prev, current: 1 }))}
        className="bg-[#de2a15] hover:bg-[#c22412] rounded-l-none h-8 w-10 px-0 flex items-center justify-center border-none shadow-sm transition-colors"
        icon={<Search className="w-4 h-4 text-white" strokeWidth={2.5} />}
        />
        </div>
        <Button 
        type="primary" 
        icon={<Plus className="w-4 h-4" />}
        className="bg-[#de2a15] hover:bg-[#c22412] text-white font-medium px-5 h-8 border-none shadow-sm transition-colors rounded-md"
        onClick={handleAddProfileClick}
        >
        ADD PROFILE
        </Button>
        </div>
        </div>

        {/* Sub Toolbar / Table Controls */}
        <div className="flex items-center justify-between px-4 py-3 bg-slate-50 dark:bg-slate-800 border-b border-slate-200 dark:border-slate-800 z-10">
        <div className="flex items-center gap-4 text-sm text-slate-600 dark:text-slate-400">
        <span className="font-bold text-slate-800 dark:text-slate-200 tracking-wide uppercase">
        PROFILES <span className="font-normal text-slate-500">({profiles.length > 0 ? (pagination.current - 1) * pagination.pageSize + 1 : 0} - {Math.min(pagination.current * pagination.pageSize, pagination.total)} of {pagination.total})</span>
        </span>
        
        <div className="flex items-center gap-2 border-l border-slate-300 dark:border-slate-600 pl-4">
        <Select className="cursor-pointer" value={pagination.pageSize.toString()}
        onChange={(val) => setPagination(prev => ({ ...prev, pageSize: Number(val), current: 1 }))}
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

        <div className="flex items-center gap-2 border-l border-slate-300 dark:border-slate-600 pl-4">
        <Button type="text" size="small" 
         disabled={pagination.current === 1} 
         onClick={() => setPagination(prev => ({ ...prev, current: prev.current - 1 }))}
         className={pagination.current === 1 ? "text-slate-400" : "text-slate-700 hover:text-slate-900"}>&larr;</Button>
        <span className="text-[#de2a15] font-bold">{pagination.current} of {Math.max(1, Math.ceil(pagination.total / pagination.pageSize))}</span>
        <Button type="text" size="small" 
         disabled={pagination.current >= Math.ceil(pagination.total / pagination.pageSize) || pagination.total === 0}
         onClick={() => setPagination(prev => ({ ...prev, current: prev.current + 1 }))}
         className={pagination.current >= Math.ceil(pagination.total / pagination.pageSize) || pagination.total === 0 ? "text-slate-400" : "text-slate-700 hover:text-slate-900"}>&rarr;</Button>
        </div>

        <div className="border-l border-slate-300 dark:border-slate-600 pl-4">
        <Button type="text" size="small" onClick={fetchProfiles} icon={<RefreshCcw className={`w-4 h-4 text-slate-500 hover:text-slate-800 ${loading ? 'animate-spin' : ''}`} />} />
        </div>
        </div>

        <div className="flex items-center gap-4 text-sm">
        <span className="flex items-center gap-1 cursor-pointer text-slate-700 hover:text-slate-900 font-medium">
        Columns (9) <ChevronDown className="w-4 h-4" />
        </span>
        <Button type="text" icon={<Filter className="w-4 h-4 text-[#de2a15]" />} className="text-[#de2a15] bg-red-50 hover:bg-red-100 rounded-full w-8 h-8 flex items-center justify-center p-0 border border-red-200 transition-colors" />
        </div>
        </div>

        {/* Table */}
        <div className="flex-1 overflow-auto border-t border-slate-200 dark:border-slate-800 z-10 relative scrollbar-hide">
        <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={profiles}
        loading={loading}
        pagination={false}
        className="custom-data-table"
        rowClassName="hover:bg-red-50 dark:hover:bg-slate-800 transition-colors cursor-pointer"
        onRow={(record) => ({
          onClick: () => handleProfileClick(record),
        })}
        />
        </div>

  {/* Profile Detail Modal */}
  <Modal
    title={
      <div className="flex items-center gap-2">
        <div className="w-8 h-8 rounded-lg bg-red-50 flex items-center justify-center">
          {selectedProfile?.family === "Apple" ? <Apple className="w-4 h-4 text-[#de2a15]" /> : <Smartphone className="w-4 h-4 text-[#de2a15]" />}
        </div>
        <div>
          <h3 className="text-lg font-bold text-slate-800 m-0">{selectedProfile?.name}</h3>
          <p className="text-xs text-slate-500 font-normal m-0">{selectedProfile?.family} Profile</p>
        </div>
      </div>
    }
    open={isProfileDetailModalVisible}
    onCancel={() => setIsProfileDetailModalVisible(false)}
    footer={
      <div className="flex justify-end items-center w-full pt-4 border-t border-slate-100 mt-4 gap-3 px-1 pb-1">
        <Button
          danger
          type="text"
          icon={<Trash2 className="w-4 h-4" />}
          onClick={handleDeleteProfile}
          className="text-red-500 hover:bg-red-50 hover:text-red-600 transition-colors h-10 px-4 rounded-lg font-medium mr-auto"
        >
          Delete
        </Button>
        {selectedProfile?.status === "active" ? (
          <Button
            onClick={() => handleUpdateStatus("draft")}
            className="h-10 px-4 rounded-lg font-medium text-amber-600 border-amber-300 hover:bg-amber-50 transition-colors"
          >
            Set Draft
          </Button>
        ) : (
          <Button
            onClick={() => handleUpdateStatus("active")}
            className="h-10 px-4 rounded-lg font-medium text-emerald-600 border-emerald-300 hover:bg-emerald-50 transition-colors"
          >
            Set Active
          </Button>
        )}
        <Button
          onClick={() => { setIsProfileDetailModalVisible(false); setIsAssignModalVisible(true); }}
          className="h-10 px-4 rounded-lg font-medium text-blue-600 border-blue-300 hover:bg-blue-50 transition-colors"
        >
          Assign
        </Button>
        <Button
          onClick={() => setIsProfileDetailModalVisible(false)}
          className="h-10 px-6 rounded-lg font-medium text-slate-600 hover:text-slate-800 hover:bg-slate-50 border-slate-200 transition-colors"
        >
          Close
        </Button>
        <Button 
          type="primary" 
          className="bg-[#de2a15] hover:bg-[#c22412] h-10 px-6 rounded-lg font-medium shadow-sm border-none transition-colors"
          onClick={async () => {
            setIsProfileDetailModalVisible(false);
            setNewProfileName(selectedProfile?.name || "");
            setEditingProfileId(selectedProfile?.id || null);
            
            message.loading({ content: 'Loading profile data...', key: 'loadingProfile' });
            
            try {
              const response = await profileService.getProfileById(selectedProfile!.id);
              if (response.is_success && response.data) {
                const detail = response.data;
                setEditingProfileSnapshot(detail);
                setNewProfileDesc("");
                
                setHasWifiConfig(!!(detail.network_config?.wifi && Object.keys(detail.network_config.wifi).length > 0));
                setHasVpnConfig(!!(detail.network_config?.vpn && Object.keys(detail.network_config.vpn).length > 0));
                setHasHttpProxyConfig(!!(detail.network_config?.proxy && Object.keys(detail.network_config.proxy).length > 0));
                setHasRestrictionsConfig(!!(detail.restrictions && Object.keys(detail.restrictions).length > 0));
                setHasPasscodeConfig(!!(detail.security_settings?.passcode && Object.keys(detail.security_settings.passcode).length > 0));
                setHasContentFilterConfig(!!(detail.content_filter && Object.keys(detail.content_filter).length > 0));
                setSelectedPlatform(detail.platform || "ios");

                // Pre-populate form values for editing
                if (detail.security_settings?.passcode) {
                  const p = detail.security_settings.passcode;
                  passcodeForm.setFieldsValue({
                    allowSimple: p.allow_simple ?? true,
                    requireAlphanumeric: p.require_alphanumeric ?? false,
                    minLength: p.min_length ?? 0,
                    minComplexChars: p.min_complex_chars ?? 0,
                    maxPasscodeAge: p.max_passcode_age,
                    autoLock: p.auto_lock ?? "none",
                    passcodeHistory: p.history,
                    gracePeriod: p.grace_period ?? "none",
                    maxFailedAttempts: p.retry_limit,
                  });
                  setPasscodeData(detail.security_settings.passcode);
                }
                if (detail.restrictions && Object.keys(detail.restrictions).length > 0) {
                  const r = detail.restrictions;
                  restrictionsForm.setFieldsValue({
                    allowCamera: r.camera_enabled ?? true,
                    allowFaceTime: r.facetime_enabled ?? true,
                    allowScreenshots: r.screenshots_enabled ?? true,
                    allowAirDrop: r.airdrop_enabled ?? true,
                    allowSiri: r.siri_enabled ?? true,
                    allowSafari: r.safari_enabled ?? true,
                    allowGameCenter: r.game_center_enabled ?? true,
                    allowiTunes: r.itunes_enabled ?? true,
                    allowNews: r.news_enabled ?? true,
                    allowPodcasts: r.podcasts_enabled ?? true,
                  });
                  setRestrictionsData(detail.restrictions);
                }
                if (detail.network_config?.wifi) {
                  const w = detail.network_config.wifi;
                  setWifiSsid(w.ssid || "");
                  setWifiSecurityType(w.security_type || "none");
                  setWifiProxySetup(w.proxy_setup || "none");
                  wifiForm.setFieldsValue({
                    autoJoin: w.auto_join ?? true,
                    hiddenNetwork: w.hidden_network ?? false,
                    disableCaptive: w.disable_captive ?? false,
                    disableMacRand: w.disable_mac_randomization ?? false,
                  });
                  setWifiData(detail.network_config.wifi);
                }
                if (detail.network_config?.vpn) {
                  const v = detail.network_config.vpn;
                  setVpnConnectionName(v.connection_name || "");
                  setVpnServer(v.server || "");
                  setVpnConnectionType(v.type || "l2tp");
                  setVpnData(detail.network_config.vpn);
                }

                setIsProfileFormModalVisible(true);
                message.success({ content: 'Profile data loaded', key: 'loadingProfile', duration: 2 });
              } else {
                message.error({ content: 'Failed to load profile details', key: 'loadingProfile', duration: 2 });
              }
            } catch (error) {
              message.error({ content: 'Error loading profile details', key: 'loadingProfile', duration: 2 });
            }
          }}
        >
          Edit Profile
        </Button>
      </div>
    }
    width={600}
    className="custom-modal"
  >
    {selectedProfile && (
      <div className="py-2 space-y-5">
        {/* Main Stats Grid */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex items-center gap-1.5 text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2">
              Status
            </div>
            <div className="flex items-center gap-2 font-bold text-slate-800 text-[15px]">
              {selectedProfile.status === 'active' ? (
                <><span className="relative flex h-2.5 w-2.5"><span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span><span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span></span> <span className="text-emerald-700">Active</span></>
              ) : (
                <><span className="w-2.5 h-2.5 rounded-full bg-amber-500"></span> <span className="text-amber-700">Draft</span></>
              )}
            </div>
          </div>
          
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex items-center gap-1.5 text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2">
              Version
            </div>
            <div className="font-bold text-slate-800 text-[15px]">{selectedProfile.version}</div>
          </div>

          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex items-center gap-1.5 text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2">
              Install Method
            </div>
            <div className="font-bold text-slate-800 text-[15px]">{selectedProfile.installMethod}</div>
          </div>

          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex items-center gap-1.5 text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2">
              Configurations
            </div>
            <div className="font-bold text-slate-800 text-[15px]">{selectedProfile.configurations} <span className="font-medium text-slate-500 text-sm">Settings</span></div>
          </div>
        </div>

        {/* Active Configurations List */}
        {selectedProfile.activeConfigs && selectedProfile.activeConfigs.length > 0 && (
          <div className="bg-slate-50/80 p-5 rounded-2xl border border-slate-100">
            <div className="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-4 flex items-center gap-2">
              <Settings2 className="w-3.5 h-3.5" /> Active Configurations
            </div>
            <div className="grid grid-cols-1 gap-2.5">
              {selectedProfile.activeConfigs.map((config) => (
                <div key={config.key} className="flex items-center gap-3 bg-white px-4 py-3 rounded-xl border border-slate-200 shadow-sm hover:border-[#de2a15]/30 transition-colors group cursor-default">
                  <div className="p-2 bg-red-50 text-[#de2a15] rounded-lg group-hover:bg-[#de2a15] group-hover:text-white transition-colors">
                    {config.icon}
                  </div>
                  <span className="font-semibold text-[14px] text-slate-700 group-hover:text-slate-900 transition-colors">{config.name}</span>
                </div>
              ))}
            </div>
          </div>
        )}
        
        {/* Assignment Info */}
        <div className="bg-slate-50/80 p-5 rounded-2xl border border-slate-100">
          <div className="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-4 flex items-center gap-2">
            <Calendar className="w-3.5 h-3.5" /> Assignment Info
          </div>
          <div className="space-y-3">
            <div className="flex items-center justify-between py-1 border-b border-slate-200/60 border-dashed pb-3">
              <span className="text-sm font-medium text-slate-500">Last Modified</span>
              <span className="text-sm font-bold text-slate-700 bg-white px-2.5 py-1 rounded-md border border-slate-200 shadow-sm">{selectedProfile.assignedDate}</span>
            </div>
            <div className="flex items-center justify-between py-1 pt-1">
              <span className="text-sm font-medium text-slate-500">Modified By</span>
              <span className="text-sm font-bold text-slate-700 flex items-center gap-1.5">
                <div className="w-5 h-5 rounded-full bg-slate-200 flex items-center justify-center text-[10px] text-slate-600">
                  {selectedProfile.assignedBy.charAt(0)}
                </div>
                {selectedProfile.assignedBy}
              </span>
            </div>
          </div>
        </div>
      </div>
    )}
  </Modal>

 {/* Create Profile Modals */}
 <Modal
 title={null}
 open={isCreateModalVisible}
 onCancel={() => setIsCreateModalVisible(false)}
 footer={null}
 width={500}
 className="dropdown-modal"
 centered
 closeIcon={<div className="bg-slate-100 hover:bg-slate-200 p-1.5 rounded-full transition-colors absolute top-2 right-2 z-50 cursor-pointer"><X className="w-4 h-4 text-slate-700" /></div>}
 >
 <div className="flex flex-col relative z-10 p-5 gap-6">
 {/* Top side: Logo */}
 <div className="flex flex-col items-center justify-center pt-4 pb-2">
 <div className="w-20 h-20 rounded-3xl bg-slate-50 flex items-center justify-center mb-4 shadow-sm border border-slate-200">
 <Apple className="w-10 h-10 text-slate-800" fill="currentColor" />
 </div>
 <div className="font-bold text-2xl text-slate-800 tracking-wide text-center leading-tight mb-1">
 Apple
 </div>
 <div className="text-sm text-slate-500 font-medium">Device Profile</div>
 </div>
 
 {/* Bottom side: Options */}
 <div className="grid grid-cols-1 gap-2 p-3 bg-slate-50 rounded-2xl border border-slate-200 shadow-inner">
 <button
 onClick={() => handleAppleClick("ios")}
 className="flex items-center gap-4 text-slate-700 hover:text-[#de2a15] hover:bg-red-50 py-3.5 px-5 rounded-xl transition-all font-semibold text-left border border-transparent hover:border-red-200 hover:shadow-md group cursor-pointer"
 >
 <div className="w-10 h-10 rounded-lg bg-red-50 text-[#de2a15] flex items-center justify-center group-hover:bg-red-100 transition-colors">
 <Smartphone className="w-5 h-5" />
 </div>
 <div className="flex flex-col">
 <span className="text-[15px]">iOS / iPadOS</span>
 <span className="text-[11px] text-slate-500 font-medium">iPhones and iPads</span>
 </div>
 </button>
 <button
 onClick={() => handleAppleClick("macos")}
 className="flex items-center gap-4 text-slate-700 hover:text-[#de2a15] hover:bg-red-50 py-3.5 px-5 rounded-xl transition-all font-semibold text-left border border-transparent hover:border-red-200 hover:shadow-md group cursor-pointer"
 >
 <div className="w-10 h-10 rounded-lg bg-red-50 text-[#de2a15] flex items-center justify-center group-hover:bg-red-100 transition-colors">
 <Monitor className="w-5 h-5" />
 </div>
 <div className="flex flex-col">
 <span className="text-[15px]">macOS User</span>
 <span className="text-[11px] text-slate-500 font-medium">Mac user level settings</span>
 </div>
 </button>
 <button
 onClick={() => handleAppleClick("macos")}
 className="flex items-center gap-4 text-slate-700 hover:text-[#de2a15] hover:bg-red-50 py-3.5 px-5 rounded-xl transition-all font-semibold text-left border border-transparent hover:border-red-200 hover:shadow-md group cursor-pointer"
 >
 <div className="w-10 h-10 rounded-lg bg-red-50 text-[#de2a15] flex items-center justify-center group-hover:bg-red-100 transition-colors">
 <Server className="w-5 h-5" />
 </div>
 <div className="flex flex-col">
 <span className="text-[15px]">macOS Device</span>
 <span className="text-[11px] text-slate-500 font-medium">Mac system level settings</span>
 </div>
 </button>
 <button
 onClick={() => handleAppleClick("ios")}
 className="flex items-center gap-4 text-slate-700 hover:text-[#de2a15] hover:bg-red-50 py-3.5 px-5 rounded-xl transition-all font-semibold text-left border border-transparent hover:border-red-200 hover:shadow-md group cursor-pointer"
 >
 <div className="w-10 h-10 rounded-lg bg-red-50 text-[#de2a15] flex items-center justify-center group-hover:bg-red-100 transition-colors">
 <Users className="w-5 h-5" />
 </div>
 <div className="flex flex-col">
 <span className="text-[15px]">Shared iPad User</span>
 <span className="text-[11px] text-slate-500 font-medium">For shared iPad environments</span>
 </div>
 </button>
 <button
 onClick={() => handleAppleClick("tvos")}
 className="flex items-center gap-4 text-slate-700 hover:text-[#de2a15] hover:bg-red-50 py-3.5 px-5 rounded-xl transition-all font-semibold text-left border border-transparent hover:border-red-200 hover:shadow-md group cursor-pointer"
 >
 <div className="w-10 h-10 rounded-lg bg-red-50 text-[#de2a15] flex items-center justify-center group-hover:bg-red-100 transition-colors">
 <Tv className="w-5 h-5" />
 </div>
 <div className="flex flex-col">
 <span className="text-[15px]">tvOS</span>
 <span className="text-[11px] text-slate-500 font-medium">Apple TV devices</span>
 </div>
 </button>
 </div>
 </div>
 </Modal>

 {/* Profile Form Modal */}
 <Modal
 title={null}
 open={isProfileFormModalVisible}
 onCancel={() => setIsProfileFormModalVisible(false)}
 footer={null}
 width={1000}
 className="custom-modal form-modal glass-modal-container"
 styles={{
 body: { padding: 0 }
 }}
 centered
 closeIcon={<div className="bg-white hover:bg-white p-2 rounded-full transition-all cursor-pointer"><X className="w-5 h-5 text-slate-700" /></div>}
 >
 <div className="flex flex-col h-full bg-white relative">
 {/* Header */}
 <div className="px-6 py-4 border-b border-slate-200 flex items-center gap-3 bg-slate-50">
 <div className="w-10 h-10 rounded-xl bg-white flex items-center justify-center shadow-sm border border-slate-200">
 <Apple className="w-6 h-6 text-slate-800" fill="currentColor" />
 </div>
 <div>
 <h2 className="text-lg font-bold text-slate-800 m-0 tracking-wide">CREATE PROFILE</h2>
 <p className="text-xs text-slate-500 font-medium m-0">Configure settings and restrictions for Apple devices</p>
 </div>
 </div>

 <Tabs
 defaultActiveKey="general"
 className="custom-tabs flex-1"
 items={[
 {
 key: "general",
 label: (
 <div className="flex items-center gap-2 px-4 font-semibold uppercase tracking-wider text-[13px]">
 <AlertCircle className="w-4 h-4 text-red-500" fill="currentColor" stroke="white" />
 GENERAL
 </div>
 ),
 children: (
 <div className="p-0 h-full overflow-hidden flex flex-col relative pb-16">
 
 <div className="flex-1 overflow-y-auto p-8 relative z-10 scrollbar-hide">
 <div className="max-w-[1200px] mx-auto pb-16">
 <Form layout="vertical" className="max-w-4xl mx-auto custom-form">
 <div className="flex flex-col md:flex-row gap-6 mb-6">
 <div className="w-full md:w-1/3">
 <div className="font-semibold text-slate-700 text-sm mb-1">Profile Name <span className="text-red-500">*</span></div>
 <div className="text-xs text-slate-500 mb-2">Tên hiển thị của cấu hình trên hệ thống</div>
 </div>
 <div className="w-full md:w-2/3">
 <Input 
 placeholder="Nhập tên cấu hình..." 
 value={newProfileName}
 onChange={(e) => setNewProfileName(e.target.value)}
 status={!newProfileName ? "error" : ""} 
 className="w-full h-10 rounded-md"
 />
 {!newProfileName && <div className="text-red-500 text-xs mt-1.5 font-medium">Tên cấu hình là bắt buộc và không được để trống.</div>}
 </div>
 </div>

 <div className="flex flex-col md:flex-row gap-6 mb-6 border-t border-slate-100 pt-6">
 <div className="w-full md:w-1/3">
 <div className="font-semibold text-slate-700 text-sm mb-1">Description</div>
 <div className="text-xs text-slate-500 mb-2">Mô tả chi tiết mục đích của cấu hình này (tuỳ chọn)</div>
 </div>
 <div className="w-full md:w-2/3">
 <Input.TextArea 
 placeholder="Nhập mô tả cấu hình..." 
 value={newProfileDesc}
 onChange={(e) => setNewProfileDesc(e.target.value)}
 rows={4}
 className="w-full rounded-md resize-none"
 />
 </div>
 </div>

 <div className="flex flex-col md:flex-row gap-6 mb-6 border-t border-slate-100 pt-6">
 <div className="w-full md:w-1/3">
 <div className="font-semibold text-slate-700 text-sm">System Info</div>
 <div className="text-xs text-slate-500 mt-1">Thông tin mặc định của hệ thống</div>
 </div>
 <div className="w-full md:w-2/3 grid grid-cols-2 gap-4">
 <div className="glass-card p-4 rounded-lg">
 <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Status</div>
 <div className="font-medium text-slate-800 flex items-center gap-2">
 <div className="w-2 h-2 rounded-full bg-slate-300"></div> N/A
 </div>
 </div>
 <div className="glass-card p-4 rounded-lg">
 <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Version</div>
 <div className="font-medium text-slate-800">1.0</div>
 </div>
 <div className="glass-card p-4 rounded-lg">
 <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Family</div>
 <div className="flex items-center gap-1.5 font-medium text-slate-800 cursor-pointer">
 <Apple className="w-4 h-4 text-slate-600" fill="currentColor" /> Apple
 </div>
 </div>
 <div className="glass-card p-4 rounded-lg">
 <div className="text-[11px] font-semibold text-slate-400 mb-1 uppercase tracking-wider">Type</div>
 <div className="flex items-center gap-1.5 font-medium text-slate-800 cursor-pointer">
 <Smartphone className="w-4 h-4 text-slate-600" /> iOS
 </div>
 </div>
 </div>
 </div>

 <div className="flex flex-col md:flex-row gap-6 border-t border-slate-100 pt-6 pb-12">
 <div className="w-full md:w-1/3">
 <div className="font-semibold text-slate-700 text-sm mb-1">Configurations</div>
 <div className="text-xs text-slate-500">Số lượng cấu hình đã thêm</div>
 </div>
 <div className="w-full md:w-2/3 flex items-center">
 <Tag className="px-4 py-1.5 text-sm border-slate-200 bg-slate-50 text-slate-700 font-medium rounded-md cursor-pointer hover:bg-blue-100 transition-colors m-0 flex items-center gap-2">
 <Settings2 className="w-4 h-4" /> {(hasPasscodeConfig ? 1 : 0) + (hasRestrictionsConfig ? 1 : 0) + (hasDomainsConfig ? 1 : 0) + (hasHttpProxyConfig ? 1 : 0) + (hasDnsProxyConfig ? 1 : 0) + (hasContentFilterConfig ? 1 : 0) + (hasCertificateTransparencyConfig ? 1 : 0) + (hasWifiConfig ? 1 : 0) + (hasVpnConfig ? 1 : 0) + (hasAirPlayConfig ? 1 : 0) + (hasAirPlaySecurityConfig ? 1 : 0) + (hasAirPrintConfig ? 1 : 0) + (hasCalendarConfig ? 1 : 0) + (hasContactsConfig ? 1 : 0) + (hasExchangeConfig ? 1 : 0) + (hasGoogleAccountConfig ? 1 : 0) + (hasLdapConfig ? 1 : 0)} Configurations Added
 </Tag>
 </div>
 </div>
 </Form>
 </div>
 </div>
 </div>
 ),
 },
 {
 key: "configurations",
 label: <div className="px-4 text-slate-500 font-medium uppercase tracking-wider text-[13px]">CONFIGURATIONS</div>,
 children: (
 <div className="p-0 h-full overflow-hidden flex flex-col relative pb-16">
 
 <div className="flex-1 overflow-y-auto p-8 relative z-10 scrollbar-hide">
 <div className="max-w-[1200px] mx-auto pb-16">
 {hasPasscodeConfig || hasRestrictionsConfig || hasDomainsConfig || hasHttpProxyConfig || hasDnsProxyConfig || hasContentFilterConfig || hasCertificateTransparencyConfig || hasWifiConfig || hasVpnConfig || hasAirPlayConfig || hasAirPlaySecurityConfig || hasAirPrintConfig || hasCalendarConfig || hasContactsConfig || hasExchangeConfig || hasGoogleAccountConfig || hasLdapConfig ? (
 <div className="flex flex-col h-full">
 <div className="flex items-center justify-between mb-4">
 <h3 className="text-[15px] font-bold text-slate-700 uppercase tracking-wide">PROFILE CONFIGURATION</h3>
 <Button 
 type="primary" 
 shape="circle"
 icon={<Plus className="w-5 h-5" strokeWidth={2.5} />} 
 className="bg-[#de2a15] hover:bg-[#c22412] flex items-center justify-center shadow-md w-10 h-10"
 onClick={() => setIsAddConfigModalVisible(true)}
 />
 </div>
 
 <div className="glass-card rounded-lg overflow-hidden">
 {/* Table Header */}
 <div className="grid grid-cols-12 gap-4 p-4 border-b border-slate-200 bg-white">
 <div className="col-span-4 font-bold text-[13px] text-slate-700 uppercase tracking-wider">NAME</div>
 <div className="col-span-8 font-bold text-[13px] text-slate-700 uppercase tracking-wider">DESCRIPTION</div>
 </div>
 
 {/* Table Row Passcode */}
 {hasPasscodeConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-white transition-colors border-b border-slate-200 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-white flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors border border-slate-200">
 <Lock className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsPasscodeConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Passcode
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Passcode
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasPasscodeConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Restrictions */}
 {hasRestrictionsConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Shield className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsRestrictionsConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Restrictions
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Restrictions
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasRestrictionsConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Domains */}
 {hasDomainsConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Globe className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsDomainsConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Domains
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Domains
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasDomainsConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Global HTTP Proxy */}
 {hasHttpProxyConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Globe className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsHttpProxyConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Global HTTP Proxy
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Global HTTP Proxy
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasHttpProxyConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row DNS Proxy */}
 {hasDnsProxyConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Server className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsDnsProxyConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 DNS Proxy
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 DNS Proxy
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasDnsProxyConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Content Filter */}
 {hasContentFilterConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Filter className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsContentFilterConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Content Filter
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Content Filter
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasContentFilterConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Certificate Transparency */}
 {hasCertificateTransparencyConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <CheckCircle2 className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsCertificateTransparencyConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Certificate Transparency
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Certificate Transparency
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasCertificateTransparencyConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Wi-Fi */}
 {hasWifiConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Wifi className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsWifiConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Wi-Fi
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Wi-Fi
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasWifiConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row VPN */}
 {hasVpnConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Server className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsVpnConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 VPN
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 VPN
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasVpnConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row AirPlay */}
 {hasAirPlayConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <MonitorUp className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsAirPlayConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 AirPlay
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 AirPlay
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasAirPlayConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row AirPlay Security */}
 {hasAirPlaySecurityConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Shield className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsAirPlaySecurityConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 AirPlay Security
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 AirPlay Security
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasAirPlaySecurityConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row AirPrint */}
 {hasAirPrintConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Printer className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsAirPrintConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 AirPrint
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 AirPrint
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasAirPrintConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Calendar */}
 {hasCalendarConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Calendar className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsCalendarConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Calendar
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Calendar
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasCalendarConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Contacts */}
 {hasContactsConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Contact className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsContactsConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Contacts
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Contacts
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasContactsConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Exchange ActiveSync */}
 {hasExchangeConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Mail className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsExchangeConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Exchange ActiveSync
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Exchange ActiveSync
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasExchangeConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Google Account */}
 {hasGoogleAccountConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Mail className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsGoogleAccountConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Google Account
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Google Account
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasGoogleAccountConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row LDAP */}
 {hasLdapConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Key className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsLdapConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 LDAP
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 LDAP
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasLdapConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Mail */}
 {hasMailConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Mail className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsMailConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Mail
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Mail
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasMailConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row macOS Server Account */}
 {hasMacOsServerConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Server className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsMacOsServerConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 macOS Server Account
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 macOS Server Account
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasMacOsServerConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row SCEP */}
 {hasScepConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Key className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsScepConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 SCEP
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 SCEP
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasScepConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Cellular */}
 {hasCellularConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Radio className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsCellularConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Cellular
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Cellular
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasCellularConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Notifications */}
 {hasNotificationsConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <AlertCircle className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsNotificationsConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Notifications
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Notifications
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasNotificationsConfig(false)}
 />
 </div>
 </div>
 )}
 {/* Table Row Conference Room Display */}
 {hasConferenceRoomConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Monitor className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsConferenceRoomConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Conference Room Display
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Conference Room Display
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasConferenceRoomConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row TV Remote */}
 {hasTvRemoteConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Tv className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsTvRemoteConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 TV Remote
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 TV Remote
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasTvRemoteConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Lock Screen Message */}
 {hasLockScreenMessageConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <MessageSquare className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsLockScreenMessageConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Lock Screen Message
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Lock Screen Message
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasLockScreenMessageConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Web Clip */}
 {hasWebClipConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <LinkIcon className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsWebClipConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Web Clip
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Web Clip
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasWebClipConfig(false)}
 />
 </div>
 </div>
 )}

 {/* Table Row Subscribed Calendar */}
 {hasSubscribedCalendarConfig && (
 <div className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0 cursor-pointer group">
 <div className="col-span-4 flex items-center gap-3">
 <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-red-50 group-hover:text-[#de2a15] transition-colors">
 <Calendar className="w-4 h-4" />
 </div>
 <button 
 onClick={() => setIsSubscribedCalendarConfigVisible(true)}
 className="text-slate-800 hover:text-[#de2a15] font-medium text-[15px] text-left transition-colors"
 >
 Subscribed Calendar
 </button>
 </div>
 <div className="col-span-7 text-slate-600 text-[15px]">
 Subscribed Calendar
 </div>
 <div className="col-span-1 flex justify-end">
 <Button 
 type="text" 
 shape="circle" 
 icon={<Trash2 className="w-5 h-5" />} 
 className="text-slate-400 hover:text-red-500 hover:bg-red-50 flex items-center justify-center"
 onClick={() => setHasSubscribedCalendarConfig(false)}
 />
 </div>
 </div>
 )}
 </div>
 </div>
 ) : (
 <div className="flex flex-col items-center justify-center py-20 text-center glass-card rounded-xl border border-slate-200 border-dashed">
 <div className="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mb-4 border border-blue-100">
 <Settings2 className="w-8 h-8 text-blue-500" />
 </div>
 <h3 className="text-lg font-semibold text-slate-800 mb-2">Chưa có cấu hình nào</h3>
 <p className="text-slate-500 max-w-md mb-6">Thêm các cấu hình thiết bị như Wifi, VPN, Passcode để áp dụng cho các thiết bị thuộc Profile này.</p>
 <Button 
 type="primary" 
 icon={<Plus className="w-4 h-4" />} 
 className="bg-[#de2a15] hover:bg-[#c22412] text-white border-none transition-colors h-10 px-6 font-medium"
 onClick={() => setIsAddConfigModalVisible(true)}
 >
 THÊM CẤU HÌNH (ADD CONFIGURATION)
 </Button>
 </div>
 )}
 </div>
 </div>
 </div>
 ),
 }
 ]}
 />

 {/* Form Footer Actions */}
 <div className="absolute bottom-0 left-0 right-0 p-4 bg-white border-t border-slate-200 flex justify-end gap-3 z-50">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
onClick={() => {
 setIsProfileFormModalVisible(false);
 setEditingProfileId(null);
 setEditingProfileSnapshot(null);
}}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors"
onClick={handleSaveProfile}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Add Configuration Modal */}
 <Modal
 title={<div className="flex items-center gap-2 text-slate-800 font-semibold"><Plus className="w-5 h-5 text-slate-700" /> THÊM CẤU HÌNH (ADD CONFIGURATION)</div>}
 open={isAddConfigModalVisible}
 onCancel={() => setIsAddConfigModalVisible(false)}
 footer={null}
 width={1000}
 className="custom-modal"
 styles={{
 body: { padding: '24px', background: 'transparent' }
 }}
 >
 <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
 {/* Security & Restrictions */}
 <div className="glass-card rounded-xl overflow-hidden h-fit">
  <div className="bg-orange-500/80 px-4 py-3 text-white font-medium flex items-center gap-2">
   <Shield className="w-4 h-4" /> Security & Restrictions
  </div>
  <div className="p-2 flex flex-col gap-1">
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsPasscodeConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-orange-100 rounded-lg text-slate-700 hover:text-orange-700 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Lock className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> Passcode
   </button>
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsRestrictionsConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-orange-100 rounded-lg text-slate-700 hover:text-orange-700 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Shield className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> Restrictions
   </button>
   {/* SCEP - not supported */}
   <div className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
    <Key className="w-4 h-4" /> SCEP
    <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
   </div>
   {/* Certificate Transparency - not supported */}
   <div className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
    <CheckCircle2 className="w-4 h-4" /> Certificate Transparency
    <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
   </div>
  </div>
 </div>

 {/* Network & Connectivity */}
 <div className="glass-card rounded-xl overflow-hidden h-fit">
  <div className="bg-slate-500 px-4 py-3 text-white font-medium flex items-center gap-2">
   <Globe className="w-4 h-4" /> Network & Connectivity
  </div>
  <div className="p-2 flex flex-col gap-1">
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsWifiConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-slate-200 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Wifi className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> Wi-Fi
   </button>
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsVpnConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-slate-200 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Server className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> VPN
   </button>
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsHttpProxyConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-slate-200 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Globe className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> Global HTTP Proxy
   </button>
   <button
    onClick={() => { setIsAddConfigModalVisible(false); setIsContentFilterConfigVisible(true); }}
    className="flex items-center gap-3 w-full p-2.5 hover:bg-slate-200 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-md group"
   >
    <Filter className="w-4 h-4 text-slate-400 group-hover:text-current transition-colors" /> Content Filter
   </button>
   {/* Cellular - not supported */}
   <div className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
    <Radio className="w-4 h-4" /> Cellular
    <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
   </div>
   {/* DNS Proxy - not supported */}
   <div className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
    <Server className="w-4 h-4" /> DNS Proxy
    <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
   </div>
   {/* Domains - not supported */}
   <div className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
    <Globe className="w-4 h-4" /> Domains
    <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
   </div>
  </div>
 </div>

 {/* Accounts & Mail - all not supported */}
 <div className="glass-card rounded-xl overflow-hidden h-fit">
  <div className="bg-rose-500/80 px-4 py-3 text-white font-medium flex items-center gap-2">
   <Mail className="w-4 h-4" /> Accounts & Mail
   <span className="ml-auto text-[10px] bg-white/20 text-white px-2 py-0.5 rounded font-normal">Coming Soon</span>
  </div>
  <div className="p-2 flex flex-col gap-1">
   {[
    { icon: <Mail className="w-4 h-4" />, label: "Mail" },
    { icon: <RefreshCcw className="w-4 h-4" />, label: "Exchange ActiveSync" },
    { icon: <Globe className="w-4 h-4" />, label: "Google Account" },
    { icon: <Server className="w-4 h-4" />, label: "macOS Server Account" },
    { icon: <Users className="w-4 h-4" />, label: "LDAP" },
    { icon: <Contact className="w-4 h-4" />, label: "Contacts" },
    { icon: <Calendar className="w-4 h-4" />, label: "Calendar" },
    { icon: <Calendar className="w-4 h-4" />, label: "Subscribed Calendars" },
   ].map((item) => (
    <div key={item.label} className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
     {item.icon} {item.label}
     <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
    </div>
   ))}
  </div>
 </div>

 {/* Media & Display - all not supported */}
 <div className="glass-card rounded-xl overflow-hidden h-fit">
  <div className="bg-teal-500/80 px-4 py-3 text-white font-medium flex items-center gap-2">
   <Monitor className="w-4 h-4" /> Media & Display
   <span className="ml-auto text-[10px] bg-white/20 text-white px-2 py-0.5 rounded font-normal">Coming Soon</span>
  </div>
  <div className="p-2 flex flex-col gap-1">
   {[
    { icon: <MonitorPlay className="w-4 h-4" />, label: "AirPlay" },
    { icon: <Shield className="w-4 h-4" />, label: "AirPlay Security" },
    { icon: <Printer className="w-4 h-4" />, label: "AirPrint" },
    { icon: <Tv className="w-4 h-4" />, label: "TV Remote" },
    { icon: <Monitor className="w-4 h-4" />, label: "Conference Room Display" },
    { icon: <MessageSquare className="w-4 h-4" />, label: "Lock Screen Message" },
    { icon: <LinkIcon className="w-4 h-4" />, label: "Web Clips" },
    { icon: <AlertCircle className="w-4 h-4" />, label: "Notifications" },
   ].map((item) => (
    <div key={item.label} className="flex items-center gap-3 w-full p-2.5 rounded-lg text-slate-400 cursor-not-allowed opacity-50 text-sm text-left font-medium select-none" title="Not supported by backend yet">
     {item.icon} {item.label}
     <span className="ml-auto text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Soon</span>
    </div>
   ))}
  </div>
 </div>
 </div>
 </Modal>

 {/* Passcode Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Lock className="w-5 h-5" /> 
 PASSCODE CONFIGURATION
 </div>
 }
 open={isPasscodeConfigVisible}
 onCancel={() => setIsPasscodeConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form form={passcodeForm} layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Allow simple value — not in backend schema */}
 <div className="relative flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="allowSimple" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox defaultChecked />
  </Form.Item>
  <div className="flex-1">
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
    Allow simple value
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
   <div className="text-sm text-slate-500 leading-relaxed">Permit the use of repeating, ascending, and descending character sequences</div>
  </div>
 </div>

 {/* Require alphanumeric value ✅ backend: require_alphanumeric */}
 <div className="flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm">
  <Form.Item name="requireAlphanumeric" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox />
  </Form.Item>
  <div>
   <div className="font-semibold text-slate-800 text-base mb-1">Require alphanumeric value</div>
   <div className="text-sm text-slate-500 leading-relaxed">Requires passcode to contain at least one letter and one number</div>
  </div>
 </div>

 {/* Minimum passcode length ✅ backend: min_length */}
 <div className="flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm">
  <Form.Item name="minLength" className="mb-0 w-24">
   <InputNumber min={0} max={16} defaultValue={0} className="w-full" controls={true} />
  </Form.Item>
  <div>
   <div className="font-semibold text-slate-800 text-base mb-1">Minimum passcode length</div>
   <div className="text-sm text-slate-500 leading-relaxed">Smallest number of passcode characters allowed</div>
  </div>
 </div>

 {/* Minimum number of complex characters — not in backend schema */}
 <div className="relative flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="minComplexChars" className="mb-0 w-24">
   <InputNumber min={0} max={4} defaultValue={0} className="w-full" controls={true} />
  </Form.Item>
  <div className="flex-1">
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
    Minimum number of complex characters
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
   <div className="text-sm text-slate-500 leading-relaxed">Smallest number of non-alphanumeric characters allowed</div>
  </div>
 </div>

 {/* Maximum passcode age — not in backend schema */}
 <div className="relative flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="maxAge" className="mb-0 w-24">
   <InputNumber min={1} max={730} placeholder="none" className="w-full" controls={true} />
  </Form.Item>
  <div className="flex-1">
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
    Maximum passcode age (1-730 days, or none)
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
   <div className="text-sm text-slate-500 leading-relaxed">Days after which passcode must be changed</div>
  </div>
 </div>

 {/* Maximum Auto-Lock ✅ backend: auto_lock */}
 <div className="flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm">
  <Form.Item name="maxAutoLock" className="mb-0 w-32">
   <Select defaultValue="none" className="cursor-pointer w-full">
    <Select.Option value="none">None</Select.Option>
    <Select.Option value="1">1 minute</Select.Option>
    <Select.Option value="2">2 minutes</Select.Option>
    <Select.Option value="3">3 minutes</Select.Option>
    <Select.Option value="4">4 minutes</Select.Option>
    <Select.Option value="5">5 minutes</Select.Option>
   </Select>
  </Form.Item>
  <div>
   <div className="font-semibold text-slate-800 text-base mb-1">Maximum Auto-Lock</div>
   <div className="text-sm text-slate-500 leading-relaxed">Longest auto-lock time available to the user</div>
  </div>
 </div>

 {/* Passcode history — not in backend schema */}
 <div className="relative flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="passcodeHistory" className="mb-0 w-24">
   <InputNumber min={1} max={50} placeholder="none" className="w-full" controls={true} />
  </Form.Item>
  <div className="flex-1">
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
    Passcode history (1-50 passcodes, or none)
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
   <div className="text-sm text-slate-500 leading-relaxed">Number of unique passcodes before reuse</div>
  </div>
 </div>

 {/* Maximum grace period — not in backend schema */}
 <div className="relative flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="gracePeriod" className="mb-0 w-32">
   <Select defaultValue="none" className="cursor-pointer w-full">
    <Select.Option value="none">None</Select.Option>
    <Select.Option value="immediate">Immediately</Select.Option>
    <Select.Option value="1">1 minute</Select.Option>
    <Select.Option value="5">5 minutes</Select.Option>
    <Select.Option value="15">15 minutes</Select.Option>
   </Select>
  </Form.Item>
  <div className="flex-1">
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
    Maximum grace period for device lock
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
   <div className="text-sm text-slate-500 leading-relaxed">Longest device lock grace period available to the user</div>
  </div>
 </div>

 {/* Maximum number of failed attempts ✅ backend: retry_limit */}
 <div className="flex items-start gap-4 p-4 bg-white rounded-lg border border-slate-200 shadow-sm">
  <Form.Item name="maxFailedAttempts" className="mb-0 w-24">
   <InputNumber min={2} max={11} placeholder="none" className="w-full" controls={true} />
  </Form.Item>
  <div>
   <div className="font-semibold text-slate-800 text-base mb-1">Maximum number of failed attempts</div>
   <div className="text-sm text-slate-500 leading-relaxed">Number of passcode entry attempts allowed before all data on device will be erased</div>
  </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsPasscodeConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
  const values = passcodeForm.getFieldsValue();
  setPasscodeData({
   allow_simple: values.allowSimple ?? true,
   require_alphanumeric: values.requireAlphanumeric ?? false,
   min_length: values.minLength ?? 0,
   min_complex_chars: values.minComplexChars ?? 0,
   max_passcode_age: values.maxPasscodeAge,
   auto_lock: values.autoLock ?? "none",
   history: values.passcodeHistory,
   grace_period: values.gracePeriod ?? "none",
   retry_limit: values.maxFailedAttempts,
  });
  setHasPasscodeConfig(true);
  setIsPasscodeConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Restrictions Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Shield className="w-5 h-5" /> 
 RESTRICTIONS CONFIGURATION
 </div>
 }
 open={isRestrictionsConfigVisible}
 onCancel={() => setIsRestrictionsConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <Tabs
 defaultActiveKey="functionality"
 className="custom-tabs"
 items={[
 {
 key: "functionality",
 label: <div className="px-4 font-semibold uppercase tracking-wider text-[13px]">Functionality</div>,
 children: (
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form form={restrictionsForm} layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-4">
 {/* Allow Camera ✅ backend: camera_enabled */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm">
  <Form.Item name="allowCamera" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox defaultChecked />
  </Form.Item>
  <div className="font-semibold text-slate-800 text-sm pt-1">Allow use of camera</div>
 </div>

 {/* Allow FaceTime — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm ml-8 opacity-50 pointer-events-none select-none">
  <Form.Item name="allowFaceTime" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox defaultChecked />
  </Form.Item>
  <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
   Allow FaceTime (supervised only)
   <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
  </div>
 </div>

 {/* Screenshots — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
  <Form.Item name="allowScreenshots" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox defaultChecked />
  </Form.Item>
  <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
   Allow screenshots and screen recording
   <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
  </div>
 </div>
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm ml-8 opacity-50 pointer-events-none select-none">
  <Form.Item name="allowAirPlayScreen" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox defaultChecked />
  </Form.Item>
  <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
   Allow AirPlay, View Screen by Classroom, and Screen Sharing
   <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
  </div>
 </div>
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm ml-16 opacity-50 pointer-events-none select-none">
  <Form.Item name="allowClassroomPrompt" valuePropName="checked" className="mb-0 pt-1">
   <Checkbox />
  </Form.Item>
  <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
   Allow Classroom to perform AirPlay and View Screen without prompting (supervised only)
   <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
  </div>
 </div>

 {/* Allow AirDrop ✅ backend: airdrop_enabled | rest not supported */}
 {[
  { name: "misc_0", label: "Allow AirDrop (supervised only)", supported: true },
  { name: "misc_1", label: "Allow iMessage (supervised only)", supported: false },
  { name: "misc_2", label: "Allow Apple Music (supervised only)", supported: false },
  { name: "misc_3", label: "Allow Radio (supervised only)", supported: false },
  { name: "misc_4", label: "Allow Live Voicemail (supervised only)", supported: false },
  { name: "misc_5", label: "Allow voice dialing while device is locked (deprecated in iOS 17)", supported: false },
 ].map((item) => (
  <div key={item.name} className={`flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm${!item.supported ? " opacity-50 pointer-events-none select-none" : ""}`}>
   <Form.Item name={item.name} valuePropName="checked" className="mb-0 pt-1">
    <Checkbox defaultChecked />
   </Form.Item>
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
    {item.label}
    {!item.supported && <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>}
   </div>
  </div>
 ))}

 {/* Siri Group — not in backend */}
 {[
  { name: "allowSiri", label: "Allow Siri", indent: "" },
  { name: "allowSiriLocked", label: "Allow Siri while device is locked", indent: " ml-8" },
  { name: "enableSiriFilter", label: "Enable Siri profanity filter (supervised only)", indent: " ml-8" },
  { name: "showUserContentSiri", label: "Show user-generated content in Siri (supervised only)", indent: " ml-8" },
  { name: "allowSiriSuggestions", label: "Allow Siri Suggestions", indent: "" },
 ].map((item) => (
  <div key={item.name} className={`flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none${item.indent}`}>
   <Form.Item name={item.name} valuePropName="checked" className="mb-0 pt-1">
    <Checkbox defaultChecked />
   </Form.Item>
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
    {item.label}
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
  </div>
 ))}
 </div>
 </Form>
 </div>
 ),
 },
 {
 key: "apps",
 label: <div className="px-4 font-semibold uppercase tracking-wider text-[13px]">Apps</div>,
 children: (
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-4">
 {/* iTunes — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <Form.Item name="allowiTunes" valuePropName="checked" className="mb-0 pt-1">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
 Allow use of iTunes Store (supervised only)
 <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 </div>
 {/* News — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <Form.Item name="allowNews" valuePropName="checked" className="mb-0 pt-1">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
 Allow use of News (supervised only)
 <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 </div>
 {/* Podcasts — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <Form.Item name="allowPodcasts" valuePropName="checked" className="mb-0 pt-1">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
 Allow use of Podcasts (supervised only)
 <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 </div>

 {/* Game Center — not in backend */}
 {[
  { name: "allowGameCenter", label: "Allow use of Game Center (supervised only)", indent: "" },
  { name: "allowMultiplayer", label: "Allow multiplayer gaming (supervised only)", indent: " ml-8" },
  { name: "allowAddFriends", label: "Allow adding Game Center friends (supervised only)", indent: " ml-8" },
 ].map((item) => (
  <div key={item.name} className={`flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none${item.indent}`}>
   <Form.Item name={item.name} valuePropName="checked" className="mb-0 pt-1">
    <Checkbox defaultChecked />
   </Form.Item>
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
    {item.label}
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
  </div>
 ))}

 {/* Safari — not in backend */}
 {[
  { name: "allowSafari", label: "Allow use of Safari (supervised only)", indent: "" },
  { name: "enableAutoFill", label: "Enable AutoFill (supervised only)", indent: " ml-8" },
  { name: "forceFraudWarning", label: "Force fraud warning", indent: " ml-8" },
  { name: "enableJS", label: "Enable JavaScript", indent: " ml-8" },
  { name: "blockPopups", label: "Block pop-ups", indent: " ml-8" },
 ].map((item) => (
  <div key={item.name} className={`flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none${item.indent}`}>
   <Form.Item name={item.name} valuePropName="checked" className="mb-0 pt-1">
    <Checkbox defaultChecked />
   </Form.Item>
   <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
    {item.label}
    <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
   </div>
  </div>
 ))}
 {/* Accept cookies — not in backend */}
 <div className="ml-8 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm mb-2">
  Accept cookies
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <Select defaultValue="always" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="always">Always</Select.Option>
 <Select.Option value="never">Never</Select.Option>
 <Select.Option value="visited">From visited sites</Select.Option>
 </Select>
 </div>

 {/* Restrict App Usage — not in backend */}
 <div className="mt-8 p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm mb-2">
  Restrict App Usage (supervised only)
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <Select defaultValue="allowAll" className="cursor-pointer w-full max-w-xs mb-4 cursor-pointer">
 <Select.Option value="allowAll">Allow all apps</Select.Option>
 <Select.Option value="allowSome">Allow some apps</Select.Option>
 <Select.Option value="hideSome">Hide some apps</Select.Option>
 </Select>
 <div className="h-32 bg-slate-50 border border-slate-200 rounded flex flex-col items-center justify-center text-slate-400">
 <div className="flex gap-2">
 <Button size="small" icon={<Plus className="w-3 h-3" />} />
 <Button size="small" disabled>-</Button>
 </div>
 </div>
 </div>
 </div>
 </Form>
 </div>
 ),
 },
 {
 key: "media",
 label: <div className="px-4 font-semibold uppercase tracking-wider text-[13px]">Media Content</div>,
 children: (
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Ratings region — not in backend */}
 <div className="p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm mb-1">
  Ratings region
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-xs text-slate-500 mb-3">Sets the region for the ratings</div>
 <Select defaultValue="us" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="us">United States</Select.Option>
 <Select.Option value="vn">Vietnam</Select.Option>
 <Select.Option value="uk">United Kingdom</Select.Option>
 </Select>
 </div>

 {/* Allowed content ratings — not in backend */}
 <div className="p-4 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm mb-1">
  Allowed content ratings
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-xs text-slate-500 mb-4">Sets the maximum allowed ratings</div>

 <div className="grid grid-cols-[100px_1fr] gap-4 items-center mb-3">
 <div className="text-right font-medium text-slate-700 text-sm">Movies:</div>
 <Select defaultValue="all" className="cursor-pointer max-w-xs cursor-pointer">
 <Select.Option value="all">Allow All Movies</Select.Option>
 <Select.Option value="none">Don&apos;t Allow Movies</Select.Option>
 </Select>
 </div>
 <div className="grid grid-cols-[100px_1fr] gap-4 items-center mb-3">
 <div className="text-right font-medium text-slate-700 text-sm">TV Shows:</div>
 <Select defaultValue="all" className="cursor-pointer max-w-xs cursor-pointer">
 <Select.Option value="all">Allow All TV Shows</Select.Option>
 <Select.Option value="none">Don&apos;t Allow TV Shows</Select.Option>
 </Select>
 </div>
 <div className="grid grid-cols-[100px_1fr] gap-4 items-center">
 <div className="text-right font-medium text-slate-700 text-sm">Apps:</div>
 <Select defaultValue="all" className="cursor-pointer max-w-xs cursor-pointer">
 <Select.Option value="all">Allow All Apps</Select.Option>
 <Select.Option value="none">Don&apos;t Allow Apps</Select.Option>
 </Select>
 </div>
 </div>

 {/* Explicit media — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <Form.Item name="allowExplicitMedia" valuePropName="checked" className="mb-0 pt-1">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
  Allow playback of explicit music, podcasts &amp; iTunes U media (supervised only)
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 </div>

 {/* Explicit books — not in backend */}
 <div className="flex items-start gap-4 p-3 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <Form.Item name="allowExplicitBooks" valuePropName="checked" className="mb-0 pt-1">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-sm pt-1">
  Allow explicit sexual content in Apple Books
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 </div>
 </div>
 </Form>
 </div>
 ),
 }
 ]}
 />

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsRestrictionsConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
  const values = restrictionsForm.getFieldsValue();
  setRestrictionsData({
   camera_enabled: values.allowCamera ?? true,
   airdrop_enabled: values.allowAirDrop ?? true,
   bluetooth_enabled: true,
   usb_debugging_enabled: false,
   external_app_install_allowed: values.allowiTunes ?? true,
   facetime_enabled: values.allowFaceTime ?? true,
   screenshots_enabled: values.allowScreenshots ?? true,
   siri_enabled: values.allowSiri ?? true,
   safari_enabled: values.allowSafari ?? true,
   game_center_enabled: values.allowGameCenter ?? true,
   itunes_enabled: values.allowiTunes ?? true,
   news_enabled: values.allowNews ?? true,
   podcasts_enabled: values.allowPodcasts ?? true,
  });
  setHasRestrictionsConfig(true);
  setIsRestrictionsConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Domains Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Globe className="w-5 h-5" /> 
 DOMAINS CONFIGURATION
 </div>
 }
 open={isDomainsConfigVisible}
 onCancel={() => setIsDomainsConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-8">
 {/* Unmarked Email Domains */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Unmarked Email Domains</div>
 <div className="text-sm text-slate-500 mb-4">Email addresses not matching any of these domains will be marked in Mail</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {unmarkedEmailDomains.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 unmarkedEmailDomains.map((domain, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedEmailDomainIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedEmailDomainIdx(idx)}
 >
 <Input 
 value={domain}
 onChange={(e) => {
 const newDomains = [...unmarkedEmailDomains];
 newDomains[idx] = e.target.value;
 setUnmarkedEmailDomains(newDomains);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 placeholder="example.com"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space to keep the look consistent */}
 {unmarkedEmailDomains.length > 0 && unmarkedEmailDomains.length < 4 && (
 Array.from({ length: 4 - unmarkedEmailDomains.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + unmarkedEmailDomains.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setUnmarkedEmailDomains([...unmarkedEmailDomains, ""]);
 setSelectedEmailDomainIdx(unmarkedEmailDomains.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedEmailDomainIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedEmailDomainIdx === null}
 onClick={() => {
 if (selectedEmailDomainIdx !== null) {
 const newDomains = unmarkedEmailDomains.filter((_, idx) => idx !== selectedEmailDomainIdx);
 setUnmarkedEmailDomains(newDomains);
 setSelectedEmailDomainIdx(null);
 }
 }}
 />
 </div>
 </div>

 {/* Managed Safari Web Domains */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Managed Safari Web Domains</div>
 <div className="text-sm text-slate-500 mb-4">URL patterns of domains from which documents will be considered managed</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {managedSafariDomains.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 managedSafariDomains.map((domain, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedSafariDomainIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedSafariDomainIdx(idx)}
 >
 <Input 
 value={domain}
 onChange={(e) => {
 const newDomains = [...managedSafariDomains];
 newDomains[idx] = e.target.value;
 setManagedSafariDomains(newDomains);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 placeholder="example.com"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {managedSafariDomains.length > 0 && managedSafariDomains.length < 4 && (
 Array.from({ length: 4 - managedSafariDomains.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + managedSafariDomains.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setManagedSafariDomains([...managedSafariDomains, ""]);
 setSelectedSafariDomainIdx(managedSafariDomains.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedSafariDomainIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedSafariDomainIdx === null}
 onClick={() => {
 if (selectedSafariDomainIdx !== null) {
 const newDomains = managedSafariDomains.filter((_, idx) => idx !== selectedSafariDomainIdx);
 setManagedSafariDomains(newDomains);
 setSelectedSafariDomainIdx(null);
 }
 }}
 />
 </div>
 </div>

 {/* Safari Password Autofill Domains */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Safari Password Autofill Domains (supervised only)</div>
 <div className="text-sm text-slate-500 mb-4">URL patterns of websites for which passwords can be saved</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {safariPasswordDomains.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 safariPasswordDomains.map((domain, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedPasswordDomainIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedPasswordDomainIdx(idx)}
 >
 <Input 
 value={domain}
 onChange={(e) => {
 const newDomains = [...safariPasswordDomains];
 newDomains[idx] = e.target.value;
 setSafariPasswordDomains(newDomains);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 placeholder="example.com"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {safariPasswordDomains.length > 0 && safariPasswordDomains.length < 4 && (
 Array.from({ length: 4 - safariPasswordDomains.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + safariPasswordDomains.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setSafariPasswordDomains([...safariPasswordDomains, ""]);
 setSelectedPasswordDomainIdx(safariPasswordDomains.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedPasswordDomainIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedPasswordDomainIdx === null}
 onClick={() => {
 if (selectedPasswordDomainIdx !== null) {
 const newDomains = safariPasswordDomains.filter((_, idx) => idx !== selectedPasswordDomainIdx);
 setSafariPasswordDomains(newDomains);
 setSelectedPasswordDomainIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsDomainsConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasDomainsConfig(true);
 setIsDomainsConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Global HTTP Proxy Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Globe className="w-5 h-5" /> 
 GLOBAL HTTP PROXY CONFIGURATION
 </div>
 }
 open={isHttpProxyConfigVisible}
 onCancel={() => setIsHttpProxyConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Proxy Type */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-2">Proxy Type</div>
 <Select defaultValue="manual" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="manual">Manual</Select.Option>
 <Select.Option value="auto">Auto</Select.Option>
 </Select>
 </div>

 {/* Proxy Server and Port */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Proxy Server and Port</div>
 <div className="text-sm text-slate-500 mb-3">Host name or IP address, and port number for the proxy server</div>
 <div className="flex items-center gap-2">
 <Input placeholder="[required]" className="flex-[3] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 <span className="text-slate-600 font-bold mx-1">:</span>
 <Input placeholder="[port]" className="flex-1 min-w-[80px] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>

 {/* User Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">User Name</div>
 <div className="text-sm text-slate-500 mb-3">User name used to connect to the proxy</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Password */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Password</div>
 <div className="text-sm text-slate-500 mb-3">Password used to authenticate with the proxy</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Allow bypassing proxy */}
 <div className="flex items-start gap-4 p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <Form.Item name="allowBypassing" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox />
 </Form.Item>
 <div className="font-semibold text-slate-800 text-base">Allow bypassing proxy to access captive networks</div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsHttpProxyConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasHttpProxyConfig(true);
 setIsHttpProxyConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* DNS Proxy Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Server className="w-5 h-5" /> 
 DNS PROXY CONFIGURATION
 </div>
 }
 open={isDnsProxyConfigVisible}
 onCancel={() => setIsDnsProxyConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* App Bundle ID */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">App Bundle ID</div>
 <div className="text-sm text-slate-500 mb-3">Bundle identifier of the app containing the DNS proxy network extention</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Provider Bundle ID */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Provider Bundle ID</div>
 <div className="text-sm text-slate-500 mb-3">Bundle identifier of the DNS proxy network extension to use</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Provider Configuration */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Provider Configuration</div>
 <div className="text-sm text-slate-500 mb-3">Vendor specific configuration values</div>
 <Input.TextArea 
 placeholder="[optional]" 
 rows={10} 
 className="bg-slate-50 hover:bg-white focus:bg-white transition-colors resize-none"
 />
 <div className="flex gap-3 mt-4">
 <Button className="bg-slate-800 hover:bg-slate-700 text-white border-none rounded-md px-4">Upload File...</Button>
 <Button className="bg-slate-800 hover:bg-slate-700 text-white border-none rounded-md px-4">Remove Configuration</Button>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsDnsProxyConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasDnsProxyConfig(true);
 setIsDnsProxyConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Content Filter Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Filter className="w-5 h-5" /> 
 CONTENT FILTER CONFIGURATION
 </div>
 }
 open={isContentFilterConfigVisible}
 onCancel={() => setIsContentFilterConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Filter Type */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-2">Filter Type</div>
 <Select 
 value={contentFilterType}
 onChange={(value) => setContentFilterType(value)}
 className="w-full max-w-sm"
 >
 <Select.Option value="limit-adult">Built-in: Limit Adult Content</Select.Option>
 <Select.Option value="specific-websites">Built-in: Specific Websites Only</Select.Option>
 <Select.Option value="plugin">Plugin (Third Party App)</Select.Option>
 </Select>
 </div>

 {/* Dynamic Content based on Filter Type */}
 {contentFilterType === 'limit-adult' && (
 <div className="space-y-6">
 {/* Allowed URLs */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Allowed URLs</div>
 <div className="text-sm text-slate-500 mb-4">Specific URLs that will be allowed</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {allowedUrls.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 allowedUrls.map((url, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedAllowedUrlIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAllowedUrlIdx(idx)}
 >
 <Input 
 value={url}
 onChange={(e) => {
 const newUrls = [...allowedUrls];
 newUrls[idx] = e.target.value;
 setAllowedUrls(newUrls);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 ))
 )}
 {allowedUrls.length > 0 && allowedUrls.length < 4 && (
 Array.from({ length: 4 - allowedUrls.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + allowedUrls.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAllowedUrls([...allowedUrls, ""]);
 setSelectedAllowedUrlIdx(allowedUrls.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAllowedUrlIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAllowedUrlIdx === null}
 onClick={() => {
 if (selectedAllowedUrlIdx !== null) {
 const newUrls = allowedUrls.filter((_, idx) => idx !== selectedAllowedUrlIdx);
 setAllowedUrls(newUrls);
 setSelectedAllowedUrlIdx(null);
 }
 }}
 />
 </div>
 </div>

 {/* Unallowed URLs */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Unallowed URLs</div>
 <div className="text-sm text-slate-500 mb-4">Additional URLs that will not be allowed</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {unallowedUrls.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 unallowedUrls.map((url, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedUnallowedUrlIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedUnallowedUrlIdx(idx)}
 >
 <Input 
 value={url}
 onChange={(e) => {
 const newUrls = [...unallowedUrls];
 newUrls[idx] = e.target.value;
 setUnallowedUrls(newUrls);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 ))
 )}
 {unallowedUrls.length > 0 && unallowedUrls.length < 4 && (
 Array.from({ length: 4 - unallowedUrls.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + unallowedUrls.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setUnallowedUrls([...unallowedUrls, ""]);
 setSelectedUnallowedUrlIdx(unallowedUrls.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedUnallowedUrlIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedUnallowedUrlIdx === null}
 onClick={() => {
 if (selectedUnallowedUrlIdx !== null) {
 const newUrls = unallowedUrls.filter((_, idx) => idx !== selectedUnallowedUrlIdx);
 setUnallowedUrls(newUrls);
 setSelectedUnallowedUrlIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 )}

 {contentFilterType === 'specific-websites' && (
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Specific Websites</div>
 <div className="text-sm text-slate-500 mb-4">Allowed domains which will be shown as bookmarks</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="grid grid-cols-2 bg-slate-100 border-b border-slate-200 font-semibold text-slate-700 text-xs px-2 py-2">
 <div>URL</div>
 <div className="border-l border-slate-300 pl-2">Name</div>
 </div>
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {specificWebsites.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 specificWebsites.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 grid grid-cols-2 cursor-pointer transition-colors ${selectedSpecificWebsiteIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedSpecificWebsiteIdx(idx)}
 >
 <div className="px-2">
 <Input 
 value={item.url}
 onChange={(e) => {
 const newSites = [...specificWebsites];
 newSites[idx].url = e.target.value;
 setSpecificWebsites(newSites);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="px-2 border-l border-slate-200">
 <Input 
 value={item.name}
 onChange={(e) => {
 const newSites = [...specificWebsites];
 newSites[idx].name = e.target.value;
 setSpecificWebsites(newSites);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 </div>
 ))
 )}
 {specificWebsites.length > 0 && specificWebsites.length < 5 && (
 Array.from({ length: 5 - specificWebsites.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + specificWebsites.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setSpecificWebsites([...specificWebsites, {url: "", name: ""}]);
 setSelectedSpecificWebsiteIdx(specificWebsites.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedSpecificWebsiteIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedSpecificWebsiteIdx === null}
 onClick={() => {
 if (selectedSpecificWebsiteIdx !== null) {
 const newSites = specificWebsites.filter((_, idx) => idx !== selectedSpecificWebsiteIdx);
 setSpecificWebsites(newSites);
 setSelectedSpecificWebsiteIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 )}

 {contentFilterType === 'plugin' && (
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Filter Name</div>
 <div className="text-sm text-slate-500 mb-2">Display name of the filter in the app and on the device</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Identifier</div>
 <div className="text-sm text-slate-500 mb-2">Bundle identifier of the app containing filter providers</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Service Address</div>
 <div className="text-sm text-slate-500 mb-2">Host name or IP address or URL for the device</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Organization</div>
 <div className="text-sm text-slate-500 mb-2">Organization parameter for use by the filter provider</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">User Name</div>
 <div className="text-sm text-slate-500 mb-2">User name for authenticating the service</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Password</div>
 <div className="text-sm text-slate-500 mb-2">Password for authenticating the service</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Certificate</div>
 <div className="text-sm text-slate-500 mb-2">Certificate for authenticating to the service</div>
 <Select 
 className="w-full"
 placeholder="Add certificates in the Certificates payload"
 >
 <Select.Option value="none">None</Select.Option>
 </Select>
 </div>

 <div className="flex flex-col gap-2 pt-2">
 <div className="flex items-center gap-3">
 <Checkbox defaultChecked />
 <span className="font-semibold text-slate-800 text-sm">Filter WebKit Traffic</span>
 </div>
 <div className="flex items-center gap-3">
 <Checkbox defaultChecked />
 <span className="font-semibold text-slate-800 text-sm">Filter Socket Traffic</span>
 </div>
 </div>

 <div className="pt-2">
 <div className="font-semibold text-slate-800 text-base mb-1">Custom Data</div>
 <div className="text-sm text-slate-500 mb-4">Custom configuration data for the filter</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="grid grid-cols-2 bg-slate-100 border-b border-slate-200 font-semibold text-slate-700 text-xs px-2 py-2">
 <div>Key</div>
 <div className="border-l border-slate-300 pl-2">Value</div>
 </div>
 <div className="h-32 bg-slate-50 flex flex-col overflow-y-auto">
 {pluginCustomData.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 bg-slate-50"></div>
 </>
 ) : (
 pluginCustomData.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 grid grid-cols-2 cursor-pointer transition-colors ${selectedPluginDataIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedPluginDataIdx(idx)}
 >
 <div className="px-2">
 <Input 
 value={item.key}
 onChange={(e) => {
 const newData = [...pluginCustomData];
 newData[idx].key = e.target.value;
 setPluginCustomData(newData);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="px-2 border-l border-slate-200">
 <Input 
 value={item.value}
 onChange={(e) => {
 const newData = [...pluginCustomData];
 newData[idx].value = e.target.value;
 setPluginCustomData(newData);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 </div>
 ))
 )}
 {pluginCustomData.length > 0 && pluginCustomData.length < 4 && (
 Array.from({ length: 4 - pluginCustomData.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + pluginCustomData.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setPluginCustomData([...pluginCustomData, {key: "", value: ""}]);
 setSelectedPluginDataIdx(pluginCustomData.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedPluginDataIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedPluginDataIdx === null}
 onClick={() => {
 if (selectedPluginDataIdx !== null) {
 const newData = pluginCustomData.filter((_, idx) => idx !== selectedPluginDataIdx);
 setPluginCustomData(newData);
 setSelectedPluginDataIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </div>
 )}
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsContentFilterConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasContentFilterConfig(true);
 setIsContentFilterConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Certificate Transparency Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <CheckCircle2 className="w-5 h-5" /> 
 CERTIFICATE TRANSPARENCY
 </div>
 }
 open={isCertificateTransparencyConfigVisible}
 onCancel={() => setIsCertificateTransparencyConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Excluded Certificates */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Excluded Certificates</div>
 <div className="text-sm text-slate-500 mb-4 leading-relaxed">
 Certificates to be excluded from Certificate Transparency enforcement. The value should be set to the SHA-256 hash of the certificate&apos;s subject public key info
 </div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {excludedCertificates.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 excludedCertificates.map((cert, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedExcludedCertificateIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedExcludedCertificateIdx(idx)}
 >
 <Input 
 value={cert}
 onChange={(e) => {
 const newCerts = [...excludedCertificates];
 newCerts[idx] = e.target.value;
 setExcludedCertificates(newCerts);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm font-mono"
 placeholder="SHA-256 hash"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {excludedCertificates.length > 0 && excludedCertificates.length < 5 && (
 Array.from({ length: 5 - excludedCertificates.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + excludedCertificates.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setExcludedCertificates([...excludedCertificates, ""]);
 setSelectedExcludedCertificateIdx(excludedCertificates.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedExcludedCertificateIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedExcludedCertificateIdx === null}
 onClick={() => {
 if (selectedExcludedCertificateIdx !== null) {
 const newCerts = excludedCertificates.filter((_, idx) => idx !== selectedExcludedCertificateIdx);
 setExcludedCertificates(newCerts);
 setSelectedExcludedCertificateIdx(null);
 }
 }}
 />
 </div>
 </div>

 {/* Excluded Domains */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Excluded Domains</div>
 <div className="text-sm text-slate-500 mb-4 leading-relaxed">
 Domains to be excluded from Certificate Transparency enforcement
 </div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {excludedDomains.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 excludedDomains.map((domain, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedExcludedDomainIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedExcludedDomainIdx(idx)}
 >
 <Input 
 value={domain}
 onChange={(e) => {
 const newDomains = [...excludedDomains];
 newDomains[idx] = e.target.value;
 setExcludedDomains(newDomains);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 placeholder="example.com"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {excludedDomains.length > 0 && excludedDomains.length < 5 && (
 Array.from({ length: 5 - excludedDomains.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + excludedDomains.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setExcludedDomains([...excludedDomains, ""]);
 setSelectedExcludedDomainIdx(excludedDomains.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedExcludedDomainIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedExcludedDomainIdx === null}
 onClick={() => {
 if (selectedExcludedDomainIdx !== null) {
 const newDomains = excludedDomains.filter((_, idx) => idx !== selectedExcludedDomainIdx);
 setExcludedDomains(newDomains);
 setSelectedExcludedDomainIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsCertificateTransparencyConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasCertificateTransparencyConfig(true);
 setIsCertificateTransparencyConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Wi-Fi Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Wifi className="w-5 h-5" /> 
 WI-FI CONFIGURATION
 </div>
 }
 open={isWifiConfigVisible}
 onCancel={() => setIsWifiConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form form={wifiForm} layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* SSID */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Service Set Identifier (SSID)</div>
 <div className="text-sm text-slate-500 mb-3">Identification of the wireless network to connect to</div>
 <Input
  value={wifiSsid}
  onChange={(e) => setWifiSsid(e.target.value)}
  placeholder="[required]"
  className="bg-slate-50 hover:bg-white focus:bg-white transition-colors"
 />
 </div>

 {/* Checkboxes */}
 <div className="flex flex-col gap-4 p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 {/* Hidden Network — not in backend */}
 <div className="flex items-start gap-4 opacity-50 pointer-events-none select-none">
 <Form.Item name="hiddenNetwork" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox />
 </Form.Item>
 <div>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base">
  Hidden Network
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500">Enable if target network is not open or broadcasting</div>
 </div>
 </div>

 {/* Auto Join — supported */}
 <div className="flex items-start gap-4">
 <Form.Item name="autoJoin" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox defaultChecked />
 </Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-base">Auto Join</div>
 <div className="text-sm text-slate-500">Automatically join this wireless network</div>
 </div>
 </div>

 {/* Disable Captive Network Detection — not in backend */}
 <div className="flex items-start gap-4 opacity-50 pointer-events-none select-none">
 <Form.Item name="disableCaptive" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox />
 </Form.Item>
 <div>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base">
  Disable Captive Network Detection
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500">Do not show the captive network assistant</div>
 </div>
 </div>

 {/* Disable MAC Randomization — not in backend */}
 <div className="flex items-start gap-4 opacity-50 pointer-events-none select-none">
 <Form.Item name="disableMacRand" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox />
 </Form.Item>
 <div>
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base">
  Disable Association MAC Randomization
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500">Connections to this Wi-Fi network will use a non-private MAC address</div>
 </div>
 </div>
 </div>

 {/* Select Options */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 {/* Proxy Setup — not in backend */}
 <div className="opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
  Proxy Setup
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-2">Configures proxies to be used with this network</div>
 <Select value={wifiProxySetup} onChange={setWifiProxySetup} className="cursor-pointer w-full max-w-xs">
 <Select.Option value="none">None</Select.Option>
 <Select.Option value="manual">Manual</Select.Option>
 <Select.Option value="auto">Auto</Select.Option>
 </Select>
 </div>

 {/* Security Type — not in backend */}
 <div className="opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
  Security Type
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-2">Wireless network encryption to use when connecting</div>
 <Select value={wifiSecurityType} onChange={setWifiSecurityType} className="cursor-pointer w-full max-w-xs">
 <Select.Option value="none">None</Select.Option>
 <Select.Option value="wep">WEP</Select.Option>
 <Select.Option value="wpa">WPA / WPA2</Select.Option>
 <Select.Option value="any">Any</Select.Option>
 </Select>
 </div>

 {/* Network Type — not in backend */}
 <div className="opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
  Network Type
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-2">Configures network to appear as legacy or Passpoint hotspot</div>
 <Select defaultValue="standard" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="standard">Standard</Select.Option>
 <Select.Option value="legacy">Legacy</Select.Option>
 <Select.Option value="passpoint">Passpoint</Select.Option>
 </Select>
 </div>

 {/* Fast Lane QoS — not in backend */}
 <div className="opacity-50 pointer-events-none select-none">
 <div className="flex items-center gap-2 font-semibold text-slate-800 text-base mb-1">
  Fast Lane QoS Marking
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <Select defaultValue="not-restrict" className="cursor-pointer w-full max-w-xs mt-1 cursor-pointer">
 <Select.Option value="not-restrict">Do not restrict QoS marking</Select.Option>
 <Select.Option value="restrict">Restrict QoS marking</Select.Option>
 </Select>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsWifiConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
  const values = wifiForm.getFieldsValue();
  setWifiData({
   ssid: wifiSsid,
   auto_join: values.autoJoin ?? true,
   hidden_network: values.hiddenNetwork ?? false,
   disable_captive: values.disableCaptive ?? false,
   disable_mac_randomization: values.disableMacRand ?? false,
   security_type: wifiSecurityType,
   proxy_setup: wifiProxySetup,
  });
  setHasWifiConfig(true);
  setIsWifiConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* VPN Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Server className="w-5 h-5" /> 
 VPN CONFIGURATION
 </div>
 }
 open={isVpnConfigVisible}
 onCancel={() => setIsVpnConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Connection Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Connection Name</div>
 <div className="text-sm text-slate-500 mb-3">Display name of the connection (displayed on the device)</div>
 <Input
  value={vpnConnectionName}
  onChange={(e) => setVpnConnectionName(e.target.value)}
  placeholder="[required]"
  className="bg-slate-50 hover:bg-white focus:bg-white transition-colors"
 />
 </div>

 {/* Connection Type */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Connection Type</div>
 <div className="text-sm text-slate-500 mb-3">Type of connection enabled by this policy</div>
 <Select value={vpnConnectionType} onChange={setVpnConnectionType} className="cursor-pointer w-full max-w-xs">
 <Select.Option value="l2tp">L2TP</Select.Option>
 <Select.Option value="pptp">PPTP</Select.Option>
 <Select.Option value="ipsec">IPSec</Select.Option>
 <Select.Option value="ikev2">IKEv2</Select.Option>
 </Select>
 </div>

 {/* Server */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Server</div>
 <div className="text-sm text-slate-500 mb-3">Host name or IP address for server</div>
 <Input
  value={vpnServer}
  onChange={(e) => setVpnServer(e.target.value)}
  placeholder="[required]"
  className="bg-slate-50 hover:bg-white focus:bg-white transition-colors"
 />
 </div>

 {/* Account — not in backend */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Account
 <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-3">User account for authenticating the connection</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* User Authentication — not in backend */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
  User Authentication
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-4">Authentication type for connection</div>

 <div className="space-y-4">
 <div className="flex items-center gap-4">
 <div className="flex items-center gap-2">
 <input type="radio" id="auth-password" name="userAuth" defaultChecked className="w-4 h-4 text-slate-700 border-slate-300 focus:ring-blue-500" />
 <label htmlFor="auth-password" className="font-medium text-slate-700">Password</label>
 </div>
 <Input.Password className="max-w-xs bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div className="flex items-center gap-2">
 <input type="radio" id="auth-rsa" name="userAuth" className="w-4 h-4 text-slate-700 border-slate-300 focus:ring-blue-500" />
 <label htmlFor="auth-rsa" className="font-medium text-slate-700">RSA SecurID</label>
 </div>

 <div className="flex items-start gap-3 pt-2">
 <Form.Item name="sendAllTraffic" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox />
 </Form.Item>
 <div className="font-semibold text-slate-800 text-sm">Send all traffic through VPN</div>
 </div>
 </div>
 </div>

 {/* Machine Authentication — not in backend */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
  Machine Authentication
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-3">Authentication type for connection</div>
 <Select defaultValue="shared-secret" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="shared-secret">Shared Secret</Select.Option>
 <Select.Option value="certificate">Certificate</Select.Option>
 </Select>
 </div>

 {/* Shared Secret — not in backend */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
  Shared Secret
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-3">Shared secret for the connection</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Proxy Configuration — not in backend */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm opacity-50 pointer-events-none select-none">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
  Proxy Setup
  <span className="text-[10px] bg-slate-200 text-slate-500 px-1.5 py-0.5 rounded font-normal">Not supported</span>
 </div>
 <div className="text-sm text-slate-500 mb-3">Configures proxies to be used with this VPN connection</div>
 <Select defaultValue="none" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="none">None</Select.Option>
 <Select.Option value="manual">Manual</Select.Option>
 <Select.Option value="auto">Auto</Select.Option>
 </Select>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsVpnConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
  setVpnData({
   connection_name: vpnConnectionName,
   server: vpnServer,
   type: vpnConnectionType,
  });
  setHasVpnConfig(true);
  setIsVpnConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* AirPlay Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <MonitorPlay className="w-5 h-5" /> 
 AIRPLAY CONFIGURATION
 </div>
 }
 open={isAirPlayConfigVisible}
 onCancel={() => setIsAirPlayConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Passwords */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-4">Passwords</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="grid grid-cols-2 bg-slate-100 border-b border-slate-200 font-semibold text-slate-700 text-xs px-2 py-2">
 <div>Device Name</div>
 <div className="border-l border-slate-300 pl-2">Password</div>
 </div>
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {airPlayPasswords.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 airPlayPasswords.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 grid grid-cols-2 cursor-pointer transition-colors ${selectedAirPlayPasswordIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAirPlayPasswordIdx(idx)}
 >
 <div className="px-2">
 <Input 
 value={item.deviceName}
 onChange={(e) => {
 const newPass = [...airPlayPasswords];
 newPass[idx].deviceName = e.target.value;
 setAirPlayPasswords(newPass);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="px-2 border-l border-slate-200">
 <Input.Password 
 value={item.password}
 onChange={(e) => {
 const newPass = [...airPlayPasswords];
 newPass[idx].password = e.target.value;
 setAirPlayPasswords(newPass);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {airPlayPasswords.length > 0 && airPlayPasswords.length < 5 && (
 Array.from({ length: 5 - airPlayPasswords.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + airPlayPasswords.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAirPlayPasswords([...airPlayPasswords, {deviceName: "", password: ""}]);
 setSelectedAirPlayPasswordIdx(airPlayPasswords.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAirPlayPasswordIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAirPlayPasswordIdx === null}
 onClick={() => {
 if (selectedAirPlayPasswordIdx !== null) {
 const newPass = airPlayPasswords.filter((_, idx) => idx !== selectedAirPlayPasswordIdx);
 setAirPlayPasswords(newPass);
 setSelectedAirPlayPasswordIdx(null);
 }
 }}
 />
 </div>
 </div>

 {/* Allowed */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-4">Allowed</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {airPlayAllowed.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 airPlayAllowed.map((macAddress, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-2 cursor-pointer transition-colors ${selectedAirPlayAllowedIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAirPlayAllowedIdx(idx)}
 >
 <Input 
 value={macAddress}
 onChange={(e) => {
 const newAllowed = [...airPlayAllowed];
 newAllowed[idx] = e.target.value;
 setAirPlayAllowed(newAllowed);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm font-mono"
 placeholder="MAC Address"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {airPlayAllowed.length > 0 && airPlayAllowed.length < 5 && (
 Array.from({ length: 5 - airPlayAllowed.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + airPlayAllowed.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAirPlayAllowed([...airPlayAllowed, ""]);
 setSelectedAirPlayAllowedIdx(airPlayAllowed.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAirPlayAllowedIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAirPlayAllowedIdx === null}
 onClick={() => {
 if (selectedAirPlayAllowedIdx !== null) {
 const newAllowed = airPlayAllowed.filter((_, idx) => idx !== selectedAirPlayAllowedIdx);
 setAirPlayAllowed(newAllowed);
 setSelectedAirPlayAllowedIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsAirPlayConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasAirPlayConfig(true);
 setIsAirPlayConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* AirPlay Security Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Shield className="w-5 h-5" /> 
 AIRPLAY SECURITY CONFIGURATION
 </div>
 }
 open={isAirPlaySecurityConfigVisible}
 onCancel={() => setIsAirPlaySecurityConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Access */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Access</div>
 <div className="text-sm text-slate-500 mb-3">Network requirement for devices that connect to Apple TV using AirPlay</div>
 <Select defaultValue="wifi" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="wifi">Devices on any Wi-Fi network</Select.Option>
 <Select.Option value="same-wifi">Devices on same Wi-Fi network</Select.Option>
 <Select.Option value="anyone">Anyone</Select.Option>
 </Select>
 </div>

 {/* Security */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Security</div>
 <div className="text-sm text-slate-500 mb-3">Security requirement for devices that connect to the Apple TV using AirPlay</div>
 <Select defaultValue="first-time" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="first-time">First-time passcode</Select.Option>
 <Select.Option value="every-time">Passcode on every connection</Select.Option>
 <Select.Option value="none">None</Select.Option>
 </Select>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsAirPlaySecurityConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasAirPlaySecurityConfig(true);
 setIsAirPlaySecurityConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* AirPrint Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Printer className="w-5 h-5" /> 
 AIRPRINT CONFIGURATION
 </div>
 }
 open={isAirPrintConfigVisible}
 onCancel={() => setIsAirPrintConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Printers */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-4">Printers</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="grid grid-cols-12 bg-slate-100 border-b border-slate-200 font-semibold text-slate-700 text-xs px-2 py-2">
 <div className="col-span-5">Host name or IP Addre...</div>
 <div className="col-span-2 border-l border-slate-300 pl-2">Use TLS</div>
 <div className="col-span-2 border-l border-slate-300 pl-2">Port</div>
 <div className="col-span-3 border-l border-slate-300 pl-2">Resource path</div>
 </div>
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {airPrintPrinters.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 airPrintPrinters.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 grid grid-cols-12 items-center cursor-pointer transition-colors ${selectedAirPrintIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAirPrintIdx(idx)}
 >
 <div className="col-span-5 px-2">
 <Input 
 value={item.host}
 onChange={(e) => {
 const newPrinters = [...airPrintPrinters];
 newPrinters[idx].host = e.target.value;
 setAirPrintPrinters(newPrinters);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="col-span-2 px-2 border-l border-slate-200 flex items-center justify-center">
 <Checkbox 
 checked={item.useTls}
 onChange={(e) => {
 const newPrinters = [...airPrintPrinters];
 newPrinters[idx].useTls = e.target.checked;
 setAirPrintPrinters(newPrinters);
 }}
 />
 </div>
 <div className="col-span-2 px-2 border-l border-slate-200">
 <Input 
 value={item.port}
 onChange={(e) => {
 const newPrinters = [...airPrintPrinters];
 newPrinters[idx].port = e.target.value;
 setAirPrintPrinters(newPrinters);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="col-span-3 px-2 border-l border-slate-200">
 <Input 
 value={item.resourcePath}
 onChange={(e) => {
 const newPrinters = [...airPrintPrinters];
 newPrinters[idx].resourcePath = e.target.value;
 setAirPrintPrinters(newPrinters);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {airPrintPrinters.length > 0 && airPrintPrinters.length < 5 && (
 Array.from({ length: 5 - airPrintPrinters.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + airPrintPrinters.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAirPrintPrinters([...airPrintPrinters, {host: "", useTls: false, port: "", resourcePath: ""}]);
 setSelectedAirPrintIdx(airPrintPrinters.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAirPrintIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAirPrintIdx === null}
 onClick={() => {
 if (selectedAirPrintIdx !== null) {
 const newPrinters = airPrintPrinters.filter((_, idx) => idx !== selectedAirPrintIdx);
 setAirPrintPrinters(newPrinters);
 setSelectedAirPrintIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsAirPrintConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasAirPrintConfig(true);
 setIsAirPrintConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Calendar Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Calendar className="w-5 h-5" /> 
 CALENDAR CONFIGURATION
 </div>
 }
 open={isCalendarConfigVisible}
 onCancel={() => setIsCalendarConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Account Description */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-3">The display name of the account (e.g. &quot;Company CalDAV Account&quot;)</div>
 <Input defaultValue="My CalDAV Account" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account Host Name and Port */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Host Name and Port</div>
 <div className="text-sm text-slate-500 mb-3">The CalDAV host name or IP address and port number</div>
 <div className="flex items-center gap-2">
 <Input placeholder="[required]" className="flex-[3] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 <span className="text-slate-600 font-bold mx-1">:</span>
 <Input defaultValue="443" className="flex-1 min-w-[80px] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>

 {/* Principal URL */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Principal URL</div>
 <div className="text-sm text-slate-500 mb-3">The Principal URL for the CalDAV account</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account User Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Account User Name
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-3">The CalDAV user name</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account Password */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Password</div>
 <div className="text-sm text-slate-500 mb-3">The CalDAV password</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Use SSL */}
 <div className="flex items-start gap-4 p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <Form.Item name="useSsl" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox defaultChecked />
 </Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-base">Use SSL</div>
 <div className="text-sm text-slate-500">Enable Secure Socket Layer communication with CalDAV server</div>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsCalendarConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasCalendarConfig(true);
 setIsCalendarConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Contacts Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Contact className="w-5 h-5" /> 
 CONTACTS CONFIGURATION
 </div>
 }
 open={isContactsConfigVisible}
 onCancel={() => setIsContactsConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Account Description */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-3">The display name of the account (e.g. &quot;Company CardDAV Account&quot;)</div>
 <Input defaultValue="My CardDAV Account" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account Host Name and Port */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Host Name and Port</div>
 <div className="text-sm text-slate-500 mb-3">The CardDAV host name or IP address and port number</div>
 <div className="flex items-center gap-2">
 <Input placeholder="[required]" className="flex-[3] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 <span className="text-slate-600 font-bold mx-1">:</span>
 <Input defaultValue="443" className="flex-1 min-w-[80px] bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>

 {/* Principal URL */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Principal URL</div>
 <div className="text-sm text-slate-500 mb-3">The Principal URL for the CardDAV account</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account User Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Account User Name
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-3">The CardDAV user name</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account Password */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Password</div>
 <div className="text-sm text-slate-500 mb-3">The CardDAV password</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Use SSL */}
 <div className="flex items-start gap-4 p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <Form.Item name="useSslContacts" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox defaultChecked />
 </Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-base">Use SSL</div>
 <div className="text-sm text-slate-500">Enable Secure Socket Layer communication with CardDAV server</div>
 </div>
 </div>

 {/* Communication Service Rules */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Communication Service Rules</div>
 <div className="text-sm text-slate-500 mb-3">Choose a default app to be used when calling contacts from this account</div>
 <div className="flex items-center gap-3">
 <Button className="bg-slate-800 hover:bg-slate-700 text-white border-none rounded-md px-4">Choose...</Button>
 <span className="text-slate-400 text-sm">Optional</span>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsContactsConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasContactsConfig(true);
 setIsContactsConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Exchange ActiveSync Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Mail className="w-5 h-5" /> 
 EXCHANGE ACTIVESYNC CONFIGURATION
 </div>
 }
 open={isExchangeConfigVisible}
 onCancel={() => setIsExchangeConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Account Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Name</div>
 <div className="text-sm text-slate-500 mb-3">Name for the Exchange ActiveSync account</div>
 <Input defaultValue="Exchange ActiveSync" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Host */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Exchange ActiveSync Host</div>
 <div className="text-sm text-slate-500 mb-3">Microsoft Exchange Server</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Use SSL */}
 <div className="flex items-start gap-4 p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <Form.Item name="useSslExchange" valuePropName="checked" className="mb-0 pt-0.5">
 <Checkbox defaultChecked />
 </Form.Item>
 <div className="font-semibold text-slate-800 text-base">Use SSL</div>
 </div>

 {/* User */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 User
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-3">The user of the account with optional domain (e.g. &quot;user&quot; or &quot;domain\user&quot;)</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Email Address */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Email Address
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-3">The address of the account (e.g. &quot;john@example.com&quot;)</div>
 <Input placeholder="[defaults to username@host]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* OAuth & Password */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-4">
 <div className="flex items-center gap-3">
 <Checkbox />
 <span className="font-semibold text-slate-800 text-sm">Use OAuth for authentication</span>
 </div>
 
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Password</div>
 <div className="text-sm text-slate-500 mb-3">The password for the account (e.g. &quot;MyP4ssw0rd!&quot;)</div>
 <Input.Password className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div className="flex items-center gap-3">
 <Checkbox />
 <span className="font-semibold text-slate-800 text-sm">Override previous password</span>
 </div>
 </div>

 {/* Sync & Auth */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Past Days of Mail to Sync</div>
 <div className="text-sm text-slate-500 mb-2">The number of past days of mail to synchronize</div>
 <Select defaultValue="1week" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="1week">1 week</Select.Option>
 <Select.Option value="2weeks">2 weeks</Select.Option>
 <Select.Option value="1month">1 month</Select.Option>
 <Select.Option value="all">No limit</Select.Option>
 </Select>
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Authentication Credential Name</div>
 <div className="text-sm text-slate-500 mb-2">Name or description for ActiveSync</div>
 <Input className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Authentication Credential</div>
 <div className="text-sm text-slate-500 mb-2">Credential for authenticating the ActiveSync account</div>
 <Select defaultValue="none" className="cursor-pointer w-full max-w-xs text-slate-400 cursor-pointer">
 <Select.Option value="none">No Value</Select.Option>
 </Select>
 </div>
 </div>

 {/* Toggles */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-4">
 <div className="flex items-start gap-3 cursor-pointer">
 <Checkbox defaultChecked className="mt-0.5" />
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow messages to be moved</div>
 <div className="text-sm text-slate-500">Messages can be moved from this account to another</div>
 </div>
 </div>
 
 <div className="flex items-start gap-3 cursor-pointer">
 <Checkbox defaultChecked className="mt-0.5" />
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow recent addresses to be synced</div>
 <div className="text-sm text-slate-500">Include this account in recent address syncing</div>
 </div>
 </div>

 <div className="flex items-start gap-3 cursor-pointer">
 <Checkbox className="mt-0.5" />
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow Mail Drop</div>
 <div className="text-sm text-slate-500">Allow Mail Drop for this account</div>
 </div>
 </div>

 <div className="flex items-start gap-3 cursor-pointer">
 <Checkbox className="mt-0.5" />
 <div>
 <div className="font-semibold text-slate-800 text-sm">Use only in Mail</div>
 <div className="text-sm text-slate-500">Send outgoing mail from this account only from Mail app</div>
 </div>
 </div>
 </div>

 {/* S/MIME */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-3">
 <div className="flex items-center gap-3"><Checkbox /><span className="font-semibold text-slate-800 text-sm">Enable S/MIME signing</span></div>
 <div className="flex items-center gap-3"><Checkbox /><span className="font-semibold text-slate-800 text-sm">Allow user to enable or disable S/MIME signing</span></div>
 <div className="flex items-center gap-3"><Checkbox /><span className="font-semibold text-slate-800 text-sm">Enable S/MIME encryption by default</span></div>
 <div className="flex items-center gap-3"><Checkbox /><span className="font-semibold text-slate-800 text-sm">Allow user to enable or disable S/MIME encryption</span></div>
 <div className="flex items-center gap-3"><Checkbox /><span className="font-semibold text-slate-800 text-sm">Enable per-message encryption switch</span></div>
 </div>

 {/* Enabled Services */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Enabled Services</div>
 <div className="text-sm text-slate-500 mb-3">Enabled services for this account. At least one of them should be enabled</div>
 <div className="space-y-2">
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Mail</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Contacts</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Calendars</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Reminders</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Notes</span></div>
 </div>
 </div>

 {/* Account Modification */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Modification</div>
 <div className="text-sm text-slate-500 mb-3">Allow users to modify the state of the following services</div>
 <div className="space-y-2">
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Mail</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Contacts</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Calendars</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Reminders</span></div>
 <div className="flex items-center gap-3"><Checkbox defaultChecked /><span className="text-slate-700">Notes</span></div>
 </div>
 </div>

 {/* Communication Service Rules */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Communication Service Rules</div>
 <div className="text-sm text-slate-500 mb-3">Choose a default app to be used when calling contacts from this account</div>
 <div className="flex items-center gap-3">
 <Button className="bg-slate-800 hover:bg-slate-700 text-white border-none rounded-md px-4">Choose...</Button>
 <span className="text-slate-400 text-sm">Optional</span>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsExchangeConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasExchangeConfig(true);
 setIsExchangeConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Google Account Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Mail className="w-5 h-5" /> 
 GOOGLE ACCOUNT CONFIGURATION
 </div>
 }
 open={isGoogleAccountConfigVisible}
 onCancel={() => setIsGoogleAccountConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Account Description */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-3">The display name of the account</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Account Name */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Account Name</div>
 <div className="text-sm text-slate-500 mb-3">The full name of the user for the account</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Email Address */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Email Address</div>
 <div className="text-sm text-slate-500 mb-3">The Google email address of the account</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 {/* Communication Service Rules */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Communication Service Rules</div>
 <div className="text-sm text-slate-500 mb-3">Choose a default app to be used when calling contacts from this account</div>
 <div className="flex items-center gap-3">
 <Button className="bg-slate-800 hover:bg-slate-700 text-white border-none rounded-md px-4">Choose...</Button>
 <span className="text-slate-400 text-sm">Optional</span>
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsGoogleAccountConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasGoogleAccountConfig(true);
 setIsGoogleAccountConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* LDAP Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Key className="w-5 h-5" /> 
 LDAP CONFIGURATION
 </div>
 }
 open={isLdapConfigVisible}
 onCancel={() => setIsLdapConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{
 body: { padding: 0 }
 }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 {/* Basic Info */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-2">The display name of the account (e.g. &quot;Company LDAP Account&quot;)</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account User Name</div>
 <div className="text-sm text-slate-500 mb-2">The user name for this LDAP account</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Password</div>
 <div className="text-sm text-slate-500 mb-2">The password for this LDAP account</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Host Name</div>
 <div className="text-sm text-slate-500 mb-2">The LDAP host name or IP address</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>

 <div className="flex items-center gap-3 pt-2">
 <Form.Item name="useSslLdap" valuePropName="checked" className="mb-0">
 <Checkbox defaultChecked />
 </Form.Item>
 <span className="font-semibold text-slate-800 text-base">Use SSL</span>
 </div>
 </div>

 {/* Search Settings */}
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Search Settings</div>
 <div className="text-sm text-slate-500 mb-4">Search settings for this LDAP server</div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="grid grid-cols-3 bg-slate-100 border-b border-slate-200 font-semibold text-slate-700 text-xs px-2 py-2">
 <div>Description</div>
 <div className="border-l border-slate-300 pl-2">Scope</div>
 <div className="border-l border-slate-300 pl-2">Search Base</div>
 </div>
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {ldapSearchSettings.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 ldapSearchSettings.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 grid grid-cols-3 items-center cursor-pointer transition-colors ${selectedLdapSearchIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedLdapSearchIdx(idx)}
 >
 <div className="px-2">
 <Input 
 value={item.description}
 onChange={(e) => {
 const newSettings = [...ldapSearchSettings];
 newSettings[idx].description = e.target.value;
 setLdapSearchSettings(newSettings);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 <div className="px-2 border-l border-slate-200">
 <Select 
 value={item.scope}
 onChange={(value) => {
 const newSettings = [...ldapSearchSettings];
 newSettings[idx].scope = value;
 setLdapSearchSettings(newSettings);
 }}
 bordered={false}
 className="w-full text-sm custom-select-no-border"
 >
 <Select.Option value="base">Base</Select.Option>
 <Select.Option value="one-level">One Level</Select.Option>
 <Select.Option value="subtree">Subtree</Select.Option>
 </Select>
 </div>
 <div className="px-2 border-l border-slate-200">
 <Input 
 value={item.searchBase}
 onChange={(e) => {
 const newSettings = [...ldapSearchSettings];
 newSettings[idx].searchBase = e.target.value;
 setLdapSearchSettings(newSettings);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {ldapSearchSettings.length > 0 && ldapSearchSettings.length < 5 && (
 Array.from({ length: 5 - ldapSearchSettings.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + ldapSearchSettings.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setLdapSearchSettings([...ldapSearchSettings, {description: "", scope: "subtree", searchBase: ""}]);
 setSelectedLdapSearchIdx(ldapSearchSettings.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedLdapSearchIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedLdapSearchIdx === null}
 onClick={() => {
 if (selectedLdapSearchIdx !== null) {
 const newSettings = ldapSearchSettings.filter((_, idx) => idx !== selectedLdapSearchIdx);
 setLdapSearchSettings(newSettings);
 setSelectedLdapSearchIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>

 {/* Footer Actions */}
 <div className="modal-footer-sticky">
 <Button 
 type="text" 
 className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6"
 onClick={() => setIsLdapConfigVisible(false)}
 >
 CANCEL
 </Button>
 <Button 
 className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8"
 onClick={() => {
 setHasLdapConfig(true);
 setIsLdapConfigVisible(false);
 }}
 >
 SAVE
 </Button>
 </div>
 </div>
 </Modal>

 {/* Mail Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Mail className="w-5 h-5" /> 
 MAIL CONFIGURATION
 </div>
 }
 open={isMailConfigVisible}
 onCancel={() => setIsMailConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-2">The display name of the account (e.g. &quot;Company Mail Account&quot;)</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div className="grid grid-cols-2 gap-4">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Type</div>
 <div className="text-sm text-slate-500 mb-2">The protocol for accessing the email account</div>
 <Select defaultValue="imap" className="cursor-pointer w-full cursor-pointer">
 <Select.Option value="imap">IMAP</Select.Option>
 <Select.Option value="pop">POP</Select.Option>
 </Select>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 invisible">Path Prefix</div>
 <div className="text-sm text-slate-500 mb-2">Path Prefix:</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">User Display Name</div>
 <div className="text-sm text-slate-500 mb-2">The display name of the user (e.g. &quot;John Appleseed&quot;)</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Email Address
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The address of the account (e.g. &quot;john@example.com&quot;)</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>

 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-3">
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="mail_move" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox defaultChecked /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow user to move messages from this account</div>
 <div className="text-xs text-slate-500">Messages can be moved from this account to another</div>
 </div>
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="mail_recent" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox defaultChecked /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow recent addresses to be synced</div>
 <div className="text-xs text-slate-500">Include this account in recent address syncing</div>
 </div>
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="mail_drop" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Allow Mail Drop</div>
 <div className="text-xs text-slate-500">Allow Mail Drop for this account</div>
 </div>
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="mail_use_only" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Use only in Mail</div>
 <div className="text-xs text-slate-500">Send outgoing mail from this account only from Mail app</div>
 </div>
 </div>
 
 <div className="pt-2 space-y-3">
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Enable S/MIME signing</span></div>
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Allow user to enable or disable S/MIME signing</span></div>
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Enable S/MIME encryption by default</span></div>
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Allow user to enable or disable S/MIME encryption</span></div>
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Enable per-message encryption switch</span></div>
 </div>
 </div>

 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <Tabs defaultActiveKey="incoming" centered className="custom-tabs">
 <Tabs.TabPane tab="Incoming Mail" key="incoming">
 <div className="space-y-5 pt-4">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Mail Server and Port</div>
 <div className="text-sm text-slate-500 mb-2">Host name or IP address, and port number for incoming mail</div>
 <div className="flex items-center gap-2">
 <Input placeholder="[required]" className="flex-[3] bg-slate-50" />
 <span className="text-slate-600 font-bold mx-1">:</span>
 <Input defaultValue="993" className="flex-1 min-w-[80px] bg-slate-50" />
 </div>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 User Name
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Account setting required on device">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The user name used to connect to the server for incoming mail</div>
 <Input placeholder="[set on device]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Authentication Type</div>
 <div className="text-sm text-slate-500 mb-2">The authentication method for the incoming mail server</div>
 <Select defaultValue="password" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="password">Password</Select.Option>
 </Select>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Password</div>
 <div className="text-sm text-slate-500 mb-2">The password for the incoming mail server</div>
 <Input.Password className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="mail_ssl" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox defaultChecked /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Use SSL</div>
 <div className="text-xs text-slate-500">Require secure connection when connecting to mail server</div>
 </div>
 </div>
 </div>
 </Tabs.TabPane>
 <Tabs.TabPane tab="Outgoing Mail" key="outgoing">
 <div className="space-y-5 pt-4">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Mail Server and Port</div>
 <div className="text-sm text-slate-500 mb-2">Host name or IP address, and port number for outgoing mail</div>
 <div className="flex items-center gap-2">
 <Input placeholder="[required]" className="flex-[3] bg-slate-50" />
 <span className="text-slate-600 font-bold mx-1">:</span>
 <Input defaultValue="587" className="flex-1 min-w-[80px] bg-slate-50" />
 </div>
 </div>
 </div>
 </Tabs.TabPane>
 </Tabs>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsMailConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasDomainsConfig(true); setIsDomainsConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* macOS Server Account Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Server className="w-5 h-5" /> 
 MACOS SERVER ACCOUNT
 </div>
 }
 open={isMacOsServerConfigVisible}
 onCancel={() => setIsMacOsServerConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 flex items-start gap-3">
 <div className="mt-0.5 text-yellow-500"><AlertCircle className="w-5 h-5" /></div>
 <div className="text-sm text-yellow-800">This payload is deprecated and will be removed in future versions of Apple Configurator</div>
 </div>

 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Account Description</div>
 <div className="text-sm text-slate-500 mb-2">The display name for the macOS Server account</div>
 <Input defaultValue="My Server Account" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Server Address
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The host name, IP address, or URL of the server</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 User Name
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The user name of the account</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Password
 <div className="w-4 h-4 rounded-full bg-yellow-400 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Optional">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The password for the account</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Documents Server Port</div>
 <div className="text-sm text-slate-500 mb-2">The port to connect to for the documents service</div>
 <Input defaultValue="8071" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsMacOsServerConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8" onClick={() => { setHasMacOsServerConfig(true); setIsMacOsServerConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* SCEP Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Key className="w-5 h-5" /> 
 SCEP CONFIGURATION
 </div>
 }
 open={isScepConfigVisible}
 onCancel={() => setIsScepConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 URL
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The base URL for the SCEP server</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Name</div>
 <div className="text-sm text-slate-500 mb-2">The name of the instance: CA-IDENT</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Subject</div>
 <div className="text-sm text-slate-500 mb-2">Representation of a X.500 name</div>
 <Input placeholder="[optional] e.g. O=Company Name/CN=Foo" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Subject Alternative Name Type</div>
 <div className="text-sm text-slate-500 mb-2">The type of a subject alternative name</div>
 <Select defaultValue="none" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="none">None</Select.Option>
 <Select.Option value="rfc822Name">RFC 822 Name</Select.Option>
 <Select.Option value="dNSName">DNS Name</Select.Option>
 <Select.Option value="uniformResourceIdentifier">URI</Select.Option>
 </Select>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Subject Alternative Name Value</div>
 <div className="text-sm text-slate-500 mb-2">The value of a subject alternative name</div>
 <Input className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">NT Principal Name</div>
 <div className="text-sm text-slate-500 mb-2">An NT principal name for use in the certificate request</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Retries</div>
 <div className="text-sm text-slate-500 mb-2">The number of times to poll the SCEP server for a signed certificate before giving up</div>
 <Input defaultValue="3" className="w-full max-w-xs bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Retry Delay</div>
 <div className="text-sm text-slate-500 mb-2">The number of seconds to wait between poll attempts</div>
 <Input defaultValue="10" className="w-full max-w-xs bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Challenge</div>
 <div className="text-sm text-slate-500 mb-2">Used as the pre-shared secret for automatic enrollment</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Key Size</div>
 <div className="text-sm text-slate-500 mb-2">Keysize in bits. 4096 is not supported in iOS 13, iPadOS 13 or tvOS 13 and earlier.</div>
 <Select defaultValue="1024" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="1024">1024</Select.Option>
 <Select.Option value="2048">2048</Select.Option>
 <Select.Option value="4096">4096</Select.Option>
 </Select>
 </div>
 
 <div className="space-y-3">
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Use as digital signature</span></div>
 <div className="flex items-center gap-3"><Checkbox /> <span className="font-semibold text-slate-800 text-sm">Use for key encipherment</span></div>
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Fingerprint</div>
 <div className="text-sm text-slate-500 mb-2">Enter hex string to be used as a fingerprint or use button to create fingerprint from Certificate</div>
 <Input className="bg-slate-50 hover:bg-white focus:bg-white transition-colors mb-3" />
 <div className="flex justify-end">
 <Button>Create from Certificate...</Button>
 </div>
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsScepConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasHttpProxyConfig(true); setIsHttpProxyConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Cellular Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Radio className="w-5 h-5" /> 
 CELLULAR CONFIGURATION
 </div>
 }
 open={isCellularConfigVisible}
 onCancel={() => setIsCellularConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Configured APN Type</div>
 <Select defaultValue="default" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="default">Default APN</Select.Option>
 </Select>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Default APN Name
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">Access Point Name for the default APN configuration</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Default APN Authentication Type</div>
 <div className="text-sm text-slate-500 mb-2">Authentication used by the default APN configuration</div>
 <Select defaultValue="pap" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="pap">PAP</Select.Option>
 <Select.Option value="chap">CHAP</Select.Option>
 </Select>
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Default APN User Name</div>
 <div className="text-sm text-slate-500 mb-2">User name for authenticating the connection</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Default APN Password</div>
 <div className="text-sm text-slate-500 mb-2">Password for authenticating the connection</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Default APN Supported IP Versions</div>
 <div className="text-sm text-slate-500 mb-2">Supported Internet Protocol versions for connections</div>
 <Select defaultValue="any" className="cursor-pointer w-full max-w-xs cursor-pointer">
 <Select.Option value="any">-</Select.Option>
 <Select.Option value="ipv4">IPv4</Select.Option>
 <Select.Option value="ipv6">IPv6</Select.Option>
 <Select.Option value="ipv4v6">IPv4/IPv6</Select.Option>
 </Select>
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsCellularConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasContentFilterConfig(true); setIsContentFilterConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Notifications Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <AlertCircle className="w-5 h-5" /> 
 NOTIFICATIONS CONFIGURATION
 </div>
 }
 open={isNotificationsConfigVisible}
 onCancel={() => setIsNotificationsConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="font-medium text-slate-600">Notifications (supervised only)</div>
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Notification Settings</div>
 <div className="text-sm text-slate-500 mb-4 leading-relaxed">
 Configure notification settings for each app
 </div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {notificationSettings.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 notificationSettings.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-3 cursor-pointer transition-colors ${selectedNotificationSettingIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedNotificationSettingIdx(idx)}
 >
 <Input 
 value={item.appBundleId}
 onChange={(e) => {
 const newSettings = [...notificationSettings];
 newSettings[idx].appBundleId = e.target.value;
 setNotificationSettings(newSettings);
 }}
 bordered={false}
 placeholder="App Bundle ID"
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {notificationSettings.length > 0 && notificationSettings.length < 5 && (
 Array.from({ length: 5 - notificationSettings.length }).map((_, i) => (
 <div key={`empty-${i}`} className={`h-8 border-b border-slate-200 ${(i + notificationSettings.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setNotificationSettings([...notificationSettings, {appBundleId: ""}]);
 setSelectedNotificationSettingIdx(notificationSettings.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedNotificationSettingIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedNotificationSettingIdx === null}
 onClick={() => {
 if (selectedNotificationSettingIdx !== null) {
 const newSettings = notificationSettings.filter((_, idx) => idx !== selectedNotificationSettingIdx);
 setNotificationSettings(newSettings);
 setSelectedNotificationSettingIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsNotificationsConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasRestrictionsConfig(true); setIsRestrictionsConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Conference Room Display Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Monitor className="w-5 h-5" /> 
 CONFERENCE ROOM DISPLAY CONFIGURATION
 </div>
 }
 open={isConferenceRoomConfigVisible}
 onCancel={() => setIsConferenceRoomConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="font-medium text-slate-600">Conference Room Display (supervised only)</div>
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Custom Message</div>
 <div className="text-sm text-slate-500 mb-2">Message displayed on-screen in Conference Room Display mode</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsConferenceRoomConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8" onClick={() => { setHasConferenceRoomConfig(true); setIsConferenceRoomConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* TV Remote Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Tv className="w-5 h-5" /> 
 TV REMOTE CONFIGURATION
 </div>
 }
 open={isTvRemoteConfigVisible}
 onCancel={() => setIsTvRemoteConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="font-medium text-slate-600">TV Remote (supervised only)</div>
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Allowed Remotes (tvOS only)</div>
 <div className="text-sm text-slate-500 mb-4 leading-relaxed">
 Add the MAC addresses of permitted iOS devices to this list to enable remote control with only these devices
 </div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {allowedRemotes.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 allowedRemotes.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-3 cursor-pointer transition-colors ${selectedAllowedRemoteIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAllowedRemoteIdx(idx)}
 >
 <Input 
 value={item}
 onChange={(e) => {
 const newSettings = [...allowedRemotes];
 newSettings[idx] = e.target.value;
 setAllowedRemotes(newSettings);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {allowedRemotes.length > 0 && allowedRemotes.length < 5 && (
 Array.from({ length: 5 - allowedRemotes.length }).map((_, i) => (
 <div key={`empty-remote-${i}`} className={`h-8 border-b border-slate-200 ${(i + allowedRemotes.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAllowedRemotes([...allowedRemotes, ""]);
 setSelectedAllowedRemoteIdx(allowedRemotes.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAllowedRemoteIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAllowedRemoteIdx === null}
 onClick={() => {
 if (selectedAllowedRemoteIdx !== null) {
 const newSettings = allowedRemotes.filter((_, idx) => idx !== selectedAllowedRemoteIdx);
 setAllowedRemotes(newSettings);
 setSelectedAllowedRemoteIdx(null);
 }
 }}
 />
 </div>
 </div>
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm">
 <div className="font-semibold text-slate-800 text-base mb-1">Allowed TVs (iOS only)</div>
 <div className="text-sm text-slate-500 mb-4 leading-relaxed">
 Add the MAC addresses of permitted tvOS devices to this list to restrict remote control to only these devices
 </div>
 
 <div className="border border-slate-200 rounded-md overflow-hidden mb-3">
 <div className="h-40 bg-slate-50 flex flex-col overflow-y-auto">
 {allowedTvs.length === 0 ? (
 <>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 border-b border-slate-200 bg-white"></div>
 <div className="h-8 border-b border-slate-200 bg-slate-50"></div>
 <div className="h-8 bg-white"></div>
 </>
 ) : (
 allowedTvs.map((item, idx) => (
 <div 
 key={idx} 
 className={`h-8 border-b border-slate-200 flex items-center px-3 cursor-pointer transition-colors ${selectedAllowedTvIdx === idx ? 'bg-blue-100' : (idx % 2 === 0 ? 'bg-white hover:bg-slate-100' : 'bg-slate-50 hover:bg-slate-100')}`}
 onClick={() => setSelectedAllowedTvIdx(idx)}
 >
 <Input 
 value={item}
 onChange={(e) => {
 const newSettings = [...allowedTvs];
 newSettings[idx] = e.target.value;
 setAllowedTvs(newSettings);
 }}
 bordered={false}
 className="p-0 h-full bg-transparent focus:shadow-none text-sm"
 />
 </div>
 ))
 )}
 {/* Fill remaining empty space */}
 {allowedTvs.length > 0 && allowedTvs.length < 5 && (
 Array.from({ length: 5 - allowedTvs.length }).map((_, i) => (
 <div key={`empty-tv-${i}`} className={`h-8 border-b border-slate-200 ${(i + allowedTvs.length) % 2 === 0 ? 'bg-white' : 'bg-slate-50'}`}></div>
 ))
 )}
 </div>
 </div>
 <div className="flex gap-2">
 <Button 
 size="small" 
 className="w-8 h-8 flex items-center justify-center bg-white border-slate-300 text-slate-600 hover:text-slate-700 hover:border-blue-600" 
 icon={<Plus className="w-4 h-4" />} 
 onClick={() => {
 setAllowedTvs([...allowedTvs, ""]);
 setSelectedAllowedTvIdx(allowedTvs.length);
 }}
 />
 <Button 
 size="small" 
 className={`w-8 h-8 flex items-center justify-center bg-white border-slate-300 ${selectedAllowedTvIdx !== null ? 'text-slate-600 hover:text-red-500 hover:border-red-500' : 'text-slate-300'}`}
 icon={<Minus className="w-4 h-4" />} 
 disabled={selectedAllowedTvIdx === null}
 onClick={() => {
 if (selectedAllowedTvIdx !== null) {
 const newSettings = allowedTvs.filter((_, idx) => idx !== selectedAllowedTvIdx);
 setAllowedTvs(newSettings);
 setSelectedAllowedTvIdx(null);
 }
 }}
 />
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsTvRemoteConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasPasscodeConfig(true); setIsPasscodeConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Lock Screen Message Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <MessageSquare className="w-5 h-5" /> 
 LOCK SCREEN MESSAGE CONFIGURATION
 </div>
 }
 open={isLockScreenMessageConfigVisible}
 onCancel={() => setIsLockScreenMessageConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="font-medium text-slate-600">Lock Screen Message (supervised only)</div>
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">&quot;If Lost, Return to...&quot; Message</div>
 <div className="text-sm text-slate-500 mb-2">Message displayed on the login window and lock screen</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Asset Tag Information</div>
 <div className="text-sm text-slate-500 mb-2">Message displayed at the bottom of the login window and lock screen</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsLockScreenMessageConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8" onClick={() => { setHasLockScreenMessageConfig(true); setIsLockScreenMessageConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Web Clip Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <LinkIcon className="w-5 h-5" /> 
 WEB CLIP CONFIGURATION
 </div>
 }
 open={isWebClipConfigVisible}
 onCancel={() => setIsWebClipConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 Label
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The name to display for the Web Clip</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 URL
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The URL to be displayed when opening the Web Clip</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 
 <div className="flex items-start gap-3 pt-2">
 <Form.Item name="webclip_removable" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox defaultChecked /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Removable</div>
 <div className="text-xs text-slate-500">Enable removal of the Web Clip</div>
 </div>
 </div>

 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Icon</div>
 <div className="text-sm text-slate-500 mb-2">The icon to use for the Web Clip</div>
 <div className="w-20 h-20 border border-slate-300 rounded-xl mb-2 bg-slate-50"></div>
 <Button size="small">Choose...</Button>
 </div>

 <div className="space-y-4 pt-2">
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="webclip_precomposed" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Precomposed Icon</div>
 <div className="text-xs text-slate-500">The icon will be displayed with no added visual effects</div>
 </div>
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="webclip_fullscreen" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Full Screen</div>
 <div className="text-xs text-slate-500">Displays the web clip as a full screen application</div>
 </div>
 </div>
 <div className="flex items-start gap-3 cursor-pointer">
 <Form.Item name="webclip_ignoremanifest" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox /></Form.Item>
 <div>
 <div className="font-semibold text-slate-800 text-sm">Ignore Manifest Scope</div>
 <div className="text-xs text-slate-500">Allow web pages that are not in the manifest to load in full screen mode</div>
 </div>
 </div>
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsWebClipConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8 transition-colors" onClick={() => { setHasDnsProxyConfig(true); setIsDnsProxyConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Subscribed Calendar Configuration Modal */}
 <Modal
 title={
 <div className="flex items-center gap-2 text-white font-medium tracking-wide text-sm uppercase">
 <Calendar className="w-5 h-5" /> 
 SUBSCRIBED CALENDAR CONFIGURATION
 </div>
 }
 open={isSubscribedCalendarConfigVisible}
 onCancel={() => setIsSubscribedCalendarConfigVisible(false)}
 footer={null}
 width={800}
 className="custom-modal-header-red config-modal"
 styles={{ body: { padding: 0 } }}
 centered
 >
 <div className="flex flex-col h-full bg-slate-50">
 <div className="p-8 flex-1 overflow-y-auto max-h-[70vh] scrollbar-hide">
 <Form layout="vertical" className="max-w-3xl mx-auto custom-form">
 <div className="space-y-6">
 <div className="p-5 bg-white rounded-lg border border-slate-200 shadow-sm space-y-5">
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Description</div>
 <div className="text-sm text-slate-500 mb-2">The description of the calendar subscription</div>
 <Input defaultValue="My Subscribed Calendar" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1 flex items-center gap-2">
 URL
 <div className="w-4 h-4 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold ml-auto cursor-help" title="Required">!</div>
 </div>
 <div className="text-sm text-slate-500 mb-2">The URL of the calendar file</div>
 <Input placeholder="[required]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">User Name</div>
 <div className="text-sm text-slate-500 mb-2">The user name for this subscription</div>
 <Input placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 <div>
 <div className="font-semibold text-slate-800 text-base mb-1">Password</div>
 <div className="text-sm text-slate-500 mb-2">The password for this subscription</div>
 <Input.Password placeholder="[optional]" className="bg-slate-50 hover:bg-white focus:bg-white transition-colors" />
 </div>
 
 <div className="flex items-start gap-3 pt-2">
 <Form.Item name="subcal_ssl" valuePropName="checked" className="mb-0 pt-0.5"><Checkbox defaultChecked /></Form.Item>
 <div className="font-semibold text-slate-800 text-sm mt-0.5">Use SSL</div>
 </div>
 </div>
 </div>
 </Form>
 </div>
 <div className="modal-footer-sticky">
 <Button type="text" className="font-semibold text-slate-600 hover:text-slate-800 hover:bg-slate-100 px-6" onClick={() => setIsSubscribedCalendarConfigVisible(false)}>CANCEL</Button>
 <Button className="font-semibold bg-[#de2a15] hover:bg-[#c22412] text-white border-none px-8" onClick={() => { setHasSubscribedCalendarConfig(true); setIsSubscribedCalendarConfigVisible(false); }}>SAVE</Button>
 </div>
 </div>
 </Modal>

 {/* Assign Profile Modal */}
 <Modal
  title={
   <div className="flex items-center gap-2">
    <div className="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center">
     <Users className="w-4 h-4 text-blue-600" />
    </div>
    <div>
     <h3 className="text-base font-bold text-slate-800 m-0">Assign Profile</h3>
     <p className="text-xs text-slate-500 font-normal m-0">{selectedProfile?.name}</p>
    </div>
   </div>
  }
  open={isAssignModalVisible}
  onCancel={() => { setIsAssignModalVisible(false); setAssignDeviceId(""); setAssignGroupId(""); }}
  footer={
   <div className="flex justify-end gap-3">
    <Button onClick={() => { setIsAssignModalVisible(false); setAssignDeviceId(""); setAssignGroupId(""); }}>
     Cancel
    </Button>
    <Button
     type="primary"
     loading={assignLoading}
     className="bg-[#de2a15] hover:bg-[#c22412] border-none"
     onClick={handleAssignProfile}
    >
     Assign
    </Button>
   </div>
  }
  width={500}
  className="custom-modal"
 >
  <div className="space-y-5 py-2">
   {/* Target Type */}
   <div>
    <div className="text-sm font-semibold text-slate-700 mb-2">Assign to</div>
    <Select
     value={assignTargetType}
     onChange={(val) => { setAssignTargetType(val); setAssignDeviceId(""); setAssignGroupId(""); }}
     className="w-full"
    >
     <Select.Option value="device">Device (by UDID)</Select.Option>
     <Select.Option value="group">Device Group</Select.Option>
    </Select>
   </div>

   {/* Device ID or Group ID */}
   {assignTargetType === "device" ? (
    <div>
     <div className="text-sm font-semibold text-slate-700 mb-2">Device UDID <span className="text-red-500">*</span></div>
     <Input
      placeholder="Enter device UDID..."
      value={assignDeviceId}
      onChange={(e) => setAssignDeviceId(e.target.value)}
     />
    </div>
   ) : (
    <div>
     <div className="text-sm font-semibold text-slate-700 mb-2">Group ID <span className="text-red-500">*</span></div>
     <Input
      placeholder="Enter group ID..."
      value={assignGroupId}
      onChange={(e) => setAssignGroupId(e.target.value)}
      type="number"
     />
    </div>
   )}

   {/* Schedule Type */}
   <div>
    <div className="text-sm font-semibold text-slate-700 mb-2">Schedule</div>
    <Select
     value={assignScheduleType}
     onChange={setAssignScheduleType}
     className="w-full"
    >
     <Select.Option value="immediate">Immediate — push to device now</Select.Option>
     <Select.Option value="scheduled">Scheduled — push later</Select.Option>
    </Select>
    {assignScheduleType === "immediate" && (
     <p className="text-xs text-slate-500 mt-2">The profile will be queued and pushed to the device via APNs right away.</p>
    )}
   </div>
  </div>
 </Modal>

 {/* Custom Styles for Data-Dense Dashboard */}
 <style jsx global>{`
 /* Hide scrollbar for a cleaner look */
 .scrollbar-hide::-webkit-scrollbar {
 display: none;
 }
 .scrollbar-hide {
 -ms-overflow-style: none;
 scrollbar-width: none;
 }
 
 /* Standard Classes */
 .glass-panel {
 background: #ffffff !important;
 border: 1px solid #e2e8f0 !important;
 box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06) !important;
 }
 .dark .glass-panel {
 background: #1e293b !important;
 border: 1px solid #334155 !important;
 }

 .glass-card {
 background: #ffffff !important;
 border: 1px solid #e2e8f0 !important;
 box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05) !important;
 transition: border-color 0.2s ease;
 }
 .glass-card:hover {
 border-color: #cbd5e1 !important;
 }

 .glass-input {
 background: #ffffff !important;
 border: 1px solid #cbd5e1 !important;
 transition: all 0.2s ease;
 }
 .glass-input:focus {
 border-color: #de2a15 !important;
 box-shadow: 0 0 0 2px rgba(222, 42, 21, 0.1) !important;
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
 
 .dark .custom-data-table .ant-table-thead > tr > th {
 background: #0f172a !important;
 color: #cbd5e1 !important;
 border-bottom-color: #334155 !important;
 }
 .dark .custom-data-table .ant-table-tbody > tr > td {
 background: #1e293b !important;
 border-bottom-color: #334155 !important;
 }
 .dark .custom-data-table .ant-table-tbody > tr:hover > td {
 background: #334155 !important;
 }

 .ant-input-group-addon {
 padding: 0 !important;
 border: none !important;
 overflow: hidden !important;
 }
 
 /* Modal Styles */
 .custom-modal .ant-modal-content {
 padding: 0 !important;
 overflow: hidden;
 border-radius: 12px;
 background: #ffffff !important;
 border: 1px solid #e2e8f0;
 box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
 }
 .custom-modal .ant-modal-header {
 padding: 16px 24px;
 margin: 0;
 background: #ffffff !important;
 border-bottom: 1px solid #e2e8f0;
 }
 .custom-modal .ant-modal-body {
 padding: 24px;
 background: #ffffff !important;
 }
 .custom-modal .ant-modal-close {
 top: 16px;
 right: 24px;
 }

 .custom-modal-header-red .ant-modal-content {
 padding: 0 !important;
 overflow: hidden;
 border-radius: 12px;
 background: #ffffff !important;
 border: 1px solid #e2e8f0;
 box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
 }
 .custom-modal-header-red .ant-modal-header {
 padding: 16px 24px;
 margin: 0;
 background: #de2a15 !important;
 border-bottom: 1px solid #c22412;
 }
 .custom-modal-header-red .ant-modal-title {
 color: white !important;
 }
 .custom-modal-header-red .ant-modal-close {
 top: 16px;
 right: 24px;
 color: white !important;
 }
 .custom-modal-header-red .ant-modal-close:hover {
 background: rgba(255, 255, 255, 0.1) !important;
 border-radius: 6px;
 }
 
 /* Dropdown menu modal */
 .dropdown-modal .ant-modal-content {
 background: #ffffff !important;
 border: 1px solid #e2e8f0;
 box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
 border-radius: 12px;
 }
 .dropdown-modal .ant-modal-body {
 padding: 0 !important;
 }
 
 /* Config modals body background */
 .config-modal .ant-modal-body > div > div:first-child {
 background: #f8fafc !important;
 }
 
 /* Update form fields inside modals */
 .config-modal .bg-white {
 background: #ffffff !important;
 border: 1px solid #e2e8f0;
 box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
 }
 .config-modal .bg-slate-50 {
 background: #f8fafc !important;
 }
 .config-modal input, .config-modal .ant-select-selector {
 background: #ffffff !important;
 border: 1px solid #cbd5e1 !important;
 }
 .config-modal input:focus, .config-modal .ant-select-focused .ant-select-selector {
 background: #ffffff !important;
 border-color: #de2a15 !important;
 box-shadow: 0 0 0 2px rgba(222, 42, 21, 0.1) !important;
 }

 /* Form Modal Tabs Customization */
 .form-modal .bg-white {
 background: #ffffff !important;
 border: 1px solid #e2e8f0;
 }
 .form-modal .bg-slate-50 {
 background: #f8fafc !important;
 }
 .form-modal input, .form-modal textarea {
 background: #ffffff !important;
 border: 1px solid #cbd5e1 !important;
 }
 .form-modal input:focus, .form-modal textarea:focus {
 background: #ffffff !important;
 border-color: #de2a15 !important;
 box-shadow: 0 0 0 2px rgba(222, 42, 21, 0.1) !important;
 }
 
 .form-modal .ant-modal-content {
 display: flex;
 flex-direction: column;
 height: 85vh; /* Fixed height for consistent layout */
 max-height: 85vh;
 border-radius: 12px;
 overflow: hidden;
 }
 .form-modal .ant-modal-body {
 flex: 1;
 overflow: hidden;
 display: flex;
 flex-direction: column;
 padding: 0 !important;
 background: #ffffff !important;
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
 background: #f8fafc;
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
 overflow-y: hidden;
 background: #ffffff;
 }
 .custom-tabs .ant-tabs-content {
 height: 100%;
 }
 .custom-tabs .ant-tabs-tabpane {
 height: 100%;
 }
 
 /* Footer sticky container */
 .modal-footer-sticky {
 border-top: 1px solid #e2e8f0;
 padding: 16px 24px;
 background: #f8fafc;
 display: flex;
 justify-content: flex-end;
 gap: 12px;
 z-index: 10;
 }
 `}</style>
 </div>
 );
}
