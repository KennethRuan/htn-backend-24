#!/bin/sh

# Wait for the database service to be ready
while ! nc -z db 5432; do
  sleep 0.1
done

# Now that db is ready, run the ETL script
python3 etl_user_skills.py
exec ./main