CREATE DATABASE aethly;

\c aethly;

-- Create Users Table
CREATE TABLE IF NOT EXISTS Users (
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(50),
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Organizations Table
CREATE TABLE IF NOT EXISTS Organizations (
    organization_id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TYPE role_enum AS ENUM ('creator', 'admin', 'member');

-- Create UserOrganizations Table
CREATE TABLE IF NOT EXISTS UserOrganizations (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) ON DELETE CASCADE,
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    role role_enum NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Folders Table
CREATE TABLE IF NOT EXISTS Folders (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    parent_folder_id UUID REFERENCES Folders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ
);

-- Create Files Table
CREATE TABLE IF NOT EXISTS Files (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    folder_id UUID REFERENCES Folders(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ
);

-- Create Fleets Table
CREATE TABLE IF NOT EXISTS Fleets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Devices Table
CREATE TABLE IF NOT EXISTS Devices (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    fleet_id UUID REFERENCES Fleets(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    ip_address VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create UserDevices Table
CREATE TABLE IF NOT EXISTS UserDevices (
    id UUID PRIMARY KEY,
    user_organization_id UUID REFERENCES UserOrganizations(id) ON DELETE CASCADE,
    device_id UUID REFERENCES Devices(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Plans Table
CREATE TABLE IF NOT EXISTS Plans (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    interval VARCHAR(50),
    storage_limit BIGINT
);

-- Create Subscriptions Table
CREATE TABLE IF NOT EXISTS Subscriptions (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    plan_id UUID REFERENCES Plans(id) ON DELETE CASCADE,
    status VARCHAR(50),
    stripe_subscription_id VARCHAR(255),
    start_date DATE,
    end_date DATE,
    paypal_subscription_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
