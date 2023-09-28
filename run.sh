result=`go version`
if [[ $result == "" ]] ; then
    echo -e "\e[31m !!! panic : golang is not installed"
    exit
fi

if [ ! -f "./smpc_k" ];then
    echo -e "\e[31m !!! panic : Must copy the smpc_k file to current path"
    exit
fi

git clone https://github.com/TpsSmpc/speedUpTx
cd cd speedUpTx/
go build
cp speedUpTx ../speedUpTool
cd ..


./speedUpTool