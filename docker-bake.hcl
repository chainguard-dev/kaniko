target "_common" {
  dockerfile = "deploy/Dockerfile"
}
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
  tags = ["${registry_host}${user}/${item.container}:${image_semvar}"]
  name = item.container
  target = item.container
  dockerfile = "./deploy/Dockerfile"
  contexts = item.context
  registry_host = "registry.gnoejuan.com"
  user = "library"
  matrix = {
    item = [
      {
        container = "kaniko-slim"
        context = {
          context =  "target:kaniko-base-slim"
        }
      },
      {
        container = "kaniko-debug"
        context = {
          context =  "target:kaniko-executor"
        }
      },
      {
        container = "kaniko-executor"
        context = {
          context =  "target:kaniko-base"
        }
      },
      {
        container = "kaniko-warmer"
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
