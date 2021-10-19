#!/bin/bash
WORKING_DIR_PATH=`dirname "$0"`

build_file_name=$1
arg=$2

if [[ $2 != "" ]]; then
    ${WORKING_DIR_PATH}/main/${build_file_name} $2
else
    ${WORKING_DIR_PATH}/main/${build_file_name}
fi

date=`date --rfc-3339=seconds`
echo ${date} run ${build_file_name} ${arg} >> ${WORKING_DIR_PATH}/main/${build_file_name}.log 



