/**
 * CSV Export Utility Functions
 * Provides functionality to export data to CSV format
 */

export interface CSVExportOptions {
    filename?: string;
    headers?: string[];
    data: any[];
    columns?: { key: string; header: string }[];
}

/**
 * Convert array of objects to CSV string
 */
function arrayToCSV(data: any[], columns?: { key: string; header: string }[]): string {
    if (data.length === 0) return '';

    let csvContent = '';

    // If columns are specified, use them; otherwise use all keys from first object
    if (columns && columns.length > 0) {
        // Add headers
        csvContent += columns.map(col => escapeCSVValue(col.header)).join(',') + '\n';

        // Add data rows
        data.forEach(row => {
            const values = columns.map(col => {
                const value = getNestedValue(row, col.key);
                return escapeCSVValue(formatValue(value));
            });
            csvContent += values.join(',') + '\n';
        });
    } else {
        // Use all keys from first object
        const keys = Object.keys(data[0]);

        // Add headers
        csvContent += keys.map(key => escapeCSVValue(key)).join(',') + '\n';

        // Add data rows
        data.forEach(row => {
            const values = keys.map(key => escapeCSVValue(formatValue(row[key])));
            csvContent += values.join(',') + '\n';
        });
    }

    return csvContent;
}

/**
 * Get nested value from object using dot notation
 */
function getNestedValue(obj: any, path: string): any {
    return path.split('.').reduce((current, key) => current?.[key], obj);
}

/**
 * Format value for CSV
 */
function formatValue(value: any): string {
    if (value === null || value === undefined) return '';
    if (value instanceof Date) return value.toISOString();
    if (typeof value === 'object') return JSON.stringify(value);
    return String(value);
}

/**
 * Escape CSV value (handle commas, quotes, newlines)
 */
function escapeCSVValue(value: string): string {
    const stringValue = String(value);

    // If value contains comma, quote, or newline, wrap in quotes and escape existing quotes
    if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
        return `"${stringValue.replace(/"/g, '""')}"`;
    }

    return stringValue;
}

/**
 * Download CSV file
 */
function downloadCSV(csvContent: string, filename: string): void {
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');

    if (link.download !== undefined) {
        const url = URL.createObjectURL(blob);
        link.setAttribute('href', url);
        link.setAttribute('download', filename);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    }
}

/**
 * Main export function
 */
export function exportToCSV(options: CSVExportOptions): void {
    const { data, columns, filename = 'export.csv' } = options;

    if (!data || data.length === 0) {
        console.warn('No data to export');
        return;
    }

    const csvContent = arrayToCSV(data, columns);
    downloadCSV(csvContent, filename);
}

/**
 * Utility function for common date formatting in exports
 */
export function formatDateForCSV(date: string | Date): string {
    if (!date) return '';
    const d = typeof date === 'string' ? new Date(date) : date;
    return d.toLocaleDateString('en-US', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
    });
}

/**
 * Utility function for currency formatting in exports
 */
export function formatCurrencyForCSV(amount: number): string {
    return amount.toFixed(2);
}
