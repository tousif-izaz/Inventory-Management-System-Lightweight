import { Select } from '../UI/Select';
import { Category } from '../../types/product';

interface ProductFiltersProps {
    categoryFilter: string;
    statusFilter: string;
    categories: Category[];
    onCategoryChange: (value: string) => void;
    onStatusChange: (value: string) => void;
}

export const ProductFilters = ({
    categoryFilter,
    statusFilter,
    categories,
    onCategoryChange,
    onStatusChange,
}: ProductFiltersProps) => {
    return (
        <div className="flex flex-col sm:flex-row gap-4">
            <Select
                value={categoryFilter}
                onChange={(e) => onCategoryChange(e.target.value)}
                options={[
                    { value: '', label: 'All Categories' },
                    ...categories.map((c) => ({ value: c.category_id.toString(), label: c.name })),
                ]}
                className="w-full sm:w-48"
            />

            <Select
                value={statusFilter}
                onChange={(e) => onStatusChange(e.target.value)}
                options={[
                    { value: '', label: 'All Status' },
                    { value: 'active', label: 'Active' },
                    { value: 'inactive', label: 'Inactive' },
                ]}
                className="w-full sm:w-48"
            />
        </div>
    );
};
