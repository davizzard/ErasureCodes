#!/bin/sh
curl http://localhost:8000/alvaro/NonExistingContainer -X GET -w %{http_code}