#!/bin/bash

USERNAME="pg-health-check"
mkdir -p "/var/lib/$USERNAME"
useradd --system "$USERNAME" --shell /sbin/nologin --home-dir "/var/lib/$USERNAME"
chown -R "$USERNAME:$USERNAME" "/var/lib/$USERNAME"