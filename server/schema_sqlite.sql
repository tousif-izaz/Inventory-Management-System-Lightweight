-- Inventory Management System Database Schema (SQLite)

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    quantity INTEGER NOT NULL,
    category TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_name ON products(name);

-- Insert sample data (optional)
-- Uncomment the lines below if you want some test data

-- Sample admin user (password is 'admin123' - you should hash this properly in production)
-- INSERT INTO users (username, password, role) VALUES ('admin', '$2a$10$...[hashed_password]', 'admin');

-- Sample products
-- INSERT INTO products (name, price, quantity, category) VALUES
-- ('Laptop', 999.99, 50, 'Electronics'),
-- ('Mouse', 29.99, 200, 'Electronics'),
-- ('Desk Chair', 199.99, 30, 'Furniture'),
-- ('Notebook', 4.99, 500, 'Office Supplies');
