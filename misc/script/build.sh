#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

wegigo_ROOT=$(dirname "${BASH_SOURCE}")/..


KUBE_VERBOSE="${KUBE_VERBOSE:-1}"
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::build_binaries "$@"
kube::golang::place_bins

go get github.com/kabukky/httpscerts
go get github.com/julienschmidt/httprouter
go get github.com/rfancn/goy2h
go build -o wegigo