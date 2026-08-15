// Rabbit 流密码移植，语义对齐 crypto-js 的 Rabbit 实现
// （https://github.com/brix/crypto-js 的 src/rabbit.js）。
// 算法：128 位密钥 + 64 位 IV，每次状态迭代产生 16 字节密钥流，明文与密钥流异或。
package encryption

import (
	"crypto/cipher"
	"fmt"
)

// rabbitState 保存 Rabbit 内部状态。
type rabbitState struct {
	x [8]uint32
	c [8]uint32
	b uint32
}

// newRabbitState 以密钥与 IV 初始化 Rabbit 状态。
func newRabbitState(key, iv []byte) (*rabbitState, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("Rabbit 密钥长度必须为 16 字节，实际 %d", len(key))
	}
	if len(iv) != 8 {
		return nil, fmt.Errorf("Rabbit IV 长度必须为 8 字节，实际 %d", len(iv))
	}

	k := [4]uint32{}
	for i := 0; i < 4; i++ {
		// crypto-js 在 _doReset 先对 K 字做字节序交换（等效小端读取）。
		k[i] = swapEndian(beUint32(key[i*4 : i*4+4]))
	}

	r := &rabbitState{}
	// 初始状态值（对齐 crypto-js _doReset）。
	r.x = [8]uint32{
		k[0], (k[3] << 16) | (k[2] >> 16),
		k[1], (k[0] << 16) | (k[3] >> 16),
		k[2], (k[1] << 16) | (k[0] >> 16),
		k[3], (k[2] << 16) | (k[1] >> 16),
	}
	// 初始计数器值。
	r.c = [8]uint32{
		(k[2] << 16) | (k[2] >> 16), (k[0] & 0xffff0000) | (k[1] & 0x0000ffff),
		(k[3] << 16) | (k[3] >> 16), (k[1] & 0xffff0000) | (k[2] & 0x0000ffff),
		(k[0] << 16) | (k[0] >> 16), (k[2] & 0xffff0000) | (k[3] & 0x0000ffff),
		(k[1] << 16) | (k[1] >> 16), (k[3] & 0xffff0000) | (k[0] & 0x0000ffff),
	}
	r.b = 0

	// 迭代系统四次。
	for i := 0; i < 4; i++ {
		r.nextState()
	}

	// 修改计数器。
	for i := 0; i < 8; i++ {
		r.c[i] ^= r.x[(i+4)&7]
	}

	// IV 设置（对齐 crypto-js _doReset 的 IV 分支）。
	iv0 := swapEndian(beUint32(iv[0:4]))
	iv1 := swapEndian(beUint32(iv[4:8]))
	i0 := iv0
	i2 := iv1
	i1 := (i0 >> 16) | (i2 & 0xffff0000)
	i3 := (i2 << 16) | (i0 & 0x0000ffff)
	r.c[0] ^= i0
	r.c[1] ^= i1
	r.c[2] ^= i2
	r.c[3] ^= i3
	r.c[4] ^= i0
	r.c[5] ^= i1
	r.c[6] ^= i2
	r.c[7] ^= i3
	for i := 0; i < 4; i++ {
		r.nextState()
	}
	return r, nil
}

// nextState 迭代一次 Rabbit 状态，产生 16 字节密钥流写入 dst。
func (r *rabbitState) nextState() [16]byte {
	var old [8]uint32
	copy(old[:], r.c[:])

	// 计数器更新（无符号比较产生进位）。
	r.c[0] += 0x4d34d34d + r.b
	r.c[1] += 0xd34d34d3 + b2u(r.c[0] < old[0])
	r.c[2] += 0x34d34d34 + b2u(r.c[1] < old[1])
	r.c[3] += 0x4d34d34d + b2u(r.c[2] < old[2])
	r.c[4] += 0xd34d34d3 + b2u(r.c[3] < old[3])
	r.c[5] += 0x34d34d34 + b2u(r.c[4] < old[4])
	r.c[6] += 0x4d34d34d + b2u(r.c[5] < old[5])
	r.c[7] += 0xd34d34d3 + b2u(r.c[6] < old[6])
	r.b = b2u(r.c[7] < old[7])

	// g 值：gh = x² 的高 32 位，gl = x² 的低 32 位，G = gh ^ gl。
	var g [8]uint32
	for i := 0; i < 8; i++ {
		sq := uint64(r.x[i]+r.c[i]) * uint64(r.x[i]+r.c[i])
		g[i] = uint32(sq>>32) ^ uint32(sq&0xffffffff)
	}

	// 新状态值（对齐 crypto-js nextState 的 X 更新）。
	r.x[0] = g[0] + rot16(g[7]) + rot16(g[6])
	r.x[1] = g[1] + rot8(g[0]) + g[7]
	r.x[2] = g[2] + rot16(g[1]) + rot16(g[0])
	r.x[3] = g[3] + rot8(g[2]) + g[1]
	r.x[4] = g[4] + rot16(g[3]) + rot16(g[2])
	r.x[5] = g[5] + rot8(g[4]) + g[3]
	r.x[6] = g[6] + rot16(g[5]) + rot16(g[4])
	r.x[7] = g[7] + rot8(g[6]) + g[5]

	// 生成四个密钥流字（对齐 crypto-js _doProcessBlock）。
	var ks [4]uint32
	ks[0] = r.x[0] ^ (r.x[5] >> 16) ^ (r.x[3] << 16)
	ks[1] = r.x[2] ^ (r.x[7] >> 16) ^ (r.x[5] << 16)
	ks[2] = r.x[4] ^ (r.x[1] >> 16) ^ (r.x[7] << 16)
	ks[3] = r.x[6] ^ (r.x[3] >> 16) ^ (r.x[1] << 16)

	var out [16]byte
	for i := 0; i < 4; i++ {
		w := swapEndian(ks[i])
		out[i*4] = byte(w >> 24)
		out[i*4+1] = byte(w >> 16)
		out[i*4+2] = byte(w >> 8)
		out[i*4+3] = byte(w)
	}
	return out
}

// rabbitStream 包装 Rabbit 状态实现 cipher.Stream。
type rabbitStream struct {
	state *rabbitState
}

// newRabbitCipher 创建 Rabbit 流密码。
func newRabbitCipher(key, iv []byte) (cipher.Stream, error) {
	s, err := newRabbitState(key, iv)
	if err != nil {
		return nil, err
	}
	return &rabbitStream{state: s}, nil
}

// XORKeyStream 用 Rabbit 密钥流异或 src 写入 dst。
func (r *rabbitStream) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("cipher: output smaller than input")
	}
	for len(src) >= 16 {
		block := r.state.nextState()
		for i := 0; i < 16; i++ {
			dst[i] = src[i] ^ block[i]
		}
		dst = dst[16:]
		src = src[16:]
	}
	if len(src) > 0 {
		block := r.state.nextState()
		for i := 0; i < len(src); i++ {
			dst[i] = src[i] ^ block[i]
		}
	}
}

// beUint32 读取大端 32 位字（crypto-js WordArray 字节序）。
func beUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// swapEndian 交换 32 位字的字节序（对齐 crypto-js 的 "Swap endian" 操作）。
func swapEndian(w uint32) uint32 {
	return (w<<24)&0xff000000 | (w<<8)&0x00ff0000 | (w>>8)&0x0000ff00 | (w>>24)&0x000000ff
}

// rot16 循环左移 16 位（等价 crypto-js 的 (x<<16)|(x>>>16)）。
func rot16(x uint32) uint32 {
	return (x << 16) | (x >> 16)
}

// rot8 循环左移 8 位（等价 crypto-js 的 (x<<8)|(x>>>24)）。
func rot8(x uint32) uint32 {
	return (x << 8) | (x >> 24)
}

// b2u 将布尔值转为 uint32 用于进位。
func b2u(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
