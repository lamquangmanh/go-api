# CI/CD Workflow

This document describes the current CI/CD flow of go-api.

## Overview

The project uses a single Dev workflow file:

- `.github/workflows/on_dev.yml`

It runs in 3 stages:

1. `ci_process`
2. `cd_process`
3. `k8s_process`

## Trigger Rules

The workflow is triggered by:

- `workflow_dispatch`
- `push` on branch `dev`
- `pull_request` on branch `dev`

## Stage Details

### 1. CI stage (`ci_process`)

Calls reusable workflow:

- `lamquangmanh/gitops-config/.github/workflows/go-ci.yml@main`

Inputs:

- `environment: dev`
- `go_version: 1.25.x`
- `working_directory: .`
- `build_command: go build -o bin/grpc ./cmd/grpc`
- `test_command: go test ./...`

### 2. CD stage (`cd_process`)

Runs after CI success (`needs: ci_process`).

Condition:

- Runs only when event is `push` or `workflow_dispatch`

Calls reusable workflow:

- `lamquangmanh/gitops-config/.github/workflows/go-cd.yml@main`

Inputs:

- `environment: dev`
- `repository_path: quangmanhlam/user-backend-dev`
- `dockerfile: Dockerfile`
- `context: .`

### 3. K8s stage (`k8s_process`)

Runs after CD success (`needs: cd_process`).

Calls reusable workflow:

- `lamquangmanh/gitops-config/.github/workflows/k8s.yml@main`

Inputs:

- `environment: dev`
- `gitops_path: microk8s/microservice/apps/user-backend/overlays/dev`
- `repository_path: quangmanhlam/user-backend-dev`

## Required Secrets

The following repository secrets are required:

- `DOCKER_USERNAME`
- `DOCKER_TOKEN`
- `GIT_PAT`

## Notes

- `cd_process` is intentionally skipped on pull requests.
- GitOps manifest update is handled in `k8s_process`.
- The current flow does not require repository variables.
