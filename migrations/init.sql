-- ===============================================
-- DATABASE DESIGN - STRUKIT
-- ===============================================

-- Enable UUID extension for PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_budget DECIMAL(15,2) DEFAULT 0,
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(7), -- hex color code
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    -- Receipt Image
    image_url TEXT NOT NULL,
    image_size INTEGER, -- in bytes
    image_width INTEGER,
    image_height INTEGER,
    
    -- Extracted Data
    merchant_name VARCHAR(255),
    total_amount DECIMAL(15,2) NOT NULL,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    transaction_date DATE NOT NULL,
    transaction_time TIME,
    
    -- AI Processing
    raw_extracted_text TEXT,
    extraction_confidence DECIMAL(3,2), -- 0.00 to 1.00
    ai_model_used VARCHAR(50) DEFAULT 'gemini',
    processing_status VARCHAR(20) DEFAULT 'pending' CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed')),
    
    -- User Verification
    is_verified BOOLEAN DEFAULT FALSE,
    verification_notes TEXT,
    
    -- Additional Info
    receipt_number VARCHAR(100),
    payment_method VARCHAR(50), -- cash, debit, credit, etc
    notes TEXT,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT positive_amount CHECK (total_amount > 0),
    CONSTRAINT valid_confidence CHECK (extraction_confidence >= 0 AND extraction_confidence <= 1)
);

CREATE TABLE receipt_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    
    -- Item Details
    item_name VARCHAR(255) NOT NULL,
    quantity INTEGER DEFAULT 1,
    unit_price DECIMAL(15,2),
    total_price DECIMAL(15,2) NOT NULL,
    
    -- Additional Info
    item_code VARCHAR(100),
    category VARCHAR(100),
    discount_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Metadata
    line_number INTEGER, -- urutan item dalam struk
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_total_price CHECK (total_price >= 0)
);

CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    allocated_amount DECIMAL(15,2) NOT NULL,
    spent_amount DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint
    UNIQUE(project_id, category_id),
    
    -- Constraints
    CONSTRAINT positive_allocated CHECK (allocated_amount >= 0),
    CONSTRAINT positive_spent CHECK (spent_amount >= 0)
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    
    -- Report Details
    title VARCHAR(255) NOT NULL,
    report_type VARCHAR(50) NOT NULL, -- 'summary', 'detailed', 'category_breakdown'
    date_from DATE,
    date_to DATE,
    
    -- Generated Files
    pdf_url TEXT,
    excel_url TEXT,
    
    -- Report Data (JSON)
    report_data JSONB,
    
    -- Metadata
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE, -- for cleanup
    download_count INTEGER DEFAULT 0
);


-- Receipt indexes (most queried table)
CREATE INDEX idx_receipts_user_id ON receipts(user_id);
CREATE INDEX idx_receipts_project_id ON receipts(project_id);
CREATE INDEX idx_receipts_transaction_date ON receipts(transaction_date DESC);
CREATE INDEX idx_receipts_merchant ON receipts(merchant_name);
CREATE INDEX idx_receipts_amount ON receipts(total_amount);
CREATE INDEX idx_receipts_created_at ON receipts(created_at DESC);
CREATE INDEX idx_receipts_user_project_date ON receipts(user_id, project_id, transaction_date);

-- Receipt items index
CREATE INDEX idx_receipt_items_receipt_id ON receipt_items(receipt_id);
CREATE INDEX idx_receipt_items_name ON receipt_items(item_name);

-- Projects index
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);

-- Reports index
CREATE INDEX idx_reports_user_id ON reports(user_id);
CREATE INDEX idx_reports_project_id ON reports(project_id);

-- Function untuk update timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_receipts_updated_at BEFORE UPDATE ON receipts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
