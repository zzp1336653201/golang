-- CMS 电商系统数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS cms DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE cms;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（加密）',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:正常 0:禁用',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 商品表
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `description` text COMMENT '商品描述',
  `price` decimal(10,2) NOT NULL COMMENT '价格（分）',
  `stock` int NOT NULL DEFAULT '0' COMMENT '库存',
  `category` varchar(100) DEFAULT NULL COMMENT '分类',
  `image_url` varchar(500) DEFAULT NULL COMMENT '图片URL',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:上架 0:下架',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_category` (`category`),
  KEY `idx_products_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 订单表
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(32) NOT NULL COMMENT '订单号',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总额',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:待支付 2:已支付 3:已发货 4:已收货 5:已完成 6:已取消',
  `receiver_name` varchar(100) DEFAULT NULL COMMENT '收货人姓名',
  `receiver_phone` varchar(20) DEFAULT NULL COMMENT '收货人电话',
  `receiver_addr` text COMMENT '收货地址',
  `remark` text COMMENT '备注',
  `paid_at` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `shipped_at` datetime(3) DEFAULT NULL COMMENT '发货时间',
  `completed_at` datetime(3) DEFAULT NULL COMMENT '完成时间',
  `canceled_at` datetime(3) DEFAULT NULL COMMENT '取消时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_orders_user_id` (`user_id`),
  KEY `idx_orders_status` (`status`),
  KEY `idx_orders_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 订单项表
CREATE TABLE IF NOT EXISTS `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `product_name` varchar(255) DEFAULT NULL COMMENT '商品名称（快照）',
  `price` decimal(10,2) NOT NULL COMMENT '购买价格',
  `quantity` int NOT NULL COMMENT '购买数量',
  `sub_total` decimal(10,2) NOT NULL COMMENT '小计',
  PRIMARY KEY (`id`),
  KEY `idx_order_items_order_id` (`order_id`),
  KEY `idx_order_items_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单项表';

-- 购物车表
CREATE TABLE IF NOT EXISTS `cart_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `quantity` int NOT NULL COMMENT '数量',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cart_user_product` (`user_id`,`product_id`),
  KEY `idx_cart_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';

-- 插入测试数据
INSERT INTO `users` (`username`, `password`, `email`, `nickname`) VALUES 
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@example.com', '管理员');

INSERT INTO `products` (`name`, `description`, `price`, `stock`, `category`, `image_url`) VALUES
('iPhone 15 Pro', '苹果旗舰手机，A17 Pro芯片', 899900, 100, '手机', 'https://example.com/iphone15.jpg'),
('MacBook Pro 14', 'M3 Pro芯片，专业级笔记本', 1499900, 50, '电脑', 'https://example.com/macbook.jpg'),
('AirPods Pro 2', '主动降噪，空间音频', 189900, 200, '配件', 'https://example.com/airpods.jpg');
