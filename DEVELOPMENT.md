# QManager Development & Build Guide

This document provides the necessary instructions to set up the development environment, compile, and deploy QManager from source.

## System Requirements

- **Go Compiler** (Version 1.21 or later)
- **Qt 5.15** (Development libraries and headers)
- **Libvirt** (Development headers for C-bindings)
- **GCC / G++** (Standard build toolchain supporting C++17)
- **Pkg-config**

## Engineering Standards & Best Practices

### 1. Balanced Multiline Formatting
We avoid both cluttered inline code and extreme vertical wrapping. The goal is professional readability.

- **Conditionals:** Body must always be on a new line. 
  - *Allowed:*
    ```go
    if err != nil {
        return err
    }
    ```
  - *Forbidden:* `if err != nil { return err }`
- **Function Calls/Signatures:** Use multiline if there are more than 3 parameters or if the line exceeds 80 characters.
- **Structs/Objects:** Initializations with multiple members must be multiline.
- **No Extreme Wrapping:** Do not wrap single simple parameters or simple conditions unless necessary for logic separation.

### 2. Semantic Naming
Every file, variable, and function must be descriptive. Avoid abbreviations like `vm`, `cfg`, `ptr`.
- **Bad:** `vm_client.go`, `var c *Lvirt`
- **Good:** `libvirt_hypervisor_connector.go`, `var hypervisorConnector *LibvirtHypervisorConnector`

### 3. Import Organization
1. Standard Library
2. Internal Project
3. Third-party Library
(Grouped with a blank line between them).

### 4. Absolute Modularity
The project is organized into atomic modules: `hypervisor`, `provisioning`, `filesystem`, `discovery`. Layers must interact through the defined `api/bridge` layer.

## Build Instructions

### 1. Automated Setup
```bash
make setup
```

### 2. Compilation (Production)
```bash
make build-production-linux
```

### 3. Testing
```bash
make test
```

## Post-Build Configuration

### Hypervisor Permissions
```bash
sudo usermod -aG libvirt $USER
newgrp libvirt
```
