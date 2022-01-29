#!/bin/bash
curl -L -o /tmp/wakeup https://github.com/r2binx/heimboard-wakeup-go/releases/download/latest/wakeup.linux-$(uname -m)
sudo install /tmp/wakeup /usr/local/bin
sudo setcap cap_net_raw+ep /usr/local/bin/wakeup
