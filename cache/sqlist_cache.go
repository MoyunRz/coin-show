package cache

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var CacheClient *Cache

// Cache 数据库连接和缓存数据
type Cache struct {
	conn *sql.DB
}

func init() {
	mkdirSrc(".coin_show")
	ch, err := NewCache(".coin_show/cache.db")
	if err != nil {
		panic(err)
	}
	CacheClient = ch
}

func mkdirSrc(dir string) {
	// 要创建的目录路径

	// 检查目录是否已经存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("创建目录失败：", err)
		} else {
			fmt.Println("创建目录成功：", dir)
		}
	} else if err != nil {
		// 如果出现其他错误，则输出错误信息
		fmt.Println("检查目录时发生错误：", err)
	}
}

// NewCache 返回一个新的缓存对象
func NewCache(dbPath string) (*Cache, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 通过执行测试语句来确保数据库可以正常工作
	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	// 创建缓存表，如果不存在的话
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS cache (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            key TEXT UNIQUE NOT NULL,
            value TEXT NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return nil, err
	}

	return &Cache{conn: db}, nil
}

// Get 从缓存中检索键的值
func (c *Cache) Get(key string) (string, error) {
	var value string
	err := c.conn.QueryRow("SELECT value FROM cache WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", errors.New("key not found")
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

// Set 将键值对保存到缓存中
func (c *Cache) Set(key, value string) error {
	_, err := c.conn.Exec(`
        INSERT OR REPLACE INTO cache (key, value) VALUES (?, ?);
    `, key, value)
	return err
}

// Delete 从缓存中删除一个键值对
func (c *Cache) Delete(key string) error {
	_, err := c.conn.Exec("DELETE FROM cache WHERE key = ?", key)
	return err
}

// Search 返回所有匹配查询字符串的键值对
func (c *Cache) Search(query string) (map[string]string, error) {
	rows, err := c.conn.Query("SELECT key, value FROM cache WHERE key LIKE ?", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		results[key] = value
	}

	return results, nil
}

// Update 更新现有键的值
func (c *Cache) Update(key, value string) error {
	_, err := c.conn.Exec(`
        UPDATE cache SET value = ? WHERE key = ?
    `, value, key)
	return err
}

// Close 关闭数据库连接
func (c *Cache) Close() error {
	return c.conn.Close()
}
