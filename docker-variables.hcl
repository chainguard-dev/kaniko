variable "image_semvar" {
  default = "latest"
}
variable "registry_host"{
  default = "docker.io/"
}
variable "user" {
  default = "library"
}
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