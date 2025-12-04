# Copyright 2020 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM bash:4.4@sha256:8fb4d13129d7d61cf4f670eed1a23db397a16352671ee758f4faed7b6e3c3f57

ARG NUM
COPY context.txt .
COPY make.sh .
SHELL ["/usr/local/bin/bash", "-c"]
RUN ./make.sh $NUM
RUN ls -al /workdir | wc
