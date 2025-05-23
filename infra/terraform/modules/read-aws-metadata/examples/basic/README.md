# Basic Example: Read AWS Account Metadata Module

## Overview
> **Note:** This example demonstrates basic usage of the `read-aws-account` module found in the parent directory (`../`). It simply calls the module to retrieve AWS account metadata like Account ID, Region, Partition, etc.

### 🔑 Key Features Demonstrated
- Calling the `read-aws-account` module.
- Using fixtures (`fixtures/*.tfvars`) to test enabled/disabled states of the module call.

### 📋 Usage Guidelines
1.  **Configure:** Use the `fixtures/default.tfvars` file (usually empty) or `fixtures/disabled.tfvars` (`is_enabled = false`).
2.  **Initialize:** Run `terraform init`.
3.  **Plan:** Run `terraform plan -var-file=fixtures/default.tfvars`.
4.  **Apply:** Run `terraform apply -var-file=fixtures/default.tfvars`. The outputs will show the retrieved AWS metadata.
5.  **Makefile:** Alternatively, use the provided `Makefile` targets (`make plan-default`, `make apply-default`).

<!-- BEGIN_TF_DOCS -->
<!-- END_TF_DOCS -->
