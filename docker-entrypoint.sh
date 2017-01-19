#!/bin/bash
if [ "$1" = "sleep" ];
    then 
    sleep $2;
else  
    /blog/journey $@;
fi