variable "username" {
    type = "string"
}

variable "password" {
    type = "string"
}

variable "base_url" {
    type = "string"
}

provider "cloudcenter" {
    username = "${var.username}"
    password = "${var.password}"
    base_url = "${var.base_url}"
}

resource "cloudcenter_user" "user" {

    first_name      = "Terraform"
    last_name       = "Plugin"
    password        = "myPassword"
    email_address   = "terraform@mydomain.com"
    company_name    = "Company"
    phone_number    = "12345"
    tenant_id       = "1"

}
