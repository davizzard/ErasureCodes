#!/bin/sh
curl http://localhost:8000/alvaro -X GET -w %{http_code}