# 🍱 餐饮大数据平台实时监控中心
(Catering Industry Data Analytics Platform)

本项目是一个功能完整的餐饮大数据抓取、清洗、存储、聚合分析及实时大屏展示平台。主要用于采集和监控全城各大行政区、各种菜系的餐饮商户数据（涵盖人均消费、综合星级、详细评分、服务设施等多维指标），帮助用户发现“高分低价”的宝藏餐厅，并为餐饮行业的市场洞察提供数据支撑。

---

## ✨ 核心特性 / Features

*   **📊 实时监控大屏 (Dashboard)**：
    基于 **Vue 3** + **ECharts** 构建极具科技感、玻璃质感（Glassmorphism）的暗黑风数据大屏。支持动态渲染：热门品类占比（玫瑰饼图）、各区商户与消费双轴对比（柱状图+折线图）、价格与评分性价比分布（海量数据散点图）。
*   **🎯 深度数据下钻**：
    在“性价比散点图”中发现高分低价的“宝藏发光点”时，点击气泡即可瞬间下钻！系统会自动唤起玻璃质感大弹窗，精准展示该商户的详细数据（包含口味/环境/服务评分、有无停车位、是否包厢等）。
*   **🕷️ 全自动爬虫流水线**：
    内置 **Python** 数据采集引擎。项目启动时自动运转，覆盖全城 10 大行政区、30+ 细分菜系。支持“双写机制”：一方面将数据安全合并更新到 MySQL 数据库（`ON DUPLICATE KEY UPDATE` 防重防漏），另一方面同步在宿主机持久化保存当日 Excel 备份。
*   **⚡ 高性能后端服务**：
    后端基于 **Golang** (Gin + GORM) 编写。为前端 ECharts 提供毫秒级的大数据聚合计算（AVG/COUNT/SUM）及多维度（行政区、品类、设施、关键字）分页搜索 API。
*   **🐳 一键式容器化部署**：
    全栈采用 **Docker & Docker Compose** 编排。无需繁杂的环境配置，一条指令即可完成数据库初始化、后端编译、前端打包及爬虫拉起，实现开箱即用。

---

## 🛠️ 技术栈 / Tech Stack

*   **前端 (Frontend)**: Vue 3, Vite, Vue Router, ECharts 5, Vanilla CSS
*   **后端 (Backend)**: Golang 1.21+, Gin Web Framework, GORM
*   **爬虫 (Crawler)**: Python 3.9, Pandas, PyMySQL, Openpyxl
*   **数据库 (Database)**: MySQL 8.0
*   **部署 (Deployment)**: Docker, Docker Compose, Nginx

---

## 📂 项目结构 / Project Structure

```text
CateringIndustryDataAnalyticsPlatform/
├── backend/                  # Golang 后端 API 服务
│   ├── cmd/main.go           # 后端服务入口
│   ├── internal/             # 业务逻辑代码
│   │   ├── api/              # RESTful API 控制器 (数据统计、商户列表等)
│   │   ├── model/            # GORM 数据库映射模型
│   │   └── router/           # 路由注册与 CORS 配置
│   ├── pkg/database/         # MySQL 数据库连接池初始化
│   ├── Dockerfile            # 后端多阶段构建配置
│   └── go.mod / go.sum       # Go 依赖管理
│
├── frontend/                 # Vue 3 大屏前端项目
│   ├── src/
│   │   ├── views/            # 大屏看板 (Dashboard.vue) 与 登录页 (Login.vue)
│   │   ├── utils/request.js  # Axios 拦截器与请求封装
│   │   ├── router/index.js   # Vue Router 路由守卫
│   │   └── assets/           # 全局样式与静态资源
│   ├── nginx.conf            # Nginx 代理与路由配置文件
│   ├── Dockerfile            # 前端打包及 Nginx 镜像构建
│   └── package.json          # Node 依赖
│
├── crawler/                  # Python 爬虫数据采集引擎
│   └── dianping_spider.py    # 抓取逻辑、数据库双写机制与 Excel 导出
│
├── data/                     # 数据持久化挂载目录
│   ├── excel/                # 爬虫生成的 Excel 每日数据备份
│   ├── mysql_data/           # MySQL 物理数据文件 (Docker Volume)
│   └── sql/                  # schema.sql 数据库自动建表脚本
│
└── docker-compose.yml        # 全局容器编排文件
```

---

## 🚀 快速启动 / Quick Start

**环境要求**：您的机器上必须已安装 **Docker** 和 **Docker Compose**。

1. **克隆项目到本地**
2. **一键启动**
   在项目根目录下，打开终端并执行以下命令：
   ```bash
   docker-compose up -d --build
   ```
   > *说明：初次运行会拉取环境镜像并进行编译打包，大概需要 1~3 分钟。*

3. **观察爬虫运行情况**
   服务启动后，Python 爬虫会自动运行并向数据库灌入初始数据。
   ```bash
   docker logs -f catering_crawler
   ```
   
4. **访问平台**
   当容器启动完毕后，在浏览器中访问：
   - **地址**：`http://localhost`
   - **账号**：admin (任意用户名)
   - **密码**：任意密码 (暂未启用严格后端强校验)

---

## 💡 API 架构概览

所有 API 的基础路径为 `/api/v1`：

| 接口路径 | 方法 | 功能说明 |
| :--- | :---: | :--- |
| `/districts` | GET | 获取所有的行政区划字典列表 |
| `/categories` | GET | 获取所有的餐饮菜系字典列表 |
| `/statistics/overview` | GET | 获取大屏核心指标卡数据（商户总数、均价等） |
| `/statistics/category-pie` | GET | 获取热门品类占比（供饼图渲染） |
| `/statistics/district-bar` | GET | 获取各区商户与消费双轴对比（供柱状图渲染） |
| `/statistics/price-rating-scatter` | GET | 获取极大数据量下的人均与评分散点坐标数据 |
| `/restaurants` | GET | 带分页与多维度筛选条件的商户详情搜索列表 |

---

## 👨‍💻 进阶开发 & 配置

- **如何调节爬虫数据量**：
  打开 `crawler/dianping_spider.py`，修改 `FETCH_COUNT_MIN` 和 `FETCH_COUNT_MAX`，重启 crawler 容器即可抓取更密集的数据。
- **重新初始化数据库**：
  如果希望清空所有数据重新开始，请先执行 `docker-compose down -v`，然后**手动删除** `data/mysql_data` 目录下的所有文件，最后再次执行 `docker-compose up -d` 即可触发 `schema.sql` 重新建表。

---

*Enjoy exploring the data and discovering the hidden gems of the city!* 🍷
