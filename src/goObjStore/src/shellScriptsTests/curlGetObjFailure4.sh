#!/bin/sh
curl http://localhost:8000/alvaroFalse/container1/objFalse -X GET -w %{http_code}