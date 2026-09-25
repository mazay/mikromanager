# Agents in MikroManager

This document describes the agent architecture and components used in MikroManager.

## Overview

MikroManager is an application designed to manage MikroTik devices, monitor network configurations, and ensure proper backup procedures. The system employs agents for various tasks including device polling, configuration exports, and S3 backup management.

## Agent Components

### 1. Device Poller Agent
- **Purpose**: Periodically polls MikroTik devices for status updates
- **Configuration**:
  - `apiPollers`: Number of concurrent API pollers (default: 2)
  - `devicePollerInterval`: Interval between polls (default: 300 seconds)
- **Functionality**: 
  - Uses MikroTik API to retrieve device information
  - Applies comment-based filtering to find management IPs
  - Updates device IP addresses when subnets change

### 2. Export Agent
- **Purpose**: Creates configuration exports for devices
- **Configuration**:
  - `exportWorkers`: Number of concurrent export workers (default: 10)
  - `deviceExportCronSchedule`: Cron schedule for exports (default: `0 * * * *`)
- **Functionality**:
  - Exports device configurations to S3 storage
  - Manages backup retention and cleanup

### 3. S3 Backup Agent
- **Purpose**: Handles all S3-related operations for backups
- **Configuration**:
  - `s3Bucket`: S3 bucket name for backups (default: `mikromanager-backups`)
  - `s3Region`: AWS region for the S3 bucket
  - `s3StorageClass`: Storage class for backups (default: `STANDARD`)
  - `s3OpsRetries`: Number of retries for S3 operations (default: 5)
- **Functionality**:
  - Uploads exports to S3 buckets
  - Manages backup lifecycle and metadata
  - Handles credentials through environment variables or config

### 4. Database Agent
- **Purpose**: Manages persistent storage for device information
- **Configuration**:
  - `dbPath`: Path to the database file (default: `database/mikromanager.db`)
  - `dbLogLevel`: GORM log level (default: `silent`)
- **Functionality**:
  - Stores device configurations and status updates
  - Manages encryption of sensitive data

## System Integration

### Docker Deployment
To run MikroManager with proper agent management:

```bash
docker run -d \
    -p 8000:8000 \ # port mapping
    -v /mikromanager:/app \ # the persistence directory, should contain the config and "database" subdirectory
    -e "AWS_ACCESS_KEY_ID=AKIAIO...MPLE" \
    -e "AWS_SECRET_ACCESS_KEY=*** \
    ghcr.io/mazay/mikromanager:main
```

## Configuration

All agents are configurable via `config.yml`:
- Sensitive data is encrypted using an encryption key
- API credentials can be provided through environment variables (preferred approach)
- Database and backup paths are configurable

## Testing Strategy

Unit tests exist for core components including:
- Device polling functionality 
- Export processing
- S3 integration
- Database operations

Integration tests cover end-to-end scenarios including:
- Device connection and status monitoring
- Configuration export and backup
- S3 upload procedures

## Security Notes

- All sensitive data (credentials, encryption keys) should be stored in environment variables
- Default credentials provided in config file should never be used in production
- Encryption key must be changed and kept persistent for proper decryption of sensitive data