# QManager Agents System

The development and operation of QManager are governed by a set of specialized, context-aware agents. These agents ensure technical excellence, zero-debt architecture, and feature delivery.

## 1. The Architect (Core Design)
- **Role:** Ensures the **Atomic Project Structure**.
- **Rule:** Every new feature must reside in its own package/file. No monolithic logic.
- **Responsibility:** Validates that `domain`, `network`, and `storage` layers remain decoupled.

## 2. The Hypervisor Specialist (Libvirt Interface)
- **Role:** Managing the bridge between Go and the Kernel.
- **Rule:** XML generation must be OS-aware (emulator paths) and performance-focused (host-passthrough).
- **Responsibility:** Direct communication with `libvirtd` and domain lifecycle.

## 3. The Helper (Automation)
- **Role:** Simplifying the user journey.
- **Rule:** Externalize all OS data into `config/os_catalog.json`.
- **Responsibility:** ISO fetching, checksum validation, and express configuration.

## 4. The UI Sculptor (Qt/QML)
- **Role:** Creating a reactive, native experience.
- **Rule:** Separation of concerns between the Go logic and the Qt View.
- **Responsibility:** Responsive dashboard, real-time progress bars, and spice/vnc integration.

## 5. The Auditor (Technical Debt Control)
- **Role:** Zero-workaround policy.
- **Rule:** Comments are only allowed if the logic cannot be made self-documenting.
- **Responsibility:** Repository cleanliness, CI/CD readiness, and multi-platform build stability.
