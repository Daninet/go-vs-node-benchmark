bombardier -c 10 -d 10s -l --print=i,r http://192.168.5.10:3001/json

sleep 10

bombardier -c 10 -d 10s -l --print=i,r http://192.168.5.10:3002/json

sleep 10

bombardier -c 100 -d 10s -l --print=i,r http://192.168.5.10:3001/json

sleep 10

bombardier -c 100 -d 10s -l --print=i,r http://192.168.5.10:3002/json

sleep 10

bombardier -c 250 -d 10s -l --print=i,r http://192.168.5.10:3001/json

sleep 10

bombardier -c 250 -d 10s -l --print=i,r http://192.168.5.10:3002/json

sleep 10
