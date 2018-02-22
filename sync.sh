TARGET=$1
if [[ -z $TARGET ]];then
  echo "Please specify sync target"
fi

case $TARGET in
    "aliyun" )
        rsync -avz --delete -e "ssh -p 8899" --exclude "wegigo" --exclude "vendor" ryan@guanxigo.com:/home/ryan/go/src/github.com/rfancn/wegigo/* .
        ;;
    "hdget" )
    	rsync -avz --delete -e "ssh -p 8899" ryan@hdget.com:/home/ryan/go/src/github.com/rfancn/wegigo/Gopkg.lock .
    	rsync -avz --delete -e "ssh -p 8899" ryan@hdget.com:/home/ryan/go/src/github.com/rfancn/wegigo/Gopkg.toml .
    	rsync -avz --delete -e "ssh -p 8899" ryan@hdget.com:/home/ryan/go/src/github.com/rfancn/wegigo/vendor .
        ;;
esac

