# heimboard-wakeup-go
This is a Go implementation of the wakeup server for my [HEIMBOARD](https://github.com/r2binx/heimboard) project to improve compatibilty and efficiency.

## Installation
To install it on your Linux machine you can simply run:

```
bash <(curl -s  https://raw.githubusercontent.com/r2binx/heimboard-wakeup-go/latest/install.sh)
```

This will install the file in `/usr/local/bin/wakeup` & set the `cap_net_raw+ep` to ron it without root.
_(This is necessary to send the WoL package)_

You can also grab a matching release for your platform [here](https://github.com/r2binx/heimboard-wakeup-go/releases/tag/latest) and run it as root/admin.


## Usage
Dowload the [config file](https://github.com/r2binx/heimboard-wakeup-go/raw/latest/.config.example) and rename it to `.config`. Set it up according to your needs.

The server will look for it in the directory it's started from.
