#!/bin/sh
curl http://localhost:8000/NonExistingAccount/NonExistingContainer -X GET -w %{http_code}