CREATE TABLE "user" (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT,
    display_name VARCHAR(32) NOT NULL,
    discord_token TEXT,
    discord_refresh_token TEXT,
    discord_user_id VARCHAR(255),
    profile VARCHAR(500),
    avatar_url TEXT,
    twitter_id TEXT,
    github_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE visibility AS ENUM (
    'public',
    'private',
    'draft'
);

CREATE TABLE work (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    visibility visibility,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorite (
    work_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (work_id, user_id)
);

CREATE TYPE assettype AS ENUM (
    'zip',
    'image',
    'video',
    'music',
    'model'
);

CREATE TABLE asset (
    id VARCHAR(255) PRIMARY KEY,
    work_id VARCHAR(255),
    asset_type assettype NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comment (
    id VARCHAR(255) PRIMARY KEY,
    content TEXT NOT NULL,
    work_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255),
    reply_at VARCHAR(255),
    visibility visibility NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE urltype AS ENUM (
    'youtube',
    'soundcloud', 
    'github',
    'sketchfab',
    'unityroom',
    'other'
);

CREATE TABLE urlinfo (
    id VARCHAR(255) PRIMARY KEY,
    work_id VARCHAR(255),
    url VARCHAR(255) NOT NULL,
    url_type urltype NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tag (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tagging (
    work_id VARCHAR(255) NOT NULL,
    tag_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (work_id, tag_id)
);

CREATE TABLE thumbnail (
    work_id VARCHAR(255) NOT NULL,
    asset_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (work_id, asset_id)
);
