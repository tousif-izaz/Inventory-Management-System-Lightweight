# Inventory Management System - Frontend UI Plan

## Design Overview

### Design System & Theme
**Current CSS Theme:**
- Font: Poppins (Google Fonts)
- Primary Color: Black (#000000)
- Secondary/Hover: Dark Gray (#434343)
- Accent: Indigo (#4F46E5 - indigo-600)
- Framework: TailwindCSS
- Form Elements: Rounded corners, ring borders, indigo focus states
- Buttons: Black background with dark gray hover states

**Design Principles:**
1. Clean, minimalist interface
2. Consistent spacing and typography
3. Clear visual hierarchy
4. Mobile-responsive design
5. Accessible color contrast
6. Professional business aesthetic

---

## Navigation Structure

### Main Navigation (Sidebar + Top Header)

**Sidebar Navigation:**
- Fixed left sidebar (collapsible on mobile)
- Width: 260px (desktop), 64px (collapsed), full-width (mobile drawer)
- Background: Black (#000000)
- Text: White with gray-400 for inactive items
- Active state: Indigo-600 background with white text
- Hover state: Dark gray (#434343)

**Top Header:**
- Fixed top bar
- Background: White with bottom border
- Contains: Page title, search bar, user menu, notifications
- Height: 64px
- Shadow: subtle drop shadow

**Navigation Sections:**

1. **Dashboard**
   - Icon: Chart/Dashboard icon
   - Route: `/dashboard`

2. **Inventory Management**
   - Icon: Box/Package icon
   - Submenu:
     - All Products (`/inventory/products`)
     - Low Stock Alert (`/inventory/low-stock`)
     - Expiring Soon (`/inventory/expiring`)
     - Categories (`/inventory/categories`)
     - Locations (`/inventory/locations`)

3. **Transactions**
   - Icon: Exchange/Transfer icon
   - Submenu:
     - Purchases (`/transactions/purchases`)
     - Sales (`/transactions/sales`)
     - Adjustments (`/transactions/adjustments`)
     - Transfers (`/transactions/transfers`)
     - Transaction History (`/transactions/history`)

4. **Parties**
   - Icon: Users icon
   - Submenu:
     - Suppliers (`/parties/suppliers`)
     - Customers (`/parties/customers`)

5. **Reports**
   - Icon: Document/Report icon
   - Submenu:
     - Stock Summary (`/reports/stock-summary`)
     - Inventory by Location (`/reports/inventory-by-location`)
     - Sales Summary (`/reports/sales-summary`)
     - Purchase Summary (`/reports/purchase-summary`)
     - Low Stock Alerts (`/reports/low-stock-alerts`)
     - Expiring Products (`/reports/expiring-products`)

6. **Settings**
   - Icon: Cog/Settings icon
   - Route: `/settings`
   - Submenu:
     - User Management (`/settings/users`)
     - Profile (`/settings/profile`)
     - System Settings (`/settings/system`)

**User Menu (Top Right):**
- User avatar/icon
- Username display
- Dropdown menu:
  - Profile
  - Settings
  - Logout

---

## Page Layouts & Components

### 1. Dashboard (`/dashboard`)

**Layout:**
- Grid layout: 4 columns on desktop, 2 on tablet, 1 on mobile
- Cards with consistent styling

**Widgets:**
1. **Quick Stats Cards** (Top Row - 4 cards)
   - Total Products Count
   - Low Stock Items Count
   - Total Inventory Value
   - Today's Sales

2. **Charts & Graphs** (Middle Section)
   - Sales trend line chart (last 30 days)
   - Top selling products (bar chart)
   - Inventory distribution by category (pie chart)
   - Purchase vs Sales comparison

3. **Alerts & Notifications** (Right Sidebar/Bottom)
   - Low stock alerts (red badge)
   - Expiring products (orange badge)
   - Pending purchases (yellow badge)
   - Recent activity feed

4. **Quick Actions** (Floating Action Button or Top Bar)
   - New Sale
   - New Purchase
   - Add Product
   - Adjust Inventory

---

### 2. Inventory Pages

#### A. All Products (`/inventory/products`)

**Layout:**
- Search bar (by name, SKU, barcode)
- Filters: Category, Status (active/inactive), Stock level
- Sort options: Name, SKU, Stock, Price
- Action buttons: Add New Product, Export CSV

**Table Columns:**
- Thumbnail (if available)
- Product Name
- SKU
- Barcode
- Category
- Total Stock
- Cost Price
- Selling Price
- Status
- Actions (View, Edit, Delete)

**Product Details Modal/Page:**
- Product information
- Inventory levels by location
- Transaction history
- Edit button (manager+)

#### B. Low Stock Alert (`/inventory/low-stock`)

**Layout:**
- Similar to All Products but filtered
- Red/Orange badges for urgency
- Reorder button per product
- Shows: Product, Current Stock, Reorder Point, Suggested Reorder Quantity

#### C. Expiring Soon (`/inventory/expiring`)

**Layout:**
- Table with expiry dates
- Color-coded by urgency (red < 7 days, orange < 30 days)
- Batch information
- Quick action: Create sale/adjustment

#### D. Categories (`/inventory/categories`)

**Layout:**
- Tree/hierarchical view
- Parent-child relationships visible
- Add, Edit, Delete actions
- Product count per category

#### E. Locations (`/inventory/locations`)

**Layout:**
- Card or list view
- Location type badges (warehouse, store, shelf, zone)
- Active/inactive status
- Product count per location
- Add New Location button

---

### 3. Transaction Pages

#### A. Purchases (`/transactions/purchases`)

**Layout:**
- List view with filters (date range, supplier, payment status)
- Summary cards: Total purchases, Pending payments, Paid amount

**Table Columns:**
- Purchase Date
- Invoice Number
- Supplier Name
- Total Amount
- Payment Status (badge)
- Payment Method
- Actions (View, Edit, Mark Paid)

**Create Purchase Form:**
- Multi-step or single-page form
- Supplier selection
- Line items table (add/remove products)
- Auto-calculation of totals
- Tax and discount inputs
- Location assignment per item
- Batch and expiry date inputs

#### B. Sales (`/transactions/sales`)

**Layout:**
- Similar to Purchases
- POS-style interface for creating sales
- Customer selection (optional)
- Real-time inventory check
- Receipt generation

**Table Columns:**
- Sale Date
- Receipt Number
- Customer Name
- Total Amount
- Payment Status
- Payment Method
- Sold By
- Actions (View, Refund)

**POS Interface:**
- Product search/barcode scan
- Shopping cart
- Quick calculation of totals
- Payment methods selection
- Print receipt

#### C. Adjustments (`/transactions/adjustments`)

**Layout:**
- Form with reason selection
- Product and location selection
- Quantity adjustment (positive or negative)
- Reason dropdown (damage, theft, correction, etc.)
- Notes field
- Performed by (auto-filled)

**History Table:**
- Date/Time
- Product
- Location
- Quantity Changed
- Reason
- Performed By
- Notes

#### D. Transfers (`/transactions/transfers`)

**Layout:**
- Transfer form with from/to location selection
- Product and quantity
- Reason field
- Visual indicator of stock levels at both locations

**Transfer History:**
- Date/Time
- Product
- From Location → To Location
- Quantity
- Performed By

---

### 4. Parties Pages

#### A. Suppliers (`/parties/suppliers`)

**Layout:**
- Card or table view
- Search and filter (active/inactive)
- Add New Supplier button

**Supplier Card/Row:**
- Name
- Contact Person
- Email, Phone
- City, Country
- Active status
- Actions (View, Edit, Deactivate)

**Supplier Details Page:**
- Full supplier information
- Purchase history
- Outstanding payments
- Performance metrics

#### B. Customers (`/parties/customers`)

**Layout:**
- Similar to Suppliers
- Loyalty points display

**Customer Details Page:**
- Customer information
- Purchase history
- Total spent
- Loyalty points
- Actions (Edit, Adjust Points)

---

### 5. Reports Pages

**Common Report Layout:**
- Date range picker
- Filter options (specific to report type)
- Export buttons (PDF, CSV, Excel)
- Data visualization (charts/graphs)
- Detailed table below

#### A. Stock Summary (`/reports/stock-summary`)
- Uses `vw_product_stock_summary` view
- Shows: Product, SKU, Total Stock, Reserved, Available, Stock Status
- Color-coded status indicators

#### B. Inventory by Location (`/reports/inventory-by-location`)
- Uses `vw_inventory_by_location` view
- Grouped by location
- Subtotals per location

#### C. Sales Summary (`/reports/sales-summary`)
- Date range filter
- Group by: Day, Week, Month
- Line chart showing trends
- Table with breakdown

#### D. Purchase Summary (`/reports/purchase-summary`)
- Date range and supplier filters
- Similar layout to Sales Summary

#### E. Low Stock Alerts (`/reports/low-stock-alerts`)
- Real-time data
- Actionable items with "Create Purchase Order" button

#### F. Expiring Products (`/reports/expiring-products`)
- Days filter (default 30)
- Color-coded urgency
- Quick actions

---

### 6. Settings Pages

#### A. User Management (`/settings/users`)
- Admin only
- List of users
- Role management
- Add, Edit, Deactivate users

#### B. Profile (`/settings/profile`)
- Current user profile
- Change password
- Preferences

#### C. System Settings (`/settings/system`)
- Admin only
- General settings
- Tax settings
- Currency settings
- Backup & restore

---

## Reusable Components

### 1. **Card Component**
- White background
- Rounded corners (rounded-lg)
- Shadow (shadow-sm)
- Padding (p-6)
- Optional header with title

### 2. **Table Component**
- Striped rows
- Hover states
- Sortable columns
- Pagination
- Selection checkboxes (for bulk actions)

### 3. **Form Components**
- Input fields (text, number, date, select)
- Consistent styling with existing SignUp form
- Validation states (error, success)
- Helper text

### 4. **Button Component**
- Primary: Black bg, white text
- Secondary: White bg, black border, black text
- Danger: Red bg, white text
- Success: Green bg, white text
- Sizes: sm, md, lg

### 5. **Modal Component**
- Overlay with backdrop blur
- Centered content
- Close button
- Responsive sizing

### 6. **Badge Component**
- Status indicators (active, inactive, pending, paid, etc.)
- Color-coded (green, red, yellow, blue, gray)
- Small, rounded

### 7. **Alert/Toast Component**
- Success, Error, Warning, Info
- Auto-dismiss option
- Position: top-right

### 8. **Search Bar Component**
- Icon on left
- Clear button on right
- Debounced search
- Keyboard shortcuts (Cmd/Ctrl + K)

### 9. **Date Picker Component**
- Calendar interface
- Range selection support
- Quick presets (Today, This Week, This Month, etc.)

### 10. **Pagination Component**
- Page numbers
- Previous/Next buttons
- Items per page selector
- Total count display

---

## Color Palette (Tailwind Classes)

### Primary Colors
- Black: `bg-black`, `text-black`
- Dark Gray: `bg-[#434343]`, `hover:bg-[#434343]`
- White: `bg-white`, `text-white`

### Accent Colors
- Indigo (Primary Actions): `bg-indigo-600`, `text-indigo-600`, `ring-indigo-600`
- Gray Shades: `text-gray-500`, `text-gray-700`, `text-gray-900`, `bg-gray-50`, `bg-gray-100`

### Status Colors
- Success: `bg-green-500`, `text-green-600`, `border-green-600`
- Error/Danger: `bg-red-500`, `text-red-600`, `border-red-600`
- Warning: `bg-yellow-500`, `text-yellow-600`, `border-yellow-600`
- Info: `bg-blue-500`, `text-blue-600`, `border-blue-600`

### Border & Ring Colors
- Default: `ring-gray-300`, `border-gray-300`
- Focus: `ring-indigo-600`, `focus:ring-2`

---

## Responsive Breakpoints

- **Mobile**: < 640px (sm)
- **Tablet**: 640px - 1024px (sm - lg)
- **Desktop**: > 1024px (lg+)

**Navigation Behavior:**
- Mobile: Hamburger menu, drawer navigation
- Tablet: Collapsible sidebar
- Desktop: Full sidebar

---

## Authentication & Authorization

### Role-Based Access Control
- **Admin**: Full access
- **Manager**: All except user management
- **Staff**: Sales, customers, basic inventory view
- **User/Viewer**: Read-only access

**UI Behavior:**
- Hide/disable actions based on user role
- Show permission denied messages when appropriate
- Redirect to login if unauthenticated

---

## Performance Considerations

1. **Lazy Loading**: Code-split routes
2. **Pagination**: Limit data fetch to 20-50 items per page
3. **Debouncing**: Search inputs debounced by 300ms
4. **Caching**: Cache frequently accessed data (categories, locations)
5. **Optimistic Updates**: Update UI before server response for better UX

---

## Accessibility

1. **Keyboard Navigation**: All interactive elements accessible via keyboard
2. **ARIA Labels**: Proper labeling for screen readers
3. **Focus States**: Clear visual focus indicators
4. **Color Contrast**: WCAG AA compliance
5. **Alt Text**: Images and icons have descriptive alt text

---

## Icons

**Icon Library**: Heroicons (matches TailwindCSS ecosystem)

**Key Icons:**
- Dashboard: ChartBarIcon
- Inventory: CubeIcon
- Transactions: ArrowsRightLeftIcon
- Suppliers/Customers: UsersIcon
- Reports: DocumentChartBarIcon
- Settings: CogIcon
- Add: PlusIcon
- Edit: PencilIcon
- Delete: TrashIcon
- Search: MagnifyingGlassIcon
- Notifications: BellIcon
- Menu: Bars3Icon

---

## Next Steps

1. **Phase 1**: Implement Navigation Structure
   - Create Sidebar component
   - Create Header component
   - Create Layout wrapper component
   - Add routing structure

2. **Phase 2**: Build Core Components Library
   - Card, Table, Form, Button, Modal, Badge, Alert

3. **Phase 3**: Implement Dashboard
   - Stats cards
   - Charts integration (recharts or chart.js)
   - Activity feed

4. **Phase 4**: Build Inventory Pages
   - Products list and CRUD
   - Categories management
   - Locations management

5. **Phase 5**: Build Transaction Pages
   - Purchase flow
   - Sales/POS interface
   - Adjustments and transfers

6. **Phase 6**: Build Parties Pages
   - Suppliers and customers management

7. **Phase 7**: Build Reports
   - Report layouts
   - Data visualization
   - Export functionality

8. **Phase 8**: Settings & User Management
   - Profile management
   - User administration
   - System settings

---

## API Integration Pattern

**Standard Pattern for All Pages:**

```typescript
// API service layer
const apiService = {
  get: (endpoint: string) => fetch(`/api${endpoint}`, {
    credentials: 'include' // for JWT cookie
  }),
  post: (endpoint: string, data: any) => fetch(`/api${endpoint}`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  put: (endpoint: string, data: any) => fetch(`/api${endpoint}`, {
    method: 'PUT',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  delete: (endpoint: string) => fetch(`/api${endpoint}`, {
    method: 'DELETE',
    credentials: 'include'
  })
};

// Error handling
const handleApiError = (error: any) => {
  if (error.status === 401) {
    // Redirect to login
  } else if (error.status === 403) {
    // Show permission denied
  } else {
    // Show generic error
  }
};
```

---

## State Management

**Approach**: React Context API + Custom Hooks

**Contexts to Create:**
1. **AuthContext**: User authentication state, role, permissions
2. **NotificationContext**: Toast messages, alerts
3. **ThemeContext**: Dark mode support (future)
4. **AppStateContext**: Global app state (sidebar collapsed, etc.)

---

## Testing Strategy

1. **Unit Tests**: Component logic, utility functions
2. **Integration Tests**: API integration, form submissions
3. **E2E Tests**: Critical user flows (create sale, create purchase)

---

## Deployment Considerations

- Environment-specific API URLs
- Error logging and monitoring
- Analytics integration
- Progressive Web App (PWA) support (future)
