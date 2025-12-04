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

FROM atomist/sdm-base:0.4.1@sha256:555a4f6fb9b0fb9d180e2e9a0bfcd04c44128cb76eca26516ccb7ba7dd304b5c

COPY package.json package-lock.json ./
RUN npm ci \
    && npm cache clean --force

COPY . ./

USER atomist:atomist
