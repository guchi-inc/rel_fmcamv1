CREATE TABLE `fm_user` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
   `username` varchar(255) NOT NULL COMMENT '用户姓名--例如张三',
   `login_name` varchar(255) NOT NULL COMMENT '登录用户名',
   `password` varchar(50) DEFAULT NULL COMMENT '登陆密码',
   `leader_flag` varchar(1) DEFAULT '0' COMMENT '是否经理，0否，1是',
   `position` varchar(200) DEFAULT NULL COMMENT '职位',
   `department` varchar(255) DEFAULT NULL COMMENT '所属部门',
   `email` varchar(100) DEFAULT NULL COMMENT '电子邮箱',
   `phonenum` varchar(100) DEFAULT NULL COMMENT '手机号码',
   `ismanager` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否为管理者 0==管理者 1==员工',
   `isystem` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否系统自带数据 ',
   `Status` tinyint(4) DEFAULT 0 COMMENT '状态，0：正常，1：删除，2封禁',
   `description` varchar(500) DEFAULT NULL COMMENT '用户描述信息',
   `remark` varchar(500) DEFAULT NULL COMMENT '备注',
   `tenant_id` bigint(20) DEFAULT NULL COMMENT '租户id',
   PRIMARY KEY (`id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=146 DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci COMMENT='用户表'