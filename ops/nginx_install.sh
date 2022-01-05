#!/bin/sh
#coding: utf-8
#author zhaodong
#site www.tshare365.com
#name: nginx_install.sh
yum install gcc openssl-devel pcre-devel zlib-devel wget -y 
id -u nginx 
if [ echo $? -ne 0 ];
then
    groupadd -r nginx
    useradd -r -g nginx -s /bin/nologin -M nginx
fi
wget http://nginx.org/download/nginx-1.7.9.tar.gz
tar xf nginx-1.7.9.tar.gz && cd nginx-1.7.9
./configure \
  --prefix=/opt/zly/nginx-1.7.9 \
  --sbin-path=/usr/sbin/nginx \
  --conf-path=/etc/nginx/nginx.conf \
  --error-log-path=/opt/logs/nginx/error.log \
  --http-log-path=/opt/logs/nginx/access.log \
  --pid-path=/var/run/nginx/nginx.pid  \
  --lock-path=/var/lock/nginx.lock \
  --user=nginx \
  --group=nginx \
  --with-http_ssl_module \
  --with-http_flv_module \
  --with-http_stub_status_module \
  --with-http_gzip_static_module \
  --http-client-body-temp-path=/var/tmp/nginx/client/ \
  --http-proxy-temp-path=/var/tmp/nginx/proxy/ \
  --http-fastcgi-temp-path=/var/tmp/nginx/fcgi/ \
  --http-uwsgi-temp-path=/var/tmp/nginx/uwsgi \
  --http-scgi-temp-path=/var/tmp/nginx/scgi \
  --with-pcre
make && make install
if [ ! -d  "/var/tmp/nginx/client/" ];
then
        mkdir -p /var/tmp/nginx/client/p
fi
echo "start nginx....."
/usr/sbin/nginx -c /etc/nginx.conf
if [ echo $? -ne 0 ];
then
    echo "#################start nginx failed!###############"
fi
