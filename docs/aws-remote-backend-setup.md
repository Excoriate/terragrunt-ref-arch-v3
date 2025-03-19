# AWS Remote Backend Setup

## Overview

This guide explains how to set up a secure AWS-based remote backend for Terraform and Terragrunt state management. A properly configured remote backend provides:

- **State Locking**: Prevents concurrent state modifications that could corrupt your infrastructure state
- **State Versioning**: Maintains a history of your state files for auditing and recovery
- **Encryption**: Ensures sensitive information in your state files is protected
- **Access Control**: Centralizes and secures access to infrastructure state

## Components

The AWS remote backend consists of two primary components:

1. **S3 Bucket**: Stores the Terraform state files
2. **DynamoDB Table**: Provides state locking to prevent concurrent modifications

## Prerequisites

- AWS CLI installed and configured with appropriate credentials
- Permissions to create and configure S3 buckets and DynamoDB tables

## Setup Instructions

### Option 1: One-Step Setup (Recommended)

The following command will create and configure all required resources with best security practices:

```bash
aws s3api create-bucket \
    --bucket terraform-state-makemyinfra \
    --region us-east-1 && \
aws s3api put-bucket-versioning \
    --bucket terraform-state-makemyinfra \
    --versioning-configuration Status=Enabled && \
aws s3api put-bucket-encryption \
    --bucket terraform-state-makemyinfra \
    --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}' && \
aws s3api put-public-access-block \
    --bucket terraform-state-makemyinfra \
    --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true" && \
aws dynamodb create-table \
    --table-name terraform-state-makemyinfra \
    --region us-east-1 \
    --billing-mode PAY_PER_REQUEST \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH
```

You may need to modify the bucket name and region to fit your requirements.

### Option 2: Step-by-Step Setup

If you prefer to understand and execute each step individually:

#### 1. Create the S3 Bucket

```bash
aws s3api create-bucket \
    --bucket terraform-state-makemyinfra \
    --region us-east-1
```

For regions other than `us-east-1`, use:

```bash
aws s3api create-bucket \
    --bucket terraform-state-makemyinfra \
    --region your-region \
    --create-bucket-configuration LocationConstraint=your-region
```

#### 2. Enable Bucket Versioning

```bash
aws s3api put-bucket-versioning \
    --bucket terraform-state-makemyinfra \
    --versioning-configuration Status=Enabled
```

#### 3. Enable Server-Side Encryption

```bash
aws s3api put-bucket-encryption \
    --bucket terraform-state-makemyinfra \
    --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}'
```

#### 4. Block Public Access

```bash
aws s3api put-public-access-block \
    --bucket terraform-state-makemyinfra \
    --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"
```

#### 5. Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
    --table-name terraform-state-makemyinfra \
    --region us-east-1 \
    --billing-mode PAY_PER_REQUEST \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH
```

## Configuration in Terragrunt

Once the backend is created, update your environment configuration to reference it. The reference architecture already includes functionality to use the backend specified in your environment variables.

In your `.env` or `.envrc` file:

```bash
# Remote State Configuration
TG_STACK_REMOTE_STATE_BUCKET_NAME="terraform-state-makemyinfra"
TG_STACK_REMOTE_STATE_LOCK_TABLE="terraform-state-makemyinfra"
TG_STACK_REMOTE_STATE_REGION="us-east-1"
```

## Security Best Practices

- **IAM Policies**: Restrict access to the S3 bucket and DynamoDB table to only authorized users/roles
- **Access Logging**: Enable S3 access logging to monitor bucket access
- **Lifecycle Policies**: Consider adding lifecycle policies to manage old state versions
- **HTTPS Only**: Enforce HTTPS-only communication with the S3 bucket

### Sample IAM Policy for Terraform/Terragrunt Access

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::terraform-state-makemyinfra",
        "arn:aws:s3:::terraform-state-makemyinfra/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": "arn:aws:dynamodb:*:*:table/terraform-state-makemyinfra"
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **Access Denied**
   - Verify AWS credentials have appropriate permissions
   - Check IAM policies attached to your user/role

2. **Bucket Already Exists**
   - S3 bucket names must be globally unique
   - Choose a different bucket name or use an existing bucket

3. **State Locking Failures**
   - Ensure DynamoDB table exists and is correctly named
   - Verify permissions include DynamoDB actions
   - Check for stale locks with `terragrunt force-unlock`

4. **Region Consistency**
   - Ensure S3 bucket and DynamoDB table are in the same AWS region
   - Verify `TG_STACK_REMOTE_STATE_REGION` matches actual resource region

## Cleanup

If you need to remove the backend infrastructure:

```bash
# First, remove all files from the bucket
aws s3 rm s3://terraform-state-makemyinfra --recursive

# Delete the bucket
aws s3api delete-bucket --bucket terraform-state-makemyinfra

# Delete the DynamoDB table
aws dynamodb delete-table --table-name terraform-state-makemyinfra
```

⚠️ **WARNING**: Deleting the remote backend will remove all Terraform state files, which could make managing existing infrastructure extremely difficult. Only do this if you're sure you want to completely reset your infrastructure management.
