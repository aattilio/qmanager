# QManager: The Hypervisor Revolution

QManager is a high-performance, professional-grade orchestrator for **QEMU/KVM**, engineered for developers and systems architects who demand absolute control with zero overhead.

By interfacing directly with the Linux Kernel via **libvirt**, QManager delivers a Type-1 hypervisor experience through a modern, reactive interface built with **Go** and **Qt/QML**.

## Core Philosophy

- **Native Execution:** Eliminate middle layers. QManager communicates directly with KVM/QEMU, ensuring near-native performance for your virtualized workloads.
- **Atomic Reliability:** Every feature is built as a decoupled, specialized module, ensuring a codebase that is robust, testable, and free of technical debt.
- **Express Automation:** Inspired by the simplicity of consumer tools but powered by professional logic. Select from over 100 operating systems, and QManager handles the dynamic ISO fetching and hardware optimization instantly.
- **Professional Scalability:** Engineered for high-concurrency environments, supporting advanced storage (QCOW2), networking (NAT/Bridge), and hardware-accelerated graphics (SPICE/VirtIO).

## Key Features

- **Dynamic ISO Scraper:** Intelligent fetching logic that resolves the latest installation media from official mirrors and university repositories.
- **Hardware Synergy:** Out-of-the-box support for `host-passthrough`, ensuring your VMs leverage the full power of your physical CPU.
- **Reactive Dashboard:** A beautiful, dark-mode interface designed for efficient management of multiple VM instances.
- **Cross-Platform Foundation:** Built with portability in mind, ready for deployment on Linux, macOS, and Windows.

## Documentation

For technical information on how to build, install, and contribute to the project, please refer to the specialized guides:

- **[Development & Build Guide](DEVELOPMENT.md)**: Instructions for compiling and setting up the development environment.
- **[Agents System](AGENTS.md)**: Overview of the architectural principles governing the project.

---
*QManager: Manage your infrastructure, not your tools.*
