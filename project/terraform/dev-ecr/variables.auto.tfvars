tags = {
  "Environment" = "Dev"
}

ecr_name = [
  "experiment/listener-service",
  "experiment/logger-service",
  "experiment/mailer-service",
  "experiment/authentication-service",
  "experiment/broker-service",
  "experiment/front-end"
]

image_mutability = "IMMUTABLE"
