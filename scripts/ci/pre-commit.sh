#! /bin/sh

pre-commit install --install-hooks
SKIP=terraform_providers_lock,terraform_validate pre-commit run --all
