#!/bin/bash

if [[ $EUID -ne 0 ]]; then
	echo "This script must be run as root" 
	exit 1
fi

LOGFILE=/var/log/alanick.log

cp alanick.conf.example /etc/alanick.conf
touch $LOGFILE
chown alanick:alanick $LOGFILE
chmod g+w $LOGFILE
