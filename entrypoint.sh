#!/bin/bash

if [[ -z ${LOG_ERRORS_TO_STDERR+default} ]]; then
  /usr/src/app/hapttic $@
else
  /usr/src/app/hapttic -logErrors $@
fi