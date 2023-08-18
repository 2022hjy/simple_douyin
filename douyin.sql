SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database douyin;
use douyin;
-- Table structure for comments
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT UNIQUE COMMENT '评论的唯一ID',
    `user_info_id` bigint unsigned DEFAULT 0 COMMENT '发表评论的用户的引用ID',
    `video_id`     bigint unsigned DEFAULT 0 COMMENT '被评论的视频的引用ID',
    `content`      varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论内容',
    `action_type`  tinyint NOT NULL DEFAULT 0 COMMENT '评论的类型，1：添加评论，2：删除评论',
    `created_at`   datetime(3) DEFAULT NULL COMMENT '评论创建时间',
    PRIMARY KEY (`id`),
    KEY            `fk_videos_comments` (`video_id`) USING BTREE,
    KEY            `fk_user_infos_comments` (`user_info_id`) USING BTREE,
    FULLTEXT INDEX `idx_content` (`content`),
    CONSTRAINT `fk_user_infos_comments` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`),
    CONSTRAINT `fk_videos_comments` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='评论表';

-- Table structure for user_favorite_video
DROP TABLE IF EXISTS `user_like_video`;
CREATE TABLE `user_like_video`
(
    `user_info_id` bigint unsigned NOT NULL COMMENT '喜欢该视频的用户的引用ID',
    `video_id`     bigint unsigned NOT NULL COMMENT '被喜欢的视频的引用ID',
    PRIMARY KEY (`user_info_id`, `video_id`),
    KEY            `fk_user_favorite_video_video` (`video_id`) USING BTREE,
    CONSTRAINT `fk_user_favorite_video_user_info` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`),
    CONSTRAINT `fk_user_favorite_video_video` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户喜欢的视频表';

---- Table structure for user_info
--DROP TABLE IF EXISTS `user_info`;
--CREATE TABLE `user_info`
--(
--    `id`             bigint unsigned NOT NULL AUTO_INCREMENT  unique  COMMENT '用户的唯一ID',
--    `username`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL unique DEFAULT '' COMMENT '用户名',
--    `follow_count`   int                                              DEFAULT 0 COMMENT '该用户关注的人数',
--    `follower_count` int                                              DEFAULT 0 COMMENT '关注该用户的人数',
--    `is_follow`      tinyint unsigned DEFAULT 0 COMMENT '标记是否被关注',
--    PRIMARY KEY (`id`)
--) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户信息表';
--
---- Table structure for user_login
--DROP TABLE IF EXISTS `user`;
--CREATE TABLE `user`
--(
--    `id`           bigint unsigned NOT NULL AUTO_INCREMENT unique  COMMENT '登录记录的唯一ID',
--    `user_info_id` bigint unsigned DEFAULT 0 COMMENT '相关用户的引用ID',
--    `username`     varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录用户名',
--    `password`     varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录密码',
--    PRIMARY KEY (`id`),
--    UNIQUE INDEX `idx_username` (`username`),
--    KEY            `fk_user_info_login` (`user_info_id`) USING BTREE,
--    CONSTRAINT `fk_user_info_login` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`)
--) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户登录信息表';

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


-- Table structure for relation
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '关系ID',
    `user_id`      bigint unsigned NOT NULL COMMENT '用户ID',
    `following_id` bigint unsigned NOT NULL COMMENT '关注的用户ID',
    `is_followed`  tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '是否已关注，1表示已关注，0表示未关注',
    `created_at`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `updated_at`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    key            `idx_relation_user_id` (`user_id`) USING BTREE,
    key            `idx_relation_following_id` (`following_id`) USING BTREE,
    CONSTRAINT `fk_relation_user` FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_relation_following` FOREIGN KEY (`following_id`) REFERENCES `user_info` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='用户关系表';

-- Table structure for video
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
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
    FULLTEXT INDEX `idx_title` (`title`),
    CONSTRAINT `fk_user_info_video` FOREIGN KEY (`user_info_id`) REFERENCES `user_info` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='视频信息表';

DROP TABLE IF EXISTS `user_chat`;
CREATE TABLE `user_chat`
(
    msg_id      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '聊天id',
    sender_id    BIGINT UNSIGNED NOT NULL COMMENT '发送者用户id',
    recipient_id BIGINT UNSIGNED NOT NULL COMMENT '接收者用户id',
    msg_content      TEXT      NOT NULL COMMENT '消息内容',
    send_time    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    PRIMARY KEY (msg_id),
    FOREIGN KEY (sender_id) REFERENCES `user_info` (`id`),
    FOREIGN KEY (recipient_id) REFERENCES `user_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存储私人聊天记录' ROW_FORMAT=DYNAMIC;