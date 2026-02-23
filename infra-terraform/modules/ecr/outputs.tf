output "repository_urls" {
  description = "ECR repo URLs by repository name"
  value = {
    for name, repo in aws_ecr_repository.this : name => repo.repository_url
  }
}
