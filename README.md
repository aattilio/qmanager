# QManager: The Hypervisor Revolution

QManager is a high-performance orchestrator for KVM, engineered for technical excellence and absolute control. It provides a native, low-overhead interface to the Linux Kernel via `libvirt`, delivering a Type-1 hypervisor management experience.

## Key Architectures

- **Go-Native Backend:** Leveraging Go's concurrency and system-level performance for robust infrastructure management.
- **Hardware-Accelerated UI:** Built with Fyne (v2), the interface is rendered via OpenGL/Vulkan for extreme responsiveness.
- **Atomic Modularity:** Every feature—from storage management to dynamic media resolution—is built as an isolated, structural module.
- **Smart Orchestration:** Professional automation for OS deployment, including dynamic mirror resolution and optimized hardware defaults.

## Project Structure

- `src/backend/hypervisor`: Direct kernel/libvirt orchestration logic.
- `src/backend/discovery`: Dynamic ISO scraping and mirror resolution engine.
- `src/backend/filesystem`: Virtual disk (QCOW2) and storage lifecycle management.
- `src/backend/provisioning`: Asynchronous media handling and VM creation.
- `config/catalog`: Modular XML-based OS metadata database.

---
*QManager: Manage your infrastructure, not your tools.*
