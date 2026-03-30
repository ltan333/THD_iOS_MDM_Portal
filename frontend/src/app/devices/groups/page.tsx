import { DeviceGroupsList } from "@/features/devices";

export default function DeviceGroupsPage() {
  return (
    <div className="flex-1 w-full bg-slate-50 min-h-[calc(100vh-64px)]">
      <DeviceGroupsList />
    </div>
  );
}
