/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3307
 Source Server Type    : MySQL
 Source Server Version : 80100
 Source Host           : localhost:3307
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 80100
 File Encoding         : 65001

 Date: 30/08/2023 12:54:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论的唯一ID',
  `user_info_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '发表评论的用户的引用ID',
  `video_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '被评论的视频的引用ID',
  `content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论内容',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '评论创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `id`(`id` ASC) USING BTREE,
  INDEX `fk_videos_comments`(`video_id` ASC) USING BTREE,
  INDEX `fk_user_infos_comments`(`user_info_id` ASC) USING BTREE,
  FULLTEXT INDEX `idx_content`(`content`)
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '评论表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for favorite
-- ----------------------------
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞视频的唯一标识',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '点赞视频的用户，关联ID',
  `video_id` bigint UNSIGNED NOT NULL COMMENT '被点赞的视频，关联ID',
  `is_favorite` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '视频是否已被点赞，1 表示已点赞，0 表示未点赞',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '记录创建时间',
  PRIMARY KEY (`user_id`, `video_id`) USING BTREE,
  UNIQUE INDEX `id`(`id` ASC) USING BTREE,
  INDEX `fk_user_favorite_video`(`video_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户点赞视频表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of favorite
-- ----------------------------

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '聊天id',
  `to_user_id` bigint UNSIGNED NOT NULL COMMENT '接收者id',
  `from_user_id` bigint UNSIGNED NOT NULL COMMENT '发送者id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '消息内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `to_user_id`(`to_user_id` ASC) USING BTREE,
  INDEX `from_user_id`(`from_user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '聊天消息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------
INSERT INTO `message` VALUES (2, 3, 2, 'Hello', '2023-08-30 12:52:37');
INSERT INTO `message` VALUES (3, 3, 2, '你好呀', '2023-08-30 12:52:39');
INSERT INTO `message` VALUES (4, 2, 3, '你拍的风景视频很好看', '2023-08-30 12:53:17');
INSERT INTO `message` VALUES (5, 3, 2, '你拍的猫猫打拳也很可爱', '2023-08-30 12:53:45');

-- ----------------------------
-- Table structure for relation
-- ----------------------------
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系ID',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `following_id` bigint UNSIGNED NOT NULL COMMENT '关注的用户ID',
  `is_followed` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已关注，1表示已关注，0表示未关注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_relation_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_relation_following_id`(`following_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户关系表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of relation
-- ----------------------------
INSERT INTO `relation` VALUES (2, 3, 2, 1, '2023-08-30 12:52:03', '2023-08-30 12:52:03');
INSERT INTO `relation` VALUES (3, 2, 3, 1, '2023-08-30 12:52:29', '2023-08-30 12:52:29');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '登录记录的唯一ID',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '登录用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '加密后的登录密码',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '用户头像',
  `background_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '用户个人页顶部大图',
  `signature` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '个人简介',
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE INDEX `idx_username`(`username` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'user1', '$2a$04$D0gotw7GO65uraamc7jVOewQez8lC1ovYSVpPw0dQkOgMwax69Nme', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/zgs.jpg', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/%E8%83%8C%E6%99%AF1.jpg', 'Hello, I am user1');
INSERT INTO `user` VALUES (2, 'user2', '$2a$04$D0gotw7GO65uraamc7jVOewQez8lC1ovYSVpPw0dQkOgMwax69Nme', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/hjy.jpg', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/%E8%83%8C%E6%99%AF2.jpg', 'Hello, I am user2');
INSERT INTO `user` VALUES (3, 'user3', '$2a$04$D0gotw7GO65uraamc7jVOewQez8lC1ovYSVpPw0dQkOgMwax69Nme', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/zx.jpg', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_PIC/%E8%BF%99%E6%98%AF%E8%B0%81%E7%9A%84%E5%B1%8F%E5%B9%95%E8%BF%99%E4%B9%88%E5%A5%BD%E7%9C%8B.jpg', '耶耶耶');

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '视频的唯一ID',
  `user_info_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '上传视频的用户的引用ID',
  `play_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频播放URL',
  `cover_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频封面图片URL',
  `favorite_count` int NULL DEFAULT 0 COMMENT '视频被喜欢的次数',
  `comment_count` int NULL DEFAULT 0 COMMENT '视频的评论数',
  `is_favorite` tinyint UNSIGNED NULL DEFAULT 0 COMMENT '标记视频是否被喜欢',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频标题',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '视频创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '视频上次更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `id`(`id` ASC) USING BTREE,
  INDEX `fk_user_info_video`(`user_info_id` ASC) USING BTREE,
  FULLTEXT INDEX `idx_title`(`title`)
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '视频信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of video
-- ----------------------------
INSERT INTO `video` VALUES (1, 1, 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/b2fe620c-1ceb-4332-b075-644a3524019a.mp4', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/b2fe620c-1ceb-4332-b075-644a3524019a.mp4?x-oss-process=video/snapshot,t_2000,m_fast', 0, 0, 0, '小熊喝水', '2023-08-28 16:53:17.794', '2023-08-30 12:49:05.804');
INSERT INTO `video` VALUES (2, 2, 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/0db23502-1003-4011-b4e0-eb13b7a4fb7d.mp4', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/0db23502-1003-4011-b4e0-eb13b7a4fb7d.mp4?x-oss-process=video/snapshot,t_2000,m_fast', 1, 0, 0, '好看的风景', '2023-08-30 12:29:24.868', '2023-08-30 12:49:13.187');
INSERT INTO `video` VALUES (3, 3, 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/b06da13f-57da-4b65-87b4-672b8aacfcd1.mp4', 'https://simple-tiktok.oss-cn-beijing.aliyuncs.com/OSS_VIDEO/b06da13f-57da-4b65-87b4-672b8aacfcd1.mp4?x-oss-process=video/snapshot,t_2000,m_fast', 0, 0, 0, '猫猫打拳', '2023-08-30 12:35:41.303', '2023-08-30 12:35:41.303');

SET FOREIGN_KEY_CHECKS = 1;
