#!/bin/sh
curl http://localhost:8000/alvaroFalse/container1/obj1 -X PUT -T "../bigFile" -w %{http_code}