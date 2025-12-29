import { useState } from 'react';

export type TimeRangeOption = 'today' | 'week' | 'custom';

export interface TimeRange {
    startDate: Date;
    endDate: Date;
    type: TimeRangeOption;
}

export interface TimeRangePickerProps {
    value: TimeRange;
    onChange: (range: TimeRange) => void;
}

export const TimeRangePicker = ({ value, onChange }: TimeRangePickerProps) => {
    const [showCustomPicker, setShowCustomPicker] = useState(false);
    const [customStart, setCustomStart] = useState('');
    const [customEnd, setCustomEnd] = useState('');

    const handleQuickSelect = (type: 'today' | 'week') => {
        const now = new Date();
        let startDate: Date;
        let endDate: Date;

        if (type === 'today') {
            startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0);
            endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59);
        } else {
            // This Week (Monday to Sunday)
            const dayOfWeek = now.getDay();
            const mondayOffset = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
            startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate() + mondayOffset, 0, 0, 0);
            endDate = new Date(startDate.getFullYear(), startDate.getMonth(), startDate.getDate() + 6, 23, 59, 59);
        }

        onChange({ startDate, endDate, type });
        setShowCustomPicker(false);
    };

    const handleCustomRange = () => {
        if (customStart && customEnd) {
            const startDate = new Date(customStart);
            startDate.setHours(0, 0, 0, 0);

            const endDate = new Date(customEnd);
            endDate.setHours(23, 59, 59, 999);

            if (startDate <= endDate) {
                onChange({ startDate, endDate, type: 'custom' });
            }
        }
    };

    const formatDateRange = (range: TimeRange): string => {
        const options: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric', year: 'numeric' };
        const start = range.startDate.toLocaleDateString('en-US', options);
        const end = range.endDate.toLocaleDateString('en-US', options);
        return `${start} - ${end}`;
    };

    return (
        <div className="bg-white shadow rounded-lg p-4">
            <h3 className="text-sm font-medium text-gray-700 mb-3">Time Range</h3>

            {/* Quick Select Buttons */}
            <div className="flex flex-wrap gap-2 mb-4">
                <button
                    onClick={() => handleQuickSelect('today')}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${
                        value.type === 'today'
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                >
                    Today
                </button>

                <button
                    onClick={() => handleQuickSelect('week')}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${
                        value.type === 'week'
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                >
                    This Week
                </button>

                <button
                    onClick={() => {
                        setShowCustomPicker(!showCustomPicker);
                        if (!showCustomPicker && value.type === 'custom') {
                            setCustomStart(value.startDate.toISOString().split('T')[0]);
                            setCustomEnd(value.endDate.toISOString().split('T')[0]);
                        }
                    }}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${
                        value.type === 'custom' || showCustomPicker
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                >
                    Custom Range
                </button>
            </div>

            {/* Custom Date Picker */}
            {showCustomPicker && (
                <div className="border-t border-gray-200 pt-4 mt-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Start Date
                            </label>
                            <input
                                type="date"
                                value={customStart}
                                onChange={(e) => setCustomStart(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                End Date
                            </label>
                            <input
                                type="date"
                                value={customEnd}
                                onChange={(e) => setCustomEnd(e.target.value)}
                                min={customStart}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                            />
                        </div>
                    </div>

                    <button
                        onClick={handleCustomRange}
                        disabled={!customStart || !customEnd}
                        className="mt-3 w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                    >
                        Apply Custom Range
                    </button>
                </div>
            )}

            {/* Current Selection Display */}
            <div className="mt-4 text-sm text-gray-600">
                <span className="font-medium">Selected:</span> {formatDateRange(value)}
            </div>
        </div>
    );
};
