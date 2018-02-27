#!/bin/sh
curl http://localhost:8000/NonExistingAccount/container1 -X GET -w %{http_code}