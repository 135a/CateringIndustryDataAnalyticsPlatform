-- 创建数据库 (指定使用 utf8mb4 字符集，完美支持中文字符和 Emoji)
CREATE DATABASE IF NOT EXISTS `catering_data` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `catering_data`;

-- 1. 用户表 (Users) - 用于平台登录与注册
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    `role` TINYINT DEFAULT 0 COMMENT '角色：0-普通用户, 1-管理员',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台用户表';

-- 2. 餐饮分类表 (Categories) - 用于存放您图片中提到的30多种分类
CREATE TABLE IF NOT EXISTS `categories` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '分类名称，如：小吃快餐, 咖啡, 火锅',
    `sort_order` INT DEFAULT 0 COMMENT '排序权重(前端展示用)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='餐饮分类字典表';

-- 3. 城市行政区表 (Districts) - 用于存放杭州市各个区域，支持地理位置过滤
CREATE TABLE IF NOT EXISTS `districts` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '行政区ID',
    `city_name` VARCHAR(50) DEFAULT '杭州市' COMMENT '所属城市',
    `district_name` VARCHAR(50) NOT NULL COMMENT '行政区名称，如：西湖区, 滨江区',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uk_city_district` (`city_name`, `district_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='城市行政区字典表';

-- 4. 餐饮商户主表 (Restaurants) - 核心表，存储爬虫抓取回来的商户明细
CREATE TABLE IF NOT EXISTS `restaurants` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '系统内部主键ID',
    `shop_id` VARCHAR(100) UNIQUE COMMENT '第三方平台(大众点评)的商户ID，用于防重和更新数据',
    `name` VARCHAR(255) NOT NULL COMMENT '商户名称',
    `category_id` INT NOT NULL COMMENT '关联分类ID',
    `district_id` INT NOT NULL COMMENT '关联行政区ID',
    `address` VARCHAR(500) COMMENT '详细地址',
    `avg_price` DECIMAL(10, 2) COMMENT '人均消费价格',
    `rating` DECIMAL(3, 1) COMMENT '综合星级评分 (如 4.5)',
    `review_count` INT DEFAULT 0 COMMENT '评价总数 (重要的人气指标)',
    `opening_hours` VARCHAR(255) COMMENT '营业时间 (如 10:00-22:00)',
    `taste_score` DECIMAL(3, 1) COMMENT '口味细分评分',
    `environment_score` DECIMAL(3, 1) COMMENT '环境细分评分',
    `service_score` DECIMAL(3, 1) COMMENT '服务细分评分',
    `crawled_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据最近一次爬取或更新的时间',
    
    -- 设置外键关联
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE RESTRICT,
    FOREIGN KEY (`district_id`) REFERENCES `districts`(`id`) ON DELETE RESTRICT,
    
    -- 建立索引：后端查询和过滤会非常依赖这些字段，加索引极大提升查询速度
    INDEX `idx_district_category` (`district_id`, `category_id`),
    INDEX `idx_avg_price` (`avg_price`),
    INDEX `idx_rating` (`rating`),
    INDEX `idx_review_count` (`review_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='餐饮商户明细表';

-- ----------------------------------------------------------------------
-- 下面是给字典表预插入的一些基础数据，方便项目冷启动
-- ----------------------------------------------------------------------

-- 插入杭州市主要行政区数据
INSERT IGNORE INTO `districts` (`city_name`, `district_name`) VALUES 
('杭州市', '西湖区'), ('杭州市', '上城区'), ('杭州市', '拱墅区'), 
('杭州市', '滨江区'), ('杭州市', '萧山区'), ('杭州市', '余杭区'), 
('杭州市', '临平区'), ('杭州市', '钱塘区'), ('杭州市', '富阳区'), 
('杭州市', '临安区');

-- 插入根据您截图识别的初步餐饮分类数据
INSERT IGNORE INTO `categories` (`name`) VALUES 
('小吃快餐'), ('咖啡'), ('自助餐'), ('面包甜点'), ('酒吧'), 
('烧烤烤串'), ('创意菜'), ('鱼鲜海鲜'), ('水果生鲜'), ('饮品'),
('特色菜'), ('地方菜系'), ('食品滋补'), ('农家菜'), ('私房菜'),
('家常菜'), ('粤菜'), ('川菜'), ('面馆'), ('江浙菜'),
('火锅'), ('茶馆'), ('西餐'), ('日式料理'), ('小龙虾'),
('烤肉'), ('湘菜'), ('东北菜'), ('北京菜'), ('韩式料理');
