-- ----------------------------im服务中心----------------------------
-- ----------------------------
-- Table structure for friend
-- ----------------------------
DROP TABLE IF EXISTS `friend`;
CREATE TABLE `friend`
(
    `id`          bigint unsigned NOT NULL COMMENT '自增主键',
    `apply_user`  bigint NOT NULL COMMENT '申请用户id',
    `accept_user` bigint NOT NULL COMMENT '接受用户id',
    `apply_reason`      char(50)   NOT NULL COMMENT '申请理由',
    `extra`       varchar(1024) NOT NULL COMMENT '附加属性',
    `status`      tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，-2：拉黑，-1：拒绝，0：申请中，1：同意',
    `create_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `k_create_time_status` (`create_time`, `status`),
    UNIQUE KEY `uk_apply_user_accept_user` (`apply_user`, `accept_user`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='好友表';

-- ----------------------------
-- Table structure for seq
-- ----------------------------
DROP TABLE IF EXISTS `seq`;
CREATE TABLE `seq`
(
    `id`          bigint unsigned NOT NULL COMMENT '自增主键',
    `object_type` tinyint  NOT NULL COMMENT '对象类型,1:用户；2：群组',
    `object_id`   bigint NOT NULL COMMENT '对象id',
    `seq`         bigint NOT NULL COMMENT '序列号',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_object` (`object_type`,`object_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='序列号';

-- ----------------------------
-- Table structure for single_msg
-- ----------------------------
DROP TABLE IF EXISTS `single_msg`;
CREATE TABLE `single_msg`
(
    `id`            bigint unsigned NOT NULL COMMENT '自增主键',
    `seq`           bigint NOT NULL COMMENT '消息序列号,每个单聊都维护一个消息序列号',
    `sender_type`   tinyint(3) NOT NULL DEFAULT '1' COMMENT '发送者类型：1朋友，2打招呼，3转发',
    `sender_id`     bigint NOT NULL COMMENT '发送者id',
    `sender_device_id` varchar(100) NOT NULL COMMENT '发送设备id',
    `receiver_id`   bigint NOT NULL COMMENT '接收者id',
    `receiver_device_id` varchar(100) NOT NULL COMMENT '接收设备id：多个设备id"，"隔开，*表示全部设备',
    `msg_type`      tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息类型：0文本、1图文、2语音、3视频、4链接',
    `content`       blob NOT NULL COMMENT '消息内容',
    `parent_id`     bigint NOT NULL COMMENT '父级id，引用功能',
    `send_time`     datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT '消息发送时间',
    `status`        tinyint(255) NOT NULL DEFAULT '0' COMMENT '消息状态：-1撤回，0未处理，1已读',
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `k_send_time_status_msg_type_sender_type` (`send_time`,`status`,`msg_type`,`sender_type`),
    UNIQUE KEY `uk_seq_sender_id_receiver_id` (`seq`, `sender_id`, `receiver_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='单聊消息表';

-- ----------------------------
-- Table structure for group_msg
-- `receiver_type` tinyint(3) NOT NULL COMMENT '接收者类型, 1普通群组，2超大群组',
-- ----------------------------
DROP TABLE IF EXISTS `group_msg`;
CREATE TABLE `group_msg`
(
    `id`            bigint unsigned NOT NULL COMMENT '自增主键',
    `seq`           bigint NOT NULL COMMENT '消息序列号,每个单聊都维护一个消息序列号',
    `sender_type`   tinyint(3) NOT NULL DEFAULT '1' COMMENT '发送者类型：1群内，2转发',
    `sender_id`     bigint NOT NULL COMMENT '发送者id',
    `sender_device_id` varchar(100) NOT NULL COMMENT '发送设备id',
    `receiver_id`   bigint NOT NULL COMMENT '接收者id, group_id',
    `receiver_device_id` varchar(100) NOT NULL COMMENT '接收设备id：多个设备id"，"隔开，*表示全部设备',
    `to_user_ids`  varchar(255) NOT NULL COMMENT '需要@的用户id列表，多个用户用@隔开',
    `msg_type`          tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息类型：0文本、1图文、2语音、3视频、地址、4链接',
    `content`       blob         NOT NULL COMMENT '消息内容',
    `parent_id`         bigint NOT NULL DEFAULT '0' COMMENT '父级id，引用功能',
    `send_time`     datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT '消息发送时间',
    `status`        tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息状态，-3接收者删除，-2发送者删除，-1撤回，0未处理，1已读',
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `k_send_time_status_msg_type_sender_type` (`send_time`,`status`,`msg_type`,`sender_type`),
    UNIQUE KEY `uk_seq_sender_id_receiver_id` (`seq`, `sender_id`, `receiver_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='群聊消息表';