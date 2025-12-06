worker_processes 1;

events {
    worker_connections 1024;
}

http {

    resolver 172.25.0.1 valid=10s;
    resolver_timeout 5s;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    gzip on;
    gzip_disable "msie6";

    server {
        listen 8198;
        server_name monitor;

        location / {
            # Using the variable makes nginx verify the host during runtime not at startup
            set $upstream_monitor http://monitor:8080;
            proxy_pass $upstream_monitor;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }

    server {
        listen 8199;
        server_name service1;

        location / {
            set $upstream_service1 http://service1:8080;
            proxy_pass $upstream_service1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }

    server {
        listen 8081;
        server_name localhost;

        location /nginx-health {
            return 200 "nginx is healthy";
        }
    }
}
