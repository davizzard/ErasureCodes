#!/bin/sh
curl http://localhost:8000/NonExistingAccount -X GET -w %{http_code}