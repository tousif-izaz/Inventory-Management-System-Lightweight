import { PlusIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';

interface ProductActionsProps {
    onAddProduct: () => void;
    onExport: () => void;
}

export const ProductActions = ({
    onAddProduct,
    onExport,
}: ProductActionsProps) => {
    return (
        <div className="flex gap-2">
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
