#!/bin/bash

set -e

# 1. 源文件名
SRC_FILE=${1}
# 2. 包名
PACKAGE=${2}
# 3. 类型
TYPE=${3}
# 4. 文件后缀名
DES=${4}
# 第一个字符大写(注意括号)
PREFIX="$(tr '[:lower:]' '[:upper:]' <<< ${TYPE:0:1})${TYPE:1}"

# 文件名
DES_FILE=$(echo ${TYPE}| tr '[:upper:]' '[:lower:]')_${DES}_gen.go

sed 's/PACKAGE_NAME/'"${PACKAGE}"'/g' ${SRC_FILE} | \
    sed 's/GENERIC_TYPE/'"${TYPE}"'/g' | \
    sed 's/GENERIC_NAME/'"${PREFIX}"'/g' > ${DES_FILE}