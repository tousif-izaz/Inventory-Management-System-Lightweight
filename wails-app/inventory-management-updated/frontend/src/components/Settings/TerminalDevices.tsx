import { useState, useEffect } from 'react';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { useToast } from '../UI/Toast';
import {
  registerTerminalDevice,
  getTerminalDevices,
  updateTerminalDevice,
  deleteTerminalDevice,
  setDefaultTerminalDevice,
} from '../../services/terminalService';
import { squareService } from '../../services/squareService';
import type { TerminalDevice, TerminalDeviceCreate } from '../../types/terminal';
import {
  PlusIcon,
  TrashIcon,
  PencilIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';

export const TerminalDevices = () => {
  const { showToast } = useToast();
  const [devices, setDevices] = useState<TerminalDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [merchantId, setMerchantId] = useState<string>('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingDevice, setEditingDevice] = useState<number | null>(null);

  // Form state
  const [formData, setFormData] = useState<TerminalDeviceCreate>({
    square_device_id: '',
    device_name: '',
    location_id: '',
  });

  useEffect(() => {
    loadMerchantAndDevices();
  }, []);

  const loadMerchantAndDevices = async () => {
    try {
      setLoading(true);
      // Get merchant ID from Square OAuth
      const oauthStatus = await squareService.getOAuthStatus();
      if (!oauthStatus.isAuthorized || !oauthStatus.merchantId) {
        showToast('warning', 'Please connect your Square account first');
        return;
      }

      setMerchantId(oauthStatus.merchantId);
      await loadDevices(oauthStatus.merchantId);
    } catch (error) {
      console.error('Failed to load merchant info:', error);
      showToast('error', 'Failed to load Square account info');
    } finally {
      setLoading(false);
    }
  };

  const loadDevices = async (merchantIdParam: string) => {
    try {
      const deviceList = await getTerminalDevices(merchantIdParam);
      setDevices(deviceList);
    } catch (error) {
      console.error('Failed to load devices:', error);
      showToast('error', 'Failed to load terminal devices');
    }
  };

  const handleAddDevice = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!merchantId) {
      showToast('error', 'Merchant ID not found');
      return;
    }

    try {
      await registerTerminalDevice(merchantId, formData);
      showToast('success', 'Terminal device registered successfully');
      setShowAddForm(false);
      setFormData({ square_device_id: '', device_name: '', location_id: '' });
      await loadDevices(merchantId);
    } catch (error: any) {
      console.error('Failed to register device:', error);
      showToast(
        'error',
        error.response?.data?.error || 'Failed to register terminal device'
      );
    }
  };

  const handleUpdateDevice = async (deviceId: number, name: string) => {
    try {
      await updateTerminalDevice(deviceId, { device_name: name });
      showToast('success', 'Device updated successfully');
      setEditingDevice(null);
      await loadDevices(merchantId);
    } catch (error: any) {
      console.error('Failed to update device:', error);
      showToast('error', error.response?.data?.error || 'Failed to update device');
    }
  };

  const handleDeleteDevice = async (deviceId: number) => {
    if (!confirm('Are you sure you want to delete this device?')) {
      return;
    }

    try {
      await deleteTerminalDevice(deviceId);
      showToast('success', 'Device deleted successfully');
      await loadDevices(merchantId);
    } catch (error: any) {
      console.error('Failed to delete device:', error);
      showToast('error', error.response?.data?.error || 'Failed to delete device');
    }
  };

  const handleSetDefault = async (deviceId: number) => {
    try {
      await setDefaultTerminalDevice(deviceId, merchantId);
      showToast('success', 'Default device updated');
      await loadDevices(merchantId);
    } catch (error: any) {
      console.error('Failed to set default device:', error);
      showToast('error', error.response?.data?.error || 'Failed to set default device');
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="text-center">Loading terminal devices...</div>
      </div>
    );
  }

  if (!merchantId) {
    return (
      <div className="p-6">
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <h3 className="font-semibold text-yellow-900">Square Account Required</h3>
          <p className="text-yellow-700 mt-2">
            Please connect your Square account in Settings before managing terminal devices.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-2xl font-bold">Square Terminal Devices</h2>
          <p className="text-gray-600 mt-1">
            Manage physical Square Terminal devices for in-person payments
          </p>
        </div>
        <Button onClick={() => setShowAddForm(!showAddForm)}>
          <div className="flex items-center gap-2">
            <PlusIcon className="w-4 h-4" />
            Add Device
          </div>
        </Button>
      </div>

      {showAddForm && (
        <div className="bg-white border border-gray-200 rounded-lg p-6 mb-6">
          <h3 className="font-semibold mb-4">Register New Terminal Device</h3>
          <form onSubmit={handleAddDevice}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Device ID"
                value={formData.square_device_id}
                onChange={(e) =>
                  setFormData({ ...formData, square_device_id: e.target.value })
                }
                placeholder="Square Device ID"
                required
              />
              <Input
                label="Device Name"
                value={formData.device_name}
                onChange={(e) =>
                  setFormData({ ...formData, device_name: e.target.value })
                }
                placeholder="e.g., Main Counter Terminal"
                required
              />
              <Input
                label="Location ID (optional)"
                value={formData.location_id}
                onChange={(e) =>
                  setFormData({ ...formData, location_id: e.target.value })
                }
                placeholder="Square Location ID"
              />
            </div>
            <div className="flex gap-2 mt-4">
              <Button type="submit">Register Device</Button>
              <Button
                variant="secondary"
                onClick={() => {
                  setShowAddForm(false);
                  setFormData({ square_device_id: '', device_name: '', location_id: '' });
                }}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
        {devices.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <p>No terminal devices registered yet.</p>
            <p className="text-sm mt-2">Click "Add Device" to register your first device.</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Device Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Device ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Default
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {devices.map((device) => (
                <tr key={device.device_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    {editingDevice === device.device_id ? (
                      <Input
                        value={device.device_name}
                        onChange={(e) => {
                          const newDevices = devices.map((d) =>
                            d.device_id === device.device_id ? { ...d, device_name: e.target.value } : d
                          );
                          setDevices(newDevices);
                        }}
                        onBlur={() => handleUpdateDevice(device.device_id, device.device_name)}
                      />
                    ) : (
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{device.device_name}</span>
                        <button
                          onClick={() => setEditingDevice(device.device_id)}
                          className="text-gray-400 hover:text-gray-600"
                        >
                          <PencilIcon className="w-4 h-4" />
                        </button>
                      </div>
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600 font-mono">
                    {device.square_device_id}
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        device.status === 'active'
                          ? 'bg-green-100 text-green-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {device.status}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    {device.is_default ? (
                      <CheckCircleIcon className="w-5 h-5 text-green-600" />
                    ) : (
                      <button
                        onClick={() => handleSetDefault(device.device_id)}
                        className="text-gray-400 hover:text-blue-600"
                        title="Set as default"
                      >
                        <div className="w-5 h-5 border-2 border-gray-300 rounded-full" />
                      </button>
                    )}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => handleDeleteDevice(device.device_id)}
                      className="text-red-600 hover:text-red-800"
                      title="Delete device"
                    >
                      <TrashIcon className="w-5 h-5" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h4 className="font-semibold text-blue-900 mb-2">How to find your Device ID:</h4>
        <ol className="list-decimal list-inside text-blue-800 space-y-1 text-sm">
          <li>Open your Square Dashboard</li>
          <li>Go to Hardware → Devices</li>
          <li>Select your Terminal device</li>
          <li>Copy the Device ID from the device details</li>
        </ol>
      </div>
    </div>
  );
};
