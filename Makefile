# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REPOSITORY := quay.io/prometheus
NAME       := golang-builder
VARIANTS   ?= base main arm powerpc mips s390x
LATEST     ?= 1.12
VERSIONS   ?= 1.11 1.12

all: build

build:
	@./build.sh "$(VERSIONS)" "$(VARIANTS)"

tag:
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-main" "$(REPOSITORY)/$(NAME):latest"
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-main" "$(REPOSITORY)/$(NAME):main"
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-arm" "$(REPOSITORY)/$(NAME):arm"
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-powerpc" "$(REPOSITORY)/$(NAME):powerpc"
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-mips" "$(REPOSITORY)/$(NAME):mips"
	docker tag "$(REPOSITORY)/$(NAME):$(LATEST)-s390x" "$(REPOSITORY)/$(NAME):s390x"

push:
	docker push "$(REPOSITORY)/$(NAME)"

.PHONY: all build tag push
