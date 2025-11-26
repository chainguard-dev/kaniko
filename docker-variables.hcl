variable "image_semvar" {
  default = "latest"
}
variable "registry_host"{
  default = "docker.io/"
}
variable "user" {
  default = "library"
}
# If the following names are updated, their corresponding alias needs to be updated in deploy/Dockerfile
# example, changing "slim_name" from "kaniko-slim" to "test-value" means you need to update 
# FROM kaniko-base-slim AS kaniko-slim
# to 
# FROM kaniko-base-slim AS test-value
variable "slim_name" {
  default = "kaniko-slim"
}
variable "debug_name" {
  default = "kaniko-debug"
}
variable "executor_name" {
  default = "kaniko-executor"
}
variable "warmer_name" {
  default = "kaniko-warmer"
}