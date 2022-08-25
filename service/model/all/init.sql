-- ----------------------------用户+认证服务中心----------------------------
-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`            bigint unsigned NOT NULL COMMENT '自增主键',
    `nickname`      char(20)   NOT NULL DEFAULT '' COMMENT '昵称',
    `realName`      char(36)   NOT NULL DEFAULT '' COMMENT '真实姓名',
    -- `phone_number`  char(20)   NOT NULL DEFAULT '0' COMMENT '手机号',
    -- `email` char(36)   NOT NULL COMMENT '邮箱',
    -- `id_card`       char(36)   NOT NULL COMMENT '身份证号',
    `password`      char(36)   NOT NULL DEFAULT '0' COMMENT '密码',
    `login_salt`    char(36)   NOT NULL DEFAULT '0' COMMENT '密码加密的盐（md5计算id+注册时间）',
    `register_device`    char(36)   NOT NULL DEFAULT '0' COMMENT '注册设备信息',
    `sex`           tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别，0:未知；1:男；2:女',
    `ico`    varchar(256)  NOT NULL DEFAULT '' COMMENT '用户头像',
    `status`        tinyint(3) NOT NULL DEFAULT '0' COMMENT '用户状态，-2:删除；-1:冻结；0：待审核；1：正常；2：第三方直接注册登入；9：超管',
    `create_time`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建/注册时间',
    `update_time`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `sex`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户';

-- ----------------------------
-- Table structure for user_auths
-- ----------------------------
DROP TABLE IF EXISTS `user_auths`;
CREATE TABLE `user_auths`
(
    `id`            bigint unsigned NOT NULL COMMENT '自增主键',
    `user_id`       bigint NOT NULL COMMENT '用户id',
    `identity_type`      char(36)   NOT NULL COMMENT '身份认证方式，手机号/邮箱/身份证/第三方',
    `identifier`       char(36)   NOT NULL COMMENT '唯一标识符，手机号/邮箱/身份证/第三方open_id',
    -- `credential`         char(36)   NOT NULL COMMENT '凭证 密码/第三方的token(user_uid+注册时间为login_salt，如user_uid-20220218)',
    `status`        tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_identifier` (`identifier`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `identity_type`, `user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户身份验证关系表';

-- ----------------------------
-- Table structure for thirdparty_bind
-- ----------------------------
DROP TABLE IF EXISTS `thirdparty_bind`;
CREATE TABLE `thirdparty_bind`
(
    `id`           bigint unsigned NOT NULL COMMENT '自增主键',
    `user_id` bigint NOT NULL COMMENT '用户id',
    `open_id` char(36) NOT NULL COMMENT '第三方唯一id',
    `source` char(36)   NOT NULL COMMENT '用户来源, 微信/qq/抖音/钉钉/...',
    `auth_scope` varchar(1000)   NOT NULL COMMENT '授权范围',
    `account` char(36)   NOT NULL COMMENT '账号',
    `nickname` char(36)   NOT NULL COMMENT '用户昵称',
    `gender` char(36)   NOT NULL COMMENT '性别',
    `avatar` char(36)   NOT NULL COMMENT '用户头像',
    `blog` char(36)   NOT NULL COMMENT '用户网址',
    `company` char(36)   NOT NULL COMMENT '所在公司',
    `location` char(36)   NOT NULL COMMENT '位置',
    `email` char(36)   NOT NULL COMMENT '邮箱',
    `remark` varchar(500)   NOT NULL COMMENT '用户备注（各平台中的用户个人介绍）',
    `extra`         varchar(1024) NOT NULL DEFAULT '' COMMENT '附加属性',
    `status`        tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_open_id` (`open_id`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `source`, `user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户第三方绑定信息';


-- ----------------------------
-- Table structure for thirdApp_use
-- ----------------------------
DROP TABLE IF EXISTS `thirdApp_use`;
CREATE TABLE `thirdApp_use`
(
    `id`           bigint unsigned NOT NULL COMMENT '自增主键',
    `app_id`      char(36) NOT NULL COMMENT 'app唯一id',
    `app_secret` char(36)   NOT NULL COMMENT 'app密码',
    `auth_scope` varchar(1000)   NOT NULL COMMENT '授权范围',
    `name` char(36)   NOT NULL COMMENT 'app名称',
    `callback_url` char(100)   NOT NULL COMMENT '回调url',
    `ico` char(100)   NOT NULL COMMENT '图标',
    `email` char(36)   NOT NULL COMMENT '邮箱',
    `phone` char(36)   NOT NULL COMMENT '电话',
    `remark` varchar(500)   NOT NULL COMMENT '备注',
    `extra`         varchar(1024) NOT NULL DEFAULT '' COMMENT '附加属性',
    `status`        tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_open_id` (`app_id`),
    KEY `indexFilter_create_time` (`create_time`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='第三方应用授权关系表';


-- ----------------------------
-- Table structure for thirdApp_use_user
-- ----------------------------
DROP TABLE IF EXISTS `thirdApp_use_user`;
CREATE TABLE `thirdApp_use_user`
(
    `id`           bigint unsigned NOT NULL COMMENT '自增主键，open_id',
    `thirdApp_use_id` bigint NOT NULL COMMENT '第三方授权应用唯一id',
    `user_id` bigint NOT NULL COMMENT '用户id',
    `auth_scope` varchar(1000)   NOT NULL COMMENT '授权范围',
    `extra`         varchar(1024) NOT NULL DEFAULT '' COMMENT '附加属性',
    `status`        tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_thirdApp_use_user_id` (`user_id`, `thirdApp_use_id`),
    KEY `indexFilter_create_time` (`create_time`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户授权第三方应用使用关系表';


-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `user_id`        bigint NOT NULL COMMENT '用户id',
    `user_agent` varchar(500) NOT NULL DEFAULT '' COMMENT '用户标示',
    `machine_type`   tinyint(3) NOT NULL COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
    `brand`          char(20) NOT NULL DEFAULT '' COMMENT '设备厂商',
    `unit_type`      char(20) NOT NULL DEFAULT '' COMMENT '设备型号',
    `system_version` char(10) NOT NULL DEFAULT '' COMMENT '系统版本',
    `browser` char(50) NOT NULL DEFAULT '' COMMENT '浏览器',
    `language` char(20) NOT NULL DEFAULT '' COMMENT '语言',
    `net_type` char(20) NOT NULL DEFAULT '' COMMENT '网络类型',
    `sdk_version`    char(10) NOT NULL DEFAULT '' COMMENT 'sdk版本',
    `conn_ip`      char(25) NOT NULL COMMENT '连接层服务器ip',
    `client_ip`    char(25) NOT NULL COMMENT '客户端ip',
    `client_addr`    char(200) NOT NULL DEFAULT '' COMMENT '客户端地址,浙江省杭州市西湖区浙江大学国家大学科技园',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '设备状态，-2删除，-1禁用，0待审核，1离线，2在线',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_user_agent_ip` (`user_agent`,`client_ip`),
    KEY `idx_user_id` (`user_id`) USING BTREE,
    KEY `indexFilter_create_time` (`create_time`, `status`, `machine_type`, `brand`, `unit_type`, `system_version`, `sdk_version`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='设备';

-- ----------------------------
-- Table structure for 用户操作记录
-- ----------------------------

-- ----------------------------组服务中心----------------------------
-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `name`           varchar(20) NOT NULL COMMENT '组名称',
    `create_user`    bigint NOT NULL COMMENT '创建者id',
    `ico`            varchar(20) NOT NULL COMMENT '组图标',
    `remark`         varchar(250) NOT NULL COMMENT '备注',
    `parent_id`      bigint NOT NULL COMMENT '父级id',
    `group_type`     tinyint(100) NOT NULL DEFAULT '1' COMMENT '类型: 1部门、2用户组、3群组、4圈子、5话题',
    `rank`           tinyint(100) NOT NULL DEFAULT '1' COMMENT '排序',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_create_user_name` (`create_user`,`name`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `group_type`, `parent_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='组织管理中心';

-- ----------------------------
-- Table structure for user_group
-- ----------------------------
DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `user_id`        bigint NOT NULL COMMENT '用户id',
    `group_id`       bigint NOT NULL COMMENT '组id',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_group_user_id` (`user_id`,`group_id`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户、组、角色关系表';

-- ----------------------------权限服务中心----------------------------
-- ----------------------------
-- Table structure for operation_auth
-- ----------------------------
DROP TABLE IF EXISTS `operation_auth`;
CREATE TABLE `operation_auth`
(
    `auth_identifier`       varchar(100) NOT NULL COMMENT '权限标识（路由）',
    `auth_name`             char(20) NOT NULL COMMENT '权限名称',
    `auth_type`       varchar(100) NOT NULL COMMENT '权限类型-接口/页面/菜单/按钮',
    `status`                tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    PRIMARY KEY (`auth_identifier`),
    KEY `index_auth_type_name` (`auth_type`, `auth_name`),
    KEY `index_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='操作权限表';

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `name`           char(20) NOT NULL COMMENT '角色名称',
    `operation_auths`           varchar(250) NOT NULL COMMENT '权限列表',
    `rank`           tinyint(100) NOT NULL DEFAULT '1' COMMENT '排序',
    `parent_id`      bigint NOT NULL COMMENT '父级id',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_name` (`name`),
    KEY `indexFilter_create_time` (`create_time`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='角色表';

-- ----------------------------
-- Table structure for user_group_role
-- ----------------------------
DROP TABLE IF EXISTS `user_group_role`;
CREATE TABLE `user_group_role`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `user_group_id`  bigint NOT NULL COMMENT '用户、组关系表id',
    `role_id`        bigint NOT NULL COMMENT '角色id',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_user_group_role_id` (`user_group_id`,`role_id`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `user_group_id`, `role_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户、组、角色关系表';

-- ----------------------------
-- Table structure for user_group_data_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_group_data_auth`;
CREATE TABLE `user_group_data_auth`
(
    `id`             bigint unsigned NOT NULL COMMENT '自增主键',
    `user_group_id`        bigint NOT NULL COMMENT '用户、组关系表id',
    `resource_identifier`      varchar(25) NOT NULL COMMENT '资源标识',
    `auth`           varchar(25) NOT NULL COMMENT '权限（all、只看自己、看协同、看同级、看同组别、看同部门、看同公司(机构)、看同产品）',
    `status`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态，-2删除，-1禁用，待审核0，启用1',
    `create_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_user_group_id_identifier` (`user_group_id`,`resource_identifier`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `resource_identifier`, `user_group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户、组、数据权限关系表';

-- ----------------------------待办服务中心----------------------------
-- ----------------------------
-- Table structure for to_do
-- ----------------------------
DROP TABLE IF EXISTS `to_do`;
CREATE TABLE `to_do`
(
    `id`                bigint unsigned NOT NULL COMMENT '自增主键',
    `content`           varchar(500) NOT NULL COMMENT '待办内容',
    `create_user`       bigint NOT NULL COMMENT '创建人',
    `execute_users`     char(250) NOT NULL COMMENT '执行者，可多个',
    `join_users`        char(250) NOT NULL COMMENT '参与者，可多个',
    `todo_type`         tinyint(3) NOT NULL DEFAULT '1' COMMENT '待办类型，审批(他人指定)、个人(自己写的)、群公告(他人指定)、项目任务(他人指定)',
    `status`            tinyint(3) NOT NULL DEFAULT '0' COMMENT '待办状态；-2：暂停，-1：删除，0：待启动，1：进行中，2：完成',
    `start_time`        datetime    COMMENT '开始时间',
    `deadline`          datetime    COMMENT '截止时间',
    `remark`            varchar(250) NOT NULL COMMENT '备注,可以分割线的形式添加进度描述',
    `priority`          tinyint(3) NOT NULL DEFAULT '3' COMMENT '优先级；',
    `parent_id`         bigint NOT NULL COMMENT '父级id',
    `attachment_url`    varchar(500) NOT NULL COMMENT '附件(文件中心的下载地址)，可多个',
    `create_time`       datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`       datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `indexFilter_create_time` (`create_time`, `status`, `create_user`, `todo_type`, `priority`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='待办表';

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
    `sender_device_id` varchar(100) NOT NULL COMMENT '发送设备id：多个设备id"，"隔开，*表示全部设备',
    `receiver_id`   bigint NOT NULL COMMENT '接收者id',
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
    `sender_device_id` varchar(100) NOT NULL COMMENT '发送设备id：多个设备id"，"隔开，*表示全部设备',
    `receiver_id`   bigint NOT NULL COMMENT '接收者id, group_id',
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
