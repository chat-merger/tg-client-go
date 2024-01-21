#!/bin/bash

scheme_file="mergerapi.proto"
tmp_file="/tmp/$RANDOM-$scheme_file"


while [ -n "$1" ]; do
  case "$1" in
  -f | --file) file="$2" ;;
  -d | --destination) destination="$2" ;;
  -b | --branch) branch="$2" ;;
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

if [[ "$branch" == "" ]]; then
  branch="main"
fi

remove_tmp_file() {
  rm -f $tmp_file
}

# exit when some command fails
set -e

# prepare
remove_tmp_file
echo "https://raw.githubusercontent.com/chat-merger/api-scheme-proto/$branch/$scheme_file"
curl "https://raw.githubusercontent.com/chat-merger/api-scheme-proto/$branch/$scheme_file" > "$tmp_file"
[ -d "$destination" ] || mkdir "$destination"
mv "$tmp_file" "$destination/$file"
remove_tmp_file
echo Scheme downloaded to "$destination/$file"

exit 0
