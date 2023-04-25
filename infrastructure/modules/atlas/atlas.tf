provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}

data "mongodbatlas_project" "main" {
  name = var.atlas_project_name
}

resource "mongodbatlas_cluster" "cluster" {
  project_id = data.mongodbatlas_project.main.id
  name       = var.atlas_cluster_name
  num_shards = 1

  replication_factor           = 3
  backup_enabled               = false
  auto_scaling_disk_gb_enabled = false
  mongo_db_major_version       = "6.0"

  provider_name               = "TENANT"
  disk_size_gb                = 2
  provider_instance_size_name = "M2"
  provider_region_name        = "US_EAST_1"
  backing_provider_name       = "AWS"
}

resource "mongodbatlas_database_user" "app-user" {
  username           = var.app_username
  password           = var.app_user_password
  project_id         = data.mongodbatlas_project.main.id
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = var.app_user_database
  }
}

output "DB_URI" {
  value = replace(
    mongodbatlas_cluster.cluster.mongo_uri_with_options,
    "mongodb://",
    "mongodb://${mongodbatlas_database_user.app-user.username}:${mongodbatlas_database_user.app-user.password}@",
  )
}

output "DB_NAME" {
  value = "Please add the folowing database in the cluster: ${var.app_user_database}"
}
