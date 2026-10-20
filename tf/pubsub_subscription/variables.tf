variable "deadletter_message_retention_duration" {
  description = <<EOT
How long to retain unacknowledged messages in the subscription's backlog,
from the moment a message is published.
Cannot be more than 31 days ('"2678400s"') or less than 10 minutes ('"600s"').
EOT
  type        = string
}

variable "subscription_message_retention_duration" {
  description = <<EOT
How long to retain unacknowledged messages in the subscription's backlog,
from the moment a message is published.
Cannot be more than 31 days ('"2678400s"') or less than 10 minutes ('"600s"').
EOT
  type        = string
}

variable "subscription_name" {
  description = "The name of the Pub/Sub subscription to create for the given topic."
  type        = string
}

variable "subscribers" {
  description = "List of subscribers for the Pub/Sub subscription."
  type = list(object({
    member = string
  }))
}

variable "topic" {
  description = "The Pub/Sub topic to subscribe to."
  type = object({
    id   = string
    name = string
  })
}

variable "viewers" {
  description = "List of viewers for the Pub/Sub subscription."
  type = list(object({
    member = string
  }))
}
