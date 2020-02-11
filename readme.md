# mail-sender

Nightingale的理念，是将告警事件扔到redis里就不管了，接下来由各种sender来读取redis里的事件并发送，毕竟发送报警的方式太多了，适配起来比较费劲，希望社区同仁能够共建。

最常见的告警发送方式是邮件，所以这里我写了一个mail-sender，供参考

## compile

```bash
cd $GOPATH/src
mkdir -p github.com/n9e
cd github.com/n9e
git clone https://github.com/n9e/mail-sender.git
cd mail-sender
go build
```

如上编译完就可以拿到二进制了。

## configuration

读取告警事件，自然要给出redis的连接地址；发送邮件，自然要给出smtp配置；直接修改etc/mail-sender.yml即可

## pack

编译完成之后可以打个包扔到线上去跑，将二进制和配置文件打包即可：

```bash
tar zcvf mail-sender.tar.gz mail-sender etc/mail.html etc/mail-sender.yml
```

## test

配置etc/mail-sender.yml，相关配置修改好，我们先来测试一下smtp是否好使， `./mail-sender -t you@example.com`，程序会自动读取etc目录下的配置文件，发一封测试邮件给`you@example.com`

## run

如果测试邮件发送没问题，扔到线上跑吧，使用systemd或者supervisor之类的托管起来，systemd的配置实例：


```
$ cat mail-sender.service
[Unit]
Description=Nightingale mail sender
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root

Type=simple
ExecStart=/home/n9e/mail-sender
WorkingDirectory=/home/n9e

Restart=always
RestartSec=1
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
```