# Logrotate configuration for MDM Platform
# Place this file in /etc/logrotate.d/mdm (Linux) or configure via Docker

# Nginx logs
/data/nginx/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 root root
    sharedscripts
    postrotate
        docker kill -s USR1 mdm-nginx-proxy 2>/dev/null || true
    endscript
}

# Application logs (if mounted to host)
/data/app/logs/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 root root
    mailfirst
}

# EMQX logs
/data/emqx/log/*.log {
    weekly
    missingok
    rotate 12
    compress
    delaycompress
    notifempty
    create 0640 root root
}
