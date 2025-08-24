-- 数据库优化脚本
-- 为统计查询添加必要的索引

-- 任务表索引优化
CREATE INDEX IF NOT EXISTS idx_todos_created_by ON todos(created_by);
CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_created_by_status ON todos(created_by, status);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_completed_at ON todos(completed_at);
CREATE INDEX IF NOT EXISTS idx_todos_created_by_completed_at ON todos(created_by, completed_at);
CREATE INDEX IF NOT EXISTS idx_todos_created_by_due_date_status ON todos(created_by, due_date, status);

-- 文章表索引优化
CREATE INDEX IF NOT EXISTS idx_articles_created_by ON articles(created_by);
CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);
CREATE INDEX IF NOT EXISTS idx_articles_created_by_status ON articles(created_by, status);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at);
CREATE INDEX IF NOT EXISTS idx_articles_created_by_created_at ON articles(created_by, created_at);
CREATE INDEX IF NOT EXISTS idx_articles_view_count ON articles(view_count);
CREATE INDEX IF NOT EXISTS idx_articles_like_count ON articles(like_count);

-- 通知表索引优化
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id_is_read ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- 用户表索引优化
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 分析表统计信息
ANALYZE todos;
ANALYZE articles;
ANALYZE notifications;
ANALYZE users;
