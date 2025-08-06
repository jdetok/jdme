# migrate to pg - copy dirs from go-api-jdeko.me
GODIR=/home/jdeto/go/github.com/jdetok
PROD=$GODIR/go-api-jdeko.me
DEV=$GODIR/go-api-jdeko.me

cd $PROD

cp -r $DEV/pgdb pgdb
go get github.com/jdetok/golib
cp -r $DEV/api $PROD/api

cp $DEV/.air.toml .dev.air.toml # verify & update manually
cp $DEV/.env dev.env # verify & update manually

rm -r getenv
rm -r applog
rm -r .vscode
rm -r mdb 
go mod tidy