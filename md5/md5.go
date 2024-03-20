package md5

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type md5Generator struct {
	str string
}

// Generator ID生成器
type Generator interface {
	// MD5 获取字符串的MD5值，32位小写
	MD5() string
	// MD5Upper 获取字符串大写Md5值， 32位大写
	MD5Upper() string
	// MD5With16 获取Md5，16位小写
	MD5With16() string
	// MD5With16Upper 获取Md5，16位大写
	MD5With16Upper() string
}

func NewMD5Generator(str string) Generator {
	return &md5Generator{str: str}
}

func (g *md5Generator) MD5() string {
	h := md5.New()
	h.Write([]byte(g.str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func (g *md5Generator) MD5Upper() string {
	return strings.ToUpper(g.MD5())
}

func (g *md5Generator) MD5With16() string {
	return g.MD5()[8:24]
}

func (g *md5Generator) MD5With16Upper() string {
	return strings.ToUpper(g.MD5With16())
}
