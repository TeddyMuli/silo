CREATE DATABASE storify;

\c storify;

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
