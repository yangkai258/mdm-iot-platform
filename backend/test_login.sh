#!/bin/sh
curl -s -X POST -H "Content-Type: application/json" --data-binary @/tmp/p http://localhost:8080/api/v1/auth/login