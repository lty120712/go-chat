DELIMITER $$

CREATE PROCEDURE GenerateTestUsers()
BEGIN
    DECLARE i INT DEFAULT 0;

    WHILE i < 1000 DO
            INSERT INTO users (
                username,           -- 用户名
                password,           -- 密码
                nickname,           -- 昵称
                `desc`,             -- 简介（使用反引号转义）
                phone,              -- 电话
                email,              -- 邮箱
                avatar,             -- 头像URL
                client_ip,          -- 客户端IP
                client_port,        -- 客户端端口
                login_time,         -- 最近登录时间
                heartbeat_time,     -- 心跳时间
                logout_time,        -- 登出时间
                status,             -- 用户状态
                online_status,      -- 用户在线状态
                device_info,        -- 客户端设备信息
                created_at,         -- 创建时间
                updated_at          -- 更新时间
            )
            VALUES (
                       CONCAT('user', i),                                  -- username: user0, user1, ...
                       CONCAT('password', i),                              -- password: password0, password1, ...
                       CONCAT('nickname', i),                              -- nickname: nickname0, nickname1, ...
                       CONCAT('This is user ', i),                         -- `desc`: 简单的描述信息
                       CONCAT('100000000', LPAD(i, 3, '0')),               -- phone: 100000000000, 100000000001, ...
                       CONCAT('user', i, '@test.com'),                     -- email: user0@test.com, user1@test.com, ...
                       CONCAT('https://example.com/avatar', i, '.jpg'),    -- avatar: 头像URL（这里使用一个示例）
                       '192.168.0.1',                                      -- client_ip: 假设一个固定的客户端IP
                       '8080',                                              -- client_port: 假设固定端口号
                       UNIX_TIMESTAMP(),                                   -- login_time: 当前时间戳
                       UNIX_TIMESTAMP(),                                   -- heartbeat_time: 当前时间戳
                       NULL,                                                -- logout_time: 初始没有登出时间
                       1,                                                   -- status: 用户有效
                       1,                                                   -- online_status: 用户在线
                       'Windows 10, Chrome',                               -- device_info: 假设的设备信息
                       NOW(),                                              -- created_at: 当前时间
                       NOW()                                               -- updated_at: 当前时间
                   );

            SET i = i + 1;
        END WHILE;
END $$

DELIMITER ;

-- 执行存储过程以生成1000个用户
CALL GenerateTestUsers();
