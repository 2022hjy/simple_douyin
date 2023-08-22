SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database douyin;
use douyin;
-- Table structure for comments
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT UNIQUE COMMENT '评论的唯一ID',
    `user_info_id` bigint unsigned DEFAULT 0 COMMENT '发表评论的用户的引用ID',
    `video_id`     bigint unsigned DEFAULT 0 COMMENT '被评论的视频的引用ID',
    `content`      varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论内容',
    `created_at`   datetime(3) DEFAULT NULL COMMENT '评论创建时间',
    PRIMARY KEY (`id`),
    KEY            `fk_videos_comments` (`video_id`) USING BTREE,
    KEY            `fk_user_infos_comments` (`user_info_id`) USING BTREE,
  FULLTEXT INDEX `idx_content`(`content`),
    CONSTRAINT `fk_user_infos_comments` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`),
    CONSTRAINT `fk_videos_comments` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='评论表';

-- Table structure for user_favorite_video
-- 如果存在，则删除表格`favorite`
DROP TABLE IF EXISTS `favorite`;
-- 创建表格`favorite`
CREATE TABLE `favorite`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT UNIQUE COMMENT '点赞视频的唯一标识',
    `user_id`      bigint unsigned NOT NULL COMMENT '点赞视频的用户，关联ID',
    `video_id`     bigint unsigned NOT NULL COMMENT '被点赞的视频，关联ID',
    `is_favorite`  tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '视频是否已被点赞，1 表示已点赞，0 表示未点赞',
    `created_at`   datetime(3) DEFAULT NULL COMMENT '记录创建时间',
    PRIMARY KEY (`user_id`, `video_id`),
    KEY            `fk_user_favorite_video` (`video_id`) USING BTREE,
    CONSTRAINT `fk_user_favorite_user_info` FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`),
    CONSTRAINT `fk_user_favorite_video` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户点赞视频表';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '聊天id',
  `from_user_id` bigint UNSIGNED NOT NULL COMMENT '发送者id',
  `to_user_id` bigint UNSIGNED NOT NULL COMMENT '接收者id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '消息内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `to_user_id`(`to_user_id` ASC) USING BTREE,
  INDEX `from_user_id`(`from_user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '聊天消息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------
INSERT INTO `message` VALUES (1, 1, 2, '测试消息', '2023-08-22 03:49:40');
INSERT INTO `message` VALUES (2, 1, 2, '测试消息2', '2023-08-22 13:17:10');
INSERT INTO `message` VALUES (3, 10, 2, '2 号用户发送消息给 8 号用户', '2023-08-22 13:31:14');
INSERT INTO `message` VALUES (6, 10, 2, '2 - 10', '2023-08-22 13:45:43');
INSERT INTO `message` VALUES (7, 2, 26, '226', '2023-08-22 13:50:39');
INSERT INTO `message` VALUES (9, 2, 26, '2-26', '2023-08-22 13:52:05');
INSERT INTO `message` VALUES (10, 2, 26, '226', '2023-08-22 14:01:05');

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
  INDEX `idx_relation_following_id`(`following_id` ASC) USING BTREE,
  CONSTRAINT `fk_relation_following` FOREIGN KEY (`following_id`) REFERENCES `user_info` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_relation_user` FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户关系表' ROW_FORMAT = DYNAMIC;

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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT  unique COMMENT '视频的唯一ID',
    `user_info_id`   bigint unsigned DEFAULT 0 COMMENT '上传视频的用户的引用ID',
    `play_url`       varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频播放URL',
    `cover_url`      varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频封面图片URL',
    `favorite_count` int                                               DEFAULT 0 COMMENT '视频被喜欢的次数',
    `comment_count`  int                                               DEFAULT 0 COMMENT '视频的评论数',
    `is_favorite`    tinyint unsigned DEFAULT 0 COMMENT '标记视频是否被喜欢',
    `title`          varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT '' COMMENT '视频标题',
    `created_at`     datetime(3) DEFAULT NULL COMMENT '视频创建时间',
    `updated_at`     datetime(3) DEFAULT NULL COMMENT '视频上次更新时间',
    PRIMARY KEY (`id`),
    KEY              `fk_user_info_video` (`user_info_id`),
  FULLTEXT INDEX `idx_title`(`title`),
  CONSTRAINT `fk_user_info_video` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '视频信息表' ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
