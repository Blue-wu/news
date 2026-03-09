-- 首先创建数据库（如果还未创建）
CREATE DATABASE IF NOT EXISTS blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE blog_db;

-- 用户表（用于文章作者管理）
CREATE TABLE IF NOT EXISTS users (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
                                     username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密存储）',
    nickname VARCHAR(50) COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像URL',
    bio TEXT COMMENT '个人简介',
    role VARCHAR(20) DEFAULT 'user' COMMENT '角色（admin/user）',
    status TINYINT DEFAULT 1 COMMENT '状态（1-活跃 0-禁用）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME COMMENT '删除时间'
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章表
CREATE TABLE IF NOT EXISTS articles (
                                        id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '文章ID',
                                        user_id BIGINT NOT NULL COMMENT '作者ID',
                                        title VARCHAR(200) NOT NULL COMMENT '标题',
    content LONGTEXT NOT NULL COMMENT '内容',
    summary VARCHAR(500) COMMENT '摘要',
    cover_image VARCHAR(255) COMMENT '封面图片URL',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    comment_count INT DEFAULT 0 COMMENT '评论次数',
    like_count INT DEFAULT 0 COMMENT '点赞次数',
    status TINYINT DEFAULT 1 COMMENT '状态（1-发布 0-草稿）',
    is_top TINYINT DEFAULT 0 COMMENT '是否置顶（1-置顶 0-普通）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME COMMENT '删除时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章分类表
CREATE TABLE IF NOT EXISTS categories (
                                          id INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
                                          name VARCHAR(50) NOT NULL UNIQUE COMMENT '分类名称',
    slug VARCHAR(50) NOT NULL UNIQUE COMMENT '分类别名（URL友好）',
    description TEXT COMMENT '分类描述',
    parent_id INT DEFAULT 0 COMMENT '父分类ID（0表示顶级分类）',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
                                    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
                                    name VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    slug VARCHAR(50) NOT NULL UNIQUE COMMENT '标签别名（URL友好）',
    article_count INT DEFAULT 0 COMMENT '关联文章数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章分类关联表
CREATE TABLE IF NOT EXISTS article_categories (
                                                  article_id BIGINT NOT NULL COMMENT '文章ID',
                                                  category_id INT NOT NULL COMMENT '分类ID',
                                                  PRIMARY KEY (article_id, category_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章标签关联表
CREATE TABLE IF NOT EXISTS article_tags (
                                            article_id BIGINT NOT NULL COMMENT '文章ID',
                                            tag_id INT NOT NULL COMMENT '标签ID',
                                            PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 评论表
CREATE TABLE IF NOT EXISTS comments (
                                        id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '评论ID',
                                        article_id BIGINT NOT NULL COMMENT '文章ID',
                                        user_id BIGINT COMMENT '用户ID（匿名评论可为空）',
                                        parent_id BIGINT DEFAULT 0 COMMENT '父评论ID（0表示顶级评论）',
                                        nickname VARCHAR(50) COMMENT '昵称（匿名评论使用）',
    email VARCHAR(100) COMMENT '邮箱（匿名评论使用）',
    content TEXT NOT NULL COMMENT '评论内容',
    status TINYINT DEFAULT 1 COMMENT '状态（1-审核通过 0-待审核）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 初始化默认数据（可选）
-- 1. 创建默认管理员用户
INSERT INTO users (username, email, password, nickname, role)
VALUES ('admin', 'admin@example.com', 'admin123456', '管理员', 'admin');

-- 2. 创建默认分类
INSERT INTO categories (name, slug, description)
VALUES
    ('技术', 'tech', '技术相关文章'),
    ('生活', 'life', '生活相关文章'),
    ('分享', 'share', '经验分享文章');

-- 3. 创建默认标签
INSERT INTO tags (name, slug)
VALUES
    ('Go', 'go'),
    ('Vue', 'vue'),
    ('MySQL', 'mysql'),
    ('前端', 'frontend'),
    ('后端', 'backend');






-- 在现有文件末尾添加以下内容

-- 文章点赞表
CREATE TABLE IF NOT EXISTS article_likes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '点赞ID',
    article_id BIGINT NOT NULL COMMENT '文章ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    UNIQUE KEY uk_article_user (article_id, user_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户关注表
CREATE TABLE IF NOT EXISTS user_follows (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关注ID',
    follower_id BIGINT NOT NULL COMMENT '关注者ID',
    followed_id BIGINT NOT NULL COMMENT '被关注者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间',
    UNIQUE KEY uk_follower_followed (follower_id, followed_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文章收藏表
CREATE TABLE IF NOT EXISTS article_collections (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '收藏ID',
    article_id BIGINT NOT NULL COMMENT '文章ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    UNIQUE KEY uk_article_user (article_id, user_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 网站设置表
CREATE TABLE IF NOT EXISTS site_settings (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '设置ID',
    key_name VARCHAR(50) NOT NULL UNIQUE COMMENT '设置键名',
    value TEXT COMMENT '设置值',
    description VARCHAR(255) COMMENT '设置描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '通知ID',
    user_id BIGINT NOT NULL COMMENT '接收通知的用户ID',
    sender_id BIGINT COMMENT '发送通知的用户ID',
    type VARCHAR(20) NOT NULL COMMENT '通知类型（comment/like/follow/system）',
    content TEXT NOT NULL COMMENT '通知内容',
    related_id BIGINT COMMENT '相关实体ID（文章ID/评论ID等）',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读（1-已读 0-未读）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统日志表
CREATE TABLE IF NOT EXISTS system_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    user_id BIGINT COMMENT '操作用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    target_type VARCHAR(50) COMMENT '操作目标类型',
    target_id BIGINT COMMENT '操作目标ID',
    details TEXT COMMENT '操作详情',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 初始化网站设置数据
INSERT INTO site_settings (key_name, value, description)
VALUES
    ('site_name', '博客系统', '网站名称'),
    ('site_description', '一个功能完善的博客系统', '网站描述'),
    ('site_keywords', '博客,技术,分享', '网站关键词'),
    ('site_logo', '', '网站Logo'),
    ('site_favicon', '', '网站图标'),
    ('homepage_articles_count', '10', '首页显示文章数量'),
    ('pagination_size', '10', '分页大小'),
    ('comment_audit', '0', '评论是否需要审核（1-需要 0-不需要）'),
    ('allow_anonymous_comment', '1', '是否允许匿名评论（1-允许 0-不允许）'),
    ('copyright_info', '© 2024 博客系统 版权所有', '版权信息'),
    ('sitemap_settings', '{"cacheEnabled":true,"cacheDuration":86400}', '站点地图配置');

-- 为现有的表添加必要的索引以提高查询性能
-- 为文章表添加索引
ALTER TABLE articles ADD INDEX idx_user_id (user_id);
ALTER TABLE articles ADD INDEX idx_status (status);
ALTER TABLE articles ADD INDEX idx_is_top (is_top);
ALTER TABLE articles ADD INDEX idx_created_at (created_at);

-- 为评论表添加索引
ALTER TABLE comments ADD INDEX idx_article_id (article_id);
ALTER TABLE comments ADD INDEX idx_user_id (user_id);
ALTER TABLE comments ADD INDEX idx_parent_id (parent_id);
ALTER TABLE comments ADD INDEX idx_status (status);

-- 为通知表添加索引
ALTER TABLE notifications ADD INDEX idx_user_id (user_id);
ALTER TABLE notifications ADD INDEX idx_is_read (is_read);
ALTER TABLE notifications ADD INDEX idx_created_at (created_at);

-- 更新分类表和标签表的外键引用（如果需要）
-- 注意：如果已经存在外键约束，这些语句可能会失败
-- 可以先删除现有的外键约束再添加，或者忽略这些语句