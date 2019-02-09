# !/bin/bash

PROJECT_ROOT="$(basename `dirname $PWD`)/$(basename $PWD)"

for file in `grep -lR "PROJECT_ROOT" .`; do
  echo "Writing ${file}"
  sed -i '' 's:PROJECT_ROOT:'${PROJECT_ROOT}':g' ${file}
done
