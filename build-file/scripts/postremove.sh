#!/bin/bash
USERNAME="pg-health-check"
userdel "$USERNAME" || true
rm -rf "/var/lib/$USERNAME"
rm -rfv /etc/pg-health-check