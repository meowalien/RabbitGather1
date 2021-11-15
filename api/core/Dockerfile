FROM golang:1.16.5

WORKDIR "/app"
COPY ./ ./

# 預設關閉DEBUG模式
ENV DEBUG="false"

# RUN go install /app/src/
# RUN go build -o /app/bin/server.out /app/src/
RUN chmod +x /app/bin/core_bin_linux

# 容器對外暴露的端口
EXPOSE 2001
EXPOSE 2002

ENTRYPOINT cd /app/ && /app/bin/core_bin_linux -debug=$DEBUG
