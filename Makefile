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
VARIANTS   ?= base main
VERSION    ?= 1.18

all: build

build:
	@./build.sh "$(VERSION)" "$(VARIANTS)"

tag:
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-base" "$(REPOSITORY)/$(NAME):base"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):latest"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):main"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):arm"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):powerpc"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):mips"
	docker tag "$(REPOSITORY)/$(NAME):$(VERSION)-main" "$(REPOSITORY)/$(NAME):s390x"

push:
	docker push -a "$(REPOSITORY)/$(NAME)"

.PHONY: all build tag push
