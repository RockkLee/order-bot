resource "aws_ecr_repository" "this" {
  for_each = toset(var.repository_names)

  name                 = each.value
  image_tag_mutability = "MUTABLE"

  # An image will be scanned the possible vulnerabilities when it's pushed to ECR.
  image_scanning_configuration {
    scan_on_push = false
  }

  # Push flow:
  # 1) Client authenticates to ECR
  # 2) Client pushes layers/manifest over HTTPS
  # 3) ECR stores the image data and encrypts it
  #
  # Pull flow:
  # 1) Client authenticates to ECR
  # 2) ECR checks pull permission
  # 3) ECR reads encrypted image data from storage
  # 4) ECR decrypts it server-side automatically
  # 5) ECR sends the image over HTTPS
  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(var.tags, {
    Name = each.value
  })
}

# resource "aws_ecr_lifecycle_policy" "this" {
#   for_each = aws_ecr_repository.this
#
#   repository = each.value.name
#   policy     = jsonencode({
#     rules = [{
#       rulePriority = 1
#       description  = "Keep only the latest 20 images"
#       selection = {
#         tagStatus   = "any"
#         countType   = "imageCountMoreThan"
#         countNumber = 20
#       }
#       action = { type = "expire" }
#     }]
#   })
# }
