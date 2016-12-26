#!/bin/bash

#BASEDIR=`dirname $0`
#PROJECT_PATH=`cd $BASEDIR; pwd`

start(){
    docker run -d \
        --hostname postgres \
        --name postgres \
        -p 5432:5432 \
        postgres
    docker ps
}

stop(){
#    docker stop $(docker ps -a -q)
    docker stop postgres
}

clear(){
#    docker rm $(docker ps -a -q)
    docker rm postgres
}

case "$1" in
        start)
            start
            ;;

        stop)
            stop
            clear
            ;;

        restart)
            stop
            clear
            start
            ;;

        clear)
            clear
            ;;
        *)
            echo $"Usage: $0 {start|stop|restart|clear}"
            exit 1
esac

