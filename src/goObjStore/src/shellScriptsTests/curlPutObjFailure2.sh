#!/bin/sh
curl http://localhost:8000/alvaro/containerFalse/obj1 -X PUT -T "../bigFile" -w %{http_code}