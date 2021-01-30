rm fund log
gofmt -w .
git add -u
go build
ps -ef | grep fund | awk '{print $2}' | xargs kill -9
nohup ./fund data/conf.json > log 2>&1 &

