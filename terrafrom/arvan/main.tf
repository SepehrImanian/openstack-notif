terraform {
  required_providers {
    arvan = {
      source = "terraform.arvancloud.ir/arvancloud/iaas"
    }
  }
}

########################
# Variables
########################

variable "api_key" {
  type        = string
  description = "Arvan IaaS Machine User API"
}

variable "region" {
  type        = string
  description = "The chosen region for resources"
  default     = "eu-west1-a"
}

variable "chosen_distro_name" {
  type        = string
  description = "The chosen distro name for image"
  default     = "ubuntu"
}

variable "chosen_name" {
  type        = string
  description = "The chosen release for image"
  default     = "24.04"
}

variable "chosen_plan_id" {
  type        = string
  description = "The chosen ID of plan (flavor)"
  default     = "std-4-4-0"
}

variable "chosen_snapshot_id" {
  type        = string
  description = "The chosen ID of snapshot"
  default     = ""
}

variable "vm_count_control" {
  type        = number
  description = "Number of VMs control to create"
  default     = 3
}

variable "vm_count_compute" {
  type        = number
  description = "Number of VMs compute to create"
  default     = 3
}

variable "ssh_key_name" {
  type        = string
  description = "SSH key name"
  default     = "Sepehr"
}

provider "arvan" {
  api_key = var.api_key
}

########################
# Data sources
########################

data "arvan_security_groups" "default_security_groups" {
  region = var.region
}

data "arvan_images" "terraform_image" {
  region     = var.region
  image_type = "distributions"
}

data "arvan_plans" "plan_list" {
  region = var.region
}

locals {
  chosen_image = try(
    [
      for image in data.arvan_images.terraform_image.distributions : image
      if image.distro_name == var.chosen_distro_name && image.name == var.chosen_name
    ][0],
    null
  )

  selected_plan = try(
    [
      for plan in data.arvan_plans.plan_list.plans : plan
      if plan.id == var.chosen_plan_id
    ][0],
    null
  )
}

########################
# Private Network
########################

resource "arvan_network" "terraform_private_network" {
  region      = var.region
  description = "Terraform-created private network"
  name        = "tf-private-network"

  cidr       = "10.0.0.0/24"
  gateway_ip = "10.0.0.1"

  dhcp_range = {
    start = "10.0.0.2"
    end   = "10.0.0.254"
  }

  dns_servers    = ["8.8.8.8", "1.1.1.1", "4.2.2.4"]
  enable_dhcp    = true
  enable_gateway = true
}


########################
# Floating IPs (allocation)
########################

# resource "arvan_floating_ip" "vm_fip" {
#   count  = var.vm_count
#   region = var.region

#   description = "Terraform-created floating IP for VM index ${count.index + 1}"
# }


########################
# Abraks (instances)
########################

resource "arvan_abrak" "built_by_terraform_control" {
  depends_on = [
    arvan_network.terraform_private_network,
    # arvan_floating_ip.vm_fip
  ]

  count  = var.vm_count_control
  region = var.region

  timeouts {
    create = "1h30m"
    update = "2h"
    delete = "20m"
    read   = "10m"
  }

  name = "openstack-control-${count.index + 1}"

  image_id  = local.chosen_image.id
  flavor_id = local.selected_plan.id
  disk_size = 30

  snapshot_id = var.chosen_snapshot_id != "" ? var.chosen_snapshot_id : null

  enable_ipv4  = true
  enable_ipv6  = false
  ssh_key_name = var.ssh_key_name

  networks = [
    {
      network_id = arvan_network.terraform_private_network.network_id
    }
  ]

  security_groups = [
    data.arvan_security_groups.default_security_groups.groups[0].id
  ]


  # floating_ip = {
  #   floating_ip_id = arvan_floating_ip.vm_fip[count.index].id
  #   network_id     = arvan_network.terraform_private_network.network_id
  # }
}


resource "arvan_abrak" "built_by_terraform_compute" {
  depends_on = [
    arvan_network.terraform_private_network,
    # arvan_floating_ip.vm_fip
  ]

  count  = var.vm_count_compute
  region = var.region

  timeouts {
    create = "1h30m"
    update = "2h"
    delete = "20m"
    read   = "10m"
  }

  # Only letters, numbers, hyphens allowed
  name = "openstack-compute-${count.index + 1}"

  image_id  = local.chosen_image.id
  flavor_id = local.selected_plan.id
  disk_size = 50

  snapshot_id = var.chosen_snapshot_id != "" ? var.chosen_snapshot_id : null

  enable_ipv4  = true
  enable_ipv6  = false
  ssh_key_name = var.ssh_key_name

  networks = [
    {
      network_id = arvan_network.terraform_private_network.network_id
    }
  ]

  security_groups = [
    data.arvan_security_groups.default_security_groups.groups[0].id
  ]


  # floating_ip = {
  #   floating_ip_id = arvan_floating_ip.vm_fip[count.index].id
  #   network_id     = arvan_network.terraform_private_network.network_id
  # }
}

########################
# Outputs
########################

output "vm_names_control" {
  value       = [for a in arvan_abrak.built_by_terraform_control : a.name]
  description = "Names of all created VMs"
}

output "vm_names_compute" {
  value       = [for a in arvan_abrak.built_by_terraform_compute : a.name]
  description = "Names of all created VMs"
}

output "vm_ids_control" {
  value       = [for a in arvan_abrak.built_by_terraform_control : a.id]
  description = "IDs of all created VMs"
}

# output "vm_ids_compute" {
#   value       = [for a in arvan_abrak.arvan_abrak.built_by_terraform_compute : a.id]
#   description = "IDs of all created VMs"
# }

# output "floating_ip_addresses" {
#   value       = [for f in arvan_floating_ip.vm_fip : f.address]
#   description = "Allocated floating IP addresses (index matches VM index)"
# }
