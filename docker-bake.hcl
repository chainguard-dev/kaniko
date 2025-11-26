variable "image_semvar" {
  default = "latest"
}
variable "registry_host"{
  default = "docker.io/"
}
variable "user" {
  default = "none"
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

target "production" {
  inherits = ["_common"]
  tags = ["${registry_host}${user}/${item.container}:${image_semvar}"]
  name = item.container
  target = item.container
  contexts = item.context
  matrix = {
    item = [
      {
        container = "${slim_name}"
        context = {
          context =  "target:kaniko-base-slim"
        }
      },
      {
        container = "${debug_name}"
        context = {
          context =  "target:${executor_name}"
        }
      },
      {
        container = "${executor_name}"
        context = {
          context =  "target:kaniko-base"
        }
      },
      {
        container = "${warmer_name}"
        context = {
          context =  "target:kaniko-base"
        }
      }
    ]
  }
}
target "kaniko-base" {
  inherits = ["_common"]
  target = "kaniko-base"
  contexts = {
    kaniko-base-slim = "target:kaniko-base-slim"
  }
}
target "kaniko-base-slim" {
  inherits = ["_common"]
  target = "kaniko-base-slim"
  context = "."
}
target "_common" {
  dockerfile = "./deploy/Dockerfile"
}