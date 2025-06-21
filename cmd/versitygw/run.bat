rem set ROOT_ACCESS_KEY=testuser
rem set ROOT_SECRET_KEY=secret

.\netzongw.exe --port :11000 --access admin --secret admin posix --notify-base-url http://127.0.0.1:8080 --notify-endpoint-path /oss/rest/restServer/creatArcTask E:\vgw