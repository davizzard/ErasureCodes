#!/bin/sh
curl http://localhost:8000/alvaro/container1False/objFalse -X GET -w %{http_code}