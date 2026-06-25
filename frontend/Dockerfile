# ==============================
# 第一阶段：前端编译阶段 (Node.js)
# ==============================
FROM node:18-alpine AS builder

WORKDIR /app

# 利用 Docker 缓存层，先安装依赖
COPY package*.json ./
# 设置淘宝镜像加速下载前端包
RUN npm config set registry https://registry.npmmirror.com && npm install

# 拷贝全部前端代码并打包
COPY . .
RUN npm run build


# ==============================
# 第二阶段：生产运行阶段 (Nginx)
# ==============================
FROM nginx:alpine

# 删掉 Nginx 默认的毫无用处的默认页面配置
RUN rm /etc/nginx/conf.d/default.conf

# 放入我们刚才写好的用来支持 Vue History 路由模式的配置文件
COPY nginx.conf /etc/nginx/conf.d/

# 将第一阶段 Node.js 编译出的 dist 静态产物夹子，全部丢进 nginx 服务器里
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露 80 端口
EXPOSE 80

# 保持 Nginx 在前台运行
CMD ["nginx", "-g", "daemon off;"]
