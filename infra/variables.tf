variable "profile" {
  type        = string
  description = "AWS profile associated with credentials"
}

variable "region" {
  type        = string
  description = "AWS Region to create your resources in"
  default     = "us-east-1"
}

variable "default_tags" {
  type        = map(string)
  description = "Map of default tag for our project"
  default = {
    "project" = "Valkyrie"
  }
}

variable "cidr_block" {
  type        = string
  description = "Default CIDR block for VPC"
  default     = "10.255.0.0/20"
}

variable "public_subnet_count" {
  type        = number
  description = "Count of public subnets to be created"
  default     = 2
}

variable "frontend_server_count" {
  type        = number
  description = "Count of EC2 instances for Frontend server"
  default     = 7
}

variable "frontend_server_key_pair" {
  type        = string
  description = "Key pair for FE server EC2 instance"
}

variable "dispatcher_key_pair" {
  type        = string
  description = "Key pair for Dispatcher EC2 instance"
}
