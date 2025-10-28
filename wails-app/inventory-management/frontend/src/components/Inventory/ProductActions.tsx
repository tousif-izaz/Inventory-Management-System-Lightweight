import { PlusIcon, ArrowDownTrayIcon, TrashIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';

interface ProductActionsProps {
    selectedCount: number;
    onAddProduct: () => void;
    onExport: () => void;
    onBulkDelete: () => void;
    isDeleting: boolean;
}

export const ProductActions = ({
    selectedCount,
    onAddProduct,
    onExport,
    onBulkDelete,
    isDeleting,
}: ProductActionsProps) => {
    return (
        <div className="flex gap-2">
            {selectedCount > 0 && (
                <Button variant="danger" onClick={onBulkDelete} isLoading={isDeleting}>
                    <TrashIcon className="h-5 w-5 mr-2" />
                    Delete Selected ({selectedCount})
                </Button>
            )}
            <Button onClick={onAddProduct}>
                <PlusIcon className="h-5 w-5 mr-2" />
                Add Product
            </Button>
            <Button variant="secondary" onClick={onExport}>
                <ArrowDownTrayIcon className="h-5 w-5 mr-2" />
                Export
            </Button>
        </div>
    );
};
