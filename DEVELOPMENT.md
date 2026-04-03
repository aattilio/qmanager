# QManager Development & Build Guide

This document provides the necessary instructions to set up the development environment, compile, and deploy QManager from source.

## System Requirements

- **Go Compiler** (Version 1.21 or later)
- **Libvirt** (Development headers for C-bindings)
- **GCC / G++** (Standard build toolchain)
- **Pkg-config**

## Engineering Standards & Git Workflow

### 1. Feature-Branch Workflow
**Direct pushes to `main` are strictly prohibited.** All development must occur on feature branches. 
- Branches must follow the naming convention: `feat/feature-name` or `fix/bug-name`.
- Integration into `main` must be performed via **Pull Request (PR)** only.

### 2. Balanced Multiline Formatting
We avoid cluttered inline code. The goal is professional readability.
- **Conditionals:** Body must always be on a new line.
- **Function Calls:** Use multiline if there are more than 3 parameters or if the line exceeds 80 characters.
- **Initializers:** Every member assignment must reside on a unique line for complex objects.

### 3. Absolute Modularity
The project is organized into atomic modules. Layers must interact through defined interfaces to prevent coupling.

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
**Mandatory Verification:** All tests must pass before a PR is opened.
```bash
make test
```
The media integrity tests validate multiple mirrors per OS and perform memory cleanup after each verification.

## Post-Build Configuration

### Hypervisor Permissions
```bash
sudo usermod -aG libvirt $USER
newgrp libvirt
```
