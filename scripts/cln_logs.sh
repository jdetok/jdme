#!/bin/bash

DIR=z_log
MAXSAVE=2

cd $DIR
ls -t | tail -n  +$(($MAXSAVE + 1)) | xargs rm --