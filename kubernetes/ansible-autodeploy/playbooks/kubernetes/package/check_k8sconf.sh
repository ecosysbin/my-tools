#!/bin/bash

if [ -f /etc/kubernetes/pki/ca.crt ]
then 
     echo "kubernetes conf exist"
else
     echo "kubernetes conf not exist"
fi
