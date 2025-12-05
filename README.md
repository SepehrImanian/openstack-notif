# openstack-notif

Opinionated toolkit for provisioning an OpenStack environment and wiring it to a lightweight VM monitoring service that pushes notifications when interesting events happen.

The repository brings together:

- **Terraform** infrastructure definitions (`terrafrom/`)
- **Kolla-Ansible** values and overrides (`ansible-kolla-values/`)
- A **Go-based VM monitor/agent** (`vm-monitor/`)

---

## Table of contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Repository structure](#repository-structure)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
  - [1. Clone the repository](#1-clone-the-repository)
  - [2. Provision infrastructure with Terraform](#2-provision-infrastructure-with-terraform)
  - [3. Deploy OpenStack with Kolla-Ansible](#3-deploy-openstack-with-kolla-ansible)
  - [4. Build and run the VM monitor](#4-build-and-run-the-vm-monitor)

---

## Overview

`openstack-notif` is intended as a practical example and starter kit for:

1. **Provisioning** the underlying infrastructure for an OpenStack deployment using Terraform.
2. **Deploying** OpenStack services via Kolla-Ansible with opinionated values.
3. **Monitoring** virtual machine lifecycle events using a small Go service.
4. **Notifying** external systems (chat, email, webhooks, etc.) about those events.

