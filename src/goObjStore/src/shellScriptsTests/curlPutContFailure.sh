#!/bin/sh
curl http://localhost:8000/alvaroFalse/container1 -X PUT -w %{http_code}