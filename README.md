# coap-over-quic-exp

## What? 
This repository is configured to performance analysis for CompressedCoAP based on [minitopo](https://github.com/qdeconinck/minitopo)

## Configuration
this configuration is performed in mininet and multipass environment. 

# Testbed Configuration

## multipass based mininet 

### Create Instance 

```zsh
multipass launch -c 4 -d 30G -m 8G -n mininet
```

## Install and configure mininet

### install 


```bash
$ sudo apt install mininet

# 설치 확인
$ mn --version
```


### install mininet essential utility


```bash
$ git clone git://github.com/mininet/mininet
$ mininet/util/install.sh -fw # sh install.sh 이런 식으로 실행하지 말것!
```

## python 

* install pip

  ```bash
  $ sudo apt install python3-pip
  ```

* install mininet python module

  ```bash
  $ sudo su
  $ pip install mininet
  ```

  * You should install this module with root 

## Enable IP Forwarding

* <https://fedingo.com/how-to-enable-ip-forwarding-in-ubuntu/>


```bash
$ sysctl -w net.ipv4.ip_forward=1
```


* or modify the file : `/etc/sysctl.conf` 



## Pre-requirements for usage of matplotlib

```
pip install matplotlib
pip install numpy
pip install pyqt5
```

