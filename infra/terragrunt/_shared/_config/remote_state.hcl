// 🌍 Terragrunt Configuration Local Variables
// This section specifies local variables utilized across various child configurations.
// These locals ensure uniformity in naming conventions and metadata throughout the infrastructure.
locals {
  // 🗄️ S3 Bucket Name for Terraform State
  bucket_name_unnormalized = get_env("TG_STACK_REMOTE_STATE_BUCKET_NAME")
  bucket_name              = lower(trimspace(local.bucket_name_unnormalized))

  // 🔒 DynamoDB Lock Table Name
  lock_table_unnormalized = get_env("TG_STACK_REMOTE_STATE_LOCK_TABLE")
  lock_table              = lower(trimspace(local.lock_table_unnormalized))

  // 🌐 AWS Region for Deployment
  region_unnormalized = get_env("TG_STACK_REMOTE_STATE_REGION", "us-east-1")
  region              = lower(trimspace(local.region_unnormalized))

  // 📦 Basename for the Terraform state object
  state_object_basename_unnormalized = get_env("TG_STACK_REMOTE_STATE_OBJECT_BASENAME", "terraform.tfstate.json")
  state_object_basename              = lower(trimspace(local.state_object_basename_unnormalized))

  // 📁 Backend Terraform File
  backend_tf_filename_unnormalized = get_env("TG_STACK_REMOTE_STATE_BACKEND_TF_FILENAME", "backend.tf")
  backend_tf_filename              = lower(trimspace(local.backend_tf_filename_unnormalized))
}
