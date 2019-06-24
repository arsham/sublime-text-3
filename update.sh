#!/bin/sh
update() {
    cd ../$1
    git pull origin master
    cd ../User
}

for dir in LSP HaoGist SublPlugs GoRename; do
    update $dir
done

cd ../GoSublime
git pull origin development
