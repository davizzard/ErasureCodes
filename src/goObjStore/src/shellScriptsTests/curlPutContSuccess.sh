#!/bin/sh
curl http://localhost:8000/alvaro/container1 -X PUT -w %{http_code}