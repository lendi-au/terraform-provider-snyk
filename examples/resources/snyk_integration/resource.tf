resource "snyk_integration" "example_integration" {
  organization = snyk_organization.example.id
  type         = "bitbucket-cloud"
  credentials {
    username = "username"
    password = "password" # Make sure your backend is encrypted - this is stored in plaintext!
  }
}