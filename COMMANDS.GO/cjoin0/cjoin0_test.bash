#!/bin/bash
#!/usr/local/bin/bash -xv # コマンド処理系の変更例
#
# test script of cjoin0
#
# usage: [<test-path>/]cjoin0.test [<command-path> [<python-version>]]
#
#            <test-path>は
#                    「現ディレクトリーからみた」本スクリプトの相対パス
#                    または本スクリプトの完全パス
#                    省略時は現ディレクトリーを仮定する
#            <command-path>は
#                    「本スクリプトのディレクトリーからみた」test対象コマンドの相対パス
#                    またはtest対象コマンドの完全パス
#                    省略時は本スクリプトと同じディレクトリーを仮定する
#                    値があるときまたは空値（""）で省略を示したときはあとにつづく<python-version>を指定できる
#            <python-version>は
#                    使用するpython処理系のversion（minor versionまで指定可）を指定する
#                    （例 python2 python2.6 phthon3 python3.4など）
#                    単にpythonとしたときは現実行環境下でのdefault versionのpythonを使用する
#                    文字列"python"は大文字/小文字の区別をしない
#                    省略時はpythonを仮定する
name=cjoin0 # test対象コマンドの名前
testpath=$(dirname $0) # 本スクリプト実行コマンドの先頭部($0)から本スクリプトのディレトリー名をとりだす
cd $testpath # 本スクリプトのあるディレクトリーへ移動
# if test "$2" = ""; # <python-version>($2)がなければ
# 	then pythonversion="python" # default versionのpythonとする
# 	else pythonversion="$2" # <python-version>($2)があれば指定versionのpythonとする
# fi
if test "$1" = ""; # <command-path>($1)がなければ
	then commandpath="." # test対象コマンドは現ディレクトリーにある
	else commandpath="$1" # <command-path>($1)があればtest対象コマンドは指定のディレクトリーにある
fi
com="go run ${commandpath}/${name}.go" # python処理系によるtest対象コマンド実行の先頭部
tmp=/tmp/$$

ERROR_CHECK(){
	[ "$(echo ${PIPESTATUS[@]} | tr -d ' 0')" = "" ] && return
	echo $1
	echo "${pythonversion} ${name}" NG
	# rm -f $tmp-*
	exit 1
}

###########################################
#TEST1

cat << FIN > $tmp-tran
0000000 浜地______ 50 F 91 59 20 76 54
0000001 鈴田______ 50 F 46 39 8  5  21
0000003 杉山______ 26 F 30 50 71 36 30
0000004 白土______ 40 M 58 71 20 10 6
0000005 崎村______ 50 F 82 79 16 21 80
FIN

cat << FIN > $tmp-master
0000004
0000001
FIN

cat << FIN > $tmp-out
0000001 鈴田______ 50 F 46 39 8 5 21
0000004 白土______ 40 M 58 71 20 10 6
FIN

${com} key=1 $tmp-master $tmp-tran > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST1 error"

###########################################
#TEST2

cat << FIN > $tmp-tran
0000000 浜地______ 50 F 91 59 20 76 54
0000001 鈴田______ 50 F 46 39 8  5  21
0000003 杉山______ 26 F 30 50 71 36 30
0000004 白土______ 40 M 58 71 20 10 6
0000005 崎村______ 50 F 82 79 16 21 80
FIN

cat << FIN > $tmp-master
0000001
0000004
FIN

cat << FIN > $tmp-out
0000001 鈴田______ 50 F 46 39 8 5 21
0000004 白土______ 40 M 58 71 20 10 6
FIN

cat << FIN > $tmp-ng
0000000 浜地______ 50 F 91 59 20 76 54
0000003 杉山______ 26 F 30 50 71 36 30
0000005 崎村______ 50 F 82 79 16 21 80
FIN

${com} +ng key=1 $tmp-master $tmp-tran > $tmp-ans 2> $tmp-ans2
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST2-1 error"
diff $tmp-ans2 $tmp-ng
[ $? -eq 0 ] ; ERROR_CHECK "TEST2-2 error"

###########################################
#TEST3

cat << FIN > $tmp-tran
AAA 001 山田
BBB 002 上田
CCC 003 太田
DDD 004 堅田
FIN

cat << FIN > $tmp-master
003 太田
002 上田
FIN

cat << FIN > $tmp-out
BBB 002 上田
CCC 003 太田
FIN

${com} key=2/3 $tmp-master $tmp-tran > $tmp-ans
diff $tmp-ans $tmp-out
[ $? -eq 0 ] ; ERROR_CHECK "TEST3 error"

rm -f $tmp-*
echo "${pythonversion} ${name}" OK
exit 0
