#!/bin/sh
curl http://localhost:8000/alvaro/container1False/obj1 -X GET -w %{http_code}