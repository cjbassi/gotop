#!/bin/bash

VERSION=v1.0

curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop > /usr/bin/gotop
chmod +x /usr/bin/gotop
