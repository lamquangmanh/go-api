param(
  [string]$DbHost = "localhost",
  [int]$DbPort = 5432,
  [string]$DbUser = "postgres",
  [string]$DbPassword = "postgres",
  [string]$DbName = "go_api"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

function Require-Cmd {
  param([string]$Name)
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Missing required command: $Name" -ForegroundColor Red
    exit 1
  }
}

Write-Host "==> Checking required tools"
Require-Cmd go
Require-Cmd protoc
Require-Cmd make
Require-Cmd psql
Require-Cmd sqlc

Write-Host "==> Running code generation and dependency sync"
make proto-tools
make proto
make sqlc
make tidy

Write-Host "==> Applying migration: db/migrations/001_create_users.sql"
$env:PGPASSWORD = $DbPassword
& psql -h $DbHost -p $DbPort -U $DbUser -d $DbName -f "db/migrations/001_create_users.sql"
Remove-Item Env:PGPASSWORD -ErrorAction SilentlyContinue

Write-Host "==> Starting service"
make run
