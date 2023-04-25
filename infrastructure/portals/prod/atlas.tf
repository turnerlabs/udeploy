# Public key obtained in the access managemnet control in the Atlas web console
variable "mongodbatlas_public_key" {}

# private key obtained in the access managemnet control in the Atlas web console
variable "mongodbatlas_private_key" {}

# Project or Context in which the Cluster will be created
# This resource will need to be created in the Atlas web console
variable "atlas_project_name" {}

# The name of the Cluster that will hold ou database. 
# This resource will be created by this terraform module
variable "atlas_cluster_name" {}

# User name that will be added in the project with Read/Wrtie access
# The access wil be given to the database named in the variable app_user_databases
variable "app_username" {}

# Password for app_username
variable "app_user_password" {}

# name of the database that the user will have access
variable "app_user_database" {}

module "env" {
  source = "github.com/turnerlabs/udeploy//infrastructure/modules/atlas?ref=v0.33.2-rc"

  mongodbatlas_public_key = var.mongodbatlas_public_key

  mongodbatlas_private_key = var.mongodbatlas_private_key

  atlas_project_name = var.atlas_project_name

  atlas_cluster_name = var.atlas_cluster_name

  app_username = var.app_username

  app_user_password = var.app_user_password

  app_user_database = var.app_user_database
}

output "DB_URI" {
  value = module.env.DB_URI
}

output "DB_NAME" {
  value = module.env.DB_NAME
}
