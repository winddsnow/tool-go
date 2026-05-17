// Package password 提供密码加密与验证功能。
// 采用 MD5 + 随机盐（Salt）的哈希方案对密码进行单向加密。
//
// 为什么需要盐（Salt）？
//   盐是在哈希前附加到密码上的随机字符串，即使两个用户密码相同，由于盐不同，
//   哈希结果也不同。这可以有效防御两种攻击：
//     - 彩虹表攻击：攻击者预计算常见密码的哈希值。加了盐后，攻击者必须为每个盐
//       值分别构建彩虹表，成本极高。
//     - 相同密码识别：如果没有盐，相同密码产生相同哈希值，攻击者可轻易识别出
//       使用相同密码的用户。
package password

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// saltLength 定义生成随机盐的字节长度。
// 16 字节 = 128 位，编码为十六进制字符串后长度为 32 个字符。
// 128 位的随机盐具有足够的熵（2^128 种可能），在计算上不可穷举。
const saltLength = 16

// GenerateSalt 生成指定长度的随机盐，返回十六进制编码字符串。
// 使用 crypto/rand 包而非 math/rand 包的原因：
//   - crypto/rand 是密码学安全的伪随机数生成器（CSPRNG），其输出不可预测，
//     适用于密钥、盐值等安全敏感场景。
//   - math/rand 是普通的伪随机数生成器，基于确定性算法，可被预测，
//     仅适用于模拟、游戏等非安全场景。
//
// io.ReadFull 确保从 rand.Reader 中读取到指定数量的字节。
// 通常情况下 rand.Reader 返回完整数据，但在极端环境（如系统熵池不足）下
// 可能部分读取，ReadFull 可保证读满 len(bytes) 字节，否则返回错误。
//
// hex.EncodeToString 将二进制字节转换为十六进制字符串，
// 方便在数据库中以字符串形式存储，也便于日志打印和调试。
func GenerateSalt() (string, error) {
	bytes := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword 使用 MD5 算法对 "salt + password" 拼接字符串进行哈希。
// MD5 是一种单向哈希函数，将任意长度的输入计算为固定 128 位（16 字节）的输出摘要。
// 编码为十六进制后长度为 32 个字符。
// 单向性意味着从哈希值无法逆向推出原始密码——这是密码存储的基本原则。
func HashPassword(password string, salt string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", salt, password)))
	return hex.EncodeToString(h.Sum(nil))
}

// CreatePassword 创建密码：生成随机盐 → 计算哈希 → 返回哈希值和盐值。
// 返回的 hash 和 salt 应一起存入数据库，验证时通过 VerifyPassword 重新计算比对。
// Go 支持函数返回多个值，这里返回三个值：哈希值、盐值、错误信息。
func CreatePassword(password string) (hash string, salt string, err error) {
	salt, err = GenerateSalt()
	if err != nil {
		return "", "", err
	}
	hash = HashPassword(password, salt)
	return hash, salt, nil
}

// VerifyPassword 验证密码：使用相同的盐重新计算哈希值，与数据库中存储的哈希值比对。
// 若计算出的哈希值和存储的哈希值相等，说明密码正确，返回 true；否则返回 false。
// 注意：每次验证只应比较哈希值，绝不能解密或存储明文密码。
func VerifyPassword(password, salt, hash string) bool {
	return HashPassword(password, salt) == hash
}
