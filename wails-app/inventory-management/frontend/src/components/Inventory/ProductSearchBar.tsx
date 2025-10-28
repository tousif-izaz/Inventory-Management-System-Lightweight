import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';

interface ProductSearchBarProps {
    searchTerm: string;
    onSearchChange: (value: string) => void;
}

export const ProductSearchBar = ({ searchTerm, onSearchChange }: ProductSearchBarProps) => {
    return (
        <div className="flex-1 relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
                type="text"
                placeholder="Search by name, SKU, or barcode..."
                value={searchTerm}
                onChange={(e) => onSearchChange(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-600 focus:border-transparent"
            />
        </div>
    );
};
