#!/bin/sh
curl http://localhost:8000/alvaro -X PUT -w %{http_code}
