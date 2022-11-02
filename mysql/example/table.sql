DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
 `name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
 `tel` varchar(20) NOT NULL DEFAULT '' COMMENT '电话号码',
 `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
 `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
 `update_at` timestamp NULL DEFAULT NULL COMMENT '修改时间',
 PRIMARY KEY (`id`) USING BTREE,
 UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
CREATE TABLE IF NOT EXISTS `user_info` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
 `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户uuid',
 `number` int(3) zerofill NOT NULL DEFAULT '000' COMMENT '编号',
 `content` json NULL DEFAULT NULL COMMENT '内容',
 `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
 `update_at` timestamp NULL DEFAULT NULL COMMENT '修改时间',
 PRIMARY KEY (`id`) USING BTREE,
 CONSTRAINT `test_a_b_c` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户详情表';
ALTER TABLE `user` CHANGE COLUMN `tel` `phone` varchar(50) NOT NULL DEFAULT '' COMMENT '手机号' AFTER `name`;
ALTER TABLE `user` MODIFY COLUMN `password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码' AFTER `phone`;
ALTER TABLE `user` ADD KEY `phone_password` (`phone`,`password`) USING BTREE;
ALTER TABLE `user_info` DROP CONSTRAINT `test_a_b_c`;
ALTER TABLE `user_info` DROP INDEX `phone_password`;
ALTER TABLE `user_info` MODIFY COLUMN `id` int unsigned NOT NULL COMMENT '主键';
ALTER TABLE `user_info` DROP PRIMARY KEY;
DROP TABLE `user_info`;
DROP TABLE `user`;
