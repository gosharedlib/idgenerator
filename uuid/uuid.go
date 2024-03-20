package uuid

import (
	guid "github.com/gofrs/uuid"
)

// https://datatracker.ietf.org/doc/html/rfc4122

type v1Generator struct {
}

type v3Generator struct {
	ns   guid.UUID
	name string
}

type v4Generator struct {
}

type v5Generator struct {
	ns   guid.UUID
	name string
}

// Generator ID生成器
type Generator interface {
	// GenID 生成字符串 Key.
	GenID() string
}

func NewV1Generator() Generator {
	return &v1Generator{}
}

// 由于安全缺陷，上游依赖已移除 V2 实现
// func NewV2Generator(domain byte) Generator {
// 	return &v2Generator{domain: domain}
// }

func NewV3Generator(ns guid.UUID, name string) Generator {
	return &v3Generator{ns: ns, name: name}
}

func NewV4Generator() Generator {
	return &v4Generator{}
}

func NewV5Generator(ns guid.UUID, name string) Generator {
	return &v5Generator{ns: ns, name: name}
}

func (g *v1Generator) GenID() string {
	return guid.Must(guid.NewV1()).String()
}

func (g *v3Generator) GenID() string {
	return guid.NewV3(g.ns, g.name).String()
}

func (g *v4Generator) GenID() string {
	return guid.Must(guid.NewV4()).String()
}

func (g *v5Generator) GenID() string {
	return guid.NewV5(g.ns, g.name).String()
}
