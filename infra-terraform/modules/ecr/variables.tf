variable "repository_names" {
  description = "ECR repositories to create"
  type        = list(string)
}

variable "tags" {
  description = "Tags to apply"
  type        = map(string)
  default     = {}
}
