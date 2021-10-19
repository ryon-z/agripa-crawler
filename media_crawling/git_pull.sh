#!/bin/bash
working_dir_path=`dirname "$0"`

function build_go_files() {
    for file_path in ${working_dir_path}/main/*.go; do
        build_file_path=`echo ${file_path} | sed "s/\.go//g"` 
        go build -o ${build_file_path} ${file_path}  
        echo "${file_path} build 완료"
    done
}

function main() {
    git pull origin master
    build_go_files
}

main
