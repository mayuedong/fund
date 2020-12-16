rm fund log
git add -u
gofmt -w .
go build
ps -ef | grep fund | awk '{print $2}' | xargs kill -9
nohup ./fund data/conf.json > log 2>&1 &

