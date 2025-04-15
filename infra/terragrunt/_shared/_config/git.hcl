locals {
  # ---------------------------------------------------------------------------------------------------------------------
  # 🌐 GIT BASE URLS
  # ---------------------------------------------------------------------------------------------------------------------
  # Purpose: Provide a centralised configuration for Git base URLs
  #
  # Key Responsibilities:
  # - Provide a centralised configuration for Git base URLs
  # - Enable consistent Git base URL usage across all modules
  # - Support cross-module Git base URL configuration sharing

  # This file is sourced by the root config.hcl file, and from there, all the units
  # will inherit the base URLs.
  # ---------------------------------------------------------------------------------------------------------------------
  git_base_urls = {
    github = "git::git@github.com:"
    gitlab = "git::gitlab.com:"
    local  = "${get_repo_root()}/infra/terraform/modules"
  }
}
