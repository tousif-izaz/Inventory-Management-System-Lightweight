import './App.css';
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Root } from "./components/Root.tsx";
import { SignUp } from "./components/SignUp.tsx";
import { Login } from "./components/Login.tsx";
import { Layout } from "./components/Layout/Layout.tsx";
import { Dashboard } from "./components/Dashboard/Dashboard.tsx";

// Inventory components
import { Products } from "./components/Inventory/Products.tsx";
import { LowStock } from "./components/Inventory/LowStock.tsx";
import { Categories } from "./components/Inventory/Categories.tsx";

// Transaction components
import { Purchases } from "./components/Transactions/Purchases.tsx";
import { Sales } from "./components/Transactions/Sales.tsx";
import { Adjustments } from "./components/Transactions/Adjustments.tsx";
import { Transfers } from "./components/Transactions/Transfers.tsx";
import { TransactionHistory } from "./components/Transactions/History.tsx";

// Parties components
import { Suppliers } from "./components/Parties/Suppliers.tsx";
import { Customers } from "./components/Parties/Customers.tsx";

// Reports components
import { StockSummary } from "./components/Reports/StockSummary.tsx";
import { InventoryByLocation } from "./components/Reports/InventoryByLocation.tsx";
import { SalesSummary } from "./components/Reports/SalesSummary.tsx";
import { PurchaseSummary } from "./components/Reports/PurchaseSummary.tsx";
import { LowStockAlerts } from "./components/Reports/LowStockAlerts.tsx";
import { ExpiringProducts } from "./components/Reports/ExpiringProducts.tsx";

// New Analytics & Reports components
import { SalesAnalytics } from "./components/Reports/SalesAnalytics.tsx";
import { ProductAnalytics } from "./components/Reports/ProductAnalytics.tsx";
import { ABCAnalysisReport } from "./components/Reports/ABCAnalysisReport.tsx";
import { InventoryTurnoverReport } from "./components/Reports/InventoryTurnoverReport.tsx";

// Settings components
import { Settings } from "./components/Settings/Settings.tsx";
import { Profile } from "./components/Settings/Profile.tsx";
import { SquareSetup } from "./components/Settings/SquareSetup.tsx";
import { TerminalDevices } from "./components/Settings/TerminalDevices.tsx";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                {/* Public routes */}
                <Route path="/" element={<Root/>}/>
                <Route path="/login" element={<Login/>}/>
                <Route path="/signup" element={<SignUp/>}/>

                {/* Protected routes with navigation layout */}
                <Route element={<Layout/>}>
                    <Route path="/dashboard" element={<Dashboard/>}/>

                    {/* Inventory routes */}
                    <Route path="/inventory/products" element={<Products/>}/>
                    <Route path="/inventory/low-stock" element={<LowStock/>}/>
                    <Route path="/inventory/categories" element={<Categories/>}/>

                    {/* Transaction routes */}
                    <Route path="/transactions/purchases" element={<Purchases/>}/>
                    <Route path="/transactions/sales" element={<Sales/>}/>
                    <Route path="/transactions/adjustments" element={<Adjustments/>}/>
                    <Route path="/transactions/transfers" element={<Transfers/>}/>
                    <Route path="/transactions/history" element={<TransactionHistory/>}/>

                    {/* Parties routes */}
                    <Route path="/parties/suppliers" element={<Suppliers/>}/>
                    <Route path="/parties/customers" element={<Customers/>}/>

                    {/* Reports routes */}
                    <Route path="/reports/stock-summary" element={<StockSummary/>}/>
                    <Route path="/reports/inventory-by-location" element={<InventoryByLocation/>}/>
                    <Route path="/reports/sales-summary" element={<SalesSummary/>}/>
                    <Route path="/reports/purchase-summary" element={<PurchaseSummary/>}/>
                    <Route path="/reports/low-stock-alerts" element={<LowStockAlerts/>}/>
                    <Route path="/reports/expiring-products" element={<ExpiringProducts/>}/>

                    {/* Advanced Analytics & Reports routes */}
                    <Route path="/reports/sales-performance" element={<SalesAnalytics/>}/>
                    <Route path="/reports/product-analytics" element={<ProductAnalytics/>}/>
                    <Route path="/reports/abc-analysis" element={<ABCAnalysisReport/>}/>
                    <Route path="/reports/inventory-turnover" element={<InventoryTurnoverReport/>}/>

                    {/* Settings routes */}
                    <Route path="/settings" element={<Settings/>}/>
                    <Route path="/settings/profile" element={<Profile/>}/>
                    <Route path="/settings/square" element={<SquareSetup/>}/>
                    <Route path="/settings/terminal-devices" element={<TerminalDevices/>}/>
                </Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
