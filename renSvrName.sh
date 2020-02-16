#!/bin/sh
if [ x$1 != x ]; then
    #...有参数
    echo "we will replace demo with parameter $1"
else
    #...没有参数
    echo "you should input a parameter to set name"
    exit -1
fi
echo "find All file contain Demo:\n"
find ./ -type f -name "*" |xargs -I {} grep "Demo" {}
echo "find All file contain demo\n"
find ./ -type f -name "*" |xargs -I {} grep "demo" {}

#
find .| xargs grep -ri "demo"
find .| xargs grep -ri "demo" -l 


echo "Replace file with demo"
grep -rl 'demo' ./  | xargs sed -i "" "s/demo/$1/g"
echo "Replace file with Demo"
grep -rl 'Demo' ./  | xargs sed -i "" "s/Demo/Svr$1/g"

exit -2

#grep -rl 'aaaModule' ./  | xargs sed -i "" "s/aaaModule/bbbName/g"

filelist=("./api/api.proto ./api/client.go ./internal/server/grpc/server.go ./internal/server/http/server.go")
for x in $filelist; do
	#sed -ig 's/原字符串/新字符串/g' /home/1.txt
	#eval sed 's/$a/$b/' filename
	#sed "s/$a/$b/" filename
	#sed 's/'$a'/'$b'/' filename
	#sed s/$a/$b/ filename

	sedcmd="s/Demo/$1/g"
	echo exec: sed -i"" $sedcmd $x
	cp $x $x.backup
	sed $sedcmd $x.backup > $x
	
done
