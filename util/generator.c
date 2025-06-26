/*
 * Copyright 2021 SAP SE
 * SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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
