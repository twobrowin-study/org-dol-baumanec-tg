#!/bin/bash
set -e

psql -U "$POSTGRES_USER" "$POSTGRES_DB" < /init.sql