// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

// Protocol represents a P2P subprotocol implementation.
// Protocol 表示一个 P2P 子协议实现。
type Protocol struct {
	// Name should contain the official protocol name,
	// often a three-letter word.
	// 名称应包含官方协议名称，
	// 通常是三个字母的单词。
	Name string

	// Version should contain the version number of the protocol.
	// 版本应该包含协议的版本号。
	Version uint

	// Length should contain the number of message codes used
	// by the protocol.
	// 长度应该包含协议使用的消息代码的数量。
	Length uint64

	// Run is called in a new goroutine when the protocol has been
	// negotiated with a peer. It should read and write messages from
	// rw. The Payload for each message must be fully consumed.
	// 当与对等方协商协议时，在新的 goroutine 中调用 Run。
	// 它应该从 rw 读取和写入消息。每条消息的有效负载必须被完全消耗。

	// The peer connection is closed when Start returns. It should return
	// any protocol-level error (such as an I/O error) that is
	// encountered.
	// 当 Start 返回时，对等连接关闭。它应该返回遇到的任何协议级错误（例如 I/O 错误）。
	Run func(peer *Peer, rw MsgReadWriter) error

	// NodeInfo is an optional helper method to retrieve protocol specific metadata
	// about the host node.
	// NodeInfo 是一个可选的辅助方法，用于检索有关主机节点的协议特定元数据。
	NodeInfo func() interface{}

	// PeerInfo is an optional helper method to retrieve protocol specific metadata
	// about a certain peer in the network. If an info retrieval function is set,
	// but returns nil, it is assumed that the protocol handshake is still running.
	// PeerInfo 是一个可选的辅助方法，用于检索有关网络中某个对等点的协议特定元数据。
	// 如果设置了信息检索函数，但返回 nil，则假定协议握手仍在运行。
	PeerInfo func(id enode.ID) interface{}

	// DialCandidates, if non-nil, is a way to tell Server about protocol-specific nodes
	// that should be dialed. The server continuously reads nodes from the iterator and
	// attempts to create connections to them.
	// DialCandidates，如果非零，是一种告诉服务器应该拨号的协议特定节点的方法。
	// 服务器不断地从迭代器中读取节点并尝试创建到它们的连接。
	DialCandidates enode.Iterator

	// Attributes contains protocol specific information for the node record.
	// 属性包含节点记录的协议特定信息。
	Attributes []enr.Entry
}

func (p Protocol) cap() Cap {
	return Cap{p.Name, p.Version}
}

// Cap is the structure of a peer capability.
// Cap是peer能力的结构。
type Cap struct {
	Name    string
	Version uint
}

func (cap Cap) String() string {
	return fmt.Sprintf("%s/%d", cap.Name, cap.Version)
}

type capsByNameAndVersion []Cap

func (cs capsByNameAndVersion) Len() int      { return len(cs) }
func (cs capsByNameAndVersion) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs capsByNameAndVersion) Less(i, j int) bool {
	return cs[i].Name < cs[j].Name || (cs[i].Name == cs[j].Name && cs[i].Version < cs[j].Version)
}
