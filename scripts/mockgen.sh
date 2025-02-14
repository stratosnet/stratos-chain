#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/sds/types/expected_keepers.go -package testutil -destination x/sds/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/pot/types/expected_keepers.go -package testutil -destination x/pot/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/register/types/expected_keepers.go -package testutil -destination x/register/testutil/expected_keepers_mocks.go
