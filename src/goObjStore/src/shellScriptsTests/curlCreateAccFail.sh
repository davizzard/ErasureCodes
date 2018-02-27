#!/bin/sh
curl http://localhost:8000/acc1 -X PUT -w %{http_code}