locals {
  public_cidr_blocks = [for i in range(var.public_subnet_count) : cidrsubnet(var.cidr_block, 4, i)]
}
