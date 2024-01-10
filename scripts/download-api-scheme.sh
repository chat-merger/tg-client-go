#!/bin/bash

tmp_dir="/tmp/api-scheme-proto-$RANDOM"
scheme_file="mergerapi.proto"


while [ -n "$1" ]; do
  case "$1" in
  -f | --file) file="$2" ;;
  -d | --destination) destination="$2" ;;
  -t | --tag) tag="$2" ;;
  esac
  shift
done

if [[ "$file" == "" ]]; then
  printf "argument -f not found! specify the name of the schema file\n"
  exit 1
fi

if [[ "$destination" == "" ]]; then
  destination="./"
fi

if [[ "$tag" != "" ]]; then
  tag="--depth 1 --branch $tag"
fi

remove_tmp_directory() {
  rm -rf $tmp_dir
}

# exit when some command fails
set -e

# prepare
remove_tmp_directory


git clone $tag git@github.com:chat-merger/api-scheme-proto.git "$tmp_dir"
[ -d "$destination" ] || mkdir "$destination"
mv "$tmp_dir/$scheme_file" "$destination/$file"
remove_tmp_directory
echo Scheme downloaded to "$destination/$file"

exit 0
