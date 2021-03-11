/*******************************************************************************
*
* Copyright 2021 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

#include <syslog.h>

//NOTE: This program can be used to generate test log lines to feed into syslog-stdout.
//      Build with `make build/syslog-generator`.

int main() {
  openlog("syslog-generator", 0, LOG_LOCAL1);
  for (int i = 0; i < 10000; ++i) {
    syslog(LOG_INFO, "test line number %05d", i);
  }
  closelog();
}
